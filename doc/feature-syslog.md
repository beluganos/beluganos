# [Feature guide] syslog

This documents shows that where is the files of system log (syslog).

## Pre-requirements

- Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.

## The path of log files

The mainly log files are saved here:

|#   | Type                | Host                    | Path                                            |
|:---|:--------------------|:------------------------|:------------------------------------------------|
|1   | General             | Beluganos's host        | `/tmp/fibc.log` or `/var/log/syslog`            |
|2   | Routing engine      | Beluganos's host        | The files under `/var/log/beluganos/<LXC-name>/`|
|3   | OpenNSL agent       | OpenNetworkLinux's host | `/var/log/gonsld.log`                           |

**Notice:** The files under `/var/log/beluganos/<LXC-name>/` at Beluganos's host are the copied files by Beluganos. The original file is located at under `/var/log/` at LXC. Please note that you DO NOT edit directly these files at Beluganos's host.

## Change log levels

### FIBC

`fibc.log.conf` at Beluganos's host is the configuration file of FIBC. FIBC is a one of main module of Beluganos.

~~~
Beluganos$ vi /etc/beluganos/fibc.log.conf

[logger_root]
level=ERROR

[handler_console]
level=ERROR

[handler_file]
level=ERROR
~~~

To apply,

~~~
Beluganos$ sudo systemctl restart fibcd
~~~

### RIBX

`ribxd.conf` at LXC is the configuration file of RIBX. RIBX is a group of daemon at LXC.

~~~
LXC$ vi /etc/beluganos/ribxd.conf

[log]
level = 0   # 0: ERROR, 5: DEBUG
dump  = 0
~~~

To apply,

~~~
Beluganos$ lxc restart <container-name>
~~~

### OpenNSL Agent

~~~
OpenNetworkLinux$ vi /etc/beluganos/gonsld.conf

# DEBUG="-v"
~~~

To apply,

~~~
OpenNetworkLinux$ /etc/init.d/gonsld stop
OpenNetworkLinux% /etc/init.d/gonsld start
~~~

