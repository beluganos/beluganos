# Operation guide

Beluganos consists of **main module** and **Linux container** which provide routing functions. You can operate both components by `beluganos` commands. In order to start Beluganos properly, you should start both main module and Linux container.

## Pre-requirements

- The installation is required in advance. Please refer [install.md](install.md) before proceeding.
- The setup of Beluganos is required in advance. Please refer [setup.md](setup.md) before proceeding.
- In case you will use OpenNSL or OF-DPA:
	- The installation of OpenNetworkLinux for your white-box switches is required in advance. Please refer [setup-hardware.md](setup-hardware.md) before proceeding.
	- The setup for ASIC API is required in advance. Please refer [setup-onsl.md](setup-onsl.md) or [setup-ofdpa.md](setup-ofdpa.md).
- The startup configuration is required in advance. Please refer [configure.md](configure.md) or [configure-ansible.md](configure-ansible.md) or [configure-netconf.md](configure-netconf.md).

## Overview

To start, you should run 3 components in a fixed order.

- Step 1. Run ASIC driver agent
- Step 2. Run Beluganos's main module
- Step 3. Run Linux containers

To stop or restart Beluganos, please stop it in the reverse order of startup.

## Step 1. Run ASIC driver agent

All commands should be executed at OpenNetworkLinux.

### OpenNSL

```
ONL> /etc/init.d/gonsl start
```

### OF-DPA

``` 
ONL> service ofdpa start
ONL> brcm-indigo-ofdpa-ofagent --controller=<BeluganosVM-IP> --dpid=<Agent-dpid>
```

### OpenFlow switch

Please refer your OpenFlow switch's documents.

## Step 2. Run Beluganos main module

You have two options to start/stop Beluganos. Generally, **daemon mode** is recommended in production. All commands should be executed at Beluganos's VM.

### Daemon mode [Recommended]

To start,

```
$ beluganos start
$ sudo systemctl start netopeer2-server  # If NETCONF will be used
$ sudo systemctl start ncm.target        # If NETCONF will be used
$ sudo systemctl start snmpd             # If SNMP will be used
$ sudo systemctl start fibsd             # If SNMP will be used
$ sudo systemctl start snmpproxyd-trap   # If SNMP will be used
$ sudo systemctl start snmpproxyd-mib    # IF SNMP will be used
```

To stop,

```
$ beluganos stop
$ sudo systemctl stop ncmd               # If NETCONF was started
$ sudo systemctl stop ncmi               # If NETCONF was started
$ sudo systemctl stop ncms               # If NETCONF was started
$ sudo systemctl stop snmpproxyd-mib     # If SNMP was started
$ sudo systemctl stop snmpproxyd-trap    # If SNMP was started
$ sudo systemctl stop fibsd              # If SNMP was started
$ sudo systemctl stop snmpd              # If SNMP was started
```

### Debug mode (CLI mode)

```
$ beluganos run
```

In this method, Beluganos can be worked in your terminal for debug. This command will snatch your standard output. <kbd><kbd>Ctrl</kbd>+<kbd>C</kbd></kbd> to stop. Note that this mode is not supported in case you will use NETCONF or SNMP feature.

### Remarks
The `beluganos start` and `beluganos run` command will make start to connect with OpenNSL or OpenFlow agent. Before executing this commands, starting ASIC driver agent is required.

## Step 3. Run Linux containers

After starting main module of Beluganos, you should add Linux containers. In general, one container should be added. Only in VRF environments like MPLS-VPN, you should add all routing instance's containers.

*Important Notice:* When **NETCONF** is used for configure, you do **NOT** have to manually add Linux containers. This step is required by only the case which you have configured by [configure.md](configure.md) or [configure-ansible.md](configure-ansible.md).

To start:

```
$ beluganos add <container-name>
```

To stop:

```
$ beluganos del <container-name>
```

- Note:
	- You forgot container name? In the **Linux style**, `MIC` is used for default container name. Besides, in **ansible**, you should have already specified this name. Please check inventory file (`etc/playbooks/hosts`) or task folder (`etc/playbooks/roles/lxd/tasks/<container-name>.yml`) to confirm your container name.

## Confirm

### Confirm main module's status

`fibcd` is one of the main component name of Beluganos.

```
$ sudo systemctl status fibcd
 fibcd.service - fib controller service
   Loaded: loaded (/etc/systemd/system/fibcd.service; disabled; vendor preset: enabled)
   Active: active (running) since Tue 2018-05-08 20:50:26 JST; 2min 7s ago
```

### Confirm Linux container's status

```
$ beluganos status <container-name>
```

### Confirm routing status

In Linux container, you can confirm about routing status. Thanks to Beluganos's architecture, the route table of containers will be synchronized to White-box switches (For more detail, please check [architecture.md](architecture.md)). At first, you may login to Linux container by following command:

```
$ beluganos con <container-name>
```

After that, you can execute any commands! For example, in case of container `sample`, you can check routing status by following commands:

```
$ beluganos con sample   # Login to containers
# vtysh                  # Open FRR console
```

If you not familiar with FRRouting (or quagga), after login `vtysh`, please hit `?` to check command reference. For example, if you want to check RIB (Routing Information Base),

```
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
```

## Appendix

### Beluganos command
The command of `beluganos` is created automatically by `create.sh`. The file (python script) is located at `/usr/local/bin/beluganos`.

```
$ beluganos -h
usage: beluganos [-h] [-p [PROFILE]] [-b [BRIDGE]] [-e EXCLUDE [EXCLUDE ...]]
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

Note that `sudo systemctl start fibcd` can be used instead of `beluganos start`.

### Interface status

In Beluganos, interface status (up/down) of physical white-box switches will be synchronized with LXC. Thus, you can monitor the status of physical interfaces by `ip link` command at container. Moreover, if you want to down it administratively, `ip link set` command is available.

```
LXC> ip link
LXC> ip link set <if-name> down
```

Note that interface status and traffic counter are available at SNMP features. For more detail, please refer [feature-snmp.md](feature-snmp.md).

### Port
These port will be occupied by Beluganos. You should not use these port by another applications.

* 830: NETCONF over ssh.
* 6633: OpenFlow for white-box switches.
* 8080: [Rest API](https://github.com/osrg/ryu/blob/master/doc/source/app/ofctl_rest.rst) by Ryu.
