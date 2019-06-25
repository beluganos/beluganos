# Document index of Beluganos

## Install

- [install.md](install.md)
	- Install Beluganos and related OSS automatically

## Use case

After installing ([install.md](install.md)), you can check our use case to try Beluganos.

- [case1.md](example/case1/case1.md)
	- **Recommendation for beginner**
	- IGP router with MPLS
	- Create automatically not only Beluganos but also environment router
- [case2.md](example/case2/case2.md)
	- MPLS-VPN PE router with 2 VRF
	- Real white-box swithes and environment routers are required to try

## Setup hardware

After installing ([install.md](install.md)), you may set up about hardware.

- [setup.md](setup.md)
	- Beluganos setup to connect with white-box switch
- [setup-hardware.md](setup-hardware.md)
	- Hardware setup
- [setup-opennetworklinux.md](setup-opennetworklinux.md)
	- Installation about OpenNetworkLinux
- [setup-onsl.md](setup-onsl.md)
	- OpenNetworkLinux setup to use OpenNSL
- [setup-ofdpa.md](setup-ofdpa.md)
	- OpenNetworkLinux setup to use OF-DPA

## Configure and operation

After installing ([install.md](install.md)), you can configure Beluganos as you like.

- [configure.md](configure.md)
	- Configuration guide of Linux style.
- [configure-ansible.md](configure-ansible.md)
	- Configuration guide of IP/MPLS router by ansible.
- [configure-netconf.md](configure-netconf.md)
	- Configuration guide of router by NETCONF.
- [configure-portmapping.md](configure-portmapping.md)
	- Interface mapping information between physical and logical port. This information is required at configuration.
- [operation.md](operation.md)
	- How to start/stop Beluganos
	- How to login routing engine's console

## NETCONF components

- [netconf/etc/openconfig](https://github.com/beluganos/netconf/tree/master/etc/openconfig)
	- Yang files of Beluganos
- [netconf/doc/examples](https://github.com/beluganos/netconf/tree/master/doc/examples)
	- The examples of XML for NETCONF `<edit-config>`
- [netconf/doc/setup-guide](https://github.com/beluganos/netconf/blob/master/doc/setup-guide.md)
	- Initial settings about network-instance modules

## Feature guide

- [feature-l3vpn.md](feature-l3vpn.md)
	- Configuration guide of MPLS-VPN PE router by ansible.
- [feature-snmp.md](feature-snmp.md) 
	- SNMP feature guide (MIB, trap)
- [feature-syslog.md](feature-syslog.md) 
	- Syslog feature guide

## General

- [README.md](../README.md)
- [Beluganos-introduction.pdf](Beluganos-introduction.pdf)
	- Presentation to introduce Beluganos
- [CONTRIBUTING.md](../CONTRIBUTING.md)
- [architecture.md](architecture.md)
	- Abstraction of Beluganos's architecture
- [function.md](function.md)
	- Function matrix of Beluganos
