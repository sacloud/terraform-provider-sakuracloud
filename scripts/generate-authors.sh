#!/usr/bin/env bash
# Copyright 2016-2019 terraform-provider-sakuracloud authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

# see also ".mailmap" for how email addresses and names are deduplicated
echo "Kazumichi Yamamoto <yamamoto.febc@gmail.com>" >> .mailmap

{
	cat <<-'EOH'
	# This file lists all individuals having contributed content to the repository.
	# For how it is generated, see `scripts/generate-authors.sh`.
	EOH
	echo
	git log --format='%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf
} > AUTHORS
