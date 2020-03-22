Name: lightstar
Version: 0.8.03
Release: 1%{?dist}
Summary: LightStar's Project Software
Group: Applications/Communications
License: GPL 3.0
URL: https://github.com/danieldin95/lightstar
BuildRequires: go
Requires: libvirt-daemon libvirt qemu-kvm qemu-img
Conflicts: lightsim

%define _source_dir ${RPM_SOURCE_DIR}/lightstar-%{version}

%description
LightStar's Project Software

%build
cd %_source_dir && make

%install
mkdir -p %{buildroot}/usr/bin
cp %_source_dir/lightstar %{buildroot}/usr/bin/lightstar
cp %_source_dir/lightpix %{buildroot}/usr/bin/lightpix

mkdir -p %{buildroot}/etc/sysconfig
cat > %{buildroot}/etc/sysconfig/lightstar.cfg << EOF
OPTIONS="-static:dir /var/lightstar/static -crt:dir /var/lightstar/ca -conf /etc/lightstar"
EOF

mkdir -p %{buildroot}/usr/lib/systemd/system
cp %_source_dir/packaging/lightstar.service %{buildroot}/usr/lib/systemd/system

mkdir -p %{buildroot}/var/lightstar
cp -R %_source_dir/resource/ca %{buildroot}/var/lightstar
cp -R %_source_dir/http/static %{buildroot}/var/lightstar

mkdir -p %{buildroot}/etc/lightstar
cp -R %_source_dir/resource/*.json.example %{buildroot}/etc/lightstar
mkdir -p %{buildroot}/lightstar/datastore/01

%pre
firewall-cmd --permanent --zone=public --add-port=10080/tcp --permanent || {
  echo "You need allowed TCP port 10080 manually."
}
firewall-cmd --reload || :

%post
[ -e '/etc/lightstar/permission.json' ] || {
  cp -rvf /etc/lightstar/permission.json.example /etc/lightstar/permission.json
}

%files
%defattr(-,root,root)
/etc/lightstar
/etc/sysconfig
/usr/bin/*
/usr/lib/systemd/system/*
/var/lightstar
/lightstar/datastore
