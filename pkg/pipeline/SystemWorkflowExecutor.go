package pipeline

import (
	"context"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/pipeline/bean"
	"go.uber.org/zap"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type SystemWorkflowExecutor interface {
	WorkflowExecutor
}

type SystemWorkflowExecutorImpl struct {
	logger  *zap.SugaredLogger
	k8sUtil *util.K8sUtil
}

func NewSystemWorkflowExecutorImpl(logger *zap.SugaredLogger, k8sUtil *util.K8sUtil) *SystemWorkflowExecutorImpl {
	return &SystemWorkflowExecutorImpl{logger: logger, k8sUtil: k8sUtil}
}

func (impl *SystemWorkflowExecutorImpl) ExecuteWorkflow(workflowTemplate bean.WorkflowTemplate) error {
	//create job template with suspended state
	jobTemplate := impl.getJobTemplate(workflowTemplate)
	clientset, err := impl.k8sUtil.GetClientSetForConfig(workflowTemplate.ClusterConfig)
	if err != nil {
		impl.logger.Errorw("error occurred while creating k8s client", "WorkflowRunnerId", workflowTemplate.WorkflowRunnerId, "err", err)
		return err
	}
	ctx := context.Background()
	createdJob, err := clientset.BatchV1().Jobs(workflowTemplate.Namespace).Create(ctx, jobTemplate, v12.CreateOptions{})
	if err != nil {
		impl.logger.Errorw("error occurred while creating k8s job", "WorkflowRunnerId", workflowTemplate.WorkflowRunnerId, "err", err)
		return err
	}

	//create cm and secrets with owner reference
	err = impl.createCmAndSecrets(workflowTemplate, createdJob)
	if err != nil {
		impl.logger.Errorw("error occurred while creating cm and secret", "WorkflowRunnerId", workflowTemplate.WorkflowRunnerId, "err", err)
		return err
	}

	//change job state to running
	_, err = clientset.BatchV1().Jobs(workflowTemplate.Namespace).Patch(ctx, createdJob.Name, types.StrategicMergePatchType, []byte(`{"spec":{"suspend": false}}`), v12.PatchOptions{})
	if err != nil {
		impl.logger.Errorw("error occurred while updating job suspended status", "WorkflowRunnerId", workflowTemplate.WorkflowRunnerId, "err", err)
		return err
	}
	return nil
}

func (impl *SystemWorkflowExecutorImpl) getJobTemplate(workflowTemplate bean.WorkflowTemplate) *v1.Job {

	workflowJob := v1.Job{
		ObjectMeta: v12.ObjectMeta{
			GenerateName: workflowTemplate.WorkflowNamePrefix + "-",
			//Annotations:  map[string]string{"workflows.argoproj.io/controller-instanceid": workflowTemplate.WfControllerInstanceID},
			Labels:     map[string]string{"devtron.ai/workflow-purpose": "cd"},
			Finalizers: []string{"foregroundDeletion"},
		},
		Spec: v1.JobSpec{
			ActiveDeadlineSeconds:   workflowTemplate.ActiveDeadlineSeconds,
			TTLSecondsAfterFinished: workflowTemplate.TTLValue,
			Template: corev1.PodTemplateSpec{
				Spec: workflowTemplate.PodSpec,
			},
			Suspend: &[]bool{true}[0],
		},
	}
	return &workflowJob
}

func (impl *SystemWorkflowExecutorImpl) getCmAndSecrets(workflowTemplate bean.WorkflowTemplate, createdJob *v1.Job) ([]corev1.ConfigMap, []corev1.Secret, error) {
	var configMaps []corev1.ConfigMap
	var secrets []corev1.Secret
	configMapDataArray := workflowTemplate.ConfigMaps
	for _, configSecretMap := range configMapDataArray {
		configDataMap, err := configSecretMap.GetDataMap()
		if err != nil {
			impl.logger.Errorw("error occurred while extracting data map", "Data", configSecretMap.Data, "err", err)
			return configMaps, secrets, err
		}
		configMapSecretDto := ConfigMapSecretDto{Name: configSecretMap.Name, Data: configDataMap, OwnerRef: impl.createJobOwnerRefVal(createdJob)}
		configMap := GetConfigMapBody(configMapSecretDto)
		configMaps = append(configMaps, configMap)
	}
	secretMaps := workflowTemplate.Secrets
	for _, secretMapData := range secretMaps {
		dataMap, err := secretMapData.GetDataMap()
		if err != nil {
			impl.logger.Errorw("error occurred while extracting data map", "Data", secretMapData.Data, "err", err)
			return configMaps, secrets, err
		}
		configMapSecretDto := ConfigMapSecretDto{Name: secretMapData.Name, Data: dataMap, OwnerRef: impl.createJobOwnerRefVal(createdJob)}
		secretBody := GetSecretBody(configMapSecretDto)
		secrets = append(secrets, secretBody)
	}
	return configMaps, secrets, nil
}

func (impl *SystemWorkflowExecutorImpl) createJobOwnerRefVal(createdJob *v1.Job) v12.OwnerReference {
	return v12.OwnerReference{UID: createdJob.UID, Name: createdJob.Name, Kind: kube.JobKind, APIVersion: "batch/v1", BlockOwnerDeletion: &[]bool{true}[0], Controller: &[]bool{true}[0]}
}

func (impl *SystemWorkflowExecutorImpl) createCmAndSecrets(template bean.WorkflowTemplate, createdJob *v1.Job) error {
	client, err := impl.k8sUtil.GetK8sClientForConfig(template.ClusterConfig)
	if err != nil {
		impl.logger.Errorw("error occurred while creating k8s client", "WorkflowRunnerId", template.WorkflowRunnerId, "err", err)
		return err
	}
	configMaps, secrets, err := impl.getCmAndSecrets(template, createdJob)
	if err != nil {
		return err
	}
	for _, configMap := range configMaps {
		_, err = impl.k8sUtil.CreateConfigMap(createdJob.Namespace, &configMap, client)
		if err != nil {
			impl.logger.Errorw("error occurred while creating cm, but ignoring", "err", err)
		}
	}
	for _, secret := range secrets {
		_, err = impl.k8sUtil.CreateSecretData(createdJob.Namespace, &secret, client)
		if err != nil {
			impl.logger.Errorw("error occurred while creating secret, but ignoring", "err", err)
		}
	}
	return nil
}
