# pci-4g

[![Build Status](https://travis-ci.org/xuybin/pci-4g.svg?branch=master)](https://travis-ci.org/xuybin/pci-4g)
PCI planning.

According to MR traffic, signal to noise ratio, working days and non-working days and other models, using genetic algorithm iterative planning, planning the whole network or local PCI.

in development now.

documents will be wrote later.

## DOWNLOAD

You could download the latest build binaries from [release page](https://github.com/xuybin/pci-4g/releases) !

## RUN

You could use **cli option** or **environment varibles** to config your pci-4g

```bash
./pci-4g --help
Options:

  -h, --help                                display help information
  -l, --*listen[=$PCI_LS]              *gateway listen host and port

```

* -l --listen **PCI_LS**, gateway listen addr, format is *host:port*, example: *0.0.0.0:1329*