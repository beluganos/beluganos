# Case 1: LDP with OSPFv2

This case is suitable for a beginner of Beluganos, because not only Beluganos's sample configurations but also network environments are prepared. In this case, Beluganos will work label switch router (LSR) with LDP and OSPFv2.

## Pre-requirements

In this case, two servers are needed. The virtual machine (VM) is acceptable. **IP reachability is needed** between server 1 and 2.

- Server 1 (for Beluganos)
	- Ubuntu 16.04 server
	- 10GB+ storage
	- At least two NICs
		- One should be connected with server 2. The other is used for your login via SSH. 
		- Please set IP addresses before following setup procedure.
- Server 2 (for OVS and other routers)
	- Ubuntu 16.04 server
	- 12GB+ storage
	- At least one NICs
		- This should be connected with server 1.
		- Please set IP address before following setup procedure.

~~~~
                          |
                        (ens3)
                 +-------------------+
                 |     server 1      |
                 |     Beluganos     |
                 +-------------------+
                        (ens4)
                          | .53
                    172.16.0.0/24
                          | .55
                        (ens3)
                 +-------------------+
                 |     server 2      |
                 |    OpenvSwitch    |
                 |         &         |
                 |    environment    |
                 +-------------------+
~~~~

## Network environment

By using scripts described following capter, you can create following environments at server 2 automatically. The router of *sample* is OpenvSwitch which will be connected with server 1 by OpenFlow. The routers of P1, P2, P3, and P4 are Linux containers installed FRRouting.

~~~~
                             |
                       172.16.1.0/24 [lxdbr1]
                             | .2
                           (eth1)
                      +--------------+
                      |      P1      |
                      | <sample-p1>  |
                      | Lo: 10.0.0.1 |
                      +--------------+
                           (eth2)
                             | .1
                       172.16.2.0/24 [lxdbr2]
                             | .2
                           (eth1)
                      +--------------+
                      |      P2      |
                      | <sample-p2>  |
                      | Lo: 10.0.0.2 |
                      +--------------+
                           (eth2)
                             | .1
                             |
                       [sample-p2.1] ---+
                       172.16.3.0/24
                        [sample.1] -----+
                             |
                             | .2
                           (eth1)
                     +--------------+
                     |  OpenvSwitch | <------- Beluganos
                     |   <sample>   |
            +--(eth2)| Lo: 10.0.0.5 |(eth3)--+
            |    .1  +--------------+   .1   |
            |            (eth4.10)           |
            |                | .1            |
            |                |               |
   +---- [sample.2]       [sample.4]      [sample.3] ----+
       172.16.4.0/24    172.16.6.0/24   172.16.5.0/24
   +-- [sample-p3.1]                    [sample-p4.1] ---+
            |                                |
            | .2                             | .2
          (eth1)                           (eth1)
     +--------------+                 +--------------+
     |      P3      |                 |      P4      |
     | <sample-p3>  |                 | <sample-p4>  |
     | Lo: 10.0.0.3 |                 | Lo: 10.0.0.4 |
     +--------------+                 +--------------+
          (eth2)                           (eth2)
            | .1                             | .1
      172.16.7.0/24 [lxdbr7]           172.16.8.0/24 [lxdbr8]
            |                                |
~~~~

## Step 1. setup server 1

### Step 1-1. Install Beluganos

The server 1 will be used as Beluganos. Please check `doc/install-guide.md` for install. Please note that `FFLOW_OFC_IFACE` should be set to the interface which is connected with server 2. If you need to change this interface name or IP address, please change `create.ini` before execute `create.sh`.

~~~~
server1$ vi create.ini
  FFLOW_MNG_IFACE=ens3             # Set your management interface name for remote login
  FFLOW_OFC_IFACE=ens4             # Set your secure channel interface name connected to switches
  FFLOW_OFC_ADDR=172.16.0.55       # (Optional) You can change FFLOW_OFC_IFACE's IP address if needed
  FFLOW_OFC_MASK=255.255.255.0     # (Optional) You can change FFLOW_OFC_IFACE's subnet mask if needed
~~~~

### Step 1-2. Settings for switches

~~~~
server1$ cd ~/beluganos/
server1$ . ./setenv.sh
server1$ cd etc/playbooks
server1$ ansible-playbook -i hosts -K dp-sample.yml
~~~~


### Step 1-3. Settings for containers

~~~~
server1$ ansible-playbook -i hosts -K lxd-sample.yml
server1$ lxc stop sample
~~~~


## Step 2. setup server 2

Following network will be created by following introduction. To quickly try Beluganos, **instead of white-box switches, OpenvSwitch is used for *sample*'s dataplane**.

~~~~

 +-----------server 2------------+
 |                               |
 |                  P3--------   |
 |                  |            |           +--server 1--+
 |  ---P1---P2---OpenvSwitch <===|==========>| Beluganos  |
 |                  |            | OpenFlow  +------------+
 |                  P4--------   |
 |                               |
 +-------------------------------+
