Name: lightutils
Version: 0.8.10
Release: 1%{?dist}
Summary: LightStar's Utility Software
Group: Applications/Communications
License: GPL 3.0
URL: https://github.com/danieldin95/lightstar
BuildRequires: python3 python-virtualenv
Requires: python3

%define _venv /var/lightstar/utils-env
%define _source_dir ${RPM_SOURCE_DIR}/lightstar-%{version}

%description
LightStar's Utility Software

%build
virtualenv -p /usr/bin/python3 %_venv
%_venv/bin/pip install -i https://pypi.tuna.tsinghua.edu.cn/simple "%_source_dir/py"

%install
mkdir -p %{buildroot}/var/lightstar
cp -R /var/lightstar/utils-env %{buildroot}/var/lightstar
mkdir -p %{buildroot}/usr/bin
ln -s /var/lightstar/utils-env/bin/lightutils %{buildroot}/usr/bin/lightutils

%files
%defattr(-,root,root)
/usr/bin/lightutils
/var/lightstar/utils-env