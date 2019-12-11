# Beluganos ZTP

## Build ONL installer with beluganos ztp-init.

### 1. Get source and apply patch.

run step 0, 1, 2, 3

```
> cd /path/to/build/onl
> cp ~/beluganos/etc/embedded/onl/onl* .
> ./onl.sh

0. edit PROXY value in ./onl.sh
1. sudo ./onl.sh install
2. ./onl.sh clone
3. cd OpenNetworkLinux
--- to build onl ---
4. sudo ./onl.sh docker
5. ./onl.sh build
--- to build opennsl library ---
# after build onl
4. ./onl.sh clone-opennsl
5. sudo ./onl.sh docker
6. ./onl.sh build-opennsl

```

### 2. Copy files

```
export ZTP_DIR=~/beluganos/etc/ztp
export ONL_DIR=`pwd`/builds/any/rootfs/jessie/common/overlay

mkdir -p ${ONL_DIR}/etc/network
mkdir -p ${ONL_DIR}/usr/bin

cp ${ZTP_DIR}/files/ma1               ${ONL_DIR}/etc/network/interfaces.d/ma1
cp ${ZTP_DIR}/files/rc.local          ${ONL_DIR}/etc/rc.local
cp ${ZTP_DIR}/files/ztp-init-kvm.conf ${ONL_DIR}/etc/ztp-init.conf
cp ${ZTP_DIR}/files/ztp-init          ${ONL_DIR}/usr/bin/ztp-init
cp ~/beluganos/bin/ffctl              ${ONL_DIR}/usr/bin/ffctl-ztp

unset ZTP_DIR
unset ONL_DIR
```

### 3. Build ONL installer

run step 4, 5


---

## DHCP/WWW server.

```
> apt install nginx isc-dhcp-server

```

### 1. Setup DHCP server.

```
# see comment in /etc/default/isc-dhcp-server
> sudo vi /etc/default/isc-dhcp-server

INTERFACES="eth0"


# see ~/beluganos/etc/ztp/server/etc/dhcp/dhcpd.conf.sample
> sudo vi /etc/dhcp/dhcpd.conf

... skip ...

option beluganos-kvm-url code 250 = text;
option beluganos-ztp-url code 251 = text;

# ex: WWW server is 172.30.0.1
subnet 172.30.0.0 netmask 255.255.255.0 {
  range 172.30.0.129 172.30.0.250;
  option routers 172.30.0.1;
  option default-url = "http://172.30.0.1/onie-installer";
  option beluganos-kvm-url = "http://172.30.0.1/beluganos-kvm-installer";
  option beluganos-ztp-url = "http://172.30.0.1/beluganos-ztp-installer";

```

### 2 Setup WWW server.

```
export ZTP_DIR=~/beluganos/etc/ztp/server/var/www/html
export WWW_DIR=/var/www/html

mkdir -p ${WWW_DIR}/beluganos-kvm-installer.d

cp ${ZTP_DIR}/beluganos-kvm-installer     ${WWW_DIR}/beluganos-kvm-installer
cp ${ZTP_DIR}/beluganos-kvm-installer.d/* ${WWW_DIR}/beluganos-kvm-installer.d/

cp /path/to/ONL-INSTALLER                 ${WWW_DIR}/onie-installer
cp /path/to/gonsl_1.0.0-1_amd64.deb       ${WWW_DIR}/beluganos-kvm-installer.d/
cp /path/to/ubuntu-wbsw.qcow2             ${WWW_DIR}/beluganos-kvm-installer.d/

cp /path/to/opennsl-accton_3.5.0.3+accton4.0-2_amd64.deb ${WWW_DIR}/beluganos-kvm-installer.d/

pushd ${WWW_DIR}/beluganos-kvm-installer.d/
md5sum *.deb *.xml *.yaml > md5sum.txt
popd

unset ZTP_DIR
unset WWW_DIR

```

## Install by ZTP

Connect Whitebox and DHCP/WWW server, then Power on.
