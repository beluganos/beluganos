# Operation Guide

Beluganos consists of **main module** and **Linux container** which provide routing functions. You can operate both components by `beluganos`. In order to start Beluganos properly, you should start both main module and Linux container.

In this page, operation about NETCONF module is also described. If you will not use NETCONF module, you can ignore some operations.

## Pre-requirements
- Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.
- If you use **ansible** to configure router settings, you have to execute playbook before starting to run Beluganos. Please refer [configure-ansible.md](configure-ansible.md) before proceeding.
- If you use **NETCONF** to configure router settings, you have to execute minimum playbook before starting to run Beluganos. Please refer [configure-netconf.md](configure-netconf.md) before proceeding.
- If you use **SNMP**, please refer [feature-snmp.md](feature-snmp.md) before proceeding.

## Step 1. Run ASIC driver agent

- OpenNSL: Login to OpenNetworkLinux. Please refer [setup-guide-onsl.md](setup-guide-onsl.md)
	- `/etc/init.d/gonsl start`
- OF-DPA: Login to OpenNetworkLinux. Please refer [setup-guide-ofdpa.md](setup-guide-ofdpa.md).
	- `service ofdpa start`
	- `brcm-indigo-ofdpa-ofagent --controller=<BeluganosVM-IP> --dpid=<Agent-dpid>`
- OpenFlow switch: Please refer OpenFlow switch's documents.

## Step 2. Run Beluganos's main module

You have two options to start/stop Beluganos. Generally, in production, because you have already registered Beluganos as a Linux service at [install-guide.md](install-guide.md), **daemon mode** is recommended.

### Daemon mode [Recommended]

To start,

```
$ sudo systemctl start fibcd
$ sudo systemctl start netopeer2-server  # If NETCONF will be used
$ sudo systemctl start ncm.target        # If NETCONF will be used
$ sudo systemctl start snmpd             # If SNMP will be used
$ sudo systemctl start fibsd             # If SNMP will be used
$ sudo systemctl start snmpproxyd-trap   # If SNMP will be used
$ sudo systemctl start snmpproxyd-mib    # IF SNMP will be used
```

To stop,

```
$ sudo systemctl stop fibcd
$ sudo systemctl stop ncmd               # If NETCONF was started
$ sudo systemctl stop ncmi               # If NETCONF was started
$ sudo systemctl stop ncms               # If NETCONF was started
$ sudo systemctl stop snmpproxyd-mib     # If SNMP was started
$ sudo systemctl stop snmpproxyd-trap    # If SNMP was started
$ sudo systemctl stop fibsd              # If SNMP was started
$ sudo systemctl stop snmpd              # If SNMP was started
```

### CLI mode

```
$ sudo beluganos run
```

In this method, Beluganos can be worked in your terminal for debug. This command will snatch your standard output. <kbd><kbd>Ctrl</kbd>+<kbd>C</kbd></kbd> to stop. Note that this mode is not supported in case you will use NETCONF or SNMP feature.

### Remarks
The `fibcd` service and `beluganos run` command will make start to connect with OpenFlow agent. Before executing this commands, starting OF-DPA apps is recommended.

## Step 3. Add Linux containers

After starting main module of Beluganos, you should add Linux containers. In general, one container should be added. Only in VRF environments like MPLS-VPN, you should add all routing instance's containers.

**Important Notice:** When NETCONF is used for configure, you do **NOT** have to manually add Linux containers. Following steps are required only the case which you have configured by [configure-ansible.md](configure-ansible.md).

To add:

```
$ beluganos add <container-name>
```

To delete:

```
$ beluganos del <container-name>
```

You forgot container name? In the procedure of [configure-ansible.md](configure-ansible.md), you should have already specified this name. Please check inventory file (`etc/playbooks/hosts`) or task folder (`etc/playbooks/roles/lxd/tasks/<container-name>.yml`) to remember your container name.

## Confirm

### Confirm main module's status

`fibcd` is one of the main componet name of Beluganos.

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

### Automatic start-up

Automatic startup is just beta version.

```
$ sudo systemctl enable fibcd
$ sudo systemctl enable netopeer2-server  #Only using NETCONF
$ sudo systemctl enable ncm.target        #Only using NETCONF
```

### Command
The command of `beluganos` is created automatically by `create.sh`. The file (python script) is located at `/usr/local/bin/beluganos`.

~~~~
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
~~~~

Note that `sudo systemctl start fibcd` can be used instead of `beluganos start`.

### Port
These port will be occupied by Beluganos. You should not use these port by another applications.

* 830: NETCONF over ssh.
* 6633: OpenFlow for white-box switches.
* 8080: [Rest API](https://github.com/osrg/ryu/blob/master/doc/source/app/ofctl_rest.rst) by Ryu.
