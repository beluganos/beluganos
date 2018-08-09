# Configure by NETCONF

## Pre-requirements
- Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.
	- In setup, you may specify switch name and `dp_id`. In this documents, the sample file's value (`name: sample_sw`, `dp_id: 153`) is assumed. If you changed this value, please change to match it.

## Step 1. Prepare for launch

### 1-1. playbooks

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
		- Router identified name. Only Beluganos's main component will use this value to identify routers. This value is only internal use.
	- datapath (`<switch-name>`) and dp_id (`<switch-dp-id>`)
		- White-box hardware settings. The value of `fibc.yml` which was edited at [setup-guide.md](setup-guide.md) should be filled.

To reflect, please execute ansible for setup.

```
$ ansible-playbook -i hosts -K lxd-netconf.yml
```

### 1-2. lxcinit
In NETCONF module of Beluganos, once you create `network-instance` by NETCONF, `/etc/lxcinit/<container-type>/lxcinit.sh` will be executed. This script contains initial settings. You should configure about this script by editing `ribxd.conf` at `/etc/lxcinit/<container-type>/conf`.

#### container-type

`<container-type>` means the type of routing instance. Generally, `std_mic` should be selected. In case of VRF-Lite or MPLS-VPN environments, please refer following tables:

| container-type | description                          |
| -------------- | ------------------------------------ |
| std_mic        | Standard network instance            |
| std_ric        | Virtual router (VRF-Lite)            |
| vpn_mic        | Standard network instance with L3VPN |
| vpn_ric        | VRF for L3VPN                        |

#### ribxd.conf

In `std_mic`, the syntax is following. There is no reflection commands.

```
# -*- coding: utf-8; mode: toml -*-

[node]
nid   = 0
reid  = "<router-entity-id>"
label = 100000
allow_duplicate_ifname = false
# nid_from_ifaddr = "eth0"

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
```
- [node]
	- reid (`<router-entity-id>`)
		- Router identified name. Only Beluganos's main component will use this value to identify routers. This value is only internal use.

