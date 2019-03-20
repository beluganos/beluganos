# Setup guide for OF-DPA
This document describes about hardware setup to use Beluganos.

## Pre-requirements
- Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.
- In this document, Edge-core AS5812-54X switch is assumed to use Beluganos. All hardware which supports OF-DPA 2.0 is acceptable, but in this case, please look it up yourself to install OF-DPA.

## 1. Install OpenNetworkLinux

### Get binary

It depends on "Deploy style". The detail of deploy style is described at "Pre-requirements" in [install-guide.md](install-guide.md).

#### separated-style

Please get binary from [OpenNetworkLinux's website](https://opennetlinux.org/binaries/). Following version is recommended:

```
ONL-2.0.0-ONL-OS-DEB8-2016-12-22.1828-604af0c-AMD64-INSTALLED-INSTALLER
```

#### embedded-style

Building installer of OpenNetworkLinux is required. Please refer "1. Building OpenNetworkLinux" at [setup-guide-embedded.md](setup-guide-embedded.md) for details.

### Install via TFTP

After getting binary, you can install OpenNetworkLinux via DHCP or TFTP. In this documents, only TFTP methods are described.

#### (Step 1) Connect console cable

In TFTP installation, the connections of console cable is required to communicate between white-box switch and your working PC. In following steps, the strings of `>` represent the console screen.

#### (Step 2) Boot hardware

Plug in power cable to boot switch.

#### (Step 3) Select "ONIE install"

In booting process, GRUB menu is appeared. Select `ONIE` -> `ONIE install` by <kbd>↑</kbd> (Up) or <kbd>↓</kbd> (Down) keys to start install.

#### (Step 4) Stop DHCP discovery

In default settings of ONIE, DHCP discovery will be started. Stop DHCP discovery by following command:

```
> onie-discovery-stop
```

#### (Step 5) Configure management address

Configure the settings of physical outbound port (management port). For example, if you want to set ``172.16.0.57/24``, following step is required:

```
> ifconfig eth0 172.16.0.57 netmask 255.255.255.0 up
```

#### (Step 6) Start to install

Start to install. For example, if your tftp server is `172.16.0.59` and ONL version is `ONL-2.0.0-ONL-OS-2018-01-09.1646-04257be-AMD64-INSTALLED-INSTALLER`, following step is required:

```
> onie-nos-install tftp://172.16.0.59/ONL-2.0.0-ONL-OS-2018-01-09.1646-04257be-AMD64-INSTALLED-INSTALLER
```

#### (Step 7) Login

Once finished a install, it will be rebooted automatically. Please log in. The default user-name to is `root` and password is `onl`.

## 2. Setup for OF-DPA

OF-DPA is a ASIC driver. After installing OpenNetworkLinux, OF-DPA settings are required.

### Get binary

Please get the binary of OF-DPA. The required version of OF-DPA is 2.0.4. The following website is published OF-DPA binary. In this document, using the binary which got from Edge-core's repository is assumed.

- [Edge-core's repository](https://github.com/edge-core/beluganos-forwarding-app)
- [Broadcom's repository](https://github.com/Broadcom-Switch/of-dpa)

### Initial settings

The following steps are required in case of the first time to use OF-DPA.

#### (Step 1) Configure management address

Configure the settings of physical outbound port (management port). For example, if you want to set ``172.16.0.57/24``, following step is required:

```
> ifconfig ma1 172.16.0.57 netmask 255.255.255.0 up
```

This settings are not permanent. To set this settings permanently, following steps are required:

```
> echo 'ip addr add 172.16.0.57/24 dev ma1' >> /mnt/onl/data/rc.boot
> chmod a+x /mnt/onl/data/rc.boot
```

#### (Step 2) Transfer OF-DPA binary

Transfer the binary to OpenNetworkLinux. For example, SCP or SFTP are assumed. Assumed file name is here:

- `ofdpa-2.0-ga_2.0.4.0+accton2.4-1_amd64.deb`

#### (Step 3) Install OF-DPA

```
$ dpkg -i --force-overwrite ofdpa-2.0-ga_2.0.4.0+accton2.4-1_amd64.deb
```

### General settings

Once you finished to do "Initial settings", other general settings are required.

#### (Step 1) Link up required ports

In default settings, almost all physical port is set to down. For example, you want to use port `30`, following step is required:

```
> echo 0 > /sys/bus/i2c/devices/30-0050/sfp_tx_disable
```

If you prefer to up all ports, following example of scripts is recommended.

```
#! /bin/bash

for ((i = 2; i <= 54; i = i + 1)); do
  echo 0 > /sys/bus/i2c/devices/$i-0050/sfp_tx_disable
done
```

#### (Step 2) Start OF-DPA and OpenFlow agent

To start, following commands are required:

```
> service ofdpa start
~~~(Please wait over 15 sec.)~~~
> brcm-indigo-ofdpa-ofagent --controller=<BeluganosVM-IP> --dpid=<Agent-dpid>
```
- `<BeluganosVM-IP>`: Specify Beluganos's IP address. Please note that you already specify this IP address in `create.ini` at [install-guide.md](install-guide.md).
- `<Agent-dpid>`: Specify OpenFlow DPID. Please note that you already specify this ID in `fibc.yml` at [setup-guide.md](setup-guide.md).

## Next steps
Please refer configure guide. You can choose two methods.

- ansible: [configure-ansible.md](configure-ansible.md)
- NETCONF over SSH: [configure-netconf.md](configure-netconf.md)