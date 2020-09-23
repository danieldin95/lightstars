#!/bin/bash

set -e

version=$(cat VERSION)
mkdir -p ~/rpmbuild/SOURCES

# update version
sed -i  -e "s/Version:.*/Version:\ ${version}/" ./packaging/light*.spec

# link source
# shellcheck disable=SC2086
rm -rf ~/rpmbuild/SOURCES/lightstar-${version}
ln -s $(pwd) ~/rpmbuild/SOURCES/lightstar-"${version}"
