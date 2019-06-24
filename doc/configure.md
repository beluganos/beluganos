# Configure guide
You can configure router options like IP addresses or settings of routing protocol. In this page, the settings method by Linux shell is mainly described. Note that **ansible** and **NETCONF** is also available to configure.

## Pre-requirements

- The installation is required in advance. Please refer [install.md](install.md) before proceeding.
- The setup of Beluganos is required in advance. Please refer [setup.md](setup.md) before proceeding.
- In case you will use OpenNSL or OF-DPA:
	- The installation of OpenNetworkLinux for your white-box switches is required in advance. Please refer [setup-hardware.md](setup-hardware.md) before proceeding.
	- The setup for ASIC API is required in advance. Please refer [setup-onsl.md](setup-onsl.md) or [setup-ofdpa.md](setup-ofdpa.md).

## Comparison of configure method

Currently, **Linux style**, **ansible**, and **NETCONF** over SSH is supported by Beluganos. There are some difference about available datastore of configuration.

| Method           | Startup   | Candidate | Running   | The path of document                         | 
|:-----------------|:----------|:----------|:----------|:---------------------------------------------|
| **Linux style**  | No        | No        | Yes       | (This page.)                                 |
| **ansible**      | Yes       | No        | No (\*1)  | [configure-ansible.md](configure-ansible.md) |
| **NETCONF** (\*2)| Yes       | Yes       | Yes (\*3) | [configure-netconf.md](configure-netconf.md) |

- (\*1) In ansible, in-operation change of configuration is not supported. If you will use ansible, the combination of ansible and Linux style is recommended.
- (\*2) In NETCONF, not all configuration can configure. The supported configuration is described at [yang files](https://github.com/beluganos/netconf/tree/master/etc/openconfig).
- (\*3) In NETCONF, running configuration cannot directly be changed. The operation "commit" is required to reflect running-config.

The **Linux style** is a most simple method to configure Beluganos. In this documents, this method will be described.

## Linux style

### First setup for Linux style

Even if you will NOT use ansible to configure router settings, the initial settings require ansible. Please edit the settings of `mkpb-sample.sh` as your environment and execute it.

```
Beluganos$ cd ~/etc/playbooks
Beluganos$ vi mkpb-sample.sh

REID="10.0.0.5"
DPID=14
DPTYPE=as5812
DAEMONS="zebra ospfd ospf6d ldpd"
```

- `REID`: Configuring IPv4 loopback address is recommended. Keep in mind that you can change loopback address after setup.
	- Detail: The router entity id. This value is used by only Beluganos's component to identify the master routing instance.
- `DPID`: The datapath ID. You already specify it at `fibc.yml` which describes at  [setup.md](setup.md). Because the default settings is `14`, generally you don't have to change it.
- `DPTYPE`: The hardware type name. You can specify each type of methods:
	- Specify hardware name
		- OpenNSL: Currently only `as5812`, `as7712`, and `as7712x4` are supported. If you want to use another switch, please refer "Specify hardware setting file".
		- OF-DPA, OpenFlow Switch, OVS: Set `openflow`.
	- Specify hardware setting file
		- Set .yaml file name. The example is at `mkpb-sample.yaml` which is located at the same directory. The filed of `port_map` should be set the mapping data of physical interface and logical interface. The information of port mapping is also described at [configure-portmapping.md](configure-portmapping.md).
- `DAEMONS`: The routing protocol stack which you want to enable. The name should be described by the name of FRRouting daemons. 
	- The supported protocol is depend on FRRouting. bgpd, ospfd, ospf6d, ripd, ripngd, isisd, ldpd, eigrpd, ... is assumed. IP multicast and VRRP is not supported by Beluganos.
	- If you cannot decide which daemon will be used, enumerating all daemons is accepted.
	- Note that **GoBGP** is always available not depend on this statement.

To reflect,

```
Beluganos$ ./mkpb-sample.sh
Beluganos$ ansible-playbook -i hosts -K lxd-MIC.yaml && lxc stop MIC
```

Note that `mkpb-sample.sh` will be created by `lxd-MIC.yaml` automatically.

- Other notes
	- The container name is `MIC`. This name may be required at [operation.md](operation.md).
	- The initial setup file is under `beluganos/etc/playbooks/roles/lxd/files/MIC/`. If you want to change it directly, please execute `ansible-playbook` command again after editing.

### Start Beluganos

Before configure router settings like IP address or routing protocol, you have to start the process of Beluganos. Required steps are following:

- [Step 1. Run ASIC driver agent](operation.md#step-1-run-asic-driver-agent)
- [Step 2. Run Beluganos main module](operation.md#step-2-run-beluganos-main-module)
- [Step 3. Add Linux containers](operation.md#step-3-add-linux-containers)

Please refer [operation.md](operation.md) for more detail. 

### Configure Beluganos

In Beluganos's architecture, the configuration of LXC will sync to the configuration of white-box switch. Thus, you should configure LXC to use Beluganos.

#### IP address, routing settings

Use FRRouting console.

```
Beluganos$ beluganos con MIC
LXC> vtysh
FRRouting> conf t
```

Cisco IOS like CLI is available. Note that start-up configuration does not automatically saved without `write memory`.For more detail, please refer [FRRouting official document](http://docs.frrouting.org/en/latest/).

If you want to use GoBGP, edit the configuration file.

```
Beluganos$ beluganos con MIC
LXC> vi /etc/frr/gobgpd.conf
```

Please note that reflection of changes in gobgpd.conf requires restarting GoBGP. For more detail, please refer [GoBGP repository](https://github.com/osrg/gobgp).

#### Other features

Please refer feature guide.

- [SNMP](feature-snmp.md): SNMP MIB, SNMP trap
- [syslog](feature-syslog.md): syslog
- [L3VPN](feature-l3vpn.md): MPLS-VPN L3VPN, VRF lite