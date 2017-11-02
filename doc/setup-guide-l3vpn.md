# Setup guide for L3VPN

Beluganos support PE router of BGP/MPLS based IP-VPN environments. To setup PE router, additional settings are needed. After setting general settings by [setup-guide](setup-guide.md), please execute following procedure.

## Container constitution

Generally, PE router needs a master instance and multiple VRF instance in order to separate the area of IP address. In Beluganos, **the separation of Linux containers will be roled as the separation of IP address tables**. We call these containers MIC and RIC.

### MIC

- **M**aster **I**nstance **C**ontainer
- Connect with P, PE, and RR routers
- MP-iBGP by GoBGP
- Required number at PE: 1

### RIC

- **R**outing **I**nstance **C**ontainer
- Connect with CE routers
- The learned route by eBGP will be redistributed to MIC's GoBGP.
- Requrired number at PE: 1 per VRF


~~~~
PE--//--P   P--//--RR
         \ /
   + - - MIC - - +
   :      |      :
   :      o      :
   :     / \     :
   + -RIC - RIC- +
       |    / \
      CE   CE CE
~~~~

## A. Settings for white-box switches

Common with [setup-guide](setup-guide.md).

## B. Settings for Linux Containers

### 1. Container name (inventry)

In inventory file (`etc/playbooks/hosts`), please note that you should set all container name not only MIC but also RIC.

~~~~
[lxd-sample-vpn]
sample-mic
sample-ric10
samplr-ric11
~~~~

### 2. Container basic settings (playbook)

Common with [setup-guide](setup-guide.md).

## C. Router settings

In sample of playbooks, the files under `etc/playbooks/roles/lxd/files/<container-name>` will be transfered each linux container. This capter will be described the difference from general routers. The files which is not mentioned in this capter is common with [setup-guide](setup-guide.md).

### 1. fibc.yml

The file of `roles/lxd/files/<container-name>/fibc.yml` should be modified like following sample.

#### Common (MIC&RIC)

- Same port cannot belong to both MIC and RIC. The value of port number should not be dupplicated.
- **The value of `re_id` and `datappath` should be same value** in both MIC and all RIC. We recommended re_id will be set to MIC's router-id.
- `<instance-number>` is the Linux container ID (integer). Please note that you have to match this value with `ribxd.conf` file. The value of `0` means MIC, and other value means RIC.
- Syntax is following.

~~~~
routers:
  - desc: <description>
    re_id: <router-entity-id>
    datapath: <switch-name>
    ports:
      - { name: <instance-number>/<lxc-interface-name>, port: <switch-interface-number> }                                                # For physical interface
      - { name: <instance-number>/<lxc-interface-name>.<lxc-interface-vlan>, port: 0, link: "<instance-number>/<lxc-interface-name>" }   # For VLAN interface
~~~~

#### MIC

The value of ports -> name should be `0/<ifname>` format. Sample is here:

~~~~
routers:
  - desc: SAMPLE-VPN-MIC
    re_id: 10.0.1.1
    datapath: whitebox1
    ports:
      - { name: 0/eth1,     port: 1 }
      - { name: 0/eth1.10,  port: 0, link: "0/eth1" }  # VLAN Device
~~~~

#### RIC

The values of ports -> name should be `<instance-number>/<ifname>` format. Sample is here:

~~~~
routers:
  - desc: SAMPLE-VPN-RIC10
    re_id: 10.0.1.1
    datapath: whitebox1
    ports:
      - { name: 10/eth1,     port: 2 }
      - { name: 10/eth2,     port: 3 }
      - { name: 10/eth2.10,  port: 0, link: "10/eth2" } # VLAN Device

~~~~

### 2. interfaces.cfg

The files of `etc/playbooks/roles/lxd/files/<container-name>/interfaces.cfg` should be modified like following sample.

#### MIC

See [setup-guide](setup-guide.md).

#### RIC

Additional configurations are needed for internal use. Please add following configurations. The bridge name of `ffbr0` will be used in `ribxd.conf`. This bridge is used for route redistribution between MIC and RIC.

~~~~
auto ffbr0
iface ffbr0 inet manual
  pre-up ip link add ffbr0 type bridge
~~~~

### 4. daemons

The file of `etc/playbooks/roles/lxd/files/<container-name>/daemons` should be modified like following sample.

