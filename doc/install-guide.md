# Install guide
This document shows how to install Beluganos in your systems. If you want to try Beluganos, you should read this page at first. Automation scripts are prepared to start easy.

## Pre-requirements

### Deploy style

Generally, network OS is installed into the white-box switches. Beluganos can deploy **into the white-box switches (embedded-style)**, but Beluganos can also deploy at **separated server (separated-style)**. If it is a first time for you to try Beluganos, we recommend **separated-style** because it is easier to deploy. 

1. embedded-style
	- Please prepare x86 physical server and virtual machine (VM) for installation, even if you prefer to select embedded-style. In this case, after installation to VM, you can move VM image (qcow2) from physical server to white-box switches.
2. separated-style
	- Please connect with separated server and the outbound port directly.

### Resources


1. Server
	- Software requirements:
		- **Ubuntu 18.04** (18.04-live-server-amd64) is strongly recommended.
	    - If you use Ubuntu 18.04.1 or later, additional settings are required before proceed. Please check **Appendix A** of this document.
	- Network requirements:
		- separated-style: **Two or more network interfaces** are required.
		- embedded-style: **One or more network interfaces** are required.
	- Storage requirements:
		- Some LXC instance will be created. More than **12GB HDD** is recommended.
		- If you have a plan to use multiple VRF, more HDD is required.
1. White-box switches
	- To use OpenNSL mode, **[OpenNSL 3.5](https://github.com/Broadcom-Switch/OpenNSL) supported switch** is required. OpenNSL agent is included in this repository. OpenNSL application in Edge-core switches is also available at [Edge-core's blog](https://support.edge-core.com/hc/en-us/sections/360002115754-OpenNSL).
	- To use OF-DPA mode, **[OF-DPA 2.0](https://github.com/Broadcom-Switch/of-dpa/) supported switch** and OpenFlow agent are required. OF-DPA application in Edge-core switches is also available at [Edge-core's repository](https://github.com/edge-core/beluganos-forwarding-app).
	- If you don't have real switches, any OpenFlow 1.3 switches are acceptable to try Beluganos. In this case, [Lagopus switch](http://www.lagopus.org/) is recommended.

## 1. Build
Using shell scripts (`create.sh`) is recommended for building Beluganos. Before starting scripts, setting file (`create.ini`) should be edited for your environments. This script will get the required resources including repository of [beluganos/netconf](https://github.com/beluganos/netconf) and [beluganos/go-opennsl](https://github.com/beluganos/go-opennsl) automatically. The internet access is required.

```
$ cd ~
$ git clone https://github.com/beluganos/beluganos/ && cd beluganos/
$ vi create.ini

  #
  # Proxy
  #
  # PROXY=http://<ip>:<port>       # (Optional) Comment out if you need internet proxy server

  #
  # Host
  #
  BELUG_OFC_IFACE=ens4             # Set your secure channel interface name connected to switches
  BELUG_OFC_ADDR=172.16.0.55/24    # (Optional) You can change BELUG_OFC_IFACE's IP address and prefix-length if needed
  # ENABLE_VIRTUALENV=yes

$ ./create.sh
```

## 2. Register as a service

Generally, registering Beluganos's main module as a Linux service is recommended.

```
$ cd ~/beluganos
$ make install-service
```

If you will use NETCONF to configure beluganos, following steps are also required.

```
$ cd ~/netconf
$ sudo make install-service
```

## Next steps

You may choose two options.

### Quick start by example
If you want to try our example cases like [case 1 (IP/MPLS router)](example/case1/case1.md) or [case 2 (MPLS-VPN PE router)](example/case2/case2.md), please get back the example documentations.

### Step-by-step procedure
You should register your white-box switches (or OpenFlow switches) to Beluganos's main module. Please refer [setup-guide.md](setup-guide.md) for more details.


---

## Appendix
### Appendix A. Additional settings at Ubuntu18.04.1 or later

In Ubuntu18.04.01 or later, some settings of apt source are removed. In this case, additional apt source is required to install Beluganos.

```
$ sudo vi /etc/apt/sources.list.d/beluganos.list

deb http://archive.ubuntu.com/ubuntu/ bionic universe
deb http://archive.ubuntu.com/ubuntu/ bionic-updates universe
deb http://archive.ubuntu.com/ubuntu/ bionic multiverse
deb http://archive.ubuntu.com/ubuntu/ bionic-updates multiverse
deb http://security.ubuntu.com/ubuntu bionic-security universe
deb http://security.ubuntu.com/ubuntu bionic-security multiverse

$ sudo apt update
```

### Appendix B. Change the connection settings of white-box switches after installation

If you want to change the white-box switch's settings which specify `BELUG_OFC_IFACE` or `BELUG_OFC_ADDR` at `create.ini` after installation, you can use netplan.

```
$ sudoedit /etc/netplan/02-beluganos.yaml

# -*- coding: utf-8 -*-
network:
  version: 2
  renderer: networkd
  ethernets:
    ens4:  ## <= In case device name was changed
      addresses:
        - 172.16.0.55/24  ## <= In case IP address was changed

```

After editing, to reflect settings, please reboot OS or issue apply command.

```
$ sudo netplan apply
```
