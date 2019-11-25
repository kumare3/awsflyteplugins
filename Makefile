include flytepropeller/Makefile

.PHONY: build_python
build_python:
	@python setup.py sdist

.PHONY: package
package:
	@package.sh

test_py:
	python setup.py test

.PHONY: release
release:
	./package.sh

.PHONY: generate_protos
generate_protos:
	@scripts/generate_protos.sh

.PHONY: generate
generate: generate_protos
	which pflags || (go get github.com/lyft/flytestdlib/cli/pflags)
	@go generate ./...

.PHONY: test_unit
test_go:
	@go test -cover ./... -race


.PHONY: test_unit_cover
test_unit_cover:
	@go test ./... -coverprofile /tmp/cover.out -covermode=count; go tool cover -func /tmp/cover.out


export AWS_PROFILE=flytedemo
export IMAGE_NAME=awsdemoplugin
PROJECT=aws
DOMAIN=development
VERSION=$(shell git rev-parse HEAD)
ACCOUNT_ID=`aws sts get-caller-identity | jq -r '.Account'`
IMAGE="${ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/${IMAGE_NAME}:${VERSION}"

.PHONY: build_docker
build_demo_docker:
	docker build -t "${IMAGE}" --build-arg DOCKER_IMAGE="${IMAGE}" .

ECR_LOGIN=$(shell aws ecr get-login --no-include-email)
.PHONY: deploy_demo_docker
deploy_demo_docker: build_demo_docker
	@${ECR_LOGIN}
	docker push "${IMAGE}"
	docker run --network host -e FLYTE_PLATFORM_URL='host.docker.internal:8089' ${IMAGE} pyflyte -p aws -d development -c flyte.config register workflows
