#! /usr/bin/env bash
[ -n "${DEBUG}" ] && set -x
set -euo pipefail

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
source "${script_dir}/vendor_package.sh"

# shellcheck disable=SC2154
golang_version=$( cat "${root_dir}/golang-release/packages/golang-1-linux/version")
stripped_go_version=$(echo "${golang_version}" | cut -d . -f -2)
echo -n "${golang_version}" > version

vendor-package golang-release golang-1-linux "${golang_version}"

echo " - updating mod files with ${stripped_go_version}"
# shellcheck disable=SC2154
find "${autoscaler_dir}" -name go.mod -type f -exec sed -i "s/^[[:space:]]*go 1.*/go ${stripped_go_version}/g" "{}" \;

echo " - updating .tool-version with ${golang_version}"
sed -i "s/golang 1.*/golang ${golang_version}/g" "${autoscaler_dir}/.tool-versions"
