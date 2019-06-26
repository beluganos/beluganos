# [Feature guide] L2 switching

Beluganos support L2 switching function provided by hardware processing. The MAC learning is also processed by hardware. This document describes the way how to configure L2 switching to Beluganos.

## Pre-requirements

- The installation is required in advance. Please refer [install.md](install.md) before proceeding.
- The setup of Beluganos is required in advance. Please refer [setup.md](setup.md) before proceeding.
- In case you will use OpenNSL or OF-DPA:
   - The installation of OpenNetworkLinux for your white-box switches is required in advance. Please refer [setup-hardware.md](setup-hardware.md) before proceeding.
   - The setup for ASIC API is required in advance. Please refer [setup-onsl.md](setup-onsl.md) or [setup-ofdpa.md](setup-ofdpa.md).

## Important notice

### Reserved VLAN range

The VLAN 3900-4094 is reserved by Beluganos. The VLAN configuration range supported by Beluganos is 2-3899 in default. But you can also change the reserved VLAN range by setting OpenNSL agent. Please check [appendix](#configure-the-reserved-vlan-range).

### Configure notice

- When you handle the IF configuration to bridge in (2), please set target IF down manually before. After complete the IF configuration to bridge, please set target IF up.
- Do NOT change L2 configuration in container without `fibc` running. It will cause container error with massive log file and make Beluganos's VM disk full. If you en-count this problem, please deal with the following procedure.
   1. Delete all IF on bridge in container
   1. Delete syslog files
   1. Reboot container
- The bridge configuration is not persisted. If you want to persist all settings, please also check [appendix](#persistence-of-bridge-configuration).

## Setup

### Add L2 switch configuration

There are two types of configuration method for L2 switch, by "[Linux command]" **or** by "[ffctl]" command. You can choose any configuration method you like.

Configuration needs three steps.

- (1) Configure Linux bridge to container
- (2) Configure L2IF to Linux bridge
- (3) Configure VLAN setting to L2IF

Sample environment is here:

```
           l2swbr0          (1) Configure Linux bridge to container
              |
     +-----+--+--+------+   (2) Configure L2IF to Linux bridge
     |     |     |      |
    eth1  eth2  eth3  eth4  (3) Configure VLAN setting to L2IF

- eth1: access port. vlan 10
- eth2: access port. vlan 20
- eth3: access port. vlan 20
- eth4: trunk  port. vlan 10,20
```

#### (1) Configure Linux bridge to container

Note that you can configure any bridge name, but you are allowed to make only one bridge. Bridge name is `l2swbr0` in sample environment.

```
[Linux command]
LXC> ip link add l2swbr0 type bridge vlan_filtering 1
LXC> ip link add l2swbr0 multicast off
LXC> ip link set l2swbr0 up


[ffctl]
LXC> ffctl bridge vlan add-br l2swbr0
```

#### (2) Configure L2IF to Linux bridge

```
[Linux command]
# eth1
LXC> ip link set eth1 down
LXC> ip link set eth1 master l2swbr0
LXC> ip link set eth1 up

# eth2
LXC> ip link set eth2 down
LXC> ip link set eth2 master l2swbr0
LXC> ip link set eth2 up

# eth3
LXC> ip link set eth3 down
LXC> ip link set eth3 master l2swbr0
LXC> ip link set eth3 up

# eth4
LXC> ip link set eth4 down
LXC> ip link set eth4 master l2swbr0
LXC> ip link set eth4 up


[ffctl]
LXC> ffctl bridge vlan add-ports l2swbr0 eth1 eth2 eth3 eth4
```

#### (3) Configure VLAN setting to L2IF

```
[Linux command]
# eth1 => access port, vlan 10
LXC> bridge vlan add vid 10 dev eth1 pvid untagged
LXC> bridge vlan del vid 1  dev eth1 pvid untagged

# eth2 => access port, vlan 20
LXC> bridge vlan add vid 20 dev eth2 pvid untagged
LXC> bridge vlan del vid 1  dev eth2 pvid untagged

# eth3 => access port, vlan 20
LXC> bridge vlan add vid 20 dev eth3 pvid untagged
LXC> bridge vlan del vid 1  dev eth3 pvid untagged

# eth4 => trunk port, vlan 10 20
LXC> bridge vlan add vid 10 dev eth4
LXC> bridge vlan add vid 20 dev eth4
LXC> bridge vlan del vid 1  dev eth4


[ffctl]
LXC> ffctl bridge vlan add-access 10 eth1
LXC> ffctl bridge vlan add-access 20 eth2 eth3
LXC> ffctl bridge vlan add-trunk eth4 10 20
```

### Delete L2 switch configuration

```
           l2swbr0         (3) Delete Linux bridge to container
              |
     +-----+--+--+------+  (2) Delete L2IF to Linux bridge
     |     |     |      |
    eth1  eth2  eth3  eth4 (1) Delete VLAN setting to L2IF
```

#### (1) Delete VLAN setting to L2IF

```
[Linux command]
# eth1
LXC> bridge vlan del vid 10 dev eth1 pvid untagged
LXC> bridge vlan add vid 1  dev eth1 pvid untagged

# eth2
LXC> bridge vlan del vid 20 dev eth2 pvid untagged
LXC> bridge vlan add vid 1  dev eth2 pvid untagged

# eth3
LXC> bridge vlan del vid 20 dev eth3 pvid untagged
LXC> bridge vlan add vid 1  dev eth3 pvid untagged

# eth4
LXC> bridge vlan del vid 10 dev eth4
LXC> bridge vlan del vid 20 dev eth4
LXC> bridge vlan add vid 1  dev eth4 pvid untagged


[ffctl]
LXC> ffctl bridge vlan del-access 10 eth1
LXC> ffctl bridge vlan del-access 20 eth2 eth3
LXC> ffctl bridge vlan del-trunk eth4 10 20
```

#### (2) Delete L2IF to Linux bridge

```
[Linux command]
# eth1
LXC> ip link set eth1 down
LXC> ip link set eth1 nomaster
LXC> ip link set eth1 up

# eth2
LXC> ip link set eth2 down
LXC> ip link set eth2 nomaster
LXC> ip link set eth2 up

# eth3
LXC> ip link set eth3 down
LXC> ip link set eth3 nomaster
LXC> ip link set eth3 up

# eth4
LXC> ip link set eth4 down
LXC> ip link set eth4 nomaster
LXC> ip link set eth4 up


[ffctl]
LXC> ffctl bridge vlan del-ports l2swbr0 eth1 eth2 eth3 eth4
```

#### (3) Delete Linux bridge to container

```
[Linux command]
LXC> ip link del l2swbr0 type bridge


[ffctl]
LXC> ffctl bridge vlan del-br l2swbr0
```

## Operation

### Confirm switching configuration

```
[Linux command]
LXC> bridge vlan


[ffctl]
LXC> ffctl bridge vlan show
```

### Confirm MAC address table

MAC address table information in white-box switch hardware is synchronized to `fdb` in the container on the Beluganos.

```
[Linux command]
# Linux command shows not only L2 unicast entry but also multicast entry, but multicast entry is NOT configured to white-box switch hardware.

LXC> bridge fdb

33:33:00:00:00:01 dev l2swbr0 self permanent
01:00:5e:00:00:01 dev l2swbr0 self permanent
33:33:ff:3b:b1:5b dev l2swbr0 self permanent
33:33:00:00:00:01 dev eth0 self permanent
01:00:5e:00:00:01 dev eth0 self permanent
00:11:11:11:11:11 dev eth1 vlan 10 master l2swbr0 permanent
00:11:11:11:11:11 dev eth2 vlan 20 master l2swbr0 permanent
00:44:44:44:44:44 dev eth4 vlan 10 master l2swbr0 permanent


[ffctl]
LXC> ffctl bridge fdb show
```

### Delete MAC address entry manually

```
[Linux command]
LXC> bridge fdb del 00:11:11:11:11:11 dev eth1 vlan 10 master

[ffctl]
LXC> ffctl bridge fdb del 00:11:11:11:11:11 10 eth1
```

Note that learning MAC address will be automatically aged. The settings of aging timer is described at [appendix](#configure-the-mac-agint-timer).

## Appendix

### Configure the reserved VLAN range

```
ONL> vi /etc/beluganos/gonsld.yaml

dpaths:
  default:
    dpid: 14
    addr: 172.16.0.1
    port: 50070
    l2sw:
      aging_sec: 3600
      sweep_sec: 3
      notify_limit: 256

    block_bcast:
	  # min, max => ports range
	  # base_vid => base vlan id for L3 routing.
	  #             VLAN IDs(L2): 2...(base_vid+min-1)
	  #             VLAN IDs(L3): (base_vid+min)...(base_vid+max)
      range: { min: 1, max: 190, base_vid: 3900 }
```

To reflect, restart `gonsld`.

### Configure the MAC aging timer

```
ONL> vi /etc/beluganos/gonsld.yaml

dpaths:
  default:
    dpid: 14
    addr: 172.16.0.1
    port: 50070
    l2sw:
      aging_sec: 3600    # => aging timeout
      sweep_sec: 3       # => do not change
      notify_limit: 256  # => do not change
```

To reflect, restart `gonsld`.

### Persistence of bridge configuration

Linux bridge configurations are not persisted in default. If you want to persisted the configurations, editing setting file is required.

```
LXC> vi /etc/beluganos/bridge_vlan.yaml
# playbook: etc/playbooks/roles/lxd/files/<container-name>/bridge_vlan.yaml

---

# network: {}

network:
  vlans:
    eth1:
      # access port
      id: 10
    eth2:
      # access port
      id: 20
    eth3:
      # trunk port
      ids: [10, 20]

  bridges:
    l2swbr0:  # => bridge name
      vlan_filtering: 1
      interfaces:
        - eth1
        - eth2
        - eth3

```

To reflect,

```
LXC> systemctl restart ribbr.service
```