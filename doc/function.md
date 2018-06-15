# Beluganos Function Matrix

Supported feature and planned feature of Beluganos is described in this page. The implimantation plan is subject to change.

## Basic

|Function           |OF-DPA   |OFSwitch  |OpenNSL|
|:------------------|:--------|:---------|:------|
|ARP                |Yes      |Yes       |Planned|
|IPv4 Routing       |Yes      |Yes       |Planned|
|Neighbor Discovery |Planned  |Planned   |Planned|
|IPv6 Routing       |Planned  |Planned   |Planned|
|IP/MPLS|Yes|Yes|Planned|
|IP/MPLS - MPLS SWAP|Yes|Yes|Planned|
|IP/MPLS - MPLS POP|Yes|Yes|Planned|
|IP/MPLS - MPLS PUSH|Yes|Yes|Planned|
|Protocol MTU|Yes|Yes|Planned|
|TTL check|Yes|Yes|Planned|

## L2

| Function | OF-DPA | OFSwitch |OpenNSL|
|:---------|:-------|:---------|:------|
|802.1q (VLAN)|Yes \*1|Yes|Planned|
|Link Aggregation|Yes \*2|Yes|Planned|
|802.3ad (LACP)|Yes|Yes|Planned|
|RSTP|Planned|Planned|Planned|
|MSTP|Planned|Planned|Planned|

- \*1: Because of OF-DPA restrictions, non-tagged packet is not worked properly.
- \*2: Because of OF-DPA restrictions, packet load-balancing is not worked properly in MPLS environments.

## L3
| Function | OF-DPA | OFSwitch | OpenNSL |
|:---------|:-------|:---------|:--------|
|BGP4|Yes|Yes|Planned|
|BGP4+|Planned|Planned|Planned|
|BGP Route Reflector|Yes|Yes|Planned|
|MP-BGP|Yes|Yes|No Roadmap|
|BMP|Yes|Yes|Planned|
|Static Routing|Yes|Yes|Planned|
|OSPFv2|Yes|Yes|Planned|
|OSPFv3|Planned|Planned|Planned|
|ISIS for IPv4|Yes|Yes|Planned|
|ISIS for IPv6|Planned|Planned|Planned|
|VRRP|Planned|Planned|Planned|
|VRRPv3|Planned|Planned|Planned|
|PIM|No Roadmap|No Roadmap|No Roadmap|

## Load-balancing
| Function | OF-DPA | OFSwitch | OpenNSL |
|:---------|:-------|:---------|:--------|
|IP ECMP|Planned|Planned|Planned|
|MPLS ECMP|No Roadmap \*1|Planned|Planned|

- \*1: Because of ASIC restrictions, MPLS ECMP functions cannot be supported.

## MPLS
| Function | OF-DPA | OFSwitch | OpenNSL |
|:---------|:-------|:---------|:--------|
|LDP (explicit-null)|Yes|Yes|Planned|
|LDP (implicit-null)|Yes|Yes|Planned|
|RSVP-TE|Planned|Planned|Planned|
|Fast Reroute|Planned|Planned|Planned|
|Segment Routing with OSPF|Planned|Planned|Planned|
|Segment Routing with ISIS|Planned|Planned|Planned|

## L3VPN
| Function | OF-DPA | OFSwitch | OpenNSL |
|:---------|:-------|:---------|:--------|
|Virtual router (VRF-Lite) |Yes|Yes \*1|No Roadmap|
|VRF for L3VPN|Yes|Yes \*1|No Roadmap|
|Intra-AS MPLS-VPN PE for IPv4|Yes|Yes|No Roadmap|
|Intra-AS MPLS-VPN 6PE/6VPE|Planned|Planned|No Roadmap|
|BGP4 between PE-CE|Yes|Yes|No Roadmap|
|OSPFv2 between PE-CE|No Roadmap \*2|No Roadmap \*2|No Roadmap|
|Direct Connections between PE-CE|Yes|Yes|No Roadmap|
|Static Routing to CE|Yes|Yes|No Roadmap|
|VPN Label Mapping per VRF|Yes|Yes|No Roadmap|
|VPN Label Mapping per Route|No Roadmap|No Roadmap|No Roadmap|

- \*1: OpenFlow metadata fields is utilized for VRF. OFSwitch needs to support metadata fields.
- \*2: No DN bits support.

## VXLAN
| Function | OF-DPA | OFSwitch | OpenNSL |
|:---------|:-------|:---------|:--------|
|VXLAN encap/decap|No Roadmap|No Roadmap|Planned|
|Ingress Replication|No Roadmap|No Roadmap|Planned|
|EVPN BGP|No Roadmap|No Roadmap|Planned|
|Multi-homing|No Roadmap|No Roadmap|Planned|

## Management
| Function | OF-DPA | OFSwitch | OpenNSL |
|:---------|:-------|:---------|:--------|
|SSH|Yes|Yes|Planned|
|CLI|Partially Yes \*1|Partially Yes \*1|Planned|
|NETCONF/YANG|Yes \*2|Yes \*2|Planned|
|OpenConfig|Partially Yes \*2|Partially Yes \*2|Planned|
|SNMP MIB|Planned|Planned|Planned|
|SNMP Trap|Planned|Planned|Planned|
|Telemetry|Planned|Planned|Planned|
|Syslog|Yes \*3|Yes \*3|Planned|
|NTP|Yes|Yes|Planned|

- \*1: Only in demonstration level. Partially "show" command is supported.
- \*2: Supported configuration is published at [netconf/etc/openconfig](https://github.com/beluganos/netconf/tree/master/etc/openconfig).
- \*3: Only in demonstration level.

## Other
| Function | OF-DPA | OFSwitch | OpenNSL |
|:---------|:-------|:---------|:--------|
|Ping|Yes|Yes|Planned|
|Traceroute|Yes|Yes|Planned|
|ACL (L2, L3, L4)|Planned|Planned|Planned|
|Policing|No Roadmap|No Roadmap|Planned|
|LLDP|Planned|Planned|Planned|
|Carrier Delay|Planned|Planned|Planned|
|SPAN (mirroring)|Planned|Planned|Planned|
