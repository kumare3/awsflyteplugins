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


export IMAGE_NAME=awsdemoplugin
PROJECT=aws
DOMAIN=development
VERSION=$(shell git rev-parse HEAD)
IMAGE="docker.io/aws/${IMAGE_NAME}:${VERSION}"

.PHONY: build_docker
build_demo_docker:
	docker build -t "${IMAGE}" --build-arg DOCKER_IMAGE="${IMAGE}" .
