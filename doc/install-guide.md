# Install guide
<!--- <img src="img/environments.png" width="320px"> --->

## 1. Pre-requirements

### Resources
You can also build Beluganos in the white-box switches, but installing another Ubuntu server is recommended at first.

- Ubuntu server
	- **Ubuntu 17.10 server** is strongly recommended.
	- **Two or more network interfaces** are required.
- White-box switches
	- **[OF-DPA 2.0](https://github.com/Broadcom-Switch/of-dpa/) switch** and OpenFlow agent are required. OF-DPA application is also available at [here] (https://github.com/edge-core/beluganos-forwarding-app).
	- If you don't have OF-DPA switches, any OpenFlow 1.3 switches are acceptable to try Beluganos. In this case, we recommend [Lagopus switch](http://www.lagopus.org/).

### LXC settings

If LXC have not configured yet, please set up LXC before starting to build Beluganos. Most of the settings may be default or may be changed if needed. However, please note that following points:

- The new bridge name should be default (lxdbr0).
- The size of new loop device is depend on number of VRF which you will configure. ( Num-of-VRF + 2 ) GB or more is recommended.

~~~~
$ lxd init
Do you want to configure a new storage pool (yes/no) [default=yes]?
Name of the new storage pool [default=default]:
Name of the storage backend to use (dir, btrfs, lvm) [default=btrfs]:
Create a new BTRFS pool (yes/no) [default=yes]?
Would you like to use an existing block device (yes/no) [default=no]?
Size in GB of the new loop device (1GB minimum) [default=15GB]: 8
Would you like LXD to be available over the network (yes/no) [default=no]?
Would you like stale cached images to be updated automatically (yes/no) [default=yes]?
Would you like to create a new network bridge (yes/no) [default=yes]?
What should the new bridge be called [default=lxdbr0]?
What IPv4 address should be used (CIDR subnet notation, “auto” or “non [default=auto]?
What IPv6 address should be used (CIDR subnet notation, “auto” or “non [default=auto]?
LXD has been successfully configured.
~~~~

## 2. Build
You can use building scripts (`create.sh`) for building Beluganos. Before starting scripts, setting file (`create.ini`) should be edited.

~~~~
$ cd ~
$ git clone https://github.com/beluganos/beluganos/ && cd beluganos/
$ vi create.ini
  BELUG_MNG_IFACE=ens3             # Set your management interface name for remote login
  BELUG_OFC_IFACE=ens4             # Set your secure channel interface name connected to switches
  BELUG_OFC_ADDR=172.16.0.55/24    # (Optional) You can change BELUG_OFC_IFACE's IP address and prefix-length if needed

$ ./create.sh
~~~~

### Note: For proxy environment
If you need to use proxy server to connect Internet, please add proxy settings to `create.ini` **before** execute `./create.sh`.

~~~~
$ vi create.ini
  PROXY=http://<server-ip>:<server-port>
~~~~

## 3. Systemctl

Generally, register Beluganos as a linux service is recommended by following commands:

~~~~
$ make install-service
~~~~