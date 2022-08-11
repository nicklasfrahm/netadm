# netadm 🔍

[![Go Report Card](https://goreportcard.com/badge/github.com/nicklasfrahm/netadm?style=flat-square)](https://goreportcard.com/report/github.com/nicklasfrahm/netadm)
[![Release](https://img.shields.io/github/release/nicklasfrahm/netadm.svg?style=flat-square)](https://github.com/nicklasfrahm/netadm/releases/latest)
[![Go Reference](https://img.shields.io/badge/Go-reference-informational.svg?style=flat-square)](https://pkg.go.dev/github.com/nicklasfrahm/netadm)

A CLI and client library to manage network devices written in [Go][website-go].

## CLI 🦾

Make sure to have [Go][website-go] installed and follow the instructions below to install it.

```bash
# go >= 1.17
# Using `go get` to install binaries is deprecated.
# The version suffix is mandatory.
go install github.com/nicklasfrahm/netadm@latest

# go < 1.17
go get github.com/nicklasfrahm/netadm
```

Below you may find the usage text of the command line interface.

```text
A command line interface to manage network devices
via different protocols.

Note:
  To achieve a consistent behavior all operations
  are executed twice and the results are merged.
  This is done to work around operations that do
  not succeed if the device needs to refresh its
  ARP cache by performing a MAC address lookup of
  the host via the host IP. This happens on the
  the first interaction or, I assume, when the
  cache expires naturally, which appears to be
  every 5 minutes or so.

Usage:
  netadm [flags]
  netadm [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  get         Read configuration keys
  help        Help about any command
  if          List network interfaces
  keys        List available configuration keys
  scan        Scan for devices
  set         Write configuration keys

Flags:
  -h, --help               display help for command
  -r, --retries uint       number of retries to perform (default 1)
  -t, --timeout duration   timeout per attempt (default 100ms)

Use "netadm [command] --help" for more information about a command.
```

## Configuration Items 🔧

| ID     | NAME                 | EXAMPLE                           |
| ------ | -------------------- | --------------------------------- |
| 0x0001 | model                | GS308E                            |
| 0x0003 | name                 | switch-0                          |
| 0x0004 | mac                  | 33:0b:c9:5e:51:3a                 |
| 0x0006 | ip                   | 192.168.0.253                     |
| 0x0007 | netmask              | 255.255.255.0                     |
| 0x0008 | gateway              | 192.168.0.254                     |
| 0x000A | password             | password                          |
| 0x000B | dhcp                 | false                             |
| 0x000D | firmware             | 1.00.10                           |
| 0x0014 | passwordencryption   | Hash64                            |
| 0x0017 | passwordnonce        | [1 2 3 4]                         |
| 0x001A | passwordhash         | [1 2 3 4]                         |
| 0x0C00 | portspeeds           | [1:1G 2:Down]                     |
| 0x1000 | portmetrics          | [1:64/32/0]                       |
| 0x1C00 | cabletestresult      | [0 0 0 0 0 119 30 183 118]        |
| 0x2000 | vlanengine           | Disabled                          |
| 0x2400 | vlansport            | [1:1+2+3+4+5+6+7+8]               |
| 0x2800 | vlans802q            | [1t1+2u3+4+5+6+7+8]               |
| 0x3000 | pvids                | [1:2 2:2 3:1 4:1 5:1 6:1 7:1 8:1] |
| 0x3400 | qosengine            | DSCP                              |
| 0x3800 | qospolicies          | [1:Normal 2:High]                 |
| 0x4C00 | bandwidthlimitsin    | [1:256Mbps 2:None]                |
| 0x5000 | bandwidthlimitsout   | [1:256Mbps 2:None]                |
| 0x5400 | broadcastfilter      | false                             |
| 0x5800 | broadcastlimits      | [1:256Mbps 2:None]                |
| 0x5C00 | portmirroring        | 1:2+3                             |
| 0x6000 | portcount            | 5                                 |
| 0x6800 | igmpsnoopingvlan     | 1                                 |
| 0x6C00 | multicastfilter      | false                             |
| 0x7000 | igmpheadervalidation | false                             |
| 0x9000 | loopdetection        | false                             |

## References 🔗

Below, you may find a list of useful references that were used for this implementation.

- [Netgear Switch Discovery Protocol on Wikipedia][wikipedia-ndsp]
- [`nmedb/nsdpy` on Github][github-nsdpy]
- [`tabacha/ProSafeLinux` on Github][github-tabacha-prosafe]
- [`rhulme/ProSafeLinux` on Github][github-rhulme-prosafe]

## License 📄

This project is licensed under the terms of the [MIT license](./LICENSE.md).

[wikipedia-ndsp]: https://en.wikipedia.org/wiki/Netgear_Switch_Discovery_Protocol
[github-nsdpy]: https://github.com/nmedb/nsdpy
[github-tabacha-prosafe]: https://github.com/tabacha/ProSafeLinux
[github-rhulme-prosafe]: https://github.com/rhulme/ProSafeLinux
[website-go]: https://go.dev
