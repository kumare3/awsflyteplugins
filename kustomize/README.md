# Kustomize
In Flyte we use kustomize to build the entire deployment of the data plane. So for plugins that use operators, it is expected that the operator deployment kustomize file is provided. A flyte
deployment will use this and in some cases overlay this
