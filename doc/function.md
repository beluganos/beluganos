# Beluganos Function Matrix

Supported feature and planned feature of Beluganos is described in this page. The implementation plan is subject to change.

## Basic

### Summary

| Function            | OF-DPA    | OFSwitch  | OpenNSL |
|:--------------------|:----------|:----------|:--------|
|L2 switching         |TBD        |TBD        |Planned  |
|IPv4 routing         |Yes        |Yes        |Yes      |
|IPv6 routing         |Planned    |Planned    |Yes      |
|IP Multicast (v4/v6) |TBD        |TBD        |TBD      |
|IP/MPLS              |Yes        |Yes        |Planned  |

### Basic - L1 / Ethernet

| Function          | OF-DPA | OFSwitch | OpenNSL |
|:------------------|:-------|:---------|:--------|
|Breakout cable     | TBD    | TBD      | Yes     |
|Auto negotiations  | Yes    | Yes      | Yes     |
|Protocol MTU       | Yes    | Yes      | Yes     |
|TTL check          | Yes    | Yes      | Yes     |
|Carrier delay      | TBD    | TBD      | Planned |

### Basic - L2 swiching

| Function                      | OF-DPA  | OFSwitch | OpenNSL  |
|:------------------------------|:--------|:---------|:---------|
|802.1q (VLAN)                  |Yes \*1  |Yes       |Yes       |
|802.1ad (Q-in-Q)               |TBD      |TBD       |Planned   |
|Link aggregation               |Yes \*2  |Yes       |Planned   |
|802.3ad (LACP)                 |Yes      |Yes       |Planned   |
|RSTP                           |TBD      |TBD       |Planned   |
|MSTP                           |TBD      |TBD       |Planned   |
|Loop avoidance (storm-control) |TBD      |TBD       |Planned   |
|LLDPv2                         |Planned  |Planned   |Planned   |
|LLDPv3                         |Planned  |Planned   |Planned   |
|IGMP                           |TBD      |TBD       |Planned   |
|IGMP snooping                  |TBD      |TBD       |Planned   |
|MLD                            |TBD      |TBD       |Planned   |
|MLD snooping                   |TBD      |TBD       |Planned   |
|Proxy ARP                      |TBD      |TBD       |Planned   |
|Proxy ND                       |TBD      |TBD       |Planned   |


- \*1: Because of OF-DPA restrictions, non-tagged packet is not worked properly.
- \*2: Because of OF-DPA restrictions, packet load-balancing is not worked properly in MPLS environments.

### Basic - L3 Routing

| Function                            | OF-DPA   | OFSwitch | OpenNSL  |
|:------------------------------------|:---------|:---------|:---------|
|ARP                                  |Yes       |Yes       |Yes       |
|Neighbor discovery                   |Planned   |Planned   |Yes       |
|BGP                                  |Yes       |Yes       |Yes       |
|BGP - BGP4+                          |Planned   |Planned   |Yes       |
|BGP - route reflector                |Yes       |Yes       |Yes       |
|BGP - MP-BGP                         |Yes       |Yes       |TBD       |
|BGP - BMP                            |Yes       |Yes       |Yes       |
|BGP - route filter                   |Yes       |Yes       |Yes       |
|BGP - BGP-flowspec                   |Yes       |Yes       |Yes       |
|BGP - BGP-LS                         |Yes       |Yes       |Yes       |
|BGP - multipath                      |TBD       |TBD       |Planned   |
|BGP - graceful restart               |TBD       |TBD       |Planned   |
|Static routing for IPv4              |Yes       |Yes       |Yes       |
|Static routing for IPv6              |Planned   |Planned   |Yes       |
|OSPFv2                               |Yes       |Yes       |Yes       |
|OSPFv2 - virtual links               |Yes       |Yes       |Yes       |
|OSPFv3                               |Planned   |Planned   |Yes       |
|ISIS for IPv4                        |Yes       |Yes       |Yes       |
|ISIS for IPv6                        |Planned   |Planned   |Yes       |
|VRRP                                 |Planned   |Planned   |Planned   |
|VRRPv3                               |Planned   |Planned   |Planned   |
|PIM-SSM for IPv4                     |TBD       |TBD       |TBD       |
|PIM-SSM for IPv6                     |TBD       |TBD       |TBD       |
|SVI (switch virtual interface)       |TBD       |TBD       |Planned   |
|Inter-VLAN routing                   |TBD       |TBD       |Planned   |

- BGP functions are depend on GoBGP.
- Other routing protocol's functions are depend on FRRouting.

### Basic - MPLS

