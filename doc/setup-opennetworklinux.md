# Setup guide for OpenNetworkLinux
This document describes about the installation method of OpenNetworkLinux. Some details may depend on your hardware, but the general procedure is described at this page.

## Pre-requirements
- If you chose embedded-style, please mind that the official binary of OpenNetworkLinux can not be applied for Beluganos. For more detail, please refer [setup-hardware.md](setup-hardware.md#1-build-opennetworklinux).

## Install

### (Step 1) Connect console cable

In installation, the connections of console cable is required to communicate between white-box switch and your working PC. In following steps, the strings of `>` represent the console screen.

### (Step 2) Boot hardware

Plug in power cable to boot switch.

### (Step 3) Select "ONIE install"

In booting process, GRUB menu is appeared. Select `ONIE` -> `ONIE install` by <kbd>↑</kbd> (Up) or <kbd>↓</kbd> (Down) keys to start install.

### (Step 4) Stop DHCP discovery

In default settings of ONIE, DHCP discovery will be started. Stop DHCP discovery by following command:

```
> onie-discovery-stop
```

### (Step 5) Configure management address

Configure the settings of physical out-bound port (management port). For example, if you want to set ``172.16.0.57/24``, following step is required:

```
ONL> ifconfig ma1 172.16.0.57 netmask 255.255.255.0 up
```

Not that this settings are not permanent. To set this settings permanently, following steps are also required:

```
ONL> echo 'ip addr add 172.16.0.57/24 dev ma1' >> /mnt/onl/data/rc.boot
ONL> chmod a+x /mnt/onl/data/rc.boot
```

### (Step 6) Start to install

Start to install. For example, if your tftp server is `172.16.0.59` and ONL version is `ONL-2.0.0-ONL-OS-2018-01-09.1646-04257be-AMD64-INSTALLED-INSTALLER`, following step is required:

```
> onie-nos-install tftp://172.16.0.59/ONL-2.0.0-ONL-OS-2018-01-09.1646-04257be-AMD64-INSTALLED-INSTALLER
```

Other protocol is also available like `http://`, `ftp://`, and `/path/to/local`. For more detail, please refer [ONIE's document](https://opencomputeproject.github.io/onie/cli/index.html#onie-nos-install).

### (Step 7) Login

Once finished install, it will be rebooted automatically. Please log in. The default user-name to is `root` and password is `onl`.