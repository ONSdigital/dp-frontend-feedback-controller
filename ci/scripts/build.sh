#!/bin/bash -eux

pushd dp-frontend-feedback-controller
  make build
  cp build/dp-frontend-feedback-controller Dockerfile.concourse ../build
popd
