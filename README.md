# awsflyteplugins
An example of Flyte Backend plugin for AWS SageMaker on Kubernetes. It uses the Amazon SageMaker operator for Kubernetes https://github.com/aws/amazon-sagemaker-operator-for-k8s to start SageMaker jobs in kubernetes.

- common/proto/sagemaker.proto specifies the custom information that is needed to execute besides the TaskTemplate
- go/sagemaker/plugin.go contains the backend plugin code that uses "Flyteplugins - plugin machinery"
- flytesagemakerplugin/sdk/..py contain the flytekit extensions that users can use to easily write a Flyte task that is executed on 
