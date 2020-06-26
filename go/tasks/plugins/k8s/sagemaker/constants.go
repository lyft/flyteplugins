package sagemaker

import commonv1 "github.com/aws/amazon-sagemaker-operator-for-k8s/api/v1/common"

const (
	trainingJobTaskPluginID = "sagemaker_training_job_task"
	trainingJobTaskType = "sagemaker_training_job_task"
)

const (
	hpoJobTaskPluginID = "sagemaker_hpo_job_task"
	hpoJobTaskType = "sagemaker_hpo_job_task"
)

const (
	BayesianSageMakerAPIHyperParameterTuningJobStrategyType commonv1.HyperParameterTuningJobStrategyType = "Bayesian"
	RandomSageMakerAPIHyperParameterTuningJobStrategyType commonv1.HyperParameterTuningJobStrategyType = "Random"
)

const (
	AutoSageMakerAPIHyperParameterScalingType               commonv1.HyperParameterScalingType = "Auto"
	LinearSageMakerAPIHyperParameterScalingType             commonv1.HyperParameterScalingType = "Linear"
	LogarithmicSageMakerAPIHyperParameterScalingType        commonv1.HyperParameterScalingType = "Logarithmic"
	ReverseLogarithmicSageMakerAPIHyperParameterScalingType commonv1.HyperParameterScalingType = "ReverseLogarithmic"
)

const (
	MinimizeSageMakerAPIHyperParameterTuningJobObjectiveType commonv1.HyperParameterTuningJobObjectiveType = "Minimize"
	MaximizeSageMakerAPIHyperParameterTuningJobObjectiveType commonv1.HyperParameterTuningJobObjectiveType = "Maximize"
)

const (
	FileSageMakerAPITrainingInputMode commonv1.TrainingInputMode = "File"
	PipeSageMakerAPITrainingInputMode commonv1.TrainingInputMode = "Pipe"
)

const (
	OffSageMakerAPITrainingJobEarlyStoppingType  commonv1.TrainingJobEarlyStoppingType = "Off"
	AutoSageMakerAPITrainingJobEarlyStoppingType commonv1.TrainingJobEarlyStoppingType = "Auto"
)

const (
	CustomSageMakerAPIAlgorithmName  string = "custom"
	XgboostSageMakerAPIAlgorithmName string = "xgboost"
)