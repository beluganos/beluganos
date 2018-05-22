# Config guide (ansible)
You can configure router options like IP addresses or settings of routing protocol by **ansible** as you like.

## Pre-requirements
- Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.
- In ansible, only startup config can be edited. If you already start Beluganos by `beluganos run` or `beluganos start` command, please stop before starting ansible playbooks.
- If you will configure L3VPN routers, please refer not only this documents but also [configure-ansible-l3vpn.md](configure-ansible-l3vpn.md).

## Config files at a glance
The files under `etc/playbooks` are configuration files. In this page, `lxd-*.yml` and the files under `roles/lxd` are described.

~~~~
beluganos/
    etc/
        playbooks/
            hosts                            # Container name (inventory)
            lxd-sample.yml                   # Sample playbook of container basic settings
            lxd-<group-name>.yml             # Container basic settings (playbook)
            roles/
                lxd/
                    vars/                    # Do NOT edit
                    handlers/                # Do NOT edit
                    tasks/                   # Do NOT edit
                    files/
                        common/              # Do NOT edit
                        sample/              # Sample files of router settings
                        <container-name>/    # Router settings
                            lxd_profile.yml
                            netplan.yaml
                            sysctl.conf
                            daemons
                            frr.conf
                            gobgp.comnf
                            gobgpd.conf
                            fibc.yml
                            ribxd.conf
~~~~

## Settings for Linux Containers

### 1. Container name (inventory)

