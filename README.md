This is a repo hosting a PoC plugin that enables AWS SageMaker in Flyte, Lyft's cloud native machine learning and data processing platform.

Disclaimer: this plugin is still under development and is currently unofficial. We do not provide support outside of Lyft's environment at the moment. 

# awsflyteplugins
An example of Flyte Backend plugin for AWS SageMaker on Kubernetes. It uses the Amazon SageMaker operator for Kubernetes https://github.com/aws/amazon-sagemaker-operator-for-k8s to start SageMaker jobs in kubernetes.

- common/proto/sagemaker.proto specifies the custom information that is needed to execute besides the TaskTemplate
- go/sagemaker/plugin.go contains the backend plugin code that uses "Flyteplugins - plugin machinery"
- flytesagemakerplugin/sdk/..py contain the flytekit extensions that users can use to easily write a Flyte task that is executed on 

## Installation

### Python plugin
To install the Python plugin, use the following command:
```
$ pip install flyte
```
Note that, since this is a PoC, the way to install the Python plugin is subject to change in the future.

