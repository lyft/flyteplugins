package tensorflow

import (
	"context"
	"fmt"
	"testing"
	"time"

	commonOp "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	"github.com/lyft/flyteplugins/go/tasks/logs"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/flytek8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/stretchr/testify/mock"

	"github.com/lyft/flytestdlib/storage"

	pluginsCore "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/utils"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core/mocks"

	pluginIOMocks "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/io/mocks"

	"github.com/golang/protobuf/jsonpb"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/plugins"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tfOp "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
)

const testImage = "image://"
const serviceAccount = "tensorflow_sa"

var (
	dummyEnvVars = []*core.KeyValuePair{
		{Key: "Env_Var", Value: "Env_Val"},
	}

	testArgs = []string{
		"test-args",
	}

	resourceRequirements = &corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:         resource.MustParse("1000m"),
			corev1.ResourceMemory:      resource.MustParse("1Gi"),
			flytek8s.ResourceNvidiaGPU: resource.MustParse("1"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:         resource.MustParse("100m"),
			corev1.ResourceMemory:      resource.MustParse("512Mi"),
			flytek8s.ResourceNvidiaGPU: resource.MustParse("1"),
		},
	}

	jobName      = "the-job"
	jobNamespace = "tensorflow-namespace"
)

func dummyTensorflowCustomObj(workers int32, psReplicas int32, chiefReplicas int32) *plugins.DistributedTensorflowTrainingTask {
	return &plugins.DistributedTensorflowTrainingTask{
		Workers: workers,
		PsReplicas: psReplicas,
		ChiefReplicas: chiefReplicas,
	}
}

func dummyTensorflowTaskTemplate(id string, tensorflowCustomObj *plugins.DistributedTensorflowTrainingTask) *core.TaskTemplate {

	tfObjJSON, err := utils.MarshalToString(tensorflowCustomObj)
	if err != nil {
		panic(err)
	}

	structObj := structpb.Struct{}

	err = jsonpb.UnmarshalString(tfObjJSON, &structObj)
	if err != nil {
		panic(err)
	}

	return &core.TaskTemplate{
		Id:   &core.Identifier{Name: id},
		Type: "container",
		Target: &core.TaskTemplate_Container{
			Container: &core.Container{
				Image: testImage,
				Args:  testArgs,
				Env:   dummyEnvVars,
			},
		},
		Custom: &structObj,
	}
}

func dummyTensorflowTaskContext(taskTemplate *core.TaskTemplate) pluginsCore.TaskExecutionContext {
	taskCtx := &mocks.TaskExecutionContext{}
	inputReader := &pluginIOMocks.InputReader{}
	inputReader.OnGetInputPrefixPath().Return(storage.DataReference("/input/prefix"))
	inputReader.OnGetInputPath().Return(storage.DataReference("/input"))
	inputReader.OnGetMatch(mock.Anything).Return(&core.LiteralMap{}, nil)
	taskCtx.OnInputReader().Return(inputReader)

	outputReader := &pluginIOMocks.OutputWriter{}
	outputReader.OnGetOutputPath().Return(storage.DataReference("/data/outputs.pb"))
	outputReader.OnGetOutputPrefixPath().Return(storage.DataReference("/data/"))
	taskCtx.OnOutputWriter().Return(outputReader)

	taskReader := &mocks.TaskReader{}
	taskReader.OnReadMatch(mock.Anything).Return(taskTemplate, nil)
	taskCtx.OnTaskReader().Return(taskReader)

	tID := &mocks.TaskExecutionID{}
	tID.OnGetID().Return(core.TaskExecutionIdentifier{
		NodeExecutionId: &core.NodeExecutionIdentifier{
			ExecutionId: &core.WorkflowExecutionIdentifier{
				Name:    "my_name",
				Project: "my_project",
				Domain:  "my_domain",
			},
		},
	})
	tID.OnGetGeneratedName().Return("some-acceptable-name")

	resources := &mocks.TaskOverrides{}
	resources.OnGetResources().Return(resourceRequirements)

	taskExecutionMetadata := &mocks.TaskExecutionMetadata{}
	taskExecutionMetadata.OnGetTaskExecutionID().Return(tID)
	taskExecutionMetadata.OnGetNamespace().Return("test-namespace")
	taskExecutionMetadata.OnGetAnnotations().Return(map[string]string{"annotation-1": "val1"})
	taskExecutionMetadata.OnGetLabels().Return(map[string]string{"label-1": "val1"})
	taskExecutionMetadata.OnGetOwnerReference().Return(v1.OwnerReference{
		Kind: "node",
		Name: "blah",
	})
	taskExecutionMetadata.OnIsInterruptible().Return(true)
	taskExecutionMetadata.OnGetOverrides().Return(resources)
	taskExecutionMetadata.OnGetK8sServiceAccount().Return(serviceAccount)
	taskCtx.OnTaskExecutionMetadata().Return(taskExecutionMetadata)
	return taskCtx
}

