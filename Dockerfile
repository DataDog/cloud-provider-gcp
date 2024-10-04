# Copyright 2018 The Kubernetes Authors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Based on https://github.com/DataDog/cloud-provider-aws/blob/master/Dockerfile
#
################################################################################
##                               BUILD ARGS                                   ##
################################################################################
# This build arg allows the specification of a custom Golang image.
ARG GOLANG_IMAGE=golang:1.21.5
# Datadog's base docker image
ARG BASE_IMAGE
################################################################################
##                              BUILD STAGE                                   ##
################################################################################
# Build the manager as a statically compiled binary so it has no dependencies
# libc, muscl, etc.
FROM ${GOLANG_IMAGE} as builder
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
WORKDIR /build
RUN mkdir -p /build/_output/
COPY go.mod go.sum ./
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY providers/ providers/
COPY vendor/ vendor/
COPY tools/ tools/
RUN ./tools/update_vendor.sh
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build \
	-trimpath \
	-ldflags="-w -s -X k8s.io/component-base/version.gitVersion=$VERSION" \
	-o /build/_output/gcp-cloud-controller-manager \
	./cmd/gcp-controller-manager/
################################################################################
##                               MAIN STAGE                                   ##
################################################################################
# Copy the manager into the distroless image.
FROM --platform=${TARGETPLATFORM} ${BASE_IMAGE}
COPY --from=builder /build/_output/gcp-cloud-controller-manager /bin/gcp-cloud-controller-manager
ENTRYPOINT [ "/bin/gcp-cloud-controller-manager" ]
