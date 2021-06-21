#!/usr/bin/env bash

name="vm0"
host="example.com"
iso="centos.iso"
datastore="/lightstar/datastore/01"


ssh-copy-id ${host}

scp ${datastore}/${iso} ${host}:${datastore}/${iso}
ssh ${host} virsh pool-create-as .${name} dir --target ${datastore}/${name} --build

virsh migrate --live ${name} qemu+ssh://${host}/system

ssh ${host} virsh dumpxml ${name} > ${name}.xml && virsh define ${name}.xml