func dummyTensorflowJobResource(tensorflowResourceHandler tensorflowOperatorResourceHandler,
	workers int32, psReplicas int32, chiefReplicas int32, conditionType commonOp.JobConditionType) *tfOp.TFJob {
	var jobConditions []commonOp.JobCondition

	now := time.Now()

	jobCreated := commonOp.JobCondition{
		Type:    commonOp.JobCreated,
		Status:  corev1.ConditionTrue,
		Reason:  "TensorflowJobCreated",
		Message: "TensorflowJob the-job is created.",
		LastUpdateTime: v1.Time{
			Time: now,
		},
		LastTransitionTime: v1.Time{
			Time: now,
		},
	}
	jobRunningActive := commonOp.JobCondition{
		Type:    commonOp.JobRunning,
		Status:  corev1.ConditionTrue,
		Reason:  "TensorflowJobRunning",
		Message: "TensorflowJob the-job is running.",
		LastUpdateTime: v1.Time{
			Time: now.Add(time.Minute),
		},
		LastTransitionTime: v1.Time{
			Time: now.Add(time.Minute),
		},
	}
	jobRunningInactive := *jobRunningActive.DeepCopy()
	jobRunningInactive.Status = corev1.ConditionFalse
	jobSucceeded := commonOp.JobCondition{
		Type:    commonOp.JobSucceeded,
		Status:  corev1.ConditionTrue,
		Reason:  "TensorflowJobSucceeded",
		Message: "TensorflowJob the-job is successfully completed.",
		LastUpdateTime: v1.Time{
			Time: now.Add(2 * time.Minute),
		},
		LastTransitionTime: v1.Time{
			Time: now.Add(2 * time.Minute),
		},
	}
	jobFailed := commonOp.JobCondition{
		Type:    commonOp.JobFailed,
		Status:  corev1.ConditionTrue,
		Reason:  "TensorflowJobFailed",
		Message: "TensorflowJob the-job is failed.",
		LastUpdateTime: v1.Time{
			Time: now.Add(2 * time.Minute),
		},
		LastTransitionTime: v1.Time{
			Time: now.Add(2 * time.Minute),
		},
	}
	jobRestarting := commonOp.JobCondition{
		Type:    commonOp.JobRestarting,
		Status:  corev1.ConditionTrue,
		Reason:  "TensorflowJobRestarting",
		Message: "TensorflowJob the-job is restarting because some replica(s) failed.",
		LastUpdateTime: v1.Time{
			Time: now.Add(3 * time.Minute),
		},
		LastTransitionTime: v1.Time{
			Time: now.Add(3 * time.Minute),
		},
	}

	switch conditionType {
	case commonOp.JobCreated:
		jobConditions = []commonOp.JobCondition{
			jobCreated,
		}
	case commonOp.JobRunning:
		jobConditions = []commonOp.JobCondition{
			jobCreated,
			jobRunningActive,
		}
	case commonOp.JobSucceeded:
		jobConditions = []commonOp.JobCondition{
			jobCreated,
			jobRunningInactive,
			jobSucceeded,
		}
	case commonOp.JobFailed:
		jobConditions = []commonOp.JobCondition{
			jobCreated,
			jobRunningInactive,
			jobFailed,
		}
	case commonOp.JobRestarting:
		jobConditions = []commonOp.JobCondition{
			jobCreated,
			jobRunningInactive,
			jobFailed,
			jobRestarting,
		}
	}

	tfObj := dummyTensorflowCustomObj(workers, psReplicas, chiefReplicas)
	taskTemplate := dummyTensorflowTaskTemplate("the job", tfObj)
	resource, err := tensorflowResourceHandler.BuildResource(context.TODO(), dummyTensorflowTaskContext(taskTemplate))
	if err != nil {
		panic(err)
	}

	return &tfOp.TFJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      jobName,
			Namespace: jobNamespace,
		},
		Spec: resource.(*tfOp.TFJob).Spec,
		Status: commonOp.JobStatus{
			Conditions:        jobConditions,
			ReplicaStatuses:   nil,
			StartTime:         nil,
			CompletionTime:    nil,
			LastReconcileTime: nil,
		},
	}
}

