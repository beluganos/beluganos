# Install guide
<img src="img/environments.png" width="320px">

## 1. Pre-requirements

You can also build Beluganos in the white-box switches, but installing another Ubuntu server is recommended at first.

- Ubuntu server
	- **Ubuntu 16.04 server** is strongly recommended.
	- **Two or more network interfaces** are required.
- White-box switches
	- **[OF-DPA 2.0](https://github.com/Broadcom-Switch/of-dpa/) switch** and OpenFlow agent are required.
	- If you don't have OF-DPA switches, any OpenFlow1.3 switches are acceptable to try Beluganos. In this case, we recommend [Lagopus switch](http://www.lagopus.org/).

## 2. Build
~~~~
$ cd ~
$ git clone https://github.com/beluganos/beluganos/ && cd beluganos/
$ vi create.ini
  FFLOW_MNG_IFACE=ens3             # Set your management interface name for remote login
  FFLOW_OFC_IFACE=ens4             # Set your secure channel interface name connected to switches
  FFLOW_OFC_ADDR=172.16.0.55       # (Optional) You can change FFLOW_OFC_IFACE's IP address if needed
  FFLOW_OFC_MASK=255.255.255.0     # (Optional) You can change FFLOW_OFC_IFACE's subnet mask if needed

$ ./create.sh
~~~~

**For proxy environment only:** If you need to use proxy server to connect Internet, please add proxy settings to `create.ini` before execute `./create.sh`.

~~~~
$ vi create.ini
  PROXY=http://<server-ip>:<server-port>
~~~~

## 3. Set environments
Before starting setup and operation, please execute `setenv.sh` script to set your environments properly. After executing this script, the strings of "`(mypython)`" will be appeard in your console. Because this settings will be cleared after logout, you should execute this script every login.

 ~~~~
 $ . ./setenv.sh
 ~~~~