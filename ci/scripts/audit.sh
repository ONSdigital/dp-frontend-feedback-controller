#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-frontend-feedback-controller
  make audit
popd  