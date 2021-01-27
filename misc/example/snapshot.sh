virsh snapshot-create-as --domain os.cmp01 --name os-queens
virsh snapshot-revert --domain os.cmp01  os-queens
