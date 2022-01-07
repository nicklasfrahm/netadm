# Netgear Switch Discovery Protocol (NSDP) 🔍

![https://goreportcard.com/report/github.com/nicklasfrahm/nsdp](https://goreportcard.com/badge/github.com/nicklasfrahm/nsdp)

A CLI and client library to manage Netgear Smart Switches written in [Go][website-go].

## CLI 🦾

Make sure to have [Go][website-go] installed and follow the instructions below to install it.

```bash
# go >= 1.17
# Using `go get` to install binaries is deprecated.
# The version suffix is mandatory.
go install github.com/nicklasfrahm/nsdp@latest

# go < 1.17
go get github.com/nicklasfrahm/nsdp
```

Below you may find the usage text of the command line interface.

```text
A command line interface to manage Netgear Smart Switches
via the UDP-based Netgear Switch Discovery Protocol (NSDP).

Usage:
  nsdp [flags]
  nsdp [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  if          List network interfaces
  scan        Scan for devices

Flags:
  -h, --help               display help for command
  -t, --timeout duration   timeout for commands (default 1s)

Use "nsdp [command] --help" for more information about a command.
```

## References 🔗

Below, you may find a list of useful references that were used for this implementation.

- [Netgear Switch Discovery Protocol on Wikipedia][wikipedia-ndsp]
- [`nmedb/nsdpy` on Github][github-nsdpy]

## License 📄

This project is licensed under the terms of the [MIT license](./LICENSE.md).

[wikipedia-ndsp]: https://en.wikipedia.org/wiki/Netgear_Switch_Discovery_Protocol
[github-nsdpy]: https://github.com/nmedb/nsdpy
[website-go]: https://go.dev
