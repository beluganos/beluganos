# Setup guide
This document describes about Beluganos's setup for white-box switches.

## Pre-requirements
- Please refer [install-guide.md](install-guide.md) before proceeding.
- If the strings of `(mypython)` is not appeared in your terminal, you should exec `. ./setenv.sh` to set your environments.

## Config files at a glance

The files under `etc/playbooks/` are configuration files. In this page, `dp-*.yml` and the files under `roles/dpath/` are described.

~~~~
beluganos/
    etc/
        playbooks/
            hosts                           # Inventry file
            lxd-sample.yml                  # Sample
            dp_sample.yml                   # Sample
            lxd-<group_name>.yml            # Playbook
            dp-<datapath_name>.yml          # Playbook
            roles/
                dpath/
                    vars/                   # DO NOT EDIT.
                    tasks/                  # DO NOT EDIT.
                    files/
                        common/             # DO NOT EDIT.
                        sample_sw/
                            fibc.yml        # Sample
                        <switch_name>/
                            fibc.yml
                lxd/
                    vars/                   # DO NOT EDIT.
                    handlers/               # DO NOT EDIT.
                    tasks/
                        main.yml            # DO NOT EDIT.
                        create.yml          # DO NOT EDIT.
                        setup.yml           # DO NOT EDIT.
                        host.yml            # DO NOT EDIT.
                    files/
                        common/             # DO NOT EDIT.
                        sample/             # Sample
                        <container name>/
                            lxd_profile.yml
                            netplan.yaml
                            sysctl.conf
                            daemons
                            frr.conf
                            gobgp.comnf
                            gobgpd.conf
                            fibc.yml
                            ribxd.conf
~~~~

## Settings for white-box switches

### 1. switch name

**The sample playbook is `etc/playbooks/dp-sample.yml`**. You may copy sample file and rename this playbook. Any switch name (dpname) are acceptable, and this name is only used for internal configurations. If your switch name is *switchA*, please edit as follows:

~~~~
$ cd beluganos/etc/playbooks
$ cp dp-sample.yml dp-switchA.yml
$ vi dp-switchA.yml

---

- hosts: hosts
  connection: local
  roles:
    - { role: dpath, dpname: switchA, mode: fibc }
  tags:
    fibc

~~~~

The syntax of this playbook is following:

~~~~
- hosts: hosts
  connection: local
  roles:
    - { role: dpath, dpname: <switch-name>, mode: fibc }
  tags:
    fibc
~~~~

- roles
	- dpname (`<switch-name>`): Your preferable switch name (A-Z,a-z,0-9,_-). Do not contain space.

### 2. switch type
**The sample configuration files are under `etc/playbooks/roles/dpath/files/sample_sw`**. The files under `roles/dpath` will determine the type of white-box switches. Beluganos will optimize the writing method to FIB for each switch type. If your switch's name is *switchA*, you may copy sample file and rename to dpname as follows:

~~~~
$ cd beluganos/etc/playbooks
$ cp -r roles/dpath/files/sample_sw roles/dpath/files/dp-switchA
$ vi roles/dpath/files/dp-switchA/fibc.yml

---

datapaths:
  - name: switchA           # dpname (A-Z,a-z,0-9,_-)
    dp_id: 14               # datapath id of your switches (integer)
    mode: generic           # "ofdpa2" or "generic" or "ovs"

~~~~

The syntax of this playbook is following:

~~~~
datapaths:
  - name: <switch-name>     # dpname (A-Z,a-z,0-9,_-)
    dp_id: <switch-dp-id>   # datapath id of your switches (integer)
    mode: <switch-type>     # "ofdpa2" or "generic" or "ovs"
~~~~


The value of `dp-id` means OpenFlow datapath ID of your switch. The value of `mode` should be edited for your switch types. You can choose following options:

- datapaths
	- name (`<switch-name>`): Your switch name which is already declared in `etc/playbooks/dp-switchA.yml`.
	- dp_id (`<switch-dp-id>`): OpenFlow datapath ID of your switch. Integer.
	- mode (`<switch-type>`): Your switch types. Currently Beluganos has three options.
 		1. `ofdpa2`: OF-DPA 2.0 switch. [https://github.com/Broadcom-Switch/of-dpa](https://github.com/Broadcom-Switch/of-dpa)
		1. `generic`: OpenFlow 1.3 compaliable switches. (e.g. Lagopus)
		1. `ovs`: Open vSwitch (Limited support).

### Reflect changes

In case `<switch-name>` is *switchA*:

~~~~
$ ansible-playbook -i etc/playbooks/hosts -K etc/playbooks/dp-switchA.yml
~~~~

## Next steps
After reflecting your changes, please refer "configure-guide". You can choose two methods to configure.

- ansible: [configure-ansible.md](configure-ansible.md)
- NETCONF over SSH: [configure-netconf.md](configure-netconf.md)