At first, you should edit [inventory file](http://docs.ansible.com/ansible/latest/intro_inventory.html) at `etc/playbooks/hosts`. The container name and group name should be described here. This name will be only used internal configurations.

In general, one-to-one correspondence between group name and container name is assumed except VRF environments. For example, if your group name is "*lxd-group*" and your container name (i.e. routing instance name) is "*master*", please edit as follows:

```
$ cd ~/beluganos/etc/playbooks/
$ vi hosts
---

[hosts]      # DO NOT EDIT.
localhost    # DO NOT EDIT.

[lxd-group]
master
```

**For VRF environments only:** To configure multiple routing-instances, you should describe multiple container name.

```
[lxd-group-vpn]
master
vrf10
vrf20
```

### 2. Container basic settings (playbook)

**The sample playbook is `etc/playbooks/lxd-sample.yml`**. You may copy this file and rename to your group name. For example, if your group name at inventory is "*lxd-group*", please edit following:

```
$ cd ~/beluganos/etc/playbooks
$ cp lxd-sample.yml lxd-group.yml
$ vi lxd-group.yml

---

- hosts: hosts
  connection: local
  vars:
    bridges: []
  roles:
    - bridge
  tags:
    - bridge

- hosts: hosts
  connection: local
  roles:
    - { role: lxd, mode: host }
  tags:
    - host

- hosts: hosts
  connection: local
  tasks:
    - include_role:
        name: lxd
      vars:
        mode: create
      with_items:      # Edit "lxd-group" to your group name in inventry file.
        - "{{ groups['lxd-group'] }}"
      loop_control:
        loop_var: lxcname
  tags:
    - create
    - lxd

- hosts: lxd-group     # Edit "lxd-group" to your setting in inventry file.
  connection: lxd
  roles:
    - { role: lxd, lxcname: "{{ inventory_hostname }}", mode: setup }
  tags:
    - setup
    - lxd
```

### Reflect changes

The reflection of "Settings for Linux Containers" should be done after editting "Router settings".

## Router settings

**The sample files are under `etc/playbboks/roles/lxd/files/sample`**. You have copied this folder in previous. If your container name is *master*, you can copy by folllowing commands:

```
$ cd ~/beluganos/etc/playbooks/roles/lxd/files
$ cp -r sample master
$ ls master
daemons fibc.yml frr.conf gobgp.conf gobgpd.conf netplan.yaml ribxd.conf sysctl.conf
```

### 1. lxd-profile.yml

In this files, you can determine **interface name of your switch** like "*eth1*". The sample is following:

```
- name: create profile
  lxd_profile:
    name: "{{ lxcname }}"
    state: present
    config: {"security.privileged": "true"}
    devices:
      eth0:                            # For internal use only.
        name: eth0
        nictype: bridged
        parent: lxdbr0
        type: nic
      eth1:                            # Interface name in Linux container
        type: nic
        name: eth1                     # Interface name in Linux container
        host_name: "{{ lxcname }}.1"   # Interface name in Host OS
        nictype: p2p
      eth2:
        type: nic
        name: eth2
        host_name: "{{ lxcname }}.2"
        nictype: p2p
      eth3:
        type: nic
        name: eth3
        host_name: "{{ lxcname }}.3"
        nictype: p2p
      eth4:
        type: nic
        name: eth4
        host_name: "{{ lxcname }}.4"
        nictype: p2p
      root:
        path: /
        pool: default
        type: disk
```

The syntax is following:

```
- name: create profile
  lxd_profile:
    name: "{{ lxcname }}"
    state: present
    config: {"security.privileged": "true"}
    devices:
      eth0:
        name: eth0
        nictype: bridged
        parent: lxdbr0
        type: nic
      <lxc-interface-name>:
        type: nic
        name: <lxc-interface-name>
        host_name: "{{ lxcname }}.<lxc-interface-number>"
        nictype: p2p

      ....
      
      root:
        path: /
        pool: default
        type: disk
```

- devices
	- `eth0`: The management interface for internal use. Do NOT edit.
	- `<lxc-interface-name>`: The interface name of Linux container. This interface name should be unique between other containers.
		- type: nic. Do NOT edit.
		- name (`<lxc-interface-name>`): The interface name of Linux container. Should be same value as this section.
		- host_name (`{{ lxcname }}.<lxc-interface-number>`): Beluganos's main module recognize container's interface by this value. `<lxc-interface-number>` should be a sequential number starting from 1 for each container.
		- nictype: p2p. Do NOT edit.
	- `root`: Strage pool settings for container. Do NOT edit.

### 2. fibc.yml: global settings of this router

In `fibc.yml`, you may set general settings. FIBC is one of the main components of Beluganos.

Especially, you can determine about **interface mapping between your white-box switches and Linux containers**. In Beluganos's architecture, the path information which was installed to Linux containers will be copied to white-box switches. This is because you have to configure interface mapping between switches and containers.

Please note that only physical interface should be described at `fibc.yml`. Logical (VLAN) settings should be described another files (`netplan.yml`).

```
routers:
  - desc: sample-router
    re_id: 10.0.1.6                                # Router entity id (=router-id)
    datapath: switchA                              # dpname
    ports:
      - { name: eth1,     port: 1 }                # nid=0, lxc-iface=eth1, dp-port=1
      - { name: eth2,     port: 2 }                # nid=0, lxc-iface=eth2, dp-port=2
```

Syntax is following:

```
routers:
  - desc: <description>
    re_id: <router-entity-id>
    datapath: <switch-name>
    ports:
      - { name: <lxc-interface-name>, port: <switch-interface-number> } # For physical interface
```

- desc (`<description>`): Your preferable descriptions.
- re_id (`<router-entity-id>`): Router entity id for internal use. In general, router-id is recomended.
	- **For VRF environments only**: Please set **master instance's** router-id for all routing instances.
- datapath (`<switch-name>`): Your switch name which is declared at `etc/playbooks/dp-sample.yml`. Please refer [setup-guide.md](setup-guide.md).
- ports
    - name (`<lxc-interface-name>`): Specify interface name in Beluganos's Linux container.
    	- Note: `<lxc-interface-name>` should be matched with the value in `lxd_profile.yml`.
    - port (`<switch-interface-number>`): Specify interface number of white-box switches.


### 3. netplan.yaml: interfaces configurations

In latest Ubuntu, netplan is used for interface configurations. This file will be located to `/etc/netplan/02-beluganos.yaml` in linux container.

For example, you can configure subinterface (VLAN) settings by this file. "eth3", "eth3.10" and "eth4.20" will be configured in following examples:

```
---

network:
  version: 2
  ethernets:
    eth3:
      mtu: 1500    # note: MTU is not supported yet by netplan.
  vlans:
    eth3.10:
      link: eth3   # Physical interface
      id: 10       # VLAN-ID
    eth4.20:
      link: eth4   # Physical interface
      id: 20       # VLAN-ID
```

For more detail, please refer [netplan design](https://wiki.ubuntu.com/Netplan/Design).

### 4. sysctl.conf: MPLS basic configurations: sysctl.conf

This file will be located under `/etc/sysctl.d/` in Linux container. If you want to enable MPLS for all interfaces, you have to disable `rp_filter` and set `net.mpls.conf` as follows:

```
# -*- coding: utf-8 -*-

# Max of MPLS label.
net.mpls.platform_labels=10240

# Disable rp filter
net.ipv4.conf.default.rp_filter=0
net.ipv4.conf.all.rp_filter=0

# Disable rp filter
net.ipv4.conf.eth1.rp_filter=0
net.ipv4.conf.eth2/100.rp_filter=0

# Enable MPLS
net.mpls.conf.eth1.input=1
net.mpls.conf.eth2/100.input=1
```

### 5. daemons: Routing protocols which you want to use

This is FRRouting setting. Please note that IP Multicast routing is not supported yet.

```
zebra=yes
bgpd=no
ospfd=yes
ospf6d=no
ripd=no
ripngd=no
isisd=no
pimd=no
ldpd=yes
nhrpd=no
```

**For MPLS-VPN environments only**: In MPLS-VPN, only GoBGP is supported. Do NOT enable `bgpd` of FRRouting.

### 6. frr.conf: Each routing protocol settings

Each daemon configuration of FRRouting.
 See [FRRouting user-guide](https://frrouting.org/user-guide/).

### 7. gobgp.conf: GoBGP daemon settings

This file is about gobgp booting, **NOT gobgpd itself**. In general, please do not edit.

```
# -*- coding: utf-8 -*-

CONF_PATH = /etc/frr/gobgpd.conf
CONF_TYPE = toml
LOG_LEVEL = debug
PPROF_OPT = --pprof-disable
API_HOSTS = 127.0.0.1:50051
```

### 8. gobgpd.conf: GoBGP itselfs configurations

GoBGP configurations. See [GoBGP configurations.md](https://github.com/osrg/gobgp/blob/master/docs/sources/configuration.md).

**For MPLS-VPN environments only**: In MPLS-VPN, some optional settings about GoBGP is required. For more detail, please check `doc/setup-guide-l3vpn.md`.

### 9. ribxd.conf: Beluganos's settings

This is the main configuration file of Beluganos itself. For example, following configuration will be assumed:

```
# -*- coding: utf-8; mode: toml -*-

[node]
nid   = 0
reid  = "10.0.1.6"
label = 100000
#allow_duplicate_ifname = false

[log]
level = 5
dump  = 0

[nla]
core  = "127.0.0.1:50061"
api   = "127.0.0.1:50062"

[ribc]
fibc  = "192.169.1.1:50070"

[ribs]
disable = true
# core = "sample:50071"
# api  = "127.0.0.1:50072"

# [ribs.bgpd]
# addr = "127.0.0.1"
# # port = 50051

# [ribs.nexthops]
# mode = "translate"
# args = "1.1.0.0/24"

[ribp]
api = "127.0.0.1:50091"
interval = 5000
```

**Except MPLS-VPN environments**, the syntax is following:

```
# -*- coding: utf-8; mode: toml -*-

[node]
nid   = 0
reid  = "<router-entity-id>"
allow_duplicate_ifname = <allow-duplicate-ifname>

[log]
level = <beluganos-log-level>
dump  = <beluganos-debug>

[nla]
core  = "127.0.0.1:50061"
api   = "127.0.0.1:50062"

[ribc]
fibc  = "192.169.1.1:50070"

[ribs]
disable = true

[ribp]
api = "127.0.0.1:50091"
interval = <beluganos-ribp-interval>
```


- [node]
	- nid (`<instance-number>`): Linux container ID. Except VRF environments, you should set to `0` to instance number.
	- reid (`<router-entity-id>`): Router entity id for internal use. In general, router-id is recommended. This value should match the value in `fibc.yml`.
	- allow\_duplicate\_ifname: The settings to allow overlapping of interface names between different containers, except "lo" and "eth0". Default settings is "false".
- [log]
	- level (`<beluganos-log-level>`): Log level. 0〜5. 5 is the most detailed value.
	- dump (`<beluganos-debug>`): Debug flag. If you don't need debug log, please set to 0.
- [nla]
	- The setting about netlink abstraction module.
	- core: set `localhost:50061` expect MPLS-VPN environments.
	- api: NLA API address. set `localhost:50062`. Do not change.
- [ribc]
	- The setting about RIB controller.
	- fibc: FIBC address. set `192.169.1.1:50070`.
- [ribs]
	- The setting about RIBS module. This module will enable route redistribution function between routing instances.
	- disable: set true. Only in MPLS-VPN environments, this module will be used.
- [ribp]
	- The setting about RIBP which is one of the components of RIB controller.
	- api: set `localhost:50091`. Do NOT change.
	- interval (`<beluganos-ribp-interval>`): The interval value which synchronize interface status between Linux container and main module of Beluganos. By our verification, `5000` is recommended value for stable operations is assumed.

**Note**: In MPLS-VPN environments, more configurations are required. Please refer to [configure-ansible-l3vpn.md](configure-ansible-l3vpn.md).

### Reflect changes

In case group name is "*lxd-group*" and container name is "*master*":

```
$ cd ~/beluganos/etc/playbooks
$ ansible-playbook -i hosts -K lxd-group.yml
$ lxc stop master
```

## Start Beluganos

When you finished the beluganos's settings, Let's start Beluganos! You can start main module of Beluganos by following commands:

```
$ beluganos start
$ beluganos add master
```

Please note that "*master*" means your container name. For more detail about `beluganos` commands and operations, please refer [operation-guide.md](operation-guide.md).

## Note

- MPLS-VPN configurations require more steps. Please refer [configure-ansible-l3vpn.md](configure-ansible-l3vpn.md).
- Using both ansible and NETCONF are not permitted currently to configure router settings.