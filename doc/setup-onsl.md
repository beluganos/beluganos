# Setup guide for OpenNSL
This document describes about hardware setup to use OpenNSL.

## Pre-requirements

- The installation is required in advance. Please refer [install.md](install.md) before proceeding.
- The setup of Beluganos is required in advance. Please refer [setup.md](setup.md) before proceeding.
- The installation of OpenNetworkLinux for your white-box switches is required in advance. Please refer [setup-hardware.md](setup-hardware.md) before proceeding.

## Install OpenNSL

OpenNSL is a ASIC driver. After installing OpenNetworkLinux, OpenNSL settings are required to use Beluganos.

### Get binary (.deb)

Please get the binary of OpenNSL. The required version of OpenNSL 3.5, and following website publish required binary.

- [Edge-core's blog](https://support.edge-core.com/hc/en-us/sections/360002115754-OpenNSL)
- [Broadcom's repository](https://github.com/Broadcom-Switch/OpenNSL)

If you cannot find OpenNSL binary (.deb), you can build this personally. Please refer [appendix F](#appendix-f-building-opennsl-deb-file) of this documents to build.

### Get OpenNSL's agent (.deb)

OpenNSL is just driver library. This is not contain agent in spite of OF-DPA case. In `create.sh`, Beluganos already created the OpenNSL agent in your server. In this step, you can set up OpenNSL agent.

```
Beluganos$ cd ~/etc/gonsl
Beluganos$ ./make.sh gonsld prepare
Beluganos$ ./make.sh gonsld check
Beluganos$ ./make.sh gonsld deb
Beluganos$ ls gonsld_1.0.0-1_amd64.deb
gonsld_1.0.0-1_amd64.deb
```

Now, you can obtain `gonsld_1.0.0-1_amd64.deb`.

### Install OpenNSL

#### (Step 1) Transfer OpenNSL binary and agent

Transfer the binary to OpenNetwork Linux. For example, SCP or SFTP are assumed. Assumed file name is here:

- `opennsl-accton_3.5.0.3+accton4.0-2_amd64.deb`
- `gonsl_1.0.0-1_amd64.deb`

#### (Step 2) Install OpenNSL and agent

```
ONL> dpkg -i opennsl-accton_3.5.0.3+accton4.0-2_amd64.deb
ONL> dpkg -i gonsld_1.0.0-1_amd64.deb
```

### Configure OpenNSL

#### (Step 1) Configure agent

```
ONL> vi /etc/beluganos/gonsld.yaml

---
dpaths:
  default:
    dpid: <Agent-dpid>
    addr: <BeluganosVM-IP>
    port: 50070
```

- `<BeluganosVM-IP>`: Specify Beluganos's IP address. Please note that you already specify this IP address in `create.ini` at [install.md](install.md).
- `<Agent-dpid>`: Specify OpenFlow DPID. Please note that you already specify this ID in `fibc.yml` at [setup.md](setup.md).

```
ONL> vi /etc/beluganos/gonsld.conf

# DEBUG="-v"
```

- `DEBUG`: Set debug flag. If you want to use, please comment out.

#### (Step2) Configure OpenNSL

`opennsl.conf` is the configuration file of OpenNSL. If it is not exist, default settings will be used. You can change the interface speed from default speed by `opennsl.conf`. For more detail, please refer [appendix B](#appendix-b-change-opennsl-configurations).

#### (Step3) Start OpenNSL and agent

```
ONL> /etc/init.d/gonsl start
ONL> /etc/init.d/gonsl status
```

If you can confirm that `gonsl` is launched properly, please stop it at once by `/etc/init.d/gonsl stop`. After configure, please re-run it.

## Next Steps

The prepare to try Beluganos is almost finished! In the next step, you should configure router settings like IP address, VLAN, and protocol settings. Please refer [configure.md](configure.md).

---

## Appendix

### (Appendix A) `create.sh` settings for OpenNSL

In default settings of `create.ini`, OpenNSL agent is compiled automatically.

```
$ cd ~
$ grep ONSL create.ini

BEL_ONSL_ENABLE=yes
BEL_ONSL_PLATFORM=as7712
BEL_ONSL_PKG="github.com/beluganos/go-opennsl"
```

If you change above settings, following step is required to apply.

```
$ cd ~
$ ./create.sh opennsl
```

### (Appendix B) Change OpenNSL configurations

```
> mkdir /etc/opennsl/
> vi /etc/opennsl/opennsl.conf
```

Note that the file `opennsl.conf` should be configure as your hardware. The sample file is available at [Edge-core's blog](https://support.edge-core.com/hc/en-us/sections/360002115754-OpenNSL) or [Broadcom's repository](https://github.com/Broadcom-Switch/OpenNSL). In Broadcom's repository, sample file is available by the name of "config.as7712" and so on.

If this file don't exists, default settings will be set by `gonsld`.

### (Appendix C) Check interface speed

You can check interface speed (1G or 10G or 40G) by following steps. Please note that following commands is available when gonsl (OpenNSL agent) is stopped.

```
> lsmod
> opennsl_setup insmod
> cd /usr/bin/opennsl-accton/examples
> ./example_drivshell

   ~~~~ (snipped) ~~~

drivshell> ps
                 ena/    speed/ link auto    STP                  lrn  inter   max  loop
          port  link    duplex scan neg?   state   pause  discrd ops   face frame  back
       xe0(  1)  down   10G  FD   SW  No   Forward  TX RX   None   FA    SFI  9412
       xe1(  2)  up     10G  FD   SW  No   Forward  TX RX   None   FA    SFI  9412
       xe2(  3)  down   10G  FD   SW  No   Forward  TX RX   None   FA    SFI  9412
       ~~~ (snipped) ~~~
```

### (Appendix D) Change interface speed

To change interface speed (1G or 10G or 40G), `opennsl.conf` should be changed. For detail, please refer [Edge-core's blog](https://support.edge-core.com/hc/en-us/articles/360010154034-OpenNSL-3-5-0-3).

### (Appendix E) Log files of gonlsd

At `/var/log/gonsld.log`.

### (Appendix F) Building OpenNSL .deb file

If you cannot find OpenNSL binary at [OpenNSL's repository](https://github.com/Broadcom-Switch/OpenNSL), you should build it personally.

#### Environments

Same as "1. Building OpenNetworkLinux" at [setup-hardware.md](setup-hardware.md#1-build-opennetworklinux). Please refer this documents.

#### Specify your hardware

Please change `OPENNSL_HWTYPE`.

```
Beluganos-server$ cd ~/etc/gonsl
Beluganos-server$ vi onl.sh

# In case of AS7712
OPENNSL_HWTYPE=as7712

# In case of AS5812
OPENNSL_HWTYPE=as5812
```

#### Build

```
Beluganos-server$ cd ~/etc/gonsl
Beluganos-server$ ./make.sh opennsl prepare
Beluganos-server$ ./make.sh opennsl check
Beluganos-server$ ./make.sh opennsl deb
```

Now, you can obtain `opennsl_*_amd64.deb`.