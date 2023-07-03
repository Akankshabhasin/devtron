package pipeline

//
//import (
//	"fmt"
//	"github.com/caarlos0/env"
//	"github.com/devtron-labs/devtron/client/gitSensor"
//	"github.com/devtron-labs/devtron/internal/sql/repository"
//	"github.com/devtron-labs/devtron/internal/sql/repository/app"
//	"github.com/devtron-labs/devtron/internal/sql/repository/appWorkflow"
//	repository2 "github.com/devtron-labs/devtron/internal/sql/repository/dockerRegistry"
//	"github.com/devtron-labs/devtron/internal/sql/repository/helper"
//	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
//	"github.com/devtron-labs/devtron/internal/util"
//	app2 "github.com/devtron-labs/devtron/pkg/app"
//	repository3 "github.com/devtron-labs/devtron/pkg/cluster/repository"
//	repository4 "github.com/devtron-labs/devtron/pkg/pipeline/history/repository"
//	"github.com/go-pg/pg"
//	"log"
//
//	"github.com/devtron-labs/devtron/pkg/attributes"
//	"github.com/devtron-labs/devtron/pkg/bean"
//	"github.com/devtron-labs/devtron/pkg/pipeline/history"
//	"github.com/devtron-labs/devtron/pkg/user"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//type Config struct {
//	Addr     string `env:"TEST_PG_ADDR" envDefault:"127.0.0.1"`
//	Port     string `env:"TEST_PG_PORT" envDefault:"5432"`
//	User     string `env:"TEST_PG_USER" envDefault:"postgres"`
//	Password string `env:"TEST_PG_PASSWORD" envDefault:"pass" secretData:"-"`
//	Database string `env:"TEST_PG_DATABASE" envDefault:"orchestrator"`
//	LogQuery bool   `env:"TEST_PG_LOG_QUERY" envDefault:"true"`
//}
//
//var (
//	db                       *pg.DB
//	ciCdPipelineOrchestrator *CiCdPipelineOrchestratorImpl
//)
//
//func getDbConn() (*pg.DB, error) {
//	if db != nil {
//		return db, nil
//	}
//	cfg := Config{}
//	err := env.Parse(&cfg)
//	if err != nil {
//		return nil, err
//	}
//	options := pg.Options{
//		Addr:     cfg.Addr + ":" + cfg.Port,
//		User:     cfg.User,
//		Password: cfg.Password,
//		Database: cfg.Database,
//	}
//	db = pg.Connect(&options)
//	return db, nil
//}
//func TestCiCdPipelineOrchestratorImpl_CreateCiConf(t *testing.T) {
//	InitClusterNoteService()
//	type args struct {
//		createRequest *bean.CiConfigRequest
//		templateId    int
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *bean.CiConfigRequest
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{
//			name: "CreateCiConf success",
//			args: args{
//				createRequest: &bean.CiConfigRequest{},
//				templateId:    123,
//			},
//			want:    &bean.CiConfigRequest{},
//			wantErr: assert.NoError,
//		},
//		{
//			name: "CreateCiConf success with payload",
//			args: args{
//				createRequest: &bean.CiConfigRequest{
//					Id:    12,
//					AppId: 20,
//				},
//				templateId: 123,
//			},
//			want: &bean.CiConfigRequest{
//				Id:    12,
//				AppId: 20,
//			},
//			wantErr: assert.NoError,
//		},
//		{
//			name: "CreateCiConf success with job payload",
//			args: args{
//				createRequest: &bean.CiConfigRequest{
//					Id:    0,
//					AppId: 21,
//					IsJob: true,
//					CiPipelines: []*bean.CiPipeline{{
//						IsExternal: true,
//						IsManual:   false,
//						AppId:      21,
//						CiMaterial: []*bean.CiMaterial{
//							{
//								Source: &bean.SourceTypeConfig{
//									Type:  "SOURCE_TYPE_BRANCH_FIXED",
//									Value: "main",
//									Regex: "",
//								},
//								GitMaterialId:   13,
//								Id:              38,
//								GitMaterialName: "devtron-test",
//								IsRegex:         false,
//							},
//						},
//						EnvironmentId: 3,
//					}},
//					UserId: 1,
//				},
//				templateId: 13,
//			},
//			want: &bean.CiConfigRequest{
//				Id:          0,
//				AppId:       21,
//				CiPipelines: []*bean.CiPipeline{{EnvironmentId: 3}},
//			},
//			wantErr: assert.NoError,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := ciCdPipelineOrchestrator.CreateCiConf(tt.args.createRequest, tt.args.templateId)
//			if !tt.wantErr(t, err, fmt.Sprintf("CreateCiConf(%v, %v)", tt.args.createRequest, tt.args.templateId)) {
//				return
//			}
//			assert.Equal(t, tt.want.AppId, got.AppId, "CreateCiConf(%v, %v)", tt.args.createRequest, tt.args.templateId)
//			assert.Equal(t, tt.want.Id, got.Id, "CreateCiConf(%v, %v)", tt.args.createRequest, tt.args.templateId)
//			assert.Equal(t, tt.want.CiPipelines[0].EnvironmentId, got.CiPipelines[0].EnvironmentId, "CreateCiConf(%v, %v)", tt.args.createRequest, tt.args.templateId)
//		})
//	}
//}
//
//func InitClusterNoteService() {
//	if ciCdPipelineOrchestrator != nil {
//		return
//	}
//	logger, err := util.NewSugardLogger()
//	if err != nil {
//		log.Fatalf("error in logger initialization %s,%s", "err", err)
//	}
//	conn, err := getDbConn()
//	if err != nil {
//		log.Fatalf("error in db connection initialization %s, %s", "err", err)
//	}
//
//	appRepository := app.NewAppRepositoryImpl(conn, logger)
//	materialRepository := pipelineConfig.NewMaterialRepositoryImpl(conn)
//	pipelineRepository := pipelineConfig.NewPipelineRepositoryImpl(conn, logger)
//	ciPipelineRepository := pipelineConfig.NewCiPipelineRepositoryImpl(conn, logger)
//	ciPipelineHistoryRepository := repository4.NewCiPipelineHistoryRepositoryImpl(conn, logger)
//	ciPipelineMaterialRepository := pipelineConfig.NewCiPipelineMaterialRepositoryImpl(conn, logger)
//	GitSensorClient, err := gitSensor.NewGitSensorClient(logger, &gitSensor.ClientConfig{})
//	ciConfig := &CiConfig{}
//	appWorkflowRepository := appWorkflow.NewAppWorkflowRepositoryImpl(logger, conn)
//	envRepository := repository3.NewEnvironmentRepositoryImpl(conn, logger, nil)
//	attributesService := attributes.NewAttributesServiceImpl(logger, nil)
//	appListingRepositoryQueryBuilder := helper.NewAppListingRepositoryQueryBuilder(logger)
//	appListingRepository := repository.NewAppListingRepositoryImpl(logger, conn, appListingRepositoryQueryBuilder)
//	appLabelsService := app2.NewAppCrudOperationServiceImpl(nil, logger, nil, nil, nil)
//	userAuthService := user.NewUserAuthServiceImpl(nil, nil, nil, nil, nil, nil, nil)
//	prePostCdScriptHistoryService := history.NewPrePostCdScriptHistoryServiceImpl(logger, nil, nil, nil)
//	prePostCiScriptHistoryService := history.NewPrePostCiScriptHistoryServiceImpl(logger, nil)
//	pipelineStageService := NewPipelineStageService(logger, nil, nil)
//	ciTemplateOverrideRepository := pipelineConfig.NewCiTemplateOverrideRepositoryImpl(conn, logger)
//	ciTemplateService := *NewCiTemplateServiceImpl(logger, nil, nil, nil)
//	gitMaterialHistoryService := history.NewGitMaterialHistoryServiceImpl(nil, logger)
//	ciPipelineHistoryService := history.NewCiPipelineHistoryServiceImpl(ciPipelineHistoryRepository, logger, ciPipelineRepository)
//	dockerArtifactStoreRepository := repository2.NewDockerArtifactStoreRepositoryImpl(conn)
//	ciCdPipelineOrchestrator = NewCiCdPipelineOrchestrator(appRepository, logger, materialRepository, pipelineRepository, ciPipelineRepository, ciPipelineMaterialRepository, GitSensorClient, ciConfig, appWorkflowRepository, envRepository, attributesService, appListingRepository, appLabelsService, userAuthService, prePostCdScriptHistoryService, prePostCiScriptHistoryService, pipelineStageService, ciTemplateOverrideRepository, gitMaterialHistoryService, ciPipelineHistoryService, ciTemplateService, dockerArtifactStoreRepository)
//}