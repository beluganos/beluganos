# [Feature guide] SNMP

SNMP (Simple Network Management Protocol) is a management protorol of network elements. Beluganos supports SNMP MIB to check interface status or traffic counter.

## Pre-requirements

- Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.

### Install SNMP functions

Before proceeding, please install SNMP functions by following steps:

```
$ cd ~/beluganos
$ make install-stats
$ sudo systemctl start snmpd
$ sudo systemctl start fibsd
```
## Overviews

### Start Beluganos

Please refer [operation-guide.md](operation-guide.md).

### Get SNMP

```
$ snmpget -v 2c -c public localhost <oid>
$ snmpwalk -v 2c -c public localhost <oid>
```

## Settings

TBD

## Supported features

TBD

|  OID                         |  Name.          | OpenFlow(PortStats) |
|------------------------------|-----------------|---------------------|
| .1.3.6.1.4.99999.31.1.1.1.1  | ifName          | ifname on container |
| .1.3.6.1.4.99999.31.1.1.1.2  | ifInOctets      | rx\_bytes           |
| .1.3.6.1.4.99999.31.1.1.1.3  | ifInUcastPkts   | rx\_packets         |
| .1.3.6.1.4.99999.31.1.1.1.4  | ifInNUcastPkts  | (always 0)          |
| .1.3.6.1.4.99999.31.1.1.1.5  | ifInDiscards    | rx\_dropped         |
| .1.3.6.1.4.99999.31.1.1.1.6  | ifInErrors      | rx\_errors          |
| .1.3.6.1.4.99999.31.1.1.1.7  | ifOutOctets     | tx\_bytes           |
| .1.3.6.1.4.99999.31.1.1.1.8  | ifOutUcastPkts  | tx\_packets         |
| .1.3.6.1.4.99999.31.1.1.1.9  | ifOutNUcastPkts | (always 0)          |
| .1.3.6.1.4.99999.31.1.1.1.10 | ifOutDiscards   | tx\_dropped         |
| .1.3.6.1.4.99999.31.1.1.1.11 | ifOutErrors     | tx\_errors          |