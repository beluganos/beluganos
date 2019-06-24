# Setup guide for OF-DPA
This document describes about hardware setup to use OF-DPA.

## Pre-requirements

- The installation is required in advance. Please refer [install.md](install.md) before proceeding.
- The setup of Beluganos is required in advance. Please refer [setup.md](setup.md) before proceeding.
- The installation of OpenNetworkLinux for your white-box switches is required in advance. Please refer [setup-hardware.md](setup-hardware.md) before proceeding.

## Install OF-DPA

OF-DPA is a ASIC driver. After installing OpenNetworkLinux, OF-DPA settings are required to use Beluganos.

### Get binary (.deb)

Please get the binary of OF-DPA. The required version of OF-DPA is 2.0.4. The following website is published OF-DPA binary. 

- [Edge-core's repository](https://github.com/edge-core/beluganos-forwarding-app)
- [Broadcom's repository](https://github.com/Broadcom-Switch/of-dpa)

In this section, using `ofdpa-2.0-ga_2.0.4.0+accton2.4-1_amd64.deb` which got from Edge-core's repository is assumed.

### Install

The following steps are required in case of the first time to use OF-DPA.

#### (Step 1) Transfer OF-DPA binary

Transfer the binary to OpenNetworkLinux. For example, SCP or SFTP are assumed. Assumed file name is here:

- `ofdpa-2.0-ga_2.0.4.0+accton2.4-1_amd64.deb`

#### (Step 2) Install OF-DPA

```
ONL> dpkg -i --force-overwrite ofdpa-2.0-ga_2.0.4.0+accton2.4-1_amd64.deb
```

## Configure OF-DPA

#### (Step 1) Link up required ports

In default settings, almost all physical port is set to down. For example, you want to use port `30`, following step is required:

```
ONL> echo 0 > /sys/bus/i2c/devices/30-0050/sfp_tx_disable
```

If you prefer to up all ports, following example of scripts is recommended.

```
#! /bin/bash

for ((i = 2; i <= 54; i = i + 1)); do
  echo 0 > /sys/bus/i2c/devices/$i-0050/sfp_tx_disable
done
```

#### (Step 2) Start OF-DPA and OpenFlow agent

To start, following commands are required:

```
ONL> service ofdpa start

~~~(Please wait over 15 sec.)~~~

ONL> brcm-indigo-ofdpa-ofagent --controller=<BeluganosVM-IP> --dpid=<Agent-dpid>
```

- `<BeluganosVM-IP>`: Specify Beluganos's IP address. Please note that you already specify this IP address in `create.ini` at [install.md](install.md).
- `<Agent-dpid>`: Specify OpenFlow DPID. Please note that you already specify this ID in `fibc.yml` at [setup.md](setup.md). The default settings is `14`.

If you can confirm that OF-DPA is launched properly, please stop it at once by <kbd><kbd>Ctrl</kbd>+<kbd>C</kbd></kbd> and `service ofdpa stop`. After configure, please re-run it.

## Next Steps

The prepare to try Beluganos is almost finished! In the next step, you should configure router settings like IP address, VLAN, and protocol settings. Please refer [configure.md](configure.md).