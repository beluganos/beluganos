# Beluganos Function Matrix

Currently supported feature of Beluganos is described in this page.

## Basic

|Function           |OF-DPA   |OFSwitch  |
|:------------------|:--------|:---------|
|ARP                |Yes      |Yes       |
|IPv4 Routing       |Yes      |Yes       |
|IPv6 Routing       |Planned  |Planned   |
|IP/MPLS|Yes|Yes|
|IP/MPLS - MPLS SWAP|Yes|Yes|
|IP/MPLS - MPLS POP|Yes|Yes|
|IP/MPLS - MPLS PUSH|Yes|Yes|
|Protocol MTU|Yes|Yes|
|TTL check|Yes|Yes|

## L2

| Function | OF-DPA | OFSwitch |
|:---------|:-------|:---------|
|802.1q (VLAN)|Yes \*1|Yes|
|Link Aggregation|Yes \*2|Yes|
|802.3ad (LACP)|Yes|Yes|

- \*1: Because of OF-DPA restrictions, non-tagged packet is not worked properly.
- \*2: Because of OF-DPA restrictions, packet load-balancing is not worked properly in MPLS environments.

## L3
| Function | OF-DPA | OFSwitch |
|:---------|:-------|:---------|
|BGP4|Yes|Yes|
|BGP4+|Planned|Planned|
|MP-BGP|Yes|Yes|
|BMP|Yes|Yes|
|Static Routing|Yes|Yes|
|OSPFv2|Yes|Yes|
|OSPFv3|Planned|Planned|
|ISIS for IPv4|Yes|Yes|
|ISIS for IPv6|Planned|Planned|
|PIM|No Roadmap|No Roadmap|

## Load-balancing
| Function | OF-DPA | OFSwitch |
|:---------|:-------|:---------|
|IP ECMP|Planned|Planned|
|MPLS ECMP|No Roadmap \*1|Planned|

- \*1: Because of ASIC restrictions, MPLS ECMP functions cannot be supported.

## MPLS
| Function | OF-DPA | OFSwitch |
|:---------|:-------|:---------|
|LDP|Yes|Yes|
|LDP (implicit-null)|Yes|Yes|
|RSVP-TE|Planned|Planned|
|FRR|Planned|Planned|
|Segment Routing with OSPF|Planned|Planned|

## L3VPN
| Function | OF-DPA | OFSwitch |
|:---------|:-------|:---------|
|Virtual router (VRF-Lite) |Yes|Yes \*1|
|VRF|Yes|Yes \*1|
|Intra-AS MPLS-VPN PE for IPv4|Yes|Yes|
|Intra-AS MPLS-VPN 6PE/6VPE|Planned|Planned|
|BGP4 between PE-CE|Yes|Yes|
|OSPFv2 between PE-CE|No Roadmap \*2|No Roadmap \*2|
|Direct Connections between PE-CE|Yes|Yes|
|Static Routing to CE|Yes|Yes|
|VPN Label Mapping per VRF|Yes|Yes|
|VPN Label Mapping per Route|No Roadmap|No Roadmap|

- \*1: OpenFlow metadata fields is utilized for VRF. OFSwitch needs to support metadata fields.
- \*2: No DN bits support.

## Management
| Function | OF-DPA | OFSwitch |
|:---------|:-------|:---------|
|SSH|Yes|Yes|
|CLI|Yes \*1|Yes \*1|
|NETCONF|Yes|Yes|
|YANG|Yes|Yes|
|OpenConfig|Partially yes|Partially yes|
|SNMP MIB|Planned|Planned|
|SNMP trap|Planned|Planned|
|Telemetry|Planned|Planned|
|Syslog|Yes \*2|Yes \*2|
|NTP|Yes|Yes|

- \*1: Only in demonstration level. Partially "show" command is supported.
- \*2: Only in demonstration level.

## Other
| Function | OF-DPA | OFSwitch |
|:---------|:-------|:---------|
|Ping|Yes|Yes|
|Traceroute|Yes|Yes|
|ACL (L2, L3, L4)|Planned|Planned|
|LLDP|Planned|Planned|
|SPAN (mirroring)|Planned|Planned|
