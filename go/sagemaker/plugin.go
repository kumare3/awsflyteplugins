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
	_ = trainPath
	_ = validatePath

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

	hpoJob := &hpojobv1.HyperparameterTuningJob{
		Spec: hpojobv1.HyperparameterTuningJobSpec{
			HyperParameterTuningJobConfig: &commonv1.HyperParameterTuningJobConfig{
				ResourceLimits: &commonv1.ResourceLimits{
					MaxNumberOfTrainingJobs: ToInt64Ptr(10),
					MaxParallelTrainingJobs: ToInt64Ptr(10),
				},
				Strategy: "Bayesian",
			},
			TrainingJobDefinition: &commonv1.HyperParameterTrainingJobDefinition{
				AlgorithmSpecification: &commonv1.HyperParameterAlgorithmSpecification{

				},
				InputDataConfig: [commonv1.Channel{

				}, commonv1.Channel{

				},],
				OutputDataConfig: &commonv1.OutputDataConfig{

				},
				ResourceConfig: &commonv1.ResourceConfig{

				},
				RoleArn: ,
				StoppingCondition: &commonv1.StoppingCondition{

				}
			}
			Region: ToStringPtr("us-east-1"),
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