| Function                 | OF-DPA   | OFSwitch | OpenNSL |
|:-------------------------|:---------|:---------|:--------|
|IP/MPLS - MPLS SWAP       |Yes       |Yes       |Planned  |
|IP/MPLS - MPLS POP        |Yes       |Yes       |Planned  |
|IP/MPLS - MPLS PUSH       |Yes       |Yes       |Planned  |
|LDP (explicit-null)       |Yes       |Yes       |Planned  |
|LDP (implicit-null)       |Yes       |Yes       |Planned  |
|RSVP-TE                   |Planned   |Planned   |Planned  |
|Fast Reroute              |Planned   |Planned   |Planned  |
|Segment Routing with OSPF |Planned   |Planned   |Planned  |
|Segment Routing with ISIS |Planned   |Planned   |Planned  |
|SR-TE                     |TBD       |TBD       |Planned  |
|MPLS TTL                  |Yes       |Yes       |Planned  |

### Bacis - L4

| Function           | OF-DPA  | OFSwitch | OpenNSL |
|:-------------------|:--------|:---------|:--------|
|DHCP server         |Planned  |Planned   |Planned  |
|DHCP relay          |Planned  |Planned   |Planned  |

### Basic - Load-balancing

| Function          | OF-DPA  | OFSwitch | OpenNSL |
|:------------------|:--------|:---------|:--------|
|IP ECMP            |Planned  |Planned   |Planned  |
|MPLS ECMP          |No \*1   |TBD       |No \*1   |
|RTAG7 hash         |Planned  |Planned   |Planned  |

- \*1: Because of ASIC restrictions, MPLS ECMP functions cannot be supported.

## Overlay

### Summary

| Function             | OF-DPA     | OFSwitch   | OpenNSL    |
|:---------------------|:-----------|:-----------|:-----------|
| MPLS-VPN (L3)        | Yes        | Yes        | TBD        |
| VPLS (L2)            | TBD        | TBD        | Planned    |
| EVPN-VXLAN (L2)      | TBD        | TBD        | Planned    |
| EVPN-VXLAN (L3)      | TBD        | TBD        | Planned    |
| IP tunneling (L3)    | TBD        | TBD        | Planned    |
| GRE tunneling (L3)   | TBD        | TBD        | Planned    |

### Overlay - L3VPN

| Function                             | OF-DPA  | OFSwitch | OpenNSL |
|:-------------------------------------|:--------|:---------|:--------|
|Virtual router (VRF-Lite)             |Yes      |Yes \*1   |Planned  |
|Route distinguisher                   |Yes      |Yes       |Planned  |
|VRF for L3VPN                         |Yes      |Yes \*1   |TBD      |
|Intra-AS MPLS-VPN PE for IPv4         |Yes      |Yes       |TBD      |
|Intra-AS MPLS-VPN 6PE/6VPE            |Planned  |Planned   |TBD      |
|BGP4 between PE-CE                    |Yes      |Yes       |TBD      |
|BGP4+ between PE-CE                   |Planned  |Planned   |TBD      |
|OSPFv2 between PE-CE                  |TBD \*2  |TBD \*2   |TBD      |
|OSPFv3 between PE-CE                  |TBD \*2  |TBD \*2   |TBD      |
|Direct connections between PE-CE IPv4 |Yes      |Yes       |TBD      |
|Direct connections between PE-CE IPv6 |Planned  |Planned   |TBD      |
|Static routing to CE IPv4             |Yes      |Yes       |TBD      |
|Static routing to CE IPv6             |Planned  |Planned   |TBD      |
|VPN label mapping per VRF             |Yes      |Yes       |TBD      |
|VPN label mapping per route           |TBD      |TBD       |TBD      |

- \*1: OpenFlow metadata fields is utilized for VRF. OFSwitch needs to support metadata fields.
- \*2: No DN bits support.

### Overlay - VXLAN

| Function              | OF-DPA  | OFSwitch | OpenNSL |
|:----------------------|:--------|:---------|:--------|
|VXLAN encap/decap      |TBD      |TBD       |Planned  |
|Ingress replication    |TBD      |TBD       |Planned  |
|EVPN BGP               |TBD      |TBD       |Planned  |
|Multi-homing           |TBD      |TBD       |Planned  |

## Management

### Summary

| Function        | OF-DPA           | OFSwitch         | OpenNSL          |
|:----------------|:-----------------|:-----------------|:-----------------|
|SSH              |Yes               |Yes               |Yes               |
|CLI              |Partially Yes \*1 |Partially Yes \*1 |Partially Yes \*1 |
|NETCONF/YANG     |Yes               |Yes               |Yes               |
|SNMP MIB         |Yes               |Yes               |Yes               |
|SNMP Trap        |Yes               |Yes               |Yes               |
|Telemetry        |Planned           |Planned           |Planned           |
|Syslog           |Yes \*2           |Yes \*2           |Yes \*2           |
|Mirroring        |Planned           |Planned           |Planned           |
|Management VRF   |Planned           |Planned           |Planned           |
|NTP              |Yes               |Yes               |Yes               |

