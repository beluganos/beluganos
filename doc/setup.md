# Setup guide
This document describes about Beluganos's setup for white-box switches.

## Pre-requirements
- The installation is required. Please refer [install.md](install.md) before proceeding.

## How to register the switch

### Edit setting file

The sample playbook is `etc/playbooks/dp-sample.yml`, and the sample configuration files are under `etc/playbooks/roles/dpath/files/sample_sw`.

The file `fibc.yml` will determine the type of white-box switches and its API of ASIC. You should change this settings like followings:

~~~~
$ cd ~/beluganos/etc/playbooks
$ vi roles/dpath/files/sample_sw/fibc.yml

---

datapaths:
  - name: sample_sw           # dpname (A-Z,a-z,0-9,_-)
    dp_id: 14                 # datapath id of your switches (integer)
    mode: generic             # "ofdpa2" or "generic" or "onsl" or "ovs"

~~~~

- `datapaths`
	- `dp_id`: Datapath ID of your switch. If you don't care this value, please do not change it.
	- `mode`: Your switch types. You can choose following options:
	   1. `onsl` : OpenNSL 3.5 switch.
	   1. `ofdpa2`: OF-DPA 2.0 switch.
	   1. `generic`: OpenFlow 1.3 compatible switches. (e.g. Lagopus)
	   1. `ovs`: Open vSwitch (Limited support).

### Reflect changes

~~~~
$ cd ~/beluganos/etc/playbooks
$ ansible-playbook -i hosts -K dp-sample.yml
~~~~

### Notice

- If you will use multiple type of hardware, copying sample files (`fibc.yml`) and changing switch name is recommended. The procedure of this is described at [Appendix A](#appendix-a-advanced-datapath-settings).

## Next steps

You should setup your hardware. 

- If you will use OpenNSL or OF-DPA switches, please refer [setup-hardware.md](setup-hardware.md) to setup.
- If you will use OpenFlow switch, please refer the switch's manual to setup. After that, please refer [configure.md](configure.md) to configure router settings.

---

## Appendix 

### (Appendix A) Advanced datapath settings

#### 1. Understand config files at a glance

The files under `etc/playbooks/` are configuration files. In this page, `dp-*.yml` and the files under `roles/dpath/` are described.

```
beluganos/
    etc/
        playbooks/
            hosts                           # Inventory file
            dp_sample.yml                   # Playbook
            roles/
                dpath/
                    vars/                   # DO NOT EDIT.
                    tasks/                  # DO NOT EDIT.
                    files/
                        common/             # DO NOT EDIT.
                        sample_sw/
                            fibc.yml        # The setting file of ASIC API
```

#### 2. Copy playbook

**The sample playbook is `etc/playbooks/dp-sample.yml`**. You may copy sample file and rename this playbook. Any switch name (dpname) are acceptable, and this name is only used for internal configurations. If your switch name is *switchA*, please edit as follows:

```
$ cd ~/beluganos/etc/playbooks
$ cp dp-sample.yml dp-switchA.yml
$ vi dp-switchA.yml

---

- hosts: hosts
  connection: local
  roles:
    - { role: dpath, dpname: switchA, mode: fibc }
  tags:
    fibc

```

The syntax of this playbook is following:

```
- hosts: hosts
  connection: local
  roles:
    - { role: dpath, dpname: <switch-name>, mode: fibc }
  tags:
    fibc
```

- roles
	- dpname (`<switch-name>`): Your preferable switch name (A-Z,a-z,0-9,_-). Do not contain space.

#### 3. Edit switch type

**The sample configuration files are under `etc/playbooks/roles/dpath/files/sample_sw`**. The files under `roles/dpath` will determine the type of white-box switches. Beluganos will optimize the writing method to FIB for each switch type. If your switch's name is *switchA*, you may copy sample file and rename to dpname as follows:

```
$ cd ~/beluganos/etc/playbooks
$ cp -r roles/dpath/files/sample_sw roles/dpath/files/switchA
$ vi roles/dpath/files/switchA/fibc.yml

---

datapaths:
  - name: switchA           # dpname (A-Z,a-z,0-9,_-)
    dp_id: 14               # datapath id of your switches (integer)
    mode: generic           # "ofdpa2" or "generic" or "ovs"

```

The syntax of `etc/playbooks/roles/dpath/files/<switch-name>/fibc.yml` is following:

```
datapaths:
  - name: <switch-name>     # dpname (A-Z,a-z,0-9,_-)
    dp_id: <switch-dp-id>   # datapath id of your switches (integer)
    mode: <switch-type>     # "ofdpa2" or "generic" or "onsl" or "ovs"
```


The value of `dp-id` means OpenFlow datapath ID of your switch. The value of `mode` should be edited for your switch types. You can choose following options:

- datapaths
	- name (`<switch-name>`): Your switch name which is already declared in `etc/playbooks/dp-switchA.yml`.
	- dp_id (`<switch-dp-id>`): OpenFlow datapath ID of your switch. Integer.
	- mode (`<switch-type>`): Your switch types. Currently Beluganos has three options.
 		1. `ofdpa2`: OF-DPA 2.0 switch.
 		1. `onsl` : OpenNSL 3.5 switch.
		1. `generic`: OpenFlow 1.3 compatible switches. (e.g. Lagopus)
		1. `ovs`: Open vSwitch (Limited support).

#### 4. Reflect changes

In case `<switch-name>` is *switchA*:

```
$ cd ~/beluganos/etc/playbooks
$ ansible-playbook -i hosts -K dp-switchA.yml
```