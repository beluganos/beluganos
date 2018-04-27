<img src="doc/img/brand-logo-h.png" width="420px" alt="beluganos-logomark">

Beluganos is a **new network OS** designed for **white-box switches**, which can apply large-scale networks.

- IP Routing (BGP, OSPF, ...)
- **IP/MPLS**, BGP/MPLS IP-VPNs
- **Interoperability** with conventional IP or IP/MPLS router
- ASIC based **full hardware processing**

The feature matrix of Beluganos is available at [doc/function.md](doc/function.md). Beluganos was named after [beluga whale](https://en.wikipedia.org/wiki/Beluga_whale).

## Architecture
Beluganos has one or more [Linux containers](https://linuxcontainers.org/). The main effort of Beluganos is that the route table which is installed to Linux containers is copied to white-box switches. To configure router settings like IP addresses or parameter of routing protocols, you may configure the settings of Linux containers by ansible or [NETCONF](https://github.com/beluganos/netconf/).

For more details, please check [doc/architecture.md](doc/architecture.md).

## Getting Started

### 1. Quick start
In order to try Beluganos quickly, **some example cases are prepared**. This example can configure automatically not only Beluganos but also other routers to connect with Beluganos. If you wish to use this, please refer to [doc/example/case1/case1.md](doc/example/case1/case1.md) instead of the following description.

### 2. Step-by-step procedure

- Step1: Build
	- Install Beluganos and related OSS automatically.
	- See [doc/install-guide.md](doc/install-guide.md).
- Step2: Setup
	- Register your white-box switches to Beluganos
	- See [doc/setup-guide.md](doc/setup-guide.md).
- Step3: Configure
	- Change router settings like IP address, VLAN, and routing protocols as you like
	- To configure by ansible, please see [doc/configure-ansible.md](doc/configure-ansible.md).
	- To configure by NETCONF, please see [doc/configure-netconf.md](doc/configure-netconf.md).
- Step4: Run
	- To launch Beluganos, please see [doc/operation-guide.md](doc/operation-guide.md).


## Support
Github issue page and e-mail are available. If you prefer to use e-mail, please contact `msf-contact [at] lab.ntt.co.jp`.

## Development & Contribution
Any contribution is encouraged. The main component is written in Go and Python. If you wish to create pull-request on github.com, please kindly create your request for **develop branch**, not master branch. If you find any issue, please kindly notify us by github issue pages.

For more details, please refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Beluganos is licensed under the **Apache 2.0** license. See [LICENSE](LICENSE).

## Project
This project is a part of [Multi-Service Fabric](https://github.com/multi-service-fabric/msf).

<img src="doc/img/multi-service-fabric.png" width="180px" alt="multi-service fabric's logomark">