~~~~

### Step 2-1. Transfer required file

Once you execute `create.sh` at server 1, some files which required in server 2 are created at home directory. You may copy these files to server 2's home directory. The required files are `frr_3.0_amd64.deb`, `ubuntu-16.04-server-cloudimg-amd64-*` and the files under `beluganos/`.

~~~~
server2$ sftp <IP-address-of-server-1>
sftp> ls
beluganos
frr
frr-dbg_3.0_amd64.deb
frr-doc_3.0_all.deb
frr_3.0_amd64.deb
go
mypython
ubuntu-16.04-server-cloudimg-amd64-lxd.tar.xz
ubuntu-16.04-server-cloudimg-amd64-root.tar.xz

sftp> put frr_3.0_amd54.deb
sftp> put ubuntu-16.04-server-cloudimg-amd64-*
sftp> put -R beluganos
sftp> exit

server2$ ls
beluganos/ frr-doc_3.0_all.deb ubuntu-16.04-server-cloudimg-amd64-lxd.tar.xz ubuntu-16.04-server-cloudimg-amd64-root.tar.xz
~~~~

### Step 2-2. pre-install required package

~~~~
server2$ cd ~/beluganos/
server2$ ./create.sh min
~~~~

### Step 2-3. setup P1 to P4

~~~~
server2$ cd ~/beluganos/
server2$ . ./setenv.sh
server2$ cd etc/playbooks
server2$ ansible-playbook -i hosts -K sample-net.yml
server2$ lxc stop sample-p1 sample-p2 sample-p3 sample-p4
~~~~

## Step 3. start this case

### Step 3-1. start environment

~~~~
server2$ sudo service openvswitch-switch start
server2$ lxc start sample-p1 sample-p2 sample-p3 sample-p4
~~~~

### Step 3-2. start Beluganos

If `(mypython)` didn't set your terminal at server 1, please execute following command to set your environments:

~~~~
server1$ cd ~/beluganos/
server1$ . ./setenv.sh
~~~~

After that, please execute following commands. You should have two terminals of server 1 because `beluganos.py run` command will take the standard input from you.

~~~~
server1$ beluganos.py run
server1$ beluganos.py add sample
~~~~

## Step 4. confirm this case

You can login Belunogas's routing engine by following commands:

~~~~
server1$ cd ~/beluganos
server1$ . ./setenv.sh
server1$ beluganos.py con sample
root@sample:~# vtysh

Hello, this is FRRouting (version 3.0-rc2).
Copyright 1996-2005 Kunihiro Ishiguro, et al.

sample# show ip ospf neighbor

Neighbor ID     Pri State           Dead Time Address         Interface            RXmtL RqstL DBsmL
10.0.0.2          1 Full/DR           36.319s 172.16.3.1      eth1:172.16.3.2          0     0     0      <---- Success to connect with P2
10.0.0.3          1 Full/DR           36.318s 172.16.4.2      eth2:172.16.4.1          0     0     0      <---- Success to connect with P3
10.0.0.4          1 Full/DR           35.942s 172.16.5.2      eth3:172.16.5.1          0     0     0      <---- Success to connect with P4

sample# show mpls ldp binding

AF   Destination          Nexthop         Local Label Remote Label  In Use
ipv4 10.0.0.1/32          10.0.0.2        16          16               yes                  <---- Success to learn about MPLS label for P1
ipv4 10.0.0.1/32          10.0.0.3        16          16                no
ipv4 10.0.0.1/32          10.0.0.4        16          16                no
ipv4 10.0.0.2/32          10.0.0.2        17          imp-null         yes
ipv4 10.0.0.2/32          10.0.0.3        17          17                no
ipv4 10.0.0.2/32          10.0.0.4        17          17                no
ipv4 10.0.0.3/32          10.0.0.2        20          18                no
ipv4 10.0.0.3/32          10.0.0.3        20          imp-null         yes
ipv4 10.0.0.3/32          10.0.0.4        20          18                no
ipv4 10.0.0.4/32          10.0.0.2        22          24                no
ipv4 10.0.0.4/32          10.0.0.3        22          24                no
ipv4 10.0.0.4/32          10.0.0.4        22          imp-null         yes
ipv4 10.0.0.5/32          10.0.0.2        imp-null    19                no
ipv4 10.0.0.5/32          10.0.0.3        imp-null    18                no
ipv4 10.0.0.5/32          10.0.0.4        imp-null    19                no
ipv4 172.16.1.0/24        10.0.0.2        18          17               yes
ipv4 172.16.1.0/24        10.0.0.3        18          19                no
ipv4 172.16.1.0/24        10.0.0.4        18          20                no
ipv4 172.16.2.0/24        10.0.0.2        19          imp-null         yes
ipv4 172.16.2.0/24        10.0.0.3        19          20                no
...

~~~~

You can also confirm by other router's console by following commands:

