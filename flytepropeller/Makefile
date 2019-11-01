export IMAGE_NAME=flytepropeller-stagemaker

.PHONY: linux_compile_propeller
propeller_linux_compile:
	mkdir -p ./bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=auto  go build -o /bin/flytepropeller ./flytepropeller/cmd/main.go

.PHONY: compile_propeller
propeller_compile:
	mkdir -p ./bin
	GO111MODULE=auto go build -o bin/flytepropeller ./flytepropeller/cmd/main.go

.PHONY: docker_build
propeller_docker_build:
	./flytepropeller/scripts/docker_build.sh

.PHONY: propeller_docker_push
propeller_docker_push:
	REGISTRY=docker.io/lyft ./flytepropeller/scripts/docker_build.sh

.PHONY: propeller_run_sandbox
propeller_run_sandbox:
	# Make sure to run (in a separate window/terminal): kubectl port-forward service/minio 9000 -n flyte
	# This will start a background process to port forward. It'll be killed when propeller process is.
	kubectl port-forward service/minio 9000 -n flyte &
	# Make sure to run (in a separate window/terminal): kubectl port-forward service/datacatalog 8089:89 -n flyte
	# This will start a background process to port forward. It'll be killed when propeller process is.
	kubectl port-forward service/datacatalog 8089:89 -n flyte &

	GO111MODULE=auto go run ./flytepropeller/cmd/main.go --config ./flytepropeller/config.yaml
