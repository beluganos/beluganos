# [Feature guide] IP tunneling

IP over IP tunneling (RFC1853) is one of VPN technology between backbone network by encapsulating IP packet with additional IP header. Beluganos support IP-IP tunneling protocol. The tunnel interface can be created by static configuration or dynamic methods.

## Pre-requirements

- The installation is required in advance. Please refer [install.md](install.md) before proceeding.
- The setup of Beluganos is required in advance. Please refer [setup.md](setup.md) before proceeding.
- In case you will use OpenNSL or OF-DPA:
   - The installation of OpenNetworkLinux for your white-box switches is required in advance. Please refer [setup-hardware.md](setup-hardware.md) before proceeding.
   - The setup for ASIC API is required in advance. Please refer [setup-onsl.md](setup-onsl.md) or [setup-ofdpa.md](setup-ofdpa.md).

## Overview

Beluganos works as tunnel initiator or terminator of IP-IP tunnel. For example, in IPv4 over IPv6 tunneling, following environment is assumed:

```

                         External
   +---------+            Host1
   |  BGP    |              |                ---+
   |  Server |              | 20.0.0.0/24       | IPv4 only network
   +---------+              |                ---+
        |          +---------------+
        |          |               |  Lo(v4): 10.0.0.1/32
        |          |  EdgeRouter1  |  Lo(v6): 2001:db8:1::1/128
        |          |  <Beluganos>  |  Lo(v6): 2001:db8:2::1/128 (only for tunnel local address)
        |          |   router-id:  |  tunnel-if
        |          |    10.0.0.1   |     +- local : 2001:db8:2::1
        |          |               |     +- remote: 2001:db8:1::4
        |          +---------------+
        |            /           \           ---+
        |           /             \             |
        +----- Backbone         Backbone        | IPv6 only network
                Router1          Router2        | OSPFv3
                    \             /             | (IP-IP tunnel section)
                     \           /           ---+
                   +---------------+
                   |               |  tunnel-if
                   |  EdgeRouter2  |     +- local : 2001:db8:1::4
                   |               |     +- remote: 2001:db8:2::1
                   +---------------+
                           |                 ---+
                           | 30.0.0.0/24        | IPv4 only network
                           |                 ---+
                        External
                         Host2

```
Note that the "EdgeRouter1" is a Beluganos switch. In this environments, the backbone network is deployed as IPv6 network, but the communication between host1 to host2 supports only IPv4 technology. In this case, the EdgeRouter1 and EdgeRouter2 provide IPv4 over IPv6 tunneling.

## Setup

In Beluganos, following two conditions are required to create IP-IP tunneling.

1. Know the route to remote address of IP-IP tunnel
1. Create IP tunnel's interface

### 1. Know the route to remote address of IP-IP tunnel

The route to remote IP (`2001:db8:1::4`, in the figure) should be installed to Beluganos. Generally, this is advertised by IGP. Please confirm by following command:

```
LXC> vtysh -C "show ipv6 route"

O>* 2001:db8:1::4/128 [100/200] via ....
```
#### Notice: route aggregation environments

If the remote IP address is aggregated, Beluganos may fail to create tunnel interface. To avoid this, the aggregation address should be configured in advance. The file of `ribxd.conf` at the container is the configuration file.

```
LXC> vtysh -C "show ipv6 route"

O>* 2001:db8:1::/64 [100/200] via ....      # <-- aggregated!

LXC> vi /etc/beluganos/ribxd.conf

[nla]
core = "127.0.0.1:50061"
api  = "127.0.0.1:50062"

### ADD FOLLOWING SETTINGS ###
[[nla.iptun]]
nid = 0
remotes = [
        "2001:db8:1::/64",
        ]
```

Note that route aggregated environments is just beta support. This restriction will be improved future release.

### 2. Create IP tunnel's interface

You have two options to create tunnel interface.

#### By Linux command

```
LXC> ip tunnel add tun1 mode ip4ip6 remote 2001:db8:1::4 local 2001:db8:2::1
LXC> ip link set tun1 up
```

#### By BGP

Dynamic tunnel creation by BGP is also supported by Beluganos. The service of `ribt` has responsibility for dynamic configuration of IP tunneling.

In this feature, **GoBGP** is used by Beluganos. 

##### Prepare for dynamic configuration

In advance, some configuration is required. At first, local address of IP tunneling should be configured.

```
# Linux style
LXC> vi /etc/beluganos/ribtd.conf

# ansible
Beluganos$ vi etc/playbooks/roles/lxd/files/<container-name>/ribtd.conf

TUNNEL_LOCAL6=2001:db8:2::/64    # <-- Address range for local address
```

##### Add to GoBGP's RIB

Please advertise IP-IP tunneling route by BGP. Note that Beluganos requires extended community of **tunnel encapsulation attribute** with value **14** in MP-BGP. For more detail, please refer [IANA documents](https://www.iana.org/assignments/bgp-parameters/bgp-parameters.xhtml#tunnel-types).

If you want to install RIB manually, you can use following commands:

```
LXC> gobgp global rib add -a ipv4 30.0.0.0/24 nexthop 2001:db8:1::4 encap ipv6
```

## Appendix

### Restrictions of IP-IP tunneling

In Beluganos's termination feature, there are some restrictions of IPv6 extended header. In EdgeRouter2, "Tunnel Encapsulation Limit" should be "none".

```
# Example of Linux command
EdgeRouter2> ip tunnel add tun1 mode ip4ip6 local 2001:db8:1::4 remote 2001:db8:2::1 encaplimit none
```

### Configuration of `ribtd.conf`

```
LXC> vi /etc/beluganos/ribtd.conf

# -*- coding: utf-8 -*-
 API_LISTEN_ADDR=localhost:50051
 ROUTE_FAMILY=ipv4-unicast
 TUNNEL_LOCAL4=127.0.0.1/32
 TUNNEL_LOCAL6=2010:2020::/64
 TUNNEL_PREFIX=tun
 TUNNEL_TYPE_IPV6=14
 TUNNEL_TYPE_FORCE=0
 TUNNEL_TYPE_DEFAULT=0
 # debug
 DEBUG="-v"
 DUMP_TABLE_TIME=0

```

- `ROUTE_FAMILY`
	- The address family which will be treated as IP tunneling route.
		- `ipv4-unicast`: for IPv4 over IPv6 tunneling
		- `ipv6-unicast`: for IPv6 over IPv6 tunneling
- `TUNNEL_LOCAL4` or `TUNNEL_LOCAL6`
	- The range of tunnel local address. CIDR.
- `TUNNEL_TYPE_IPV6`
	- The value of tunnel encapsulation attribute. Set `14`.