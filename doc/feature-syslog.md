# [Feature guide] syslog

This documents shows that where is the files of system log (syslog).

## Pre-requirements

- Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.

## Log files

The mainly log files are saved here:

| Type                | Server                  | Path                                       |
|:--------------------|:------------------------|:-------------------------------------------|
| General             | Beluganos's host        | `/tmp/fibc.log` or under `/var/log/syslog` |
| About routing       | Beluganos's host        | Under `/var/log/beluganos/<LXC-name>/`     |
| About OpenNSL agent | OpenNetworkLinux's host | `/var/log/gonsld.log`                      |

## Note

The files under `/var/log/beluganos/<LXC-name>/` at Beluganos's host are the copied files by Beluganos. The original file is located at under `/var/log/`. Please note that you DO NOT edit the files under `/var/log/beluganos/<LXC-name>/`.