~~~~
server2$ lxc exec sample-p3 -- /bin/bash
root@sample-p3:~# vtysh

Hello, this is FRRouting (version 3.0-rc2).
Copyright 1996-2005 Kunihiro Ishiguro, et al.

sample-p3# show ip route
Codes: K - kernel route, C - connected, S - static, R - RIP,
       O - OSPF, I - IS-IS, B - BGP, P - PIM, N - NHRP, T - Table,
       v - VNC, V - VNC-Direct,
       > - selected route, * - FIB route

K>* 0.0.0.0/0 via 192.169.1.1, eth0
O>* 10.0.0.1/32 [110/30] via 172.16.4.1, eth1, 02:08:49
O>* 10.0.0.2/32 [110/20] via 172.16.4.1, eth1, 02:08:49
O   10.0.0.3/32 [110/0] is directly connected, lo, 02:38:20
C>* 10.0.0.3/32 is directly connected, lo
O>* 10.0.0.4/32 [110/20] via 172.16.4.1, eth1, 02:08:20
O>* 10.0.0.5/32 [110/10] via 172.16.4.1, eth1, 02:08:49
O>* 172.16.1.0/24 [110/40] via 172.16.4.1, eth1, 02:08:49
O>* 172.16.2.0/24 [110/30] via 172.16.4.1, eth1, 02:08:49
...

sample-p3# ping 10.0.0.1
PING 10.0.0.1 (10.0.0.1) 56(84) bytes of data.
64 bytes from 10.0.0.1: icmp_seq=1 ttl=62 time=0.137 ms
64 bytes from 10.0.0.1: icmp_seq=2 ttl=62 time=0.081 ms
64 bytes from 10.0.0.1: icmp_seq=3 ttl=62 time=0.092 ms
64 bytes from 10.0.0.1: icmp_seq=4 ttl=62 time=0.092 ms
64 bytes from 10.0.0.1: icmp_seq=5 ttl=62 time=0.145 ms
^C
--- 10.0.0.1 ping statistics ---
5 packets transmitted, 5 received, 0% packet loss, time 3999ms
rtt min/avg/max/mdev = 0.081/0.109/0.145/0.027 ms
~~~~

## Step 5. stop this case

If there is any trouble unfortunately, please stop and restart this case.

### Step 5-1. stop Beluganos

~~~~
server1$ cd ~/beluganos/
server1$ . ./setenv.sh
server1$ beluganos.py del sample
~~~~

Note that the script `beluganos.py run` can be stopped by `Ctrl-c`. You should stop this script after executing `beluganos.py del sample`.

### Step 5-2. stop environments

~~~~
server2$ sudo service openvswitch-switch stop
server2$ lxc stop sample-p1 sample-p2 sample-p3 sample-p4
~~~~

## For advanced settings

### Use real white-box switches

In this sample, we use OpenvSwitch instead of white-box switches. This is just for simplicity, so that you can use real white-box switches.

In this case, you don't need prepare for server 2. Only server 1 is needed. Moreover, please change `fibc.yml` by following commands:

~~~~
server1$ cd ~/beluganos/
server1$ vi etc/playbooks/roles/dpath/files/sample/fibc.yml
 datapaths:
   - name: sample_sw
     dp_id: <dp-id>
     mode: <switch-type>
~~~~

In order to apply this change, please execute playbook of `dp-sample.yml` after stopping Beluganos.

~~~~
server1$ cd ~/beluganos/
server1$ . ./setenv.sh
server1$ cd etc/playbooks
server1$ ansible-playbook -i hosts -K dp-sample.yml
~~~~

### Change Beluganos's settings

You can edit Belugnos's settings if you want. For example, if you want to change interface address, please edit following file.

~~~~
server1$ cd ~/beluganos
server1$ vi etc/playbooks/roles/lxd/files/sample/frr.conf
 interface eth1
  ip address 172.16.3.2/24
 !
 interface eth2
  ip address 172.16.4.1/24
 !
 interface eth3
  ip address 172.16.5.1/24
 !
 interface eth4.10
  ip address 172.16.6.1/24
 !
 router ospf
  network 10.0.0.5/32 area 0.0.0.0
  network 172.16.3.0/24 area 0.0.0.0
  network 172.16.4.0/24 area 0.0.0.0
  network 172.16.5.0/24 area 0.0.0.0
  network 172.16.6.0/24 area 0.0.0.0
 !
~~~~

If you change the files under `roles/lxd/files/sample/`, you may execute playbook of `lxd-sample.yml` after stopping Beluganos.

~~~~
server1$ ansible-playbook -i hosts -K lxd-sample.yml
server1$ lxc stop sample
~~~~

For more detail, see `doc/setup-guide.md`.

### Start Beluganos as daemon

When you execute `beluganos.py run`, you will be token the standard input by this script. This mode is useful for debugging, but some trouble will occur in production. Of course you can use daemon mode of Beluganos. For using, you should edit `.service` file. For more detail, see `doc/operation-guide.md`.
