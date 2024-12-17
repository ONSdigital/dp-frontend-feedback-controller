#!/bin/bash -eux

pushd dp-frontend-feedback-controller
  go install github.com/kevinburke/go-bindata/go-bindata
  make test-component
popd
