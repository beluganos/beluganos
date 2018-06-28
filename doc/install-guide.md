# Install guide
This document shows how to install Beluganos in your systems. Automation scripts are prepared.

## Pre-requirements

### Resources
Generally, network OS is installed into white-box switches. In Beluganos, you can also build Beluganos into the white-box switches. However, at first, preparing extra server (or VM) and installing here are recommended at first.

- Ubuntu server
	- **Ubuntu 18.04** (18.04-live-server-amd64) is strongly recommended.
	- **Two or more network interfaces** are required.
	- Some LXC instance will be created. More than **12GB HDD** is recommended.
- White-box switches
	- **[OF-DPA 2.0](https://github.com/Broadcom-Switch/of-dpa/) switch** and OpenFlow agent are required. OF-DPA application in Edge-core switches is also available at [Edge-core's repository](https://github.com/edge-core/beluganos-forwarding-app).
	- If you don't have OF-DPA switches, any OpenFlow 1.3 switches are acceptable to try Beluganos. In this case, [Lagopus switch](http://www.lagopus.org/) is recommended.

### LXC settings

If LXC have not configured yet, please set up LXC before starting to build Beluganos. Most of the settings may be default or may be changed if needed. However, please note that following points:

- The new bridge name should be default ( `lxdbr0` ).
- The size of new loop device is depend on number of VRF which you will configure at most. ( Num-of-VRF + 2 ) GB or more is required.

```
$ sudo lxd init
Would you like to use LXD clustering? (yes/no) [default=no]:
Do you want to configure a new storage pool? (yes/no) [default=yes]:
Name of the new storage pool [default=default]:
Name of the storage backend to use (btrfs, dir, lvm) [default=btrfs]:
Create a new BTRFS pool? (yes/no) [default=yes]:
Would you like to use an existing block device? (yes/no) [default=no]:
Size in GB of the new loop device (1GB minimum) [default=15GB]: 8
Would you like to connect to a MAAS server? (yes/no) [default=no]:
Would you like to create a new network bridge? (yes/no) [default=yes]:
What should the new bridge be called? [default=lxdbr0]:
What IPv4 address should be used? (CIDR subnet notation, “auto” or “none”) [default=auto]:
What IPv6 address should be used? (CIDR subnet notation, “auto” or “none”) [default=auto]:
Would you like LXD to be available over the network? (yes/no) [default=no]:
Would you like stale cached images to be updated automatically? (yes/no) [default=yes]
Would you like a YAML "lxd init" preseed to be printed? (yes/no) [default=no]:
```

## 1. Build
Using shell scripts (`create.sh`) is recommended for building Beluganos. Before starting scripts, setting file (`create.ini`) should be edited for your environments. This script will get the required resources including repository of [beluganos/netconf](https://github.com/beluganos/netconf) automatically.

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
  BELUG_MNG_IFACE=ens3             # Set your management interface name for remote login
  BELUG_OFC_IFACE=ens4             # Set your secure channel interface name connected to switches
  BELUG_OFC_ADDR=172.16.0.55/24    # (Optional) You can change BELUG_OFC_IFACE's IP address and prefix-length if needed

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
You should register your white-box switches (or OpenFlow switches) to Beluganos's main module. Please refer [setup-guide.md](setup-guide.md) for more details.
