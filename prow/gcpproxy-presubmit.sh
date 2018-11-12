#!/bin/bash

# Copyright 2018 Google LLC

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Presubmit script triggered by Prow.

# Fail on any error.
set -e
set -u

WD=$(dirname "$0")
WD=$(cd "$WD"; pwd)
ROOT=$(dirname "$WD")
export PATH=$PATH:$GOPATH/bin

cd "${ROOT}"

# golang test
echo '======================================================'
echo '=====================   Go test  ====================='
echo '======================================================'
if [ ! -d "$GOPATH/bin" ]; then
  mkdir $GOPATH/bin
fi
make tools
make depend.install
make test
make clean

# c++ test
echo '======================================================'
echo '===================== Bazel test ====================='
echo '======================================================'
bazel test //src/...