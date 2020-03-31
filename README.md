# awsflyteplugins
An example of Flyte Backend plugin for AWS SageMaker on Kubernetes. It uses the Amazon SageMaker operator for Kubernetes https://github.com/aws/amazon-sagemaker-operator-for-k8s to start SageMaker jobs in kubernetes.

- common/proto/sagemaker.proto specifies the custom information that is needed to execute besides the TaskTemplate
- go/sagemaker/plugin.go contains the backend plugin code that uses "Flyteplugins - plugin machinery"
- flytesagemakerplugin/sdk/..py contain the flytekit extensions that users can use to easily write a Flyte task that is executed on 

This example provides a canonical example to write backend plugins. Refer to the docs for more details.

## Types of Plugins.
This repo serves as an example of how a `backend` Flyte plugin can be developed. Flyte can be extended with 2 types of plugins,
 - an SDK only plugin
 - a backend plugin

Since Flyte executes containers by default, an SDK plugin can essentially do anything that can be written in Python. It could call a remote service, launch a job and wait for it to complete or perform an inline single machine computation. SDK plugins are easy to write and by far the most useful and easily introduced. But, there are a few disadvantages of having SDK only plugins. From the point of view of Flyte backend the execution is purely a container execution. Any side-effects are not easily observed and actionable. (more on this later)

A backend plugin has extended capabilties, including improved visibility, custom information and progress indicators. It also provides better management primitives as compared to writing a SDK only (or within container plugin). 


## What are the parts of a backend plugin?

### Custom model defintion
 Does the plugin need any additional information that cannot be covered by the default definition of Flyte TaskTemplte https://github.com/lyft/flyteidl/blob/master/protos/flyteidl/core/tasks.proto#L85
If there is additional information to be passed, like in the case of spark plugin, the type of spark application, spark configuration and hadoop configuration are additional values to be passed, then
define the model in some specification language that both the backend plugin on Golang and the sdk plugin (in python) and if required the UI/Console in typescript can understand.
In case of Spark we decided to use Protobuf. But this is not a requirement. Plugin writers can use other alternatives like open-api spec, json spec etc.

If using protobuf, the protobuf is declared in the common.protos file and some helper scripts/generate_protos.sh make it easy to generate

### Backend Plugin in golang
.....


### SDK plugin flytekit.
The default flyte sdk is flytekit (python). In this example we implement a python flyte SagemakerXgBoostOptimizer task.


### Bundling dependencies
Finally if the backend plugin requires any operators, they are specified using a kustomize script


#### To test the demo workflow
1. First build python package for the plugin
```
make build_python
```

2. build docker image
```
make build_demo_docker
```
3. you need to upload the docker image or change the tag such that it is the same as the one registered and the image should be available to your sandbox cluster
```
upload the image to ecr?
```
4. enter your image and run register.sh
```
docker run --rm -ti 8cdec6176236 bash
```
