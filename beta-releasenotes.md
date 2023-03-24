## Bugs
- fix: eks nodegroup label added (#3184)
- fix:http status fix for access to jobs (#3176)
- fix:wire issue fixed (#3179)
- fix:notes.txt throws error when charts require special KubeVersion (#3170)
- fix: Gitops validate and update with empty token (#3168)
- fix:Manifest output throws error when charts require special KubeVersion (#3162)
- fix: onlyDevtronCharts flag changed type from boolean to integer (#3161)
- fix: Optimize app grouping apis (#3125)
- fix: log api panic (#3156)
- fix: kubernetes external secret not accessible (#3143)
- fix: apps in progressing state indefinitely (#3137)
- fix: ci artifacts not coming for linked CI pipeline (#3134)
- fix: Helm repository deleted from argocd-cm when deleted from UI (Github Issue #1399) (#2970)
- fix: pg prom metrics not getting exported when pg query logging is disabled (#3124)
- fix: Send webhook data(source value and source type) in pipeline api response (#3120)
- fix: updated condition for adding kubectl apply synced timeline (#3115)
- fix: json unmarshal panic error temperory fix. (#3095)
- fix: added otel for deployment history info api (#3107)
- fix:showing right image imagetags in multiplecolons in registry url (#3103)
- fix: cd metrics nil ptr fix (#3099)
- fix: getting all environmets for super admin only (#3096)
- fix: logs issue for init containers (#3076)
- fix: pod resources not visible even if user have access on those pods. (#3071)
- fix: unable to create container registry with deleted name (#2963)
- fix: helm app deployment failure (#3060)
- fix: logs file path issue fix and docker file update for non root user (#3024)
- fix: urls pipeline fetching bug fix (#3063)
- fix: app grouping appid appname (#3058)
## Enhancements
- feat: Using server url in application object (#3175)
- feat: Add timer/counter telemetry for GitOps (#3119)
- feat: show notes txt for  helmapps deployed by  helm (#2966)
- feat: enable auto deployment trigger option after deployment app change (#3110)
- feat: jobs feature (#3074)
- feat: Disabling global secrets for application environment. (#3126)
- feat: Review config changes before deployment (#3077)
- feat: global secrets for pre/post cd (#3073)
- feat: export pg query metrics to prometheus (#3118)
- feat: Add timer/counter telemetry for CI process (#3081)
- feat: Git, GitOps, Container Registries and SSO login token hide from dashboard (#2952)
- feat:show notes for gitops app (#3082)
- feat: autoselect node,error messaging improvement , node group and custom shell support (#2925)
- feat: Option to run post-ci scripts even if build fails (#3065)
- feat: branch name field added in API response for a cd Artifact material (#3064)
- feat: github PR updater plugin added (#3051)
## Documentation
- docs: added jobs section (#3097)
- docs: added K8s client doc (#3028)
- docs: doc for GCP external secret (#3029)
## Others
- chore: ADO sync action changes (#3167)
- chore: Enterprise repo sync (#3146)
- chore: version upgrade of github action (ado-sync) (#3160)
- chore: ado-sync workflow (#3153)
- chore: ado-sync workflow changes  (#3151)
- chore: ADO-sync github action changes (#3149)
- chore: Helm lint mechanism and azure ADO sync github action (#3138)
- chore:  enterprise-repo-sync.yaml (#3127)
- task: restricted deployment status updation cron to fetch pipelines deployed within hours (#3104)
- chore: upgrade common-lib dependency (#3052)
- feat - Use cluster-name instead of server url in Argo cd application objects  (#2958)