In other case (`std_ric`, `vpn_mic`, and `vpn_ric`), please refer [configure-ansible.md](https://github.com/beluganos/beluganos/blob/master/doc/configure-ansible.md#8-ribxdconf-beluganoss-settings).

## Step 2. Launch components

Actually, you have already finished to launch Beluganos! The configuring by NETCONF will be enabled after launching. Please start Beluganos by following commands.

```
$ sudo systemctl start fibcd
$ sudo systemctl start netopeer2-server
$ sudo systemctl start ncm.target
```

You can also use `beluganos start` instead of `systemctl start fibcd`. For more detail about `beluganos` commands and operations, please refer [operation-guide.md](operation-guide.md). Note that above commands are describes at step 1 in this document. Moreover, in NETCONF case, step 2 is not required.

## Step 3. Configure by NETCONF

After starting Beluganos, you can use NETCONF commands. In this sections, NETCONF operations (like `<edit-config>`) and the methods how to create proper XML.

### NETCONF operations

NETCONF over ssh will utilize TCP 830 port. NETCONF session will be started by following commands.

```
$ ssh -s <server-ip> -p 830 netconf
```

After exchanging `<hello>` message, you can operate `<get-config>` or `<edit-config>` operations. The example log is located at the bottom of this page.

Moreover, you can also use NETOPEER2's CLI like following:

```
$ netopeer2-cli
```

### Yang moudles and configuration XML

The yang modules of Beluganos are published under [netconf/etc/openconfig](https://github.com/beluganos/netconf/tree/master/etc/openconfig). Currently, Beluganos support three modules ([network-instances](https://github.com/beluganos/netconf/blob/master/etc/openconfig/beluganos-network-instance.yang), [interfaces](https://github.com/beluganos/netconf/blob/master/etc/openconfig/beluganos-interfaces.yang), [routing-policy](https://github.com/beluganos/netconf/blob/master/etc/openconfig/beluganos-routing-policy.yang)). Note that the sample NETCONF XML are available at [netconf/doc/examples](https://github.com/beluganos/netconf/tree/master/doc/examples).

#### network-instance

At least one **network-instance** is required. Even if you will configure general IP routers, single network-instance exists is assumed. **The important point is that the type and route-target of network-instance is used for Beluganos's settings.**

```
module: beluganos-network-instance
    +--rw network-instances
       +--rw network-instance* [name]
          +--rw name          -> ../config/name
          +--rw config
          |  +--rw type?                  identityref (*)
          |  +--rw route-target?          oc-ni-types:route-distinguisher (*)
          |  +--rw ....
```

The supported network instance (i.e. Linux container) type is following:

| Type             | Route-target | Description                          |
| ---------------- | -------------| ------------------------------------ |
| DEFAULT_INSTANCE | No           | Standard network instance            |
| L3VRF            | No           | Virtual router (VRF-Lite)            |
| DEFAULT_INSTANCE | Yes(*1)      | Standard network instance with L3VPN |
| L3VRF            | Yes          | VRF for L3VPN                        |

(*1) As for any value, please fill it. This value is not used by Beluganos.

For more detail about network-instance, please refer [netconf/setup-guide.md](https://github.com/beluganos/netconf/blob/master/doc/setup-guide.md).

#### interfaces

Currently Beluganos support only routed port. In routed port, `<subinterface>` describes VID (VLAN ID). For example, if you want to use `eth2` as a subinterface with VID 10, you should configure both `eth2` and `eth2.10`.

```
module: beluganos-interfaces
    +--rw interfaces
       +--rw interface* [name]
          +--rw name             -> ../config/name
          +--rw config
          |  +--rw name?          string
          |  +--rw type           identityref
          |  +--rw mtu?           uint16
          |  +--rw description?   string
          |  +--rw enabled?       boolean
          +--rw state
          +--rw subinterfaces
             +--rw subinterface* [index]
                +--rw index     -> ../config/index
                +--rw config
                |  +--rw index?         uint32
                |  +--rw description?   string
                |  +--rw enabled?       boolean
                +--rw state
```

The naming rules of interface are `eth<n>` (`<n>` is a interface index). The maximum number of `<n>` is defined at `<maximum-physical-port>` in `lxd-netconf.yml`.

## Appendix

### Restrictions

- Datastores
	- For proper network operation, you cannot edit running-configuration datastore directly. You should use candidate-configuration datastore by `<commit>` operations.
	- The startup-configuration is automatically copied from running-configurations. You don't have to operate `<copy-config>`.
- Operations
	- The operation of `<validate>` will do nothing. Actually, the minimum verification will be done in `<edit-config>` to candidate-configuration. This is because `<validate>` has no meanings in Beluganos's implementation.
- Others
	- The settings of white-box switches (like `dp_id`) cannot changed by NETCONF. You should use ansible.
	- The settings of sub-IF (VLAN at routed port) should be operated at the same time of physical IF.
	- The interface settings under beluganos-interfaced module should be added **before** adding interface under beluganos-network-instances.
	- The settings of network instance's name should be matched at `ribxd.conf` which is located at `/etc/lxcinit/<container-type>`. This restrictions will be removed at next release.

### Sample operations

```

beluganos@beluganos:~$ sudo systemctl status ncm.target
● ncm.target - Beluganos netconf services
   Loaded: loaded (/etc/systemd/system/ncm.target; static; vendor preset: enabled)
   Active: active since Thu 2018-05-10 15:24:19 JST; 16min ago

May 10 15:24:19 beluganos020-180507 systemd[1]: Reached target Beluganos netconf services.
beluganos@beluganos:~$
beluganos@beluganos:~$ sudo systemctl status fibcd
● fibcd.service - fib controller service
   Loaded: loaded (/etc/systemd/system/fibcd.service; disabled; vendor preset: enabled)
   Active: active (running) since Thu 2018-05-10 13:42:10 JST; 1h 58min ago
  Process: 13519 ExecStartPre=/bin/sleep ${START_DELAY_SEC} (code=exited, status=0/SUCCESS)
 Main PID: 13529 (ryu-manager)
    Tasks: 1 (limit: 4915)
   Memory: 52.7M
      CPU: 1.425s
   CGroup: /system.slice/fibcd.service
           mq13529 /usr/bin/python /usr/local/bin/ryu-manager ryu.app.ofctl_rest fabricflow.fibc.app.fibcapp --co


beluganos@beluganos:~$ ssh -s localhost -p 830 netconf
Interactive SSH Authentication
Type your password:
Password:
<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><capabilities><capability>urn:ietf:params:netconf:base:1.0</capability><capability>urn:ietf:params:netconf:base:1.1</capability><capability>urn:ietf:params:netconf:capability:writable-running:1.0</capability><capability>urn:ietf:params:netconf:capability:candidate:1.0</capability><capability>urn:ietf:params:netconf:capability:rollback-on-error:1.0</capability><capability>urn:ietf:params:netconf:capability:validate:1.1</capability><capability>urn:ietf:params:netconf:capability:startup:1.0</capability><capability>urn:ietf:params:netconf:capability:xpath:1.0</capability><capability>urn:ietf:params:netconf:capability:with-defaults:1.0?basic-mode=explicit&amp;also-supported=report-all,report-all-tagged,trim,explicit</capability><capability>urn:ietf:params:netconf:capability:notification:1.0</capability><capability>urn:ietf:params:netconf:capability:interleave:1.0</capability><capability>urn:ietf:params:xml:ns:yang:ietf-yang-metadata?module=ietf-yang-metadata&amp;revision=2016-08-05</capability><capability>urn:ietf:params:xml:ns:yang:1?module=yang&amp;revision=2017-02-20</capability><capability>urn:ietf:params:xml:ns:yang:ietf-inet-types?module=ietf-inet-types&amp;revision=2013-07-15</capability><capability>urn:ietf:params:xml:ns:yang:ietf-yang-types?module=ietf-yang-types&amp;revision=2013-07-15</capability><capability>urn:ietf:params:xml:ns:yang:ietf-yang-library?module=ietf-yang-library&amp;revision=2016-06-21&amp;module-set-id=37</capability><capability>urn:ietf:params:xml:ns:yang:ietf-netconf-acm?module=ietf-netconf-acm&amp;revision=2012-02-22</capability><capability>urn:ietf:params:xml:ns:netconf:base:1.0?module=ietf-netconf&amp;revision=2011-06-01&amp;features=writable-running,candidate,rollback-on-error,validate,startup,xpath</capability><capability>urn:ietf:params:xml:ns:yang:ietf-netconf-notifications?module=ietf-netconf-notifications&amp;revision=2012-02-06</capability><capability>urn:ietf:params:xml:ns:netconf:notification:1.0?module=notifications&amp;revision=2008-07-14</capability><capability>urn:ietf:params:xml:ns:netmod:notification?module=nc-notifications&amp;revision=2008-07-14</capability><capability>http://example.net/turing-machine?module=turing-machine&amp;revision=2013-12-27</capability><capability>urn:ietf:params:xml:ns:yang:ietf-interfaces?module=ietf-interfaces&amp;revision=2014-05-08</capability><capability>urn:ietf:params:xml:ns:yang:iana-if-type?module=iana-if-type&amp;revision=2014-05-08</capability><capability>urn:ietf:params:xml:ns:yang:ietf-ip?module=ietf-ip&amp;revision=2014-06-16</capability><capability>http://openconfig.net/yang/openconfig-ext?module=openconfig-extensions&amp;revision=2017-04-11</capability><capability>http://openconfig.net/yang/openconfig-types?module=openconfig-types&amp;revision=2018-01-16</capability><capability>http://openconfig.net/yang/types/yang?module=openconfig-yang-types&amp;revision=2017-07-30</capability><capability>http://openconfig.net/yang/types/inet?module=openconfig-inet-types&amp;revision=2017-08-24</capability><capability>http://openconfig.net/yang/mpls-types?module=openconfig-mpls-types&amp;revision=2017-08-24</capability><capability>http://openconfig.net/yang/bgp-types?module=openconfig-bgp-types&amp;revision=2018-03-20</capability><capability>http://openconfig.net/yang/ospf-types?module=openconfig-ospf-types&amp;revision=2017-08-24</capability><capability>http://openconfig.net/yang/policy-types?module=openconfig-policy-types&amp;revision=2017-07-14</capability><capability>http://openconfig.net/yang/network-instance-types?module=openconfig-network-instance-types&amp;revision=2017-08-24</capability><capability>https://github.com/beluganos/beluganos/yang/interfaces?module=beluganos-interfaces&amp;revision=2017-10-20</capability><capability>https://github.com/beluganos/beluganos/yang/interfaces/ip?module=beluganos-if-ip&amp;revision=2017-07-14</capability><capability>https://github.com/beluganos/beluganos/yang/interfaces/ethernet?module=beluganos-if-ethernet&amp;revision=2018-01-05</capability><capability>https://github.com/beluganos/beluganos/yang/ldp?module=beluganos-mpls-ldp&amp;revision=2017-10-20</capability><capability>https://github.com/beluganos/beluganos/yang/mpls?module=beluganos-mpls&amp;revision=2017-10-20</capability><capability>https://github.com/beluganos/beluganos/yang/routing-policy?module=beluganos-routing-policy&amp;revision=2017-10-20</capability><capability>https://github.com/beluganos/beluganos/yang/bgp?module=beluganos-bgp&amp;revision=2017-10-20</capability><capability>https://github.com/beluganos/beluganos/yang/ospfv2?module=beluganos-ospfv2&amp;revision=2017-10-20</capability><capability>https://github.com/beluganos/beluganos/yang/local-routing?module=beluganos-local-routing&amp;revision=2017-10-20</capability><capability>https://github.com/beluganos/beluganos/yang/network-instance?module=beluganos-network-instance&amp;revision=2017-10-20</capability><capability>https://github.com/beluganos/beluganos/yang/bgp-policy?module=beluganos-bgp-policy&amp;revision=2017-10-20</capability><capability>urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring?module=ietf-netconf-monitoring&amp;revision=2010-10-04</capability><capability>urn:ietf:params:xml:ns:yang:ietf-netconf-with-defaults?module=ietf-netconf-with-defaults&amp;revision=2011-06-01</capability></capabilities><session-id>5</session-id></hello>]]>]]>

<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><capabilities><capability>urn:ietf:params:netconf:base:1.1</capability></capabilities></hello>]]>]]>

#166
<?xml version="1.0" encoding="UTF-8"?><rpc message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><get-config><source><running/></source></get-config></rpc>
##

#249
<rpc-reply message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><data xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><routing-policy xmlns="https://github.com/beluganos/beluganos/yang/routing-policy"></routing-policy></data></rpc-reply>
##

#405
<?xml version="1.0" encoding="UTF-8"?><rpc message-id="102" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><edit-config><target><candidate/></target><config><network-instances xmlns="https://github.com/beluganos/beluganos/yang/network-instance"><network-instance><name>master-instance</name><config><name>master-instance</name></config></network-instance></network-instances></config></edit-config></rpc>
##

#93
<rpc-reply message-id="102" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>
##

#1472
<?xml version="1.0" encoding="UTF-8"?><rpc message-id="105" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><edit-config><target><candidate/></target><config><interfaces xmlns="https://github.com/beluganos/beluganos/yang/interfaces"></interfaces><network-instances xmlns="https://github.com/beluganos/beluganos/yang/network-instance"><network-instance><name>master-instance</name><config><name>master-instance</name><type xmlns:oc-ni-types="http://openconfig.net/yang/network-instance-types">oc-ni-types:DEFAULT_INSTANCE</type><router-id>1.1.1.1</router-id><route-distinguisher>10:1</route-distinguisher><route-target>10:1</route-target></config><interfaces></interfaces><protocols><protocol><identifier xmlns:oc-pol-types="http://openconfig.net/yang/policy-types">oc-pol-types:OSPF</identifier><name>MSF-infra</name><config><identifier xmlns:oc-pol-types="http://openconfig.net/yang/policy-types">oc-pol-types:OSPF</identifier><name>MSF-infra</name></config><ospfv2><global><config><router-id>1.1.1.1</router-id></config></global><areas><area><identifier>0.0.0.0</identifier><config><identifier>0.0.0.0</identifier></config><interfaces><interface><id>lo</id><config><id>lo</id><metric>10</metric><passive>true</passive></config><interface-ref><config><interface>lo</interface><subinterface>0</subinterface></config></interface-ref></interface></interfaces></area></areas></ospfv2></protocol></protocols></network-instance></network-instances></config></edit-config></rpc>
##

#93
<rpc-reply message-id="105" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>
##

#123
<?xml version="1.0" encoding="UTF-8"?><rpc message-id="106" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><commit/></rpc>
##

#93
<rpc-reply message-id="106" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>
##

#166
<?xml version="1.0" encoding="UTF-8"?><rpc message-id="107" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><get-config><source><running/></source></get-config></rpc>
##

#1022
<rpc-reply message-id="107" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><data xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><routing-policy xmlns="https://github.com/beluganos/beluganos/yang/routing-policy"></routing-policy><network-instances xmlns="https://github.com/beluganos/beluganos/yang/network-instance"><network-instance><name>master-instance</name><config><name>master-instance</name><type xmlns:oc-ni-types="http://openconfig.net/yang/network-instance-types">oc-ni-types:DEFAULT_INSTANCE</type><router-id>1.1.1.1</router-id><route-distinguisher>10:1</route-distinguisher><route-target>10:1</route-target></config><protocols><protocol><identifier xmlns:oc-pol-types="http://openconfig.net/yang/policy-types">oc-pol-types:OSPF</identifier><name>MSF-infra</name><config><identifier xmlns:oc-pol-types="http://openconfig.net/yang/policy-types">oc-pol-types:OSPF</identifier><name>MSF-infra</name></config><ospfv2><global><config><router-id>1.1.1.1</router-id></config></global><areas><area><identifier>0.0.0.0
#744
</identifier><config><identifier>0.0.0.0</identifier></config><interfaces><interface><id>lo</id><config><id>lo</id><metric>10</metric><passive>true</passive></config><interface-ref><config><interface>lo</interface><subinterface>0</subinterface></config></interface-ref><timers></timers></interface></interfaces></area></areas></ospfv2><bgp><global></global><zebra><config></config></zebra></bgp></protocol></protocols><mpls><global></global><signaling-protocols><ldp><global><address-families><ipv4><config><label-policy><advertise></advertise></label-policy></config></ipv4></address-families><discovery><interfaces></interfaces></discovery></global></ldp></signaling-protocols></mpls></network-instance></network-instances></data></rpc-reply>
##

#429
<?xml version="1.0" encoding="UTF-8"?><rpc message-id="108" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><edit-config><target><candidate/></target><default-operation>replace</default-operation><config><interfaces xmlns="https://github.com/beluganos/beluganos/yang/interfaces"></interfaces><network-instances xmlns="https://github.com/beluganos/beluganos/yang/network-instance"></network-instances></config></edit-config></rpc>
##

#93
<rpc-reply message-id="108" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>
##

#123
<?xml version="1.0" encoding="UTF-8"?><rpc message-id="106" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><commit/></rpc>
##

#93
<rpc-reply message-id="106" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><ok/></rpc-reply>
##

#166
<?xml version="1.0" encoding="UTF-8"?><rpc message-id="107" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><get-config><source><running/></source></get-config></rpc>
##

#249
<rpc-reply message-id="107" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><data xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><routing-policy xmlns="https://github.com/beluganos/beluganos/yang/routing-policy"></routing-policy></data></rpc-reply>
##

```


