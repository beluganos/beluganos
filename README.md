## Beluganos
Beluganos is a **new network OS** designed for **white-box switches**, which can apply large-scale networks.

- IP Routing (BGP, OSPF, IPv6, ...) and L2 switching
- **IP/MPLS**, BGP/MPLS IP-VPNs, IP-IP tunneling
- **Interoperability** with conventional IP or IP/MPLS router
- ASIC based **hardware packet processing**

The feature matrix of Beluganos is available at [doc/function.md](doc/function.md). Beluganos was named after [beluga whale](https://en.wikipedia.org/wiki/Beluga_whale).

>note
>* The commercial version of "Beluganos" released on March 31,2023 does not use the OSS technology published on theGitHub.
>* After March 31, 2023, NTT's registered trademark "Beluganos"will be used for the commercial version of "Beluganos" and willnot be used for the OSS version on the GitHub.
>* For details of the commercial version of “Beluganos”, please click [here](https://group.ntt/en/newsrelease/2023/03/28/230328b.html).

## Architecture
Beluganos has one or more **[Linux containers](https://linuxcontainers.org/)**. The main effort of Beluganos is that the route table which is installed to Linux containers is copied to white-box switches. If you will configure router settings like IP addresses or parameter of routing protocols, you may configure the settings of Linux containers by ansible or [NETCONF](https://github.com/beluganos/netconf/). Moreover, in order to control white-box switches, **OpenNSL** or **OF-DPA** is used.

For more details, please check [doc/architecture.md](doc/architecture.md).

## Getting Started

### 1. Quick start by example case
In order to try Beluganos quickly, **some example cases are prepared**. This example can configure automatically not only Beluganos but also other routers to connect with Beluganos. If you wish to use this, please refer to [doc/example/case1/case1.md](doc/example/case1/case1.md) instead of the following description.

### 2. Step-by-step procedure

<img src="doc/img/environments.png" width="350px" alt="beluganos-install-environments">

- Step 1: Build
	- Install Beluganos and related OSS automatically.
		- Check [doc/install.md](doc/install.md).
- Step 2: Setup
	- Register your white-box switch to Beluganos.
		- Check [doc/setup.md](doc/setup.md).
 	- Initial setup of your white-box switch.
	 	- Check [doc/setup-hardware.md](doc/setup-hardware.md)
	- Initial setup of ASIC API.
		- If you use OpenNSL switch, check [doc/setup-onsl.md](doc/setup-onsl.md).
	 	- If you use OF-DPA switch, check [doc/setup-ofdpa.md](doc/setup-ofdpa.md).
- Step 3: Configure
	- Change router settings like IP address, VLAN, and routing protocols as you like.
		- Check [doc/configure.md](doc/configure.md).
   - Some advanced configuration technology is also supported.
		- To configure by ansible, check [doc/configure-ansible.md](doc/configure-ansible.md).
		- To configure by NETCONF, check [doc/configure-netconf.md](doc/configure-netconf.md).
- Step 4: Operation
	- Start Beluganos. Monitor Beluganos.
		- Check [doc/operation.md](doc/operation.md)

Other document is listed at [document index page](doc/README.md).

## Support
Github issue page and e-mail are available. If you prefer to use e-mail, please contact `msf-contact-ml [at] hco.ntt.co.jp`.

## Development & Contribution
Any contribution is encouraged. The main component is written in Go and Python. For more details, please refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Beluganos is licensed under the **Apache 2.0** license. Check [LICENSE](LICENSE).

## Project
This project is a part of [Multi-Service Fabric](https://github.com/multi-service-fabric/msf).

<img src="doc/img/multi-service-fabric.png" width="180px" alt="multi-service fabric's logomark">
