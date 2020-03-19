#!/bin/bash

set -e 
set -v

version=$(cat VERSION)
mkdir -p ~/rpmbuild/SOURCES

# update version
sed -i  -e "s/Version:.*/Version:\ ${version}/" ./packaging/lightsim.spec
sed -i  -e "s/Version:.*/Version:\ ${version}/" ./packaging/lightstar.spec

# link source
rm -rf ~/rpmbuild/SOURCES/lightstar-${version}
ln -s $(pwd) ~/rpmbuild/SOURCES/lightstar-${version}
