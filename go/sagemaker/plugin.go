package sagemaker

import (
	"context"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"

	pluginsCore "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/flytek8s"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/k8s"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	commonv1 "go.amzn.com/sagemaker/sagemaker-k8s-operator/api/v1/common"
	hpojobv1 "go.amzn.com/sagemaker/sagemaker-k8s-operator/api/v1/hyperparametertuningjob"

	"github.com/kumare3/awsflyteplugins/gen/pb-go/proto"
	. "go.amzn.com/sagemaker/sagemaker-k8s-operator/controllers/controllertest"
)

const KindSagemakerHPOJob = "HyperparameterTuningJob"

const (
	pluginID          = "aws_sagemaker_hpo"
	sagemakerTaskType = "aws_sagemaker_hpo"
)

// Sanity test that the plugin implements method of k8s.Plugin
var _ k8s.Plugin = mySamplePlugin{}

type mySamplePlugin struct {
}

func (m mySamplePlugin) BuildIdentityResource(ctx context.Context, taskCtx pluginsCore.TaskExecutionMetadata) (k8s.Resource, error) {
	// TODO This should return the type of the kubernetes resource. As golang has no generics it is hard to gather the type information automatically
	// Example of a pod
	return &hpojobv1.HyperparameterTuningJob{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindSagemakerHPOJob,
			APIVersion: commonv1.GroupVersion.String(),
		},
	}, nil
}

func (m mySamplePlugin) BuildResource(ctx context.Context, taskCtx pluginsCore.TaskExecutionContext) (k8s.Resource, error) {
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
	trainPath, ok := inputLiterals["train"]
	if !ok {
		return nil, errors.Errorf("train input not specified")
	}
	validatePath, ok := inputLiterals["validation"]
	if !ok {
		return nil, errors.Errorf("validation input not specified")
	}

	outputPath := taskCtx.OutputWriter().GetOutputPrefixPath().String()

	// TODO if we have special information to be marshalled through from the python SDK then it can be retrieved using the util
	sagemakerJob := proto.SagemakerHPOJob{}
	err = utils.UnmarshalStruct(taskTemplate.GetCustom(), &sagemakerJob)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid task specification for taskType [%s]", sagemakerTaskType)
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
								S3Uri:      ToStringPtr(trainPath.GetScalar().GetBlob().GetUri()),
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
								S3Uri:      ToStringPtr(validatePath.GetScalar().GetBlob().GetUri()),
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
			Region: ToStringPtr("us-east-2"),
		},
	}

	return hpoJob, nil
}

func (m mySamplePlugin) GetTaskPhase(ctx context.Context, pluginContext k8s.PluginContext, resource k8s.Resource) (pluginsCore.PhaseInfo, error) {
	// TODO observe the applicate state from the passed in resource and return the PhaseInfo
	// E.g we will consider the resource to be a pod
	//p := resource.(*v1.Pod)
	p := resource.(*hpojobv1.HyperparameterTuningJob)
	_ = p

	return pluginsCore.PhaseInfoRunning(1, &pluginsCore.TaskInfo{}), nil
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
			// TODO Type of the k8s resource, e.g. Pod
			ResourceToWatch: &v1.PersistentVolume{},
			Plugin:          mySamplePlugin{},
			IsDefault:       false,
		})
}
