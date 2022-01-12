# Netgear Switch Discovery Protocol (NSDP) 🔍

[![Go Report Card](https://goreportcard.com/badge/github.com/nicklasfrahm/nsdp)](https://goreportcard.com/report/github.com/nicklasfrahm/nsdp)

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
  nsdp [flags]
  nsdp [command]

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
