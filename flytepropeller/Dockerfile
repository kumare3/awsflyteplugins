# Using go1.13.3
FROM golang:1.13.3-alpine3.10 as builder
RUN apk add git openssh-client make curl dep

# COPY only the go module files for efficient caching
COPY go.* /go/src/github.com/aws/awsflyteplugins/
WORKDIR /go/src/github.com/aws/awsflyteplugins

# Pull dependencies
RUN go mod download

# COPY the rest of the source code
COPY . /go/src/github.com/aws/awsflyteplugins/

# This 'propeller_linux_compile' target should compile binaries to the /artifacts directory
# The main entrypoint should be compiled to /artifacts/flytepropeller
RUN make propeller_linux_compile

# update the PATH to include the /artifacts directory
ENV PATH="/bin:${PATH}"

# This will eventually move to centurylink/ca-certs:latest for minimum possible image size
FROM alpine:3.10
COPY --from=builder /bin /bin

RUN apk --update add ca-certificates

CMD ["flytepropeller"]
