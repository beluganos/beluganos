# Operation Guide

Beluganos consists of **main module** and **Linux container** which provide routing functions. You can operate both components by `bin/beluganos`. In order to start Beluganos properly, you should start both main module and Linux container.

In this page, operation about NETCONF module is not described. If you'll use NETCONF to configure Beluganos, please refer [operation-guide-netconf.md](operation-guide-netconf.md) instead of following description.

## Pre-requirements
- Please refer [install-guide.md](install-guide.md) before proceeding.
- Please refer [setup-guide.md](setup-guide.md) before proceeding.

## Run main module

You have two options to start/stop Beluganos. 

### Option A: Daemon mode

To start,

~~~~
$ beluganos start
$ sudo systemctl start netopeer2-server
$ sudo systemctl start ncm.target
~~~~

To stop,

~~~~
$ beluganos stop
$ sudo systemctl stop ncmd
$ sudo systemctl stop ncmi
$ sudo systemctl stop ncms
~~~~

### Option B: CLI mode

This mode should be used for only debug. In production, using daemon mode is recommended.

In this method, Beluganos can be worked in your terminal. 

~~~~
$ beluganos run
~~~~

`Ctrl-c` to stop. 

## Run Linux containers

After starting main module of Beluganos, you should start Linux containers. In general, one container should be added. Only in VRF environments like MPLS-VPN, you should add all routing instance's containers.


To add:

~~~~
$ beluganos add <container-name>
~~~~

To delete:

~~~~
$ beluganos del <container-name>
~~~~

You forgot container name? In the procedure of `doc/setup-guide.md`, you should have already specified this name. Please check inventory file (`etc/playbooks/hosts`) or task foloder (`etc/playbooks/roles/lxd/tasks/<container-name>.yml`) to remenber your container name.

## Confirm

### Confirm main module's status

~~~~
$ beluganos status
~~~~

### Confirm linux container's status

~~~~
$ beluganos status <container-name>
~~~~

### Confirm routing status

In Linux container, you can confirm about routing status. Thanks to Beluganos's architecture, the route table of containers will be synchronized to White-box switches (For more detail, please check `doc/architecture.md`). At first, you may login to Linux container by following command:

~~~~
$ beluganos con <container-name>
~~~~

After that, you can execute any commands! For example, in case of container `sample`, you can check routing status by following commands:

~~~~
$ beluganos con sample   # Login to containers
# vtysh                     # Open FRR console
~~~~

If you not familiar with FRRouting or quagga, after login `vtysh`, please hit `?` to check command reference. For example, if you want to check RIB,

~~~~
$ beluganos con sample
# vtysh

Hello, this is FRRouting (version 3.0-rc2).
Copyright 1996-2005 Kunihiro Ishiguro, et al.

sample# enable
sample# show ip route
Codes: K - kernel route, C - connected, S - static, R - RIP,
       O - OSPF, I - IS-IS, B - BGP, P - PIM, N - NHRP, T - Table,
       v - VNC, V - VNC-Direct,
       > - selected route, * - FIB route

K>* 0.0.0.0/0 via 192.169.1.1, eth0
O   10.0.1.6/32 [110/0] is directly connected, lo, 00:08:06
C>* 10.0.1.6/32 is directly connected, lo
O   10.10.1.4/30 [110/100] is directly connected, eth1, 00:08:06
C>* 10.10.1.4/30 is directly connected, eth1
O   10.10.2.4/30 [110/200] is directly connected, eth2.100, 00:08:06
C>* 10.10.2.4/30 is directly connected, eth2.100

sample# show ip ?
  access-list           List IP access lists
  as-path-access-list   List AS path access lists
  bgp                   BGP information
  community-list        List community-list
  extcommunity-list     List extended-community list
  forwarding            IP forwarding status
  igmp                  IGMP information
  large-community-list  List large-community list
  mroute                IP multicast routing table
  msdp                  MSDP information
  multicast             Multicast global information
  nhrp                  NHRP information
  nht                   IP nexthop tracking table
  ospf                  OSPF information
  pim                   PIM information
  prefix-list           Build a prefix list
  protocol              IP protocol filtering status
  rib                   IP unicast routing table
  rip                   Show RIP routes
  route                 IP routing table
  rpf                   Display RPF information for multicast source
  ssmpingd              ssmpingd operation
sample#
~~~~

## Remarks

### Help
`bin/beluganos` is created automatically by `create.sh`.

~~~~
$ beluganos -h
usage: beluganos.py [-h] [-p [PROFILE]] [-b [BRIDGE]] [-e EXCLUDE [EXCLUDE ...]]
               [-v]
               {run,start,stop,add,del,status,con,init,clear} [container]

positional arguments:
  {run,start,stop,add,del,status,con,init,clear}
  container

optional arguments:
  -h, --help            show this help message and exit
  -p [PROFILE], --profile [PROFILE]
  -b [BRIDGE], --bridge [BRIDGE]
  -e EXCLUDE [EXCLUDE ...], --exclude EXCLUDE [EXCLUDE ...]
  -v, --verbose
~~~~

### Port
These port will be utilized by Beluganos.

* 6633: OpenFlow for white-box switches.
* 8080: [Rest API](https://github.com/osrg/ryu/blob/master/doc/source/app/ofctl_rest.rst) by Ryu.
