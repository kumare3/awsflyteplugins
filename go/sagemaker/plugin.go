package sagemaker

import (
	"context"
	"fmt"
	"time"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	pluginsCore "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/flytek8s"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/k8s"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/utils"

	"github.com/aws/aws-sdk-go/service/sagemaker"
	commonv1 "go.amzn.com/sagemaker/sagemaker-k8s-operator/api/v1/common"
	hpojobv1 "go.amzn.com/sagemaker/sagemaker-k8s-operator/api/v1/hyperparametertuningjob"

	"github.com/kumare3/awsflyteplugins/gen/pb-go/proto"
	taskError "github.com/lyft/flyteplugins/go/tasks/errors"
	. "go.amzn.com/sagemaker/sagemaker-k8s-operator/controllers/controllertest"
)

const KindSagemakerHPOJob = "HyperparameterTuningJob"

const (
	pluginID          = "aws_sagemaker_hpo"
	sagemakerTaskType = "aws_sagemaker_hpo"
)

// Sanity test that the plugin implements method of k8s.Plugin
var _ k8s.Plugin = awsSagemakerPlugin{}

type awsSagemakerPlugin struct {
}

func (m awsSagemakerPlugin) BuildIdentityResource(ctx context.Context, taskCtx pluginsCore.TaskExecutionMetadata) (k8s.Resource, error) {
	return &hpojobv1.HyperparameterTuningJob{}, nil
}

func convertStaticHyperparamsLiteralToSpecType(hyperparamLiteral *core.Literal) ([]*commonv1.KeyValuePair, error) {
	var retValue []*commonv1.KeyValuePair
	hyperFields := hyperparamLiteral.GetScalar().GetGeneric().GetFields()
	for k, v := range hyperFields {
		var newElem = commonv1.KeyValuePair{
			Name:  k,
			Value: v.GetStringValue(),
		}
		retValue = append(retValue, &newElem)
	}
	return retValue, nil
}

func (m awsSagemakerPlugin) BuildResource(ctx context.Context, taskCtx pluginsCore.TaskExecutionContext) (k8s.Resource, error) {
	// TODO build the actual spec of the k8s resource from the taskCtx Some helpful code is already added
	taskTemplate, err := taskCtx.TaskReader().Read(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to fetch task specification")
	} else if taskTemplate == nil {
		return nil, errors.Errorf("nil task specification")
	}

	taskInput, err := taskCtx.InputReader().Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to fetch task inputs")
	}

	// Get inputs from literals
	inputLiterals := taskInput.GetLiterals()
	trainPathLiteral, ok := inputLiterals["train"]
	if !ok {
		return nil, errors.Errorf("train input not specified")
	}
	validatePathLiteral, ok := inputLiterals["validation"]
	if !ok {
		return nil, errors.Errorf("validation input not specified")
	}
	staticHyperparamsLiteral, ok := inputLiterals["static_hyperparameters"]
	if !ok {
		return nil, errors.Errorf("static hyperparameters input not specified")
	}

	outputPath := taskCtx.OutputWriter().GetOutputPrefixPath().String()
	outputPath = fmt.Sprintf("%s/hpo_outputs", outputPath)

	sagemakerJob := proto.SagemakerHPOJob{}
	err = utils.UnmarshalStruct(taskTemplate.GetCustom(), &sagemakerJob)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid task specification for taskType [%s]", sagemakerTaskType)
	}

	// Convert the hyperparameters to the spec value
	staticHyperparams, err := convertStaticHyperparamsLiteralToSpecType(staticHyperparamsLiteral)
	if err != nil {
		return nil, errors.Wrapf(err, "could not convert static hyperparameters to spec type")
	}

	// If the container is part of the task template you can access it here
	container := taskTemplate.GetContainer()

	// When adding env vars there are some default env vars that are available, you can pass them through
	envVars := flytek8s.DecorateEnvVars(ctx, flytek8s.ToK8sEnvVar(container.GetEnv()), taskCtx.TaskExecutionMetadata().GetTaskExecutionID())
	_ = envVars

	taskName := taskCtx.TaskExecutionMetadata().GetTaskExecutionID().GetID().NodeExecutionId.GetExecutionId().GetName()

	hpoJob := &hpojobv1.HyperparameterTuningJob{
		Spec: hpojobv1.HyperparameterTuningJobSpec{
			HyperParameterTuningJobName: &taskName,
			HyperParameterTuningJobConfig: &commonv1.HyperParameterTuningJobConfig{
				ResourceLimits: &commonv1.ResourceLimits{
					MaxNumberOfTrainingJobs: ToInt64Ptr(10),
					MaxParallelTrainingJobs: ToInt64Ptr(5),
				},
				Strategy: "Bayesian",
				HyperParameterTuningJobObjective: &commonv1.HyperParameterTuningJobObjective{
					Type:       "Minimize",
					MetricName: ToStringPtr("validation:error"),
				},
				ParameterRanges: &commonv1.ParameterRanges{
					IntegerParameterRanges: []commonv1.IntegerParameterRange{
						commonv1.IntegerParameterRange{
							Name:        ToStringPtr("num_round"),
							MinValue:    ToStringPtr("10"),
							MaxValue:    ToStringPtr("20"),
							ScalingType: "Linear",
						},
					},
				},
				TrainingJobEarlyStoppingType: "Auto",
			},
			TrainingJobDefinition: &commonv1.HyperParameterTrainingJobDefinition{
				StaticHyperParameters: staticHyperparams,
				AlgorithmSpecification: &commonv1.HyperParameterAlgorithmSpecification{
					TrainingImage:     &sagemakerJob.AlgorithmSpecification.TrainingImage,
					TrainingInputMode: commonv1.TrainingInputMode(sagemakerJob.AlgorithmSpecification.TrainingInputMode),
				},
				InputDataConfig: []commonv1.Channel{
					commonv1.Channel{
						ChannelName: ToStringPtr("train"),
						DataSource: &commonv1.DataSource{
							S3DataSource: &commonv1.S3DataSource{
								S3DataType: "S3Prefix",
								S3Uri:      ToStringPtr(trainPathLiteral.GetScalar().GetBlob().GetUri()),
							},
						},
						ContentType: ToStringPtr("text/csv"),
						InputMode:   "File",
					},
					commonv1.Channel{
						ChannelName: ToStringPtr("validation"),
						DataSource: &commonv1.DataSource{
							S3DataSource: &commonv1.S3DataSource{
								S3DataType: "S3Prefix",
								S3Uri:      ToStringPtr(validatePathLiteral.GetScalar().GetBlob().GetUri()),
							},
						},
						ContentType: ToStringPtr("text/csv"),
						InputMode:   "File",
					},
				},
				OutputDataConfig: &commonv1.OutputDataConfig{
					S3OutputPath: ToStringPtr(outputPath),
				},
				ResourceConfig: &commonv1.ResourceConfig{
					InstanceType:   sagemakerJob.ResourceConfig.InstanceType,
					InstanceCount:  &sagemakerJob.ResourceConfig.InstanceCount,
					VolumeSizeInGB: &sagemakerJob.ResourceConfig.VolumeSizeInGB,
					VolumeKmsKeyId: &sagemakerJob.ResourceConfig.VolumeKmsKeyId,
				},
				RoleArn: &sagemakerJob.RoleArn,
				StoppingCondition: &commonv1.StoppingCondition{
					MaxRuntimeInSeconds: &sagemakerJob.StoppingCondition.MaxRuntimeInSeconds,
				},
			},
			Region: &sagemakerJob.Region,
		},
	}

	return hpoJob, nil
}