In general settings, any software like quagga, FRRouting, and GoBGP are acceptable for routing engine. However, in PE router's environments, Beluganos has some restrictions about routing engine. By this restrictions, sample file of `daemons` is provided.

#### MIC

- Allow to enable: zebra, ospfd, isisd, ldpd
- Deny to enable: other protocols (**including bgpd**)

~~~~
zebra=yes
ospfd=yes
ldpd=yes
bgpd=no
~~~~

#### RIC

- Allow to enable: zebra
- Deny to enable: other protocols (**including ldpd, bgpd**)

~~~~
zebra=yes
ospfd=no
ldpd=no
bgpd=no
~~~~

### 6. gobgpd.conf

The file of `etc/playbooks/roles/lxd/files/<container-name>/gobgpd.conf` should be modified like following sample.

#### MIC

~~~~
[global.config]
  as = 65001                                            # Self AS number
  router-id = "10.0.1.1"                                # Lo

  [global.apply-policy.config]
    export-policy-list = ["policy-nexthop-self"]
    default-export-policy = "accept-route"


[[neighbors]]
  [neighbors.config]
    neighbor-address = "10.0.0.254"                     # BGP peer address
    peer-as = 65001                                     # BGP peer's AS number

  [neighbors.transport.config]
    local-address = "10.0.1.1"                          # Lo address

  [[neighbors.afi-safis]]
    [neighbors.afi-safis.config]
      afi-safi-name = "l3vpn-ipv4-unicast"


[[defined-sets.neighbor-sets]]
  neighbor-set-name = "ns-rr"
  neighbor-info-list = ["10.0.0.254"]                    # BGP peer address

[[policy-definitions]]
  name = "policy-nexthop-self"
  [[policy-definitions.statements]]
    [policy-definitions.statements.conditions.match-neighbor-set]
      neighbor-set = "ns-rr"
    [policy-definitions.statements.actions.bgp-actions]
      set-next-hop = "10.0.1.1"                          # Nexthop: Lo address (iBGP)
      set-local-pref = 100                               # Local Pref: set to lower than RIC
    [policy-definitions.statements.actions]
      route-disposition = "accept-route"
~~~~

- Beluganos supports only inter-AS MPLS-VPN environments. You should set same AS number at global.config -> as and neighbors.config -> peer-as.
- In generally inter-AS MPLS-VPN environments, the next-hop address which are advertised by MP-iBGP should be modified to PE router itselfs loopback address, instead of CE router's one. The applied policy of "policy-nexthop-self" will be acted as this behaivor. The value of "set-next-hop" at "policy-nexthop-self" is PE router itselfs loopback address.
- In "policy-nexthop-self", LOCAL\_PREF will be set as 100. This settings is needed by Beluganos. This is restrictions in currently version, and we are working to fix.
	- Because of technical limitations, the route decision by administrative distance (AD) value will not be worked correctly only in PE route's case. Therefore LOCAL\_PREF will be an alternative to AD value.
   - By this restrictions, **LOCAL\_PREF should be set to the lower value than RIC's LOCAL\_PREF**. This is because eBGP (from CE) should be have high priority than iBGP (from PE).
