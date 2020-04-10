module github.com/kumare3/awsflyteplugins

go 1.13

require (
	cloud.google.com/go v0.47.0 // indirect
	github.com/Azure/azure-sdk-for-go v10.2.1-beta+incompatible // indirect
	github.com/Azure/go-autorest v13.3.0+incompatible // indirect
	github.com/GoogleCloudPlatform/spark-on-k8s-operator v0.0.0-20191028162909-4990c026d087 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/aws/amazon-sagemaker-operator-for-k8s v1.0.0
	github.com/aws/aws-sdk-go v1.25.36
	github.com/benlaurie/objecthash v0.0.0-20180202135721-d1e3d6079fc1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cheekybits/is v0.0.0-20150225183255-68e9c0620927 // indirect
	github.com/coocood/freecache v1.1.0 // indirect
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/go-redis/redis v6.15.6+incompatible // indirect
	github.com/go-test/deep v1.0.4 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/golang/protobuf v1.3.3
	github.com/graymeta/stow v0.0.0-20190522170649-903027f87de7 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/influxdata/influxdb v1.7.9 // indirect
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lyft/datacatalog v0.1.1 // indirect
	github.com/lyft/flyteidl v0.14.1
	github.com/lyft/flyteplugins v0.2.0
	github.com/lyft/flytepropeller v0.1.13-0.20191112060948-c2d2fde4537c
	github.com/lyft/flytestdlib v0.2.28
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/ncw/swift v1.0.49-0.20190728102658-a24ef33bc9b7 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/common v0.7.0 // indirect
	github.com/prometheus/procfs v0.0.5 // indirect
	github.com/satori/uuid v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.5.1
	go.opencensus.io v0.22.1 // indirect
	go.uber.org/zap v1.12.0 // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20191030232956-1e24073be82c // indirect
	golang.org/x/vgo v0.0.0-20180912184537-9d567625acf4 // indirect
	google.golang.org/genproto v0.0.0-20191028173616-919d9bdd9fe6 // indirect
	google.golang.org/grpc v1.24.0 // indirect
	k8s.io/api v0.0.0-20191031065753-b19d8caf39be
	k8s.io/apiextensions-apiserver v0.0.0-20191028232452-c47e10e6d5a3 // indirect
	k8s.io/apimachinery v0.0.0-20191030190112-bb31b70367b7
	k8s.io/client-go v11.0.1-0.20190918222721-c0e3722d5cf0+incompatible
	sigs.k8s.io/controller-runtime v0.3.1-0.20191029211253-40070e2a1958
	sigs.k8s.io/testing_frameworks v0.1.2 // indirect
)

replace k8s.io/api => k8s.io/api v0.0.0-20190918155943-95b840bb6a1f

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
