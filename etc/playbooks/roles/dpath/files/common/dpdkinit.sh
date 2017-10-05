#! /bin/bash
# -*- coding: utf-8 -*-

. /etc/lagopus/dpdk.conf

export RTE_SDK
export RTE_TARGET

# unbind ifaces
for IFACE in $IFACES; do
    sudo ${RTE_SDK}/tools/dpdk-devbind.py -u $IFACE
done

# unload mosules
sudo rmmod igb_uio
sudo rmmod rte_kni
sudo rmmod uio

sudo umount /mnt/huge

sudo modprobe uio
sudo insmod ${RTE_SDK}/${RTE_TARGET}/kmod/igb_uio.ko
sudo insmod ${RTE_SDK}/${RTE_TARGET}/kmod/rte_kni.ko

sudo sh -c "echo $HUGE_PAGE >  /sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages"
sudo mkdir -p /mnt/huge
sudo mount -t hugetlbfs nodev /mnt/huge

sudo ${RTE_SDK}/tools/dpdk-devbind.py -b igb_uio $IFACES
echo "-- current status ---"
sudo ${RTE_SDK}/tools/dpdk-devbind.py --status

