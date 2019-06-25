# Setup guide for hardware
This document describes about white-box switch's setup for Beluganos. The setup guide about embedded-style is mainly described at here.

## Pre-requirements
- The installation is required in advance. Please refer [install.md](install.md) before proceeding.
- The setup of Beluganos is required in advance. Please refer [setup.md](setup.md) before proceeding.

## Index

The required step is different depend on your deploy style.

### embedded-style

- 1. Build OpenNetworkLinux
- 2. Install OpenNetworkLinux to white-box switches
- 3. Move Beluganos VM to your switch

### separated-style

- 2. Install OpenNetworkLinux to white-box switches

## 1. Build OpenNetworkLinux

**This step is required by only embedded-style**. If you chose separated-style, you can skip this section.

In embedded-style, OpenNetworkLinux (ONL)'s installer which KVM option is enabled is required. If you already had the proper installer, this step may be skipped. However, generally, you should build it in your own.

#### Required files at a glance

```
beluganos/
    etc/
        embedded/
            onl/
                onl.sh            # build tool
                onl.patch         # patch files to activate KVM option
                opennsl_control   # for OpenNSL
                opennsl_postinst  # for OpenNSL
                opennsl_prerm     # for OpenNSL
                opennsl_setup     # for OpenNSL
```

#### Environments

Please refer "Build Hosts and Environments" at [OpenNetworkLinux website](https://opennetlinux.org/docs/build). If you need the proxy to connect internet, please comment out proxy settings at `onl.sh`.

```
$ cd ~/beluganos/etc/embedded/onl/
$ vi onl.sh
# PROXY=http://172.16.0.1:8080
```

#### Building

The installer will be created at `OpenNetworkLinux/RELEASE` by following commands.

```
$ sudo ./onl.sh install
$ ./onl.sh clone
$ cd OpenNetworkLinux
$ sudo ./onl.sh docker
docker> ./onl.sh build

... (building messages) ...

$ ls RELEASE/jessie/amd64/*_INSTALLED_INSTALLER
ONL-2.0.0_ONL-OS_YYYY-MM-DD.XXXX-04257be_AMD64_INSTALLED_INSTALLER

```

Please use `*_INSTALLED_INSTALLER` files to install white-box switches. 

## 2. Install OpenNetworkLinux

The installation of the OpenNetworkLinux will be performed to your white-box switches.

### Get binary

In **embedded-style**, please use binary which was build at the section of [1. Build OpenNetworkLinux](#1-build-opennetworklinux).

In **separated-style**, please get binary from [OpenNetworkLinux's website](https://opennetlinux.org/binaries/). Following version is recommended:

```
ONL-2.0.0-ONL-OS-DEB8-2016-12-22.1828-604af0c-AMD64-INSTALLED-INSTALLER
```

### Install

Depend on your hardware. Please refer your hardware manual. The general procedure is also described at [setup-opennetworklinux.md](setup-opennetworklinux.md).

## 3. Move Beluganos VM to your switch

**This step is required by only embedded-style**. If you chose separated-style, you can skip this section.

In embedded-style, you should move Beluganos VM image to white-box switch from x86 servers. On white-box switch, Beluganos's VM image (.qcow2) will be deployed by KVM.

#### Required files at a glance

The required file is already located your Beluganos VM.

```
beluganos/
    etc/
        embedded/
            kvm/
                kvm.sh        # setting tools
                domain.xml    # KVM Setting file of guest OS
                networks.xml  # KVM Setting file of network
```

#### Install KVM

```
ONL> apt install libvirt-bin
ONL> reboot

... After boot ...
ONL> virsh list
```

### Prepare image files

Please prepare VM image files. This files should be .qcow2 format. Note that deleting network configuration file on Beluganos VM is strongly recommended to avoid the mismatch of the network settings.

```
Beluganos$ sudo rm /etc/netplan/02-beluganos.yaml
```

After that, please confirm that DHCP is enabled at interface `ens3`.

```
Beluganos$ sudo cat /etc/netplan/50-cloud-init.yaml

network:
    ethernets:
        ens3:
            dhcp4: true
    version: 2
```

If not, please create the file of netplan configuration to enable DHCP at interface `ens3`.

### Transfer required files

Please transmit VM image files and other required files to ONL on white-box switches to following directory:

```
/mnt/
    onl/
        data/
            kvm.sh
            domain.xml
            networks.xml
            vmimages/
                ubuntu-wbsw.qcow2
```

Note that the .qcow2 file name should be `ubuntu-wbsw.qcow2`.

### Run

```
ONL> cd /mnt/onl/data/
ONL> ./kvm.sh network add
ONL> ./kvm.sh domain add
```
Note that if `./kvm.sh network add` failed, please execute `./kvm.sh network del` in advance to initialize network.

### Check IP address

After finishing to boot VM, IP address is shown by following commands:

```
ONL> ./kvm.sh domain list

1448180977 XX:XX:XX:XX:XX:XX 192.168.122.12 * *
```

To login to the Beluganos VM on the white-box switches, please use this IP address.

## Next steps

The next steps is setup for ASIC API. You should install OpenNSL or OF-DPA daemon to OpenNetworkLinux (NOT Beluganos VM).

- If you will use **OpenNSL**, please refer [setup-onsl.md](setup-onsl.md).
- If you will use **OF-DPA**, please refer [setup-ofdpa.md](setup-ofdpa.md).