func TestBuildResourceTensorflow(t *testing.T) {
	tensorflowResourceHandler := tensorflowOperatorResourceHandler{}

	tfObj := dummyTensorflowCustomObj(100, 1,1)
	taskTemplate := dummyTensorflowTaskTemplate("the job", tfObj)

	resource, err := tensorflowResourceHandler.BuildResource(context.TODO(), dummyTensorflowTaskContext(taskTemplate))
	assert.NoError(t, err)
	assert.NotNil(t, resource)

	tensorflowJob, ok := resource.(*tfOp.TFJob)
	assert.True(t, ok)
	assert.Equal(t, int32(100), *tensorflowJob.Spec.TFReplicaSpecs[tfOp.TFReplicaTypeWorker].Replicas)
	assert.Equal(t, int32(1), *tensorflowJob.Spec.TFReplicaSpecs[tfOp.TFReplicaTypePS].Replicas)
	assert.Equal(t, int32(1), *tensorflowJob.Spec.TFReplicaSpecs[tfOp.TFReplicaTypeChief].Replicas)

	for _, replicaSpec := range tensorflowJob.Spec.TFReplicaSpecs {
		var hasContainerWithDefaultTensorflowName = false

		for _, container := range replicaSpec.Template.Spec.Containers {
			if container.Name == tfOp.DefaultContainerName {
				hasContainerWithDefaultTensorflowName = true
			}

			assert.Equal(t, resourceRequirements.Requests, container.Resources.Requests)
			assert.Equal(t, resourceRequirements.Limits, container.Resources.Limits)
		}

		assert.True(t, hasContainerWithDefaultTensorflowName)
	}
}

func TestGetTaskPhase(t *testing.T) {
	tensorflowResourceHandler := tensorflowOperatorResourceHandler{}
	ctx := context.TODO()

	dummyTensorflowJobResourceCreator := func(conditionType commonOp.JobConditionType) *tfOp.TFJob {
		return dummyTensorflowJobResource(tensorflowResourceHandler, 2,1, 1, conditionType)
	}

	taskPhase, err := tensorflowResourceHandler.GetTaskPhase(ctx, nil, dummyTensorflowJobResourceCreator(commonOp.JobCreated))
	assert.NoError(t, err)
	assert.Equal(t, pluginsCore.PhaseQueued, taskPhase.Phase())
	assert.NotNil(t, taskPhase.Info())
	assert.Nil(t, err)

	taskPhase, err = tensorflowResourceHandler.GetTaskPhase(ctx, nil, dummyTensorflowJobResourceCreator(commonOp.JobRunning))
	assert.NoError(t, err)
	assert.Equal(t, pluginsCore.PhaseRunning, taskPhase.Phase())
	assert.NotNil(t, taskPhase.Info())
	assert.Nil(t, err)

	taskPhase, err = tensorflowResourceHandler.GetTaskPhase(ctx, nil, dummyTensorflowJobResourceCreator(commonOp.JobSucceeded))
	assert.NoError(t, err)
	assert.Equal(t, pluginsCore.PhaseSuccess, taskPhase.Phase())
	assert.NotNil(t, taskPhase.Info())
	assert.Nil(t, err)

	taskPhase, err = tensorflowResourceHandler.GetTaskPhase(ctx, nil, dummyTensorflowJobResourceCreator(commonOp.JobFailed))
	assert.NoError(t, err)
	assert.Equal(t, pluginsCore.PhaseRetryableFailure, taskPhase.Phase())
	assert.NotNil(t, taskPhase.Info())
	assert.Nil(t, err)

	taskPhase, err = tensorflowResourceHandler.GetTaskPhase(ctx, nil, dummyTensorflowJobResourceCreator(commonOp.JobRestarting))
	assert.NoError(t, err)
	assert.Equal(t, pluginsCore.PhaseRetryableFailure, taskPhase.Phase())
	assert.NotNil(t, taskPhase.Info())
	assert.Nil(t, err)
}

func TestGetLogs(t *testing.T) {
	assert.NoError(t, logs.SetLogConfig(&logs.LogConfig{
		IsKubernetesEnabled: true,
		KubernetesURL:       "k8s.com",
	}))

	workers := int32(2)
	psReplicas := int32(1)
	chiefReplicas := int32(1)

	tensorflowResourceHandler := tensorflowOperatorResourceHandler{}
	jobLogs, err := getLogs(dummyTensorflowJobResource(tensorflowResourceHandler, workers, psReplicas, chiefReplicas, commonOp.JobRunning),
		workers, psReplicas, chiefReplicas)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(jobLogs))
	assert.Equal(t, fmt.Sprintf("k8s.com/#!/log/%s/%s-worker-0/pod?namespace=tensorflow-namespace", jobNamespace, jobName), jobLogs[0].Uri)
	assert.Equal(t, fmt.Sprintf("k8s.com/#!/log/%s/%s-worker-1/pod?namespace=tensorflow-namespace", jobNamespace, jobName), jobLogs[1].Uri)
	assert.Equal(t, fmt.Sprintf("k8s.com/#!/log/%s/%s-psReplica-0/pod?namespace=tensorflow-namespace", jobNamespace, jobName), jobLogs[2].Uri)
	assert.Equal(t, fmt.Sprintf("k8s.com/#!/log/%s/%s-chiefReplica-0/pod?namespace=tensorflow-namespace", jobNamespace, jobName), jobLogs[3].Uri)
}