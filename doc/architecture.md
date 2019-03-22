# Architecture

Beluganos consists from some layer. The routing engine layer create routing information base (RIB), and hardware control layer may install forwarding information base (FIB) to hardware. The forwarding abstraction layer which converts from RIB to FIB is the main component of Beluganos.

<img src="img/high-level.png" alt="High-level architecture" width="600px">


## Routing engine layer

Routing engine layer creates RIB by routing protocols. This layer is placed on [Linux container](https://linuxcontainers.org/). In this layer, you can use any routing protocols as you like. The applicable requirement for routing protocols is that best paths will be installed correctly to Linux kernel. Beluganos get the path information from [netlink](https://tools.ietf.org/html/rfc3549) and install to white-box switches.

Actually, Beluganos does NOT contain any routing protocol stack. In verification, [FRRouting](https://frrouting.org/) or [GoBGP](https://osrg.github.io/gobgp/) is used. If you use simple setup script to build Beluganos, both FRRouting and GoBGP will be installed automatically. Please note, in MPLS-VPN case, only GoBGP is supported as MP-BGP stack in order to deal with vpnv4 route.

## Hardware control layer

Hardware control layer has ASIC control functions. In Beluganos, **[OpenNSL](https://jp.broadcom.com/products/ethernet-connectivity/software/opennsl/)** and **[OF-DPA](https://www.broadcom.com/products/ethernet-connectivity/software/of-dpa)** are supported to control ASIC. OpenNSL and OF-DPA provides open access for ASIC, including. IP/MPLS functions, wire-rate packet forwarding, and TCAM based longest match.

**OpenNSL** (SAI) is an open API for controlling ASIC. We developed OpenNSL agent which call this API, and use this. Moreover, we also released OpenNSL library for Go lang at [beluganos/go-opennsl](https://github.com/beluganos/go-opennsl).

**OF-DPA** is also an open API, but it is developed for OpenFlow controllers. Because OF-DPA agent is provided by some vendors, Beluganos developed as OpenFlow controller.

For OpenFlow switches: Because OF-DPA is compliant with [OpenFlow](https://www.opennetworking.org/sdn-resources/openflow)1.3, any OpenFlow switches may be applied for Beluganos's data-plane if you want. If you don't have white-box switches hardware yet, [Lagopus](http://www.lagopus.org/) is recommended for verifying Beluganos. Lagopus is one of the most compatible switches for OpenFlow 1.3.

## Forwarding abstraction layer

<img src="img/flow.png" alt="Flow modification architecture" width="420px">

This figure shows the principals of forwarding abstraction layer. Note that this figure shows it in only OF-DPA case. Forwarding abstraction layer converts path information. This layer has three important components: **NLA, RIBC, FIBC**.

**NLA (NetLink Abstraction)** is the parser of [netlink](https://tools.ietf.org/html/rfc3549). In Beluganos's architecture, kernel's path information is the original data of RIB. The main efforts of NLA is getting path information from Linux kernel.

**RIBC (Routing Information Base Controller)** and **FIBC (Forwarding Information Base Controller)** create FIB entry from path information. The main role of RIBC is creating base information for each pipelines of ASIC processing.

Since OF-DPA support OpenFlow Table Type Pattern (TTP), RIBC will send message which is separated into units of OpenFlow TTP. Yet in FIBC, creating OpenFlow entry and converting IF name to hardware are the main efforts. [Ryu](https://osrg.github.io/ryu/) is used as OpenFlow controller. Therefore, FIBC is just Ryu app. In terms of OpenNSL, almost same separation is required by OpenNSL API.

The reason why RIBC and FIBC are separated is to support MPLS-VPN. In MPLS-VPN case, multiple linux container will be launched as multiple VRF, and multiple RIBC will send message to single FIBC. Note that RIBS (RIB sync) process will be worked to redistribute route information in MPLS-VPN case.


## Management layer

The management function will be required for network OS. NETCONF is already available at [https://github.com/beluganos/netconf](https://github.com/beluganos/netconf). [ansible](configure-ansible.md) is also available for setup. In addition, SNMP and syslog feature is supported.