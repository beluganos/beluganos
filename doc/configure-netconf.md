# Configure by NETCONF

## Pre-requirements
- Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.
	- In setup, you may specify switch name and dp_id. In this documents, the sample file's value (`name: sample_sw`, `dp_id: 153`) is assumed. If you changed this value, please change to match it.

## Step 1. Prepare for launch

You have to execute minimum playbooks even if you configure only by NETCONF before configuring.

```
$ cd ~/beluganos/etc/playbooks
$ vi lxd-netconf.yml
---

- hosts: hosts
  connection: local
  become: true
  roles:
    - { role: lxd, mode: netconf }
  vars:
    port_num: 5
    re_id: 10.0.1.6
    datapath: sample_sw
    dp_id: 153
    bridge: dp0

$ ansible-playbook -i hosts -K lxd-netconf.yml
```

The syntax of `etc/playbooks/lxd-netconf.yml` is following:

```
- hosts: hosts
  connection: local
  become: true
  roles:
    - { role: lxd, mode: netconf }
  vars:
    port_num: <maximum-physical-port>
    re_id: <router-entity-id>
    datapath: <switch-name>
    dp_id: <switch-dp-id>
    bridge: dp0
```

- vars
	- port_num (`<maximum-physical-port>`)
		- The maximum physical interface number of your router. For example, if you have 48x10G and 6x40G port in a switch, the value of `port_num` may be `54`.
	- re_id (`<router-entity-id>`)
		- Router identified name. Only Beluganos's main component will use this value to identify routers. This value is not used for router settings, like loopback address or router-id.
	- datapath (`<switch-name>`) and dp_id (`<switch-dp-id>`)
		- White-box hardware settings. The value of `fibc.yml` which was edited at [setup-guide.md](setup-guide.md) should be filled.

## Step 2. Launch components

You have already finished to launch Beluganos! The configuring by NETCONF will be enabled after launching. Please start Beluganos by following commands.

```
$ beluganos start
$ sudo systemctl start netopeer2-server
$ sudo systemctl start ncm.target
```

For more detail about `beluganos` commands and operations, please refer [operation-guide.md](operation-guide.md).

## Step 3. Configure by NETCONF

After starting Beluganos, you can use NETCONF commands.

The yam of Beluganos is published under [netconf/etc/openconfig](https://github.com/beluganos/netconf/tree/master/etc/openconfig). Furthermore, sample NETCONF operations are available at [netconf/doc/examples](https://github.com/beluganos/netconf/tree/master/doc/examples). 
