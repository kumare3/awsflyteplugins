# flytek8ssparkplugin
An example of Flyte Backend plugin for Spark on Kubernetes. It uses the Spark Operator https://github.com/GoogleCloudPlatform/spark-on-k8s-operator to start Spark jobs in kubernetes.

- common/proto/spark.proto specifies the custom information that is needed to execute besides the TaskTemplate
- go/plugin/spark.go contains the backend plugin code that uses "Flyteplugins - plugin machinery"
- flytek8ssparkplugin/sdk/..py contain the flytekit extensions that users can use to easily write a Flyte task that is executed on Spark

This example provides a canonical example to write backend plugins. Refer to the docs for more details.

## Types of Plugins.
This repo serves as an example of how a `backend` Flyte plugin can be developed. Flyte can be extended with 2 types of plugins,
 - an SDK only plugin
 - a backend plugin

Since Flyte executes containers by default, an SDK plugin can essentially do anything that can be written in Python. It could call a remote service, launch a job and wait for it to complete or perform
an inline single machine computation. SDK plugins are easy to write and by far the most useful and easily introduced. But, there are a few disadvantages of having SDK only plugins. From the point of
view of Flyte backend the execution is purely a container execution. Any side-effects are not easily observed and actionable. (more on this later)

A backend plugin has extended capabilties, including improved visibility, custom information and progress indicators. It also provides better management primitives as compared to writing a SDK only (or within container plugin). 

### Example when having a simple container plugin is enough
 - A plugin that can execute a jupyter notebook and adapt it to the Flyte task and type system. This can be expressed as a python task.
 - Any additional logic that enhances some existing core execution plugin. For example, in some environment for every task (python code) the user launches, the platform desires a sidecar container to
   be executed (maybe for additional security, networking etc). This plugin simply provides a new task like `@improved_python_task` which decorates the already existing sidecar plugin: https://lyft.github.io/flyte/user/tasktypes/sidecar.html

### Example when a backend plugin is preferable
 - A Kubernetes CRD is launched
 - A remote service called that has some advanced visibility options, like a separate log link a visualization link etc. Also if the remote service results in some expensive computation.
 - Abstraction of the remote execution is desirable

`Expensive computation` Expensive here refers to expensive in dollar terms. The reason where it matters is when an execution is aborted, the abort should propagate to the downstream system. If the
remote service or execution was invoked from the container only plugin (as a side-effect) FlyteBackend today does not have the capability of cleaning up or propagating the abort to the downstream
system.


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
The default flyte sdk is flytekit (python) and in this example we implement a python flyte spark task.


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