- Other BGP attirbute may be used as you like. For more detail, see [GoBGP configurations.md](https://github.com/osrg/gobgp/blob/master/docs/sources/configuration.md).

#### RIC

~~~~
[global.config]
  as = 65001                                       # self AS number
  router-id = "10.0.1.1"                           # Lo address

  [global.apply-policy.config]
    export-policy-list = ["policy-nexthop-self"]
    default-export-policy = "accept-route"
    import-policy-list = ["policy-local-pref"]
    default-import-policy = "accept-route"


[[neighbors]]
  [neighbors.config]
    neighbor-address = "192.168.1.2"               # BGP peer address
    peer-as = 10                                   # BGP peer's AS number

  [[neighbors.afi-safis]]
    [neighbors.afi-safis.config]
      afi-safi-name = "ipv4-unicast"


[zebra]                                            # Zebra collaborations
  [zebra.config]
    enabled = true
    version = 4
    url = "unix:/var/run/frr/zserv.api"
    # redistribute-route-type-list = ["connect"]


[[defined-sets.neighbor-sets]]
  neighbor-set-name = "ns-ce1"
  neighbor-info-list = ["192.168.1.2"]             # CE1


[[policy-definitions]]
  name = "policy-nexthop-self"
  [[policy-definitions.statements]]
    [policy-definitions.statements.conditions.match-neighbor-set]
      neighbor-set = "ns-ce1"
    [policy-definitions.statements.actions.bgp-actions]
      set-next-hop = "192.168.1.1"                # Nexthop: interface address (eBGP).
    [policy-definitions.statements.actions]
      route-disposition = "accept-route"


[[policy-definitions]]
  name = "policy-local-pref"
  [[policy-definitions.statements]]
    [policy-definitions.statements.conditions.match-neighbor-set]
      neighbor-set = "ns-ce1"
    [policy-definitions.statements.actions.bgp-actions]
      set-local-pref = 110                        # Local Pref: set to higher than MIC
    [policy-definitions.statements.actions]
      route-disposition = "accept-route"
~~~~

- The collaborations of Zebra (FRRouting) should be enabled.
- Currently, route redistribution to MIC's MP-BGP daemon will be performed without configurations.
- Other notes is same as MIC.

### 8. ribxd.conf

The file of `etc/playbooks/roles/lxd/files/<container-name>/ribxd.conf` should be modified like following sample.

The value which is not described here should be set as [setup-guide](setup-guide.md).

- [node]
	- nid (`<instance-number>`): Linux container ID. Please refer `fibc.yml` to match number. MIC's `<instance-number>` should be 0, and RIC's `<instance-number>` should be 1-254.
	- label (`<vpn-start-label>`): The start values of MPLS label to identify VPN. Generally you don't have to change this. Note that this settings should be commented out only in PE router's environments.
- [nla]
	- Settings about netlink abstruction module.
	- core (`<container-name-mic>:50061`): The settings for route redistribution. The linux container name of MIC should be set. The port number cannot be changed.
	- api: Set `127.0.0.1:50062` only in MIC case. **Comment out in RIC's configuration**.
- [ribs]
	- Settings about RIBS module. This module will enable route redistribution function between routing instances.
	- core (`<container-name-mic>:50071`): The settings for route redistribution. The linux container name of MIC should be set. The port number cannot be changed.
- [ribs.nexthops]
	- Set only in MIC case.
	- mode: Set `translate`.
	- args (`<translation-address>`): The IP address ranges which you cannot use your network. Unfortunately Beluganos will be allocated some IP address for internal use, only in PE router's environments.
- [ribs.vrf]
	- Set only in RIC case.
	- rt: Set import and export routing target (RT) value.
	- rd: Set routing distinguisher (RD) value.
	- iface: Set `ffbr0`.

#### sample (MIC)

```
# -*- coding: utf-8; mode: toml -*-

[node]
nid   = 0                        # 0
reid  = "10.0.1.1"               # <router-entity-id> at fibc.yml
label = 100000

[log]
level = 5
dump  = 0

[nla]
core  = "sample-mic:50061"       # <container-name-mic>:50061
api   = "127.0.0.1:50062"        # enable in MIC

[ribc]
fibc  = "192.169.1.1:50070"
# disable = true                 # RIBc module should be enabled in MIC

[ribs]
core = "sample-mic:50071"        # <container-name-mic>:50071
api  = "127.0.0.1:50072"

[ribs.bgpd]
addr = "127.0.0.1"

[ribs.nexthops]
mode = "translate"
args = "1.1.0.0/24"              # <translation-address> for internal dedicated IP address

[ribp]
api = "127.0.0.1:50091"
interval = 5000
```

#### sample (RIC)

```
# -*- coding: utf-8; mode: toml -*-

[node]
nid = 10                        # <instance-number> in fibc.yml
reid  = "10.0.1.1"              # <router-entiry-id> in fibc.yml
label = 100000

[log]
level = 5
dump  = 0

[nla]
core  = "sample-mic:50061"      # <container-name-mic>:50061
# api   = "127.0.0.1:50062"     # Comment in RIC

[ribc]
disable = true                  # RIBc module should be disabled in RIC

[ribs]
core = "sample-mic:50071"       # <container-name-mic>:50071
api  = "127.0.0.1:50072"

[ribs.bgpd]
addr = "127.0.0.1"

[ribs.vrf]
rt = "1:10"                     # Routing target
rd = "1:2010"                   # Routing distinguisher
iface = "ffbr0"

[ribp]
api = "127.0.0.1:50091"
interval = 5000
```