func getEventInfoForHPOJob(job *hpojobv1.HyperparameterTuningJob) (*pluginsCore.TaskInfo, error) {
	var taskLogs []*core.TaskLog
	customInfoMap := make(map[string]string)

	customInfo, err := utils.MarshalObjToStruct(customInfoMap)
	if err != nil {
		return nil, err
	}

	return &pluginsCore.TaskInfo{
		Logs:       taskLogs,
		CustomInfo: customInfo,
	}, nil
}

func (m awsSagemakerPlugin) GetTaskPhase(ctx context.Context, pluginContext k8s.PluginContext, resource k8s.Resource) (pluginsCore.PhaseInfo, error) {
	job := resource.(*hpojobv1.HyperparameterTuningJob)
	info, err := getEventInfoForHPOJob(job)
	if err != nil {
		return pluginsCore.PhaseInfoUndefined, err
	}

	occurredAt := time.Now()
	switch job.Status.HyperParameterTuningJobStatus {
	case sagemaker.HyperParameterTuningJobStatusFailed:
		execError := &core.ExecutionError{
			Message: job.Status.Additional,
		}
		return pluginsCore.PhaseInfoFailed(pluginsCore.PhasePermanentFailure, execError, info), nil
	case sagemaker.HyperParameterTuningJobStatusStopped:
		reason := fmt.Sprintf("HPO Job Stopped")
		return pluginsCore.PhaseInfoRetryableFailure(taskError.DownstreamSystemError, reason, info), nil
	case sagemaker.HyperParameterTuningJobStatusCompleted:
		return pluginsCore.PhaseInfoSuccess(info), nil
	case "":
		return pluginsCore.PhaseInfoQueued(occurredAt, pluginsCore.DefaultPhaseVersion, "job submitted"), nil
	}

	return pluginsCore.PhaseInfoRunning(pluginsCore.DefaultPhaseVersion, info), nil
}

// TODO we should register the plugin
func init() {
	if err := commonv1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}

	pluginmachinery.PluginRegistry().RegisterK8sPlugin(
		k8s.PluginEntry{
			ID:                  pluginID,
			RegisteredTaskTypes: []pluginsCore.TaskType{sagemakerTaskType},
			ResourceToWatch:     &hpojobv1.HyperparameterTuningJob{},
			Plugin:              awsSagemakerPlugin{},
			IsDefault:           false,
		})
}
