#!/bin/bash -eux

go install github.com/kevinburke/go-bindata/v4/...@v4.0.2

pushd dp-frontend-feedback-controller
  make build
  cp build/dp-frontend-feedback-controller Dockerfile.concourse ../build
popd
