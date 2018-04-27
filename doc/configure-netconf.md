# Configure by NETCONF

## Pre-requirements
- Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.
	- In setup, you may specify switch name and dp_id. In this documents, the sample file's value (`name: sample_sw`, `dp_id: 153`) is assumed. If you changed this value, please change to match it.


## Prepare before starting

Before issuing NETCONF message, you need following steps.

```
$ cd etc/playbooks
$ vi lxd-netconf.yml
---

- hosts: hosts
  connection: local
  become: true
  roles:
    - { role: lxd, mode: netconf }
  vars:
    port_num: 5
    re_id: 10.0.1.6
    datapath: sample_sw
    dp_id: 153
    bridge: dp0

$ ansible-playbook -i hosts -K lxd-netconf.yml
```

The details of settings value in `lxd-netconf.yml` is following.

- `port_num`
	- The maximum physical interface number of your router.
- `re_id`
	- Router identified name. Only Beluganos's main component will use this value to identify routers. This value is not used for router settings, like loopback address or router-id.
- `datapath` and `dp_id`
	- White-box hardware settings. The value of `fibc.yml` should be filled. Please refer setup-guide.md for more details.
- `bridge`
	- Do not changed. Only for debugs.

## Configure after starting

After starting Beluganos by [operation-guide.md](operation-guide.md), you can use NETCONF commands. The yang file of Beluganos is published under [github](https://github.com/beluganos/netconf/etc/openconfig).

To try quickly, sample NETCONF operations are available at [https://github.com/beluganos/netconf/](https://github.com/beluganos/netconf/). 
