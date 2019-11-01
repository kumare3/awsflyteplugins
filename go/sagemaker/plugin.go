package sagemaker

import (
	"context"
	"fmt"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/flytek8s"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"

	pluginsCore "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/k8s"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kumare3/awsflyteplugins/gen/pb-go/proto"
)

const (
	pluginID          = "aws_sagemaker"
	sagemakerTaskType = "aws_sagemaker"
)

// Sanity test that the plugin implements method of k8s.Plugin
var _ k8s.Plugin = mySamplePlugin{}

type mySamplePlugin struct {
}

func (m mySamplePlugin) BuildIdentityResource(ctx context.Context, taskCtx pluginsCore.TaskExecutionMetadata) (k8s.Resource, error) {
	// TODO This should return the type of the kubernetes resource. As golang has no generics it is hard to gather the type information automatically
	// Example of a pod
	return &v1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: v1.SchemeGroupVersion.String(),
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

	// TODO if we have special information to be marshalled through from the python SDK then it can be retrieved using the util
	sagemakerJob := proto.SagemakerJob{}
	err = utils.UnmarshalStruct(taskTemplate.GetCustom(), &sagemakerJob)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid task specification for taskType [%s]", sagemakerTaskType)
	}

	// If the container is part of the task template you can access it here
	container := taskTemplate.GetContainer()

	// When adding env vars there are some default env vars that are available, you can pass them through
	envVars := flytek8s.DecorateEnvVars(ctx, flytek8s.ToK8sEnvVar(container.GetEnv()), taskCtx.TaskExecutionMetadata().GetTaskExecutionID())
	fmt.Print(envVars)
	panic("Implement me")
}

func (m mySamplePlugin) GetTaskPhase(ctx context.Context, pluginContext k8s.PluginContext, resource k8s.Resource) (pluginsCore.PhaseInfo, error) {
	// TODO observe the applicate state from the passed in resource and return the PhaseInfo
	// E.g we will consider the resource to be a pod
	//p := resource.(*v1.Pod)
	panic("implement me")
}

// TODO we should register the plugin
func init() {
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
