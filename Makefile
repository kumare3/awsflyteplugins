include flytepropeller/Makefile

.PHONY: build_python
build_python:
	@python setup.py sdist bdist_wheel

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
export IMAGE_NAME=awsflyte
PROJECT=aws
DOMAIN=development
VERSION=$(shell git rev-parse HEAD)
ACCOUNT_ID=`aws sts get-caller-identity | jq -r '.Account'`
#IMAGE="${ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/${IMAGE_NAME}:${VERSION}"
IMAGE="bnsblue/${IMAGE_NAME}:${VERSION}"

.PHONY: build_docker_manual
build_docker_manual:
	docker build -t "${IMAGE}" --build-arg DOCKER_IMAGE="${IMAGE}" .

.PHONY: deploy_docker_manual
deploy_docker_manual: build_docker_manual
	docker push "${IMAGE}"