- \*1: Only in demonstration level. Partially "show" command is supported.
- \*2: Only in demonstration level.

### Management - Configuration

| Function                            | OF-DPA           | OFSwitch         | OpenNSL          |
|:------------------------------------|:-----------------|:-----------------|:-----------------|
|Config backup/rollback               |Planned           |Planned           |Planned           |
|Linux style                          |Yes               |Yes               |Yes               |
|ansible for initial configuration    |Yes               |Yes               |Yes               |
|NETCONF/YANG                         |Yes \*1           |Yes \*1           |Yes \*1           |
|Candidate config at NETCONF          |Yes               |Yes               |Yes               |
|OpenConfig                           |Partially Yes \*1 |Partially Yes \*1 |Partially Yes \*1 |

- \*1: Supported configuration is published at [netconf/etc/openconfig](https://github.com/beluganos/netconf/tree/master/etc/openconfig).

### Management - SNMP MIB

| Function | OF-DPA | OFSwitch | OpenNSL |
|:---------|:-------|:---------|:--------|
| Traffic counter per physical interfafce | Yes | Yes | Yes |
| Traffic counter per VLAN | TBD | TBD | Planned |
| Interface status | Yes | Yes | Yes |
| System status (CPU, etc.) | Planned | Planned | Planned |

- Supported OID is published at [beluganos/doc/feature-snmp.md](https://github.com/beluganos/beluganos/tree/master/doc/feature-snmp.md).

### Management - Misc

| Function                | OF-DPA  | OFSwitch | OpenNSL |
|:------------------------|:--------|:---------|:--------|
|Ping                     |Yes      |Yes       |Yes      |
|Traceroute               |Yes      |Yes       |Yes      |
|Interface admin shutdown |Planned  |Planned   |Planned  |
|sflow                    |TBD      |TBD       |Planned  |
|SPAN (mirroring)         |Planned  |Planned   |Planned  |
|RSPAN (mirroring)        |TBD      |TBD       |Planned  |

## Others

### Security

| Function                     | OF-DPA  | OFSwitch | OpenNSL |
|:-----------------------------|:--------|:---------|:--------|
|ACL                           |Planned  |Planned   |Planned  |
|ACL - S/D MAC                 |Planned  |Planned   |Planned  |
|ACL - S/D IP                  |Planned  |Planned   |Planned  |
|ACL - S/D L4                  |Planned  |Planned   |Planned  |
|SSH public key authentication |Yes      |Yes       |Yes      |

### QoS

| Function                | OF-DPA  | OFSwitch | OpenNSL |
|:------------------------|:--------|:---------|:--------|
|Policing                 |TBD      |TBD       |Planned  |
|Shaping                  |TBD      |TBD       |TBD      |
|Diffserv - classify      |No       |No        |TBD      |
|Diffserv - scheduling    |No       |No        |TBD      |
|Diffserv - marking       |No       |No        |Planned  |

### Hardware

| Type          | OF-DPA   | OFSwitch | OpenNSL  |
|:--------------|:---------|:---------|:---------|
|Trident+       |Verified  |-         |No        |
|Trident II     |Verified  |-         |Verified  |
|Trident II+    |Verified  |-         |Verified  |
|Tridnet III    |No        |-         |Planned   |
|Tomahawk       |No        |-         |Verified  |
|Tomahawk II    |No        |-         |TBD       |
|DNX series     |No        |-         |Planned   |
|x86            |-         |Verified  |-         |

- OF-DPA or OpenNSL support is required to use above ASICs.
- In x86 case, verification is executed by Lagopus.

## Appendix

#### The difference of hardware

In Beluganos, some features are supported only specific hardware because of some limitations. For example, there is some specification difference between **OF-DPA** and **OpenNSL** which is a open ASIC API provided by Broadcom. In this document, supported feature is described for each API. 

Moreover, x86 **OFSwitch** (OpenFlow switch) is also supported by Beluganos. From this perspective, Lagopus is used for verification. In OFSwitch, almost all supported feature is same as OF-DPA, but there is some difference.

#### Legend of this function table

| Value        | Description |
|:------------:|:------------|
|Yes           | Supported.|
|Partially Yes | Supported, but some restriction is remained.|
|Planned       | NOT supported yet, but there is a plan for implimentation.  (SUBJECT TO CHANGE)|
|TBD           | NOT supported yet.|
|No            | NOT supported, and will not supported because of technical limitation.|


