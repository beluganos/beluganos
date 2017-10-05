# Operation Guide

Beluganos consists of **main module** and **Linux container** which provide routing functions. You can operate both components by `bin/beluganos.py`. In order to start Beluganos properly, you should start main module and Linux container.

## Environments
Before starting setup and operation, please execute `setenv.sh` script to set your environments properly. After executing this script, the strings of "`(mypython)`" will be appeard in your console. Because this settings will be cleared after logout, you should execute this script every login.

 ~~~~
 $ . ./setenv.sh
 ~~~~

## Run main module

You have two options to start/stop Beluganos.

### Option A: Daemon mode

To use daemon mode, some preparation is required.

#### Prepare-1. Edit `.service` file
```
$ cd etc/systemd
$ vi fibcd.service
```

In this file, you have to change about "virtualenv" settings. The "virtualenv" path (`ExecStart`) should be modified as your Linux [username].

```
ExecStart=/home/[username]/mypython/bin/ryu-manager ryu.app.ofctl_rest ${FIBC_APP} --config-file ${FIBC_CONF} --log-config-file ${LOG_CONF}
```

e.g.) If your user name when you operated `create.sh` is "user1",

```
ExecStart=/home/user1/mypython/bin/ryu-manager ryu.app.ofctl_rest ${FIBC_APP} --config-file ${FIBC_CONF} --log-config-file ${LOG_CONF}
```

#### Prepare-2. Copy files

```
$ cd etc/systemd
$ sudo cp fibcd.conf /etc/beluganos/
$ sudo cp fibcd.service /etc/systemd/system/
```

#### Prepare-3. Reload daemon files

```
$ sudo systemctl daemon-reload
$ systemctl status fibcd

● fibcd.service - fib controller service
   Loaded: loaded (/etc/systemd/system/fibcd.service; disabled; vendor preset: enabled)
   Active: inactive (dead)
```

#### Operation

To start,

```
$ beluganos.py start
```

To stop,

```
$ beluganos.py stop
```

### Option B: CLI mode

In this method, Beluganos can be worked in your terminal. `Ctrl-c` to stop.

```
$ beluganos.py run
```

This mode should be used for only debug. In production, using daemon mode is recommended.

## Run Linux containers

After starting main module of Beluganos, you should start Linux containers. In general, one container should be added. Only in VRF environments like MPLS-VPN, you should add all routing instance's containers.


To add:

```
$ beluganos.py add <container-name>
```

To delete:

```
$ beluganos.py del <container-name>
```

You forgot container name? In the procedure of `doc/setup-guide.md`, you should have already specified this name. Please check inventory file (`etc/playbooks/hosts`) or task foloder (`etc/playbooks/roles/lxd/tasks/<container-name>.yml`) to remenber your container name.

## Confirm

### Confirm main module's status

```
$ beluganos.py status
```

### Confirm linux container's status

```
$ beluganos.py status <container-name>
```

### Confirm routing status

In Linux container, you can confirm about routing status. Thanks to Beluganos's architecture, the route table of containers will be synchronized to White-box switches (For more detail, please check `doc/architecture.md`). At first, you may login to Linux container by following command:

```
$ beluganos.py con <container-name>
```

After that, you can execute any commands! For example, in case of container `sample`, you can check routing status by following commands:

```
$ beluganos.py con sample   # Login to containers
# vtysh                     # Open FRR console
```

If you not familiar with FRRouting or quagga, after login `vtysh`, please hit `?` to check command reference. For example, if you want to check RIB,

```
$ beluganos.py con sample
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
```

## Remarks

### beluganos.py
```
$ beluganos.py -h
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
```
`bin/beluganos.py` is created automatically by `create.sh`.

### Port
These port will be utilized by Beluganos.

* 6633: OpenFlow for white-box switches.
* 8080: [Rest API](https://github.com/osrg/ryu/blob/master/doc/source/app/ofctl_rest.rst) by Ryu.
