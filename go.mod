module github.com/kumare3/awsflyteplugins

go 1.13

// Pin the version of client-go to something that's compatible with katrogan's fork of api and apimachinery
// Type the following
//   replace k8s.io/client-go => k8s.io/client-go kubernetes-1.16.2
// and it will be replaced with the 'sha' variant of the version

replace (
	github.com/GoogleCloudPlatform/spark-on-k8s-operator => github.com/lyft/spark-on-k8s-operator v0.1.3
	github.com/aws/amazon-sagemaker-operator-for-k8s => ./amazon-sagemaker-operator-for-k8s
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1
	gopkg.in/fsnotify.v1 => github.com/fsnotify/fsnotify v1.4.7
	k8s.io/api => github.com/lyft/api v0.0.0-20191031200350-b49a72c274e0
	k8s.io/apimachinery => github.com/lyft/apimachinery v0.0.0-20191031200210-047e3ea32d7f
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200204173128-addea2498afe
)

require (
	cloud.google.com/go v0.56.0 // indirect
	github.com/Azure/azure-sdk-for-go v41.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/aws/amazon-sagemaker-operator-for-k8s v1.1.0
	github.com/aws/aws-sdk-go v1.30.7
	github.com/golang/protobuf v1.3.5
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/lyft/flyteidl v0.17.27
	github.com/lyft/flyteplugins v0.3.20
	github.com/lyft/flytepropeller v0.2.25
	github.com/lyft/flytestdlib v0.3.3
	github.com/mitchellh/mapstructure v1.2.2
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/pelletier/go-toml v1.7.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1 // indirect
	github.com/prometheus/procfs v0.0.11 // indirect
	github.com/sirupsen/logrus v1.5.0 // indirect
	github.com/spf13/cobra v0.0.7 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3 // indirect
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71 // indirect
	golang.org/x/sys v0.0.0-20200409092240-59c9f1ba88fa // indirect
	google.golang.org/api v0.21.0 // indirect
	google.golang.org/genproto v0.0.0-20200409111301-baae70f3302d // indirect
	google.golang.org/grpc v1.28.1 // indirect
	gopkg.in/ini.v1 v1.55.0 // indirect
	k8s.io/apiextensions-apiserver v0.18.1 // indirect
	k8s.io/client-go v11.0.1-0.20190918222721-c0e3722d5cf0+incompatible
	k8s.io/kube-openapi v0.0.0-20200403204345-e1beb1bd0f35 // indirect
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
	sigs.k8s.io/controller-runtime v0.5.2 // indirect
)
