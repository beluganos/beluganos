# Case 2: BGP/MPLS IP-VPN

The one of advantage using Beluganos is MPLS-VPN. In this case, the sample configurations of providor edge (PE) routers in BGP/MPLS based IP-VPN environments.

## Pre-requirements

- Ubuntu 17.10 server
	- 10GB+ strage
	- At least two NICs
		- One of this should be connected with white-box switch. Please set IP address before following setup procedure.
		- Another one is used for your login via SSH.
- White-box switch

## Network environments

This sample case will create only "..SAMPLE-VPN.." zone described following. The other routers like P1 and CE should be prepared in advance.

~~~~

             +--------------+                     +----------------+
             |      P1      |---------+----\\-----|      PE2       |
             | Lo: 10.0.0.1 |         |           +----------------+
             +--------------+         |           +----------------+
                   | .1               +----\\-----|      RR1       |
          LDP      |                              | Lo: 10.0.0.254 |
          OSPFv2   |                              | AS: 65001      |
             172.16.2.0/30                        +----------------+
   +--      [sample-mic.1]
   :               |
   :               | .2
   :             (eth1)
   :     +---------------------------------------------------------+
   :     |  +--------------+                        PE1            |
   :     |  |  sample-mic  |                    (Beluganos)        |
   S     |  | Lo: 10.0.1.1 |                     AS: 65001         |
   A     |  +--------------+                                       |
   M     |       (eth0)                                            |
   P     |         | .1                                            |
   L     |  192.169.0.0/24 [lxdbr0]                                |
   E     |         +-----------------------------+                 |
   |     |         | .X(DHCP)                    | .Y(DHCP)        |
   V     |      (eth0)                        (eth0)               |
   P     |  +---------------+             +---------------+        |
   N     |  | sample-ric10  |             | sample-ric11  |        |
   :     |  | Lo:  10.0.0.5 |             | Lo:  10.0.0.5 |        |
   :     |  | VRF: 10       |             | VRF: 11       |        |
   :     |  | RT:  1:10     |             | RT:  1:11     |        |
   :     |  | RD:  1:2010   |             | RD:  1:2011   |        |
   :     |  +---------------+             +---------------+        |
   :     |       (eth1)                 (eth1)         (eth2)      |
   :     +---------------------------------------------------------+
   :               | .1                   | .1           | .1
   :               |                      |              |
   +--     [sample-ric10.1]      [sample-ric11.1]  [sample-ric11.2]
            192.168.1.0/24        192.168.2.0/24    192.168.3.0/24
           eBGP    |             eBGP     |       eBGP   |
                   | .2                   | .2           | .2
            +--------------+    +--------------+    +--------------+
            |     CE1      |    |     CE2      |    |     CE3      |
            | Lo: 20.0.0.1 |    | Lo: 20.0.0.2 |    | Lo: 20.0.0.3 |
            | AS: 10       |    | AS: 20       |    | AS; 20       |
            | VRF:10       |    | VRF:11       |    | VRF:11       |
            +--------------+    +--------------+    +--------------+
~~~~


## Step 1. setup

### Step 1-1. Install Beluganos

Please check `doc/install-guide.md` for install.

### Step 1-2. Settings for switches

~~~~
$ cd ~/beluganos/etc/playbooks
$ vi roles/dpath/files/whitebox1/fibc.yml

datapaths:
  - name: whitebox1         # dpname (A-Z,a-z,0-9,_-)
    dp_id: 14               # datapath id of your switches (integer)
    mode: ofdpa2            # "ofdpa2" or "generic" or "ovs"

$ ansible-playbook -i hosts -K dp-whitebox1.yml
~~~~

This settings should be modified as your white-box switches. For more detail, see `doc/setup-guide.md`.

### Step 1-3. Settings for containers

~~~~
$ cd ~/beluganos/etc/playbooks
$ ansible-playbook -i hosts -K lxd-sample-vpn.yml
$ lxc stop sample-mic sample-ric10 sample-ric11
~~~~

## Step 2. start this case

Please execute following commands. You should have two terminals because `beluganos.py run` command will take the standard input from you.

~~~~
$ beluganos run
$ beluganos add sample
$ beluganos add sample-ric10
$ beluganos add sample-ric11
~~~~

## Step 3. confirm this case

You can login Belunogas's routing engine by following commands. Note that you have three containers which is `sample-mic`, `sample-ric10`, and `sample-ric11`.  The container which has "ric" in the name means VRFs, and the "sample-mic" is the master instance.

~~~~
$ beluganos con sample-mic
sample-mic> vtysh
~~~~