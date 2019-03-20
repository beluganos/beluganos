# Setup guide for embedded-style

This document should be refereed at embedded-style which Beluganos works into white-box switches, not separated servers.

## 1. Building OpenNetworkLinux

In this section, OpenNetworkLinux (ONL) installer which KVM option is enabled will be created.

### Required files at a glance

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

### Environments

Please refer "Build Hosts and Environments" at [OpenNetworkLinux website](https://opennetlinux.org/docs/build). If you need the proxy to connect internet, please comment out proxy settings at "onl.sh".

```
$ cd ~/beluganos/etc/embedded/onl/
$ vi onl.sh
# PROXY=http://172.16.0.1:8080
```

### Build

The installer will be created at `OpenNetworkLinux/RELEASE` by following commands.

```
$ sudo ./onl.sh install
$ ./onl.sh clone
$ cd OpenNetworkLinux
$ sudo ./onl.sh docker
docker> ./onl.sh build

... (building messages) ...

$ ls RELEASE/jessie/amd64/
ONL-2.0.0_ONL-OS_YYYY-MM-DD.XXXX-04257be_AMD64_INSTALLED_INSTALLER
ONL-2.0.0_ONL-OS_YYYY-MM-DD.XXXX-04257be_AMD64_INSTALLED_INSTALLER.md5sum
ONL-2.0.0_ONL-OS_YYYY-MM-DD.XXXX-04257be_AMD64.swi
ONL-2.0.0_ONL-OS_YYYY-MM-DD.XXXX-04257be_AMD64_SWI_INSTALLER
ONL-2.0.0_ONL-OS_YYYY-MM-DD.XXXX-04257be_AMD64_SWI_INSTALLER.md5sum
ONL-2.0.0_ONL-OS_YYYY-MM-DD.XXXX-04257be_AMD64.swi.md5sum

# Please use `*_INSTALLED_INSTALLER` files.

```

## 2. Deploying KVM images at ONL

In this section, deploying methods of KVM images which installs Beluganos are described.

### Required files at a glance

```
beluganos/
    etc/
        embedded/
            kvm/
                kvm.sh        # setting tools
                domain.xml    # KVM Setting file of guest OS
                networks.xml  # KVM Setting file of network
```

### Install KVM

At OpenNetworkLinux, please install KVM.

```
$ apt install libvirt-bin
$ sudo reboot

... After boot ...
$ virt list
```

### Prepare image files

Please prepare VM image files. This files should be qcow2 format. Please note that following files should be deleted at VM.

```
$ sudo rm /etc/netplan/02-beluganos.yaml
```

### Transmit image files

Please transmit VM image files and other required files to ONL on white-box switches.

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

### Run

```
$ ./kvm.sh network add
$ ./kvm.sh domain add
```

### (Appendix) IP address

After finishing to boot VM, IP address is shown by following commands:

```
$ ./kvm.sh domain list

1448180977 XX:XX:XX:XX:XX:XX 192.168.122.12 * *
```