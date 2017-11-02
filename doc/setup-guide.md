# Setup guide
This document describes about Beluganos's configuration. You can configure router options like IP addresses or settings of routing protocol as you like. You can use  **ansible** to configure Beluganos.

## Config files at a glance

The files under `etc/playbooks/` are configuration files.

~~~~
beluganos/
    etc/
        playbooks/
            hosts                           # Inventry file
            lxd-sample.yml                  # Sample
            dp_sample.yml                   # Sample
            lxd-<group_name>.yml            # Playbook
            dp-<datapath_name>.yml          # Playbook
            roles/
                dpath/
                    vars/                   # DO NOT EDIT.
                    tasks/                  # DO NOT EDIT.
                    files/
                        common/             # DO NOT EDIT.
                        sample_sw/
                            fibc.yml        # Sample
                        <switch_name>/
                            fibc.yml
                lxd/
                    vars/                   # DO NOT EDIT.
                    handlers/               # DO NOT EDIT.
                    tasks/
                        main.yml            # DO NOT EDIT.
                        create.yml          # DO NOT EDIT.
                        setup.yml           # DO NOT EDIT.
                        host.yml            # DO NOT EDIT.
                        sample.yml          # Sample
                        <container name>.yml
                    files/
                        common/             # DO NOT EDIT.
                        sample/             # Sample
                        <container name>/
                            interfaces.cfg
                            sysctl.conf
                            daemons
                            frr.conf
                            gobgp.comnf
                            gobgpd.conf
                            fibc.yml
                            ribxd.conf
~~~~

## A. Settings for white-box switches

### 1. switch name

**The sample playbook is `etc/playbooks/dp-sample.yml`**. You may copy sample file and rename this playbook. Any switch name (dpname) are acceptable, and this name is only used for internal configurations. If your switch name is *switchA*, please edit as follows:

~~~~
$ cd beluganos/etc/playbooks
$ cp dp-sample.yml dp-switchA.yml
$ vi dp-switchA.yml

---

- hosts: hosts
  connection: local
  roles:
    - { role: dpath, dpname: switchA, mode: fibc }
  tags:
    fibc

~~~~

The syntax of this playbook is following:

~~~~
- hosts: hosts
  connection: local
  roles:
    - { role: dpath, dpname: <switch-name>, mode: fibc }
  tags:
    fibc
~~~~

- roles
	- dpname (`<switch-name>`): Your preferable switch name (A-Z,a-z,0-9,_-). Do not contain space.

### 2. switch type
**The sample configuration files are under `etc/playbooks/roles/dpath/files/sample_sw`**. The files under `roles/dpath` will determine the type of white-box switches. Beluganos will optimize the writing method to FIB for each switch type. If your switch's name is *switchA*, you may copy sample file and rename to dpname as follows:

~~~~
$ cd beluganos/etc/playbooks
$ cp -r roles/dpath/files/sample_sw roles/dpath/files/dp-switchA
$ vi roles/dpath/files/dp-switchA/fibc.yml

---

datapaths:
  - name: switchA           # dpname (A-Z,a-z,0-9,_-)
    dp_id: 14               # datapath id of your switches (integer)
    mode: generic           # "ofdpa2" or "generic" or "ovs"

~~~~

The syntax of this playbook is following:

~~~~
datapaths:
  - name: <switch-name>     # dpname (A-Z,a-z,0-9,_-)
    dp_id: <switch-dp-id>   # datapath id of your switches (integer)
    mode: <switch-type>     # "ofdpa2" or "generic" or "ovs"
~~~~


The value of `dp-id` means OpenFlow datapath ID of your switch. The value of `mode` should be edited for your switch types. You can choose following options:

