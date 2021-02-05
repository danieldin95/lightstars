Name: lightsim
Version: 0.8.38
Release: 1%{?dist}
Summary: LightStar's Project Software
Group: Applications/Communications
License: GPL 3.0
URL: https://github.com/danieldin95/lightstar
BuildRequires: go
Conflicts: lightstar

%define _source_dir ${RPM_SOURCE_DIR}/lightstar-%{version}

%description
LightStar's Project Software

%build
cd %_source_dir && make bin

%install
mkdir -p %{buildroot}/usr/bin
cp %_source_dir/build/lightstar %{buildroot}/usr/bin/lightstar

mkdir -p %{buildroot}/usr/lib/systemd/system
cp %_source_dir/packaging/lightsim.service %{buildroot}/usr/lib/systemd/system

mkdir -p %{buildroot}/var/lightstar
cp -R %_source_dir/build/cert/lightstar/cert %{buildroot}/var/lightstar
cp -R %_source_dir/src/http/static %{buildroot}/var/lightstar

mkdir -p %{buildroot}/etc/lightstar
cp -R %_source_dir/packaging/resource/*.json.example %{buildroot}/etc/lightstar

%pre
/usr/bin/firewall-cmd --permanent --zone=public --add-port=10080/tcp --permanent || {
  echo "YOU NEED ALLOWED TCP PORT:10080."
}
/usr/bin/firewall-cmd --reload || :

%post
[ -e '/etc/lightstar/permission.json' ] || {
  cp -rvf /etc/lightstar/permission.json.example /etc/lightstar/permission.json
}

[ -e '/etc/lightstar/auth.json' ] || {
cat > /etc/lightstar/auth.json <<EOF
{
  "admin": {
    "type": "admin",
    "password": "$(/usr/bin/mkpasswd -l 16)"
  },
  "guest": {
    "type": "guest",
    "password": "$(/usr/bin/mkpasswd -l 16)"
  }
}
EOF
}
[ -e '/etc/sysconfig/lightstar.cfg' ] || {
cat > /etc/sysconfig/lightstar.cfg << EOF
OPTIONS="-static:dir /var/lightstar/static -crt:dir /var/lightstar/cert -conf /etc/lightstar"
EOF
}


%files
%defattr(-,root,root)
/etc/lightstar
/usr/bin/*
/usr/lib/systemd/system/*
/var/lightstar
