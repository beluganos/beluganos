# Setup guide for OpenNSL
This document describes about hardware setup to use Beluganos.

## Pre-requirements
- Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.
- In this document, Edge-core AS7712-32X switch is assumed to use Beluganos. All hardware which supports OpenNSL 3.5 is acceptable, but in this case, please look it up yourself to install OpenNSL.

## 1. Install OpenNetworkLinux

### Get binary

Please get binary from [OpenNetworkLinux's website](https://opennetlinux.org/binaries/). Following version is recommended:

```
ONL-2.0.0-ONL-OS-2018-01-09.1646-04257be-AMD64-INSTALLED-INSTALLER
```
After getting binary, you can install OpenNetworkLinux via DHCP or TFTP. In this documents, only TFTP methods are described.

### Install via TFTP

#### (Step 1) Connect console cable

In TFTP installation, the connections of console cable is required to communicate between white-box switch and your working PC. In following steps, the strings of `>` represent the console screen.

#### (Step 2) Boot hardware

Plug in power cable to boot switch.

#### (Step 3) Select "ONIE install"

In booting process, GRUB menu is apperrd. Select `ONIE` -> `ONIE install` by <kbd>↑</kbd> (Up) or <kbd>↓</kbd> (Down) keys to start install.

#### (Step 4) Stop DHCP discovery

In default settings of ONIE, DHCP discovery will be started. Stop DHCP discovery by following command:

```
> onie-discovery-stop
```

#### (Step 5) Configure management address

Configure the settings of physical outbound port (management port). For example, if you want to set ``172.16.0.57/24``, following step is required:

```
> ifconfig eth0 172.16.0.57 netmask 255.255.255.0 up
```

#### (Step 6) Start to install

Start to install. For example, if your tftp server is `172.16.0.59` and ONL version is `ONL-2.0.0-ONL-OS-2018-01-09.1646-04257be-AMD64-INSTALLED-INSTALLER`, following step is required:

```
> onie-nos-install tftp://172.16.0.59/ONL-2.0.0-ONL-OS-2018-01-09.1646-04257be-AMD64-INSTALLED-INSTALLER
```

#### (Step 7) Login

Once finished a install, it will be rebooted automatically. Please log in. The default user-name to is `root` and password is `onl`.

## 2. Setup for OpenNSL

OpenNSL is a ASIC driver. After installing OpenNetworkLinux, OpenNSL settings are required.

### Get binary

Please get the binary of OpenNSL. The required version of OpenNSL 3.5, and following website publish required binary.

- [Edge-core's blog](https://support.edge-core.com/hc/en-us/sections/360002115754-OpenNSL)
- [Broadcom's repository](https://github.com/Broadcom-Switch/OpenNSL)

In this documents, `opennsl-accton_3.5.0.3+accton4.0-2_amd64.deb` from Edge-core's blog is used to describe following steps.

### Compile OpenNSL's agent

OpenNSL is just driver library. This is not contain agent in spite of OF-DPA case. In `create.sh`, beluganos already created the OpenNSL agent in your server. In this step, in order to set up OpenNSL agent, required files are described.

The required files should be located at `~/beluganos/etc/gonsl` of **not "OpenNetworkLinux" but "Beluganos server"**.

Next, transfer required files to OpenNetworkLinux via SCP or SFTP. The required files are described here:

```
beluganos/
  etc/
    gonsl/
      gonsld                # Copy from ~/go/bin/gonnsl from Beluganos (*1)
      opennsl.conf          # Get and copy OpenNSL files (*2)
                            #   Ex) config.as7712
      libopennsl.so.1       # Get and copy OpenNSL files (*2)
      linux-kernel-bde.ko   # Get and copy OpenNSL files (*2)
      linux-user-bde.ko     # Get and copy OpenNSL files (*2)
      linux-bcm-knet.ko     # Get and copy OpenNSL files (*2)
      make.sh
      files/
        gonsld.initd
        gonsld.conf
        gonsld.yaml
```

- (\*1): By `create.sh`, the binary which is `~/go/bin/gonsld` is created.
- (\*2): In AS7712-32X, there is at `bin/as7712/` of [Broadcom's repository](https://github.com/Broadcom-Switch/OpenNSL).

To compile agent, `make.sh` is prepared.

```
Beluganos-server$ cd etc/gonsl
Beluganos-server$ ./make.sh deb
Beluganos-server$ ls gonsld_1.0.0-1_amd64.deb
gonsld_1.0.0-1_amd64.deb
```

### Initial settings

The following steps are required in case of the first time to use OF-DPA.

#### (Step 1) Configure management address

Configure the settings of physical outbound port (management port). For example, if you want to set ``172.16.0.57/24``, following step is required:

```
> ifconfig ma1 172.16.0.57 netmask 255.255.255.0 up
```

This settings are not permanent. To set this settings permanently, following steps are required:

```
> echo 'ip addr add 172.16.0.57/24 dev ma1' >> /mnt/onl/data/rc.boot
> chmod a+x /mnt/onl/data/rc.boot
```

#### (Step 2) Transfer OpenNSL binary and agent

Transfer the binary to OpenNetwork Linux. For example, SCP or SFTP are assumed. Assumed file name is here:

- `opennsl-accton_3.5.0.3+accton4.0-2_amd64.deb`
- `gonsl_1.0.0-1_amd64.deb`

#### (Step 3) Setup agent settings

```
> vi /etc/beluganos/gonsld.yaml

---
dpaths:
  default:
    dpid: <Agent-dpid>
    addr: <BeluganosVM-IP>
    port: 50070
```

- `<BeluganosVM-IP>`: Specify Beluganos's IP address. Please note that you already specify this IP address in `create.ini` at [install-guide.md](install-guide.md).
- `<Agent-dpid>`: Specify OpenFlow DPID. Please note that you already specify this ID in `fibc.yml` at [setup-guide.md](setup-guide.md).


#### (Step 4) Install OpenNSL and agent

```
> dpkg -i opennsl-accton_3.5.0.3+accton4.0-2_amd64.deb
> dpkg -i gonsld_1.0.0-1_amd64.deb
```

### General settings

Once you finished to do "Initial settings", other general settings are required.

#### (Step 1) Start OpenNSL agent

```
> /etc/init.d/gonsl start
> /etc/init.d/gonsl status
```

## Next steps
After reflecting your changes, please refer configure guide. You can choose two methods.

- ansible: [configure-ansible.md](configure-ansible.md)
- NETCONF over SSH: [configure-netconf.md](configure-netconf.md)

## Appendix

### `create.sh` settings

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

### Check interface speed

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

### Change interface speed

To change interface speed (1G or 10G or 40G), `opennsl.conf` should be changed. For detail, please refer [Edge-core's blog](https://support.edge-core.com/hc/en-us/articles/360010154034-OpenNSL-3-5-0-3).

### Log files of gonlsd

At `/var/log/gonsld.log`.