- datapaths
	- name (`<switch-name>`): Your switch name which is already declared in `etc/playbooks/dp-switchA.yml`.
	- dp_id (`<switch-dp-id>`): OpenFlow datapath ID of your switch. Integer.
	- mode (`<switch-type>`): Your switch types. Currently Beluganos has three options.
 		1. `ofdpa2`: OF-DPA 2.0 switch. [https://github.com/broadcom/ofdpa/](https://github.com/broadcom/ofdpa/)
		1. `generic`: OpenFlow 1.3 compaliable switches. (e.g. Lagopus)
		1. `ovs`: Open vSwitch (Limited support).

### Reflect changes

In case `<switch-name>` is *switchA*:

~~~~
$ ansible-playbook -i etc/playbooks/hosts -K etc/playbooks/dp-switchA.yml
~~~~

## B. Settings for Linux Containers

### 1. Container name (inventory)

At first, you should edit [inventory file](http://docs.ansible.com/ansible/latest/intro_inventory.html) at `etc/playbooks/hosts`. The container name and group name should be described here. This name will be only used internal configurations.

In general, one-to-one correspondence between group name and container name is assumed except VRF environments. For example, if your group name is *lxd-group* and your container name (i.e. routing instance name) is *master-instance*, please edit as follows:

```
$ cd beluganos/etc/playbooks/
$ vi hosts
---

[hosts]      # DO NOT EDIT.
localhost    # DO NOT EDIT.

[lxd-group]
master-instance
```

**For VRF environments only:** To configure multiple routing-instances, you should describe multiple container name.

```
[lxd-group-vpn]
master-instance
vrf10
vrf20
```

### 2. Container basic settings (playbook)

**The sample playbook is `lxd-sample.yml`**. You may copy this file and rename to your group name. For example, if your group name at inventory is *lxd-group*, please edit following:


```
$ cd etc/playbooks
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
      with_items:  # Edit "lxd-group" to your group name in inventry file.
        - "{{ groups['lxd-group'] }}"
      loop_control:
        loop_var: lxcname
  tags:
    - create
    - lxd

- hosts: lxd-group  # Edit "lxd-group" to your setting in inventry file.
  connection: lxd
  roles:
    - { role: lxd, lxcname: "{{ inventory_hostname }}", mode: setup }
  tags:
    - setup
    - lxd
```

### 3. Container interface names (playbook)

**The sample configuration file is `roles/lxd/tasks/sample.yml`**. At first, you should create task files for ansible. In task files, you can determine about **interface name of Linux container**.

If your container name is *"master-instance"*, please edit task files as follows:

~~~~
$ cd etc/playbooks
$ cp roles/lxd/tasks/sample.yml roles/lxd/tasks/master-instance.yml
$ vi roles/lxd/tasks/master-instance.yml

---

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

~~~~

**For VRF environments only:** The task files should be created per routing-instance.

The syntax is following:

~~~~
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
      <lxc-interface-name>:
        type: nic
        name: <lxc-interface-name>
        host_name: "{{ lxcname }}.<lxc-interface-number>"
        nictype: p2p
      ....
~~~~

- devices
	- `eth0`: The management interface for internal use. Do NOT edit.
	- `<lxc-intarface-name>`: The interface name of container. If your switch have 10 interfaces, you should add devices from eth1 to eth10. In Beluganos's architecture, each linux container's interface is corresponding to white-box's interface.
		- host_name (`{{ lxcname }}.<lxc-interface-number>`): Beluganos's main module recognize container's interface by this value. Any value is accepted  in `<lxc-interface-number>` unless there is overlap other interface. For example, if the value of `<lxc-intarface-name>` is *eth5*, `<lxc-interface-number>` will be *5* is assumed.

### Reflect changes

The reflection of "B. Settings for Linux Containers" should be done after editting "C. Router settings".

## C. Router settings

**The sample files are under `roles/lxd/files/sample`**. You may copy and rename these files. If your container name is *master-instance*, you can copy by folllowing commands:

```
$ cd etc/playbooks
$ cp -r roles/lxd/files/sample roles/lxd/files/master-instance
$ ls roles/lxd/files/master-instance
daemons fibc.yml frr.conf gobgp.conf gobgpd.conf interfaces.cfg ribxd.conf sysctl.conf
```

### 1. fibc.yml: global settings of this router

In `fibc.yml`, you may set general settings. FIBC is one of the main components of Beluganos.

Especially, you can determine about **interface mapping between your white-box switches and Linux containers**. In Beluganos's architecture, the path information which was installed to Linux containers will be copied to white-box switches. This is because you have to configure interface mapping between switches and containers.

~~~~
$ vi roles/lxd/files/master-instance/fibc.yml

routers:
  - desc: sample-router
    re_id: 10.0.1.6                                    # Router entity id (=router-id)
    datapath: switchA                                  # dpname
    ports:
      - { name: 0/eth1,     port: 1 }                  # nid=0, lxc-iface=eth1, dp-port=1
      - { name: 0/eth2,     port: 2 }                  # nid=0, lxc-iface=eth2, dp-port=2
      - { name: 0/eth2.100, port: 0,  link: "0/eth2" } # vlan iface of "0/eth2"
~~~~

Syntax is following:

~~~~
routers:
  - desc: <description>
    re_id: <router-entity-id>
    datapath: <switch-name>
    ports:
      - { name: <instance-number>/<lxc-interface-name>, port: <switch-interface-number> }                                                # For physical interface
      - { name: <instance-number>/<lxc-interface-name>.<lxc-interface-vlan>, port: 0, link: "<instance-number>/<lxc-interface-name>" }   # For VLAN interface
~~~~

- desc (`<description>`): Your preferable descriptions.
- re_id (`<router-entity-id>`): Router entiry id for internal use. In general, router-id is recomended.
	- **For VRF environments only**: Please set **master-instance's** router-id for all routing instances.
- datapath (`<switch-name>`): Your switch name which is declared previously.
- ports
    - name (`<instance-number>/<lxc-interface-name>`): Specify interface name in Beluganos.
	    -  `<instance-number>` should be integer. Except VRF environments, you should set to `0` to instance number. In VRF environments, you shoud set to VRF number per routing instance.
	    -  `<lxc-interface-name>` should be matched with the value in `roles/lxd/tasks/<container-name>.yml`. For example, *eth1*.
    - port (`<switch-interface-number>`): Specify interface number in White-box switches. Note that in the case of VLAN interface you should set to `0`.
	- link (`<instance-number>/<lxc-interface-name>`): Specify container's interface name which vlan-id belong to. Note that in the case of pyshical interface you should NOT set this value.


### 2. interfaces.cfg: VLAN configurations

This file will be located under `/etc/network/interfaces.d/` in linux container. In Beluganos, only when you want to configure VLAN interface, you have to edit this file.

~~~~
# -*- coding: utf-8 -*-

auto eth2.100
iface eth2.100 inet manual
~~~~

### 3. sysctl.conf: MPLS basic configurations: sysctl.conf

This file will be located under `/etc/sysctl.d/` in Linux container. If you want to enable MPLS for all interfaces, you have to disable `rp_filter` and set `net.mpls.conf` as follows:

~~~~
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
~~~~

### 4. daemons: Routing protocols which you want to use

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

### 5. frr.conf: Each routing protocol settings

Each daemon configuration of FRRouting.
 See [FRRouting user-guide](https://frrouting.org/user-guide/).

### 6. gobgp.conf: GoBGP daemon settings

This file is about gobgp booting, **NOT gobgpd itself**. In general, please do not edit.

~~~~
# -*- coding: utf-8 -*-

CONF_PATH = /etc/frr/gobgpd.conf
CONF_TYPE = toml
LOG_LEVEL = debug
PPROF_OPT = --pprof-disable
API_HOSTS = 127.0.0.1:50051
~~~~

### 7. gobgpd.conf: GoBGP itselfs configurations

GoBGP configurations. See [GoBGP configurations.md](https://github.com/osrg/gobgp/blob/master/docs/sources/configuration.md).

**For MPLS-VPN environments only**: In MPLS-VPN, some optional settings about GoBGP is required. For more detail, please check `doc/setup-guide-l3vpn.md`.

### 8. ribxd.conf: Beluganos's settings

This is the main configuration file of Beluganos itself. For example, following configuration will be assumed:

~~~~
# -*- coding: utf-8; mode: toml -*-

[node]
nid   = 0
reid  = "10.0.1.6"

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

[ribp]
api = "127.0.0.1:50091"
interval = 5000
~~~~

**Except MPLS-VPN environments**, the syntax is following:

~~~~
# -*- coding: utf-8; mode: toml -*-

[node]
nid   = 0
reid  = "<router-entity-id>"

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
~~~~


- [node]
	- nid (`<instance-number>`): Linux container ID. Except VRF environments, you should set to `0` to instance number.
	- reid (`<router-entity-id>`): Router entity id for internal use. In general, router-id is recommended. This value should match the value in `fibc.yml`.
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
	- api: set `localhost:50091`.
	- interval (`<beluganos-ribp-interval>`): The interval value which synchronize interface status between Linux container and main module of Beluganos. Curentlly, `5000` is recommended for stable operations.

**Note**: In MPLS-VPN environments, more configurations are required. Please refer to `doc/setup-guide-l3vpn.md`.

### Reflect changes

In case group name is *lxd-group* and container name is *master-instance*:

~~~~
$ ansible-playbook -i etc/playbooks/hosts -K etc/playbooks/lxd-group.yml
$ lxc stop master-instance
~~~~

Note that the reflection of "C. Router settings" will be done at the same time as "B. Settings for Linux Containers".