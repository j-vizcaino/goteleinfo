# goteleinfo: monitor your house electrical current

[![Build Status](https://travis-ci.org/j-vizcaino/goteleinfo.svg)](https://travis-ci.org/j-vizcaino/goteleinfo)

`goteleinfo` is a daemon capable of decoding Teleinfo frames sent by an electrical meter,
and inserting them into a database.

Teleinfo is the protocol used by ERDF (*Electricité Réseau Distribution France*) electrical meters. Using this
protocol, the electrical meter sends status frames periodically (~2 per second) that gives
information about your electrical current consumption and status.

The Teleinfo protocol is described [in this document](http://www.enedis.fr/sites/default/files/ERDF-NOI-CPT_02E.pdf)

In order to convert read the frames, the [Micro Teleinfo](https://www.tindie.com/products/Hallard/micro-teleinfo-v11/) or [PiTinfo](https://www.tindie.com/products/Hallard/pitinfo/) is needed.


## Features

`goteleinfo` supports the following exporters:

* `hphc.json`: extract *Heures pleines*, *Heures creuses* counters, instant power consumption and print as JSON on stdout
* `hphc.influx`: extract *Heures pleines*, *Heures creuses* counters, instant power consumption and insert into an InfluxDB database

NOTE: an InfluxDB database schema for the `hphc.influx` exporter is available in the `deploy` subdirectory. Use `new-database.sh` to generate the schema with a custom database name.


## Installation

### Dependencies

In order to build `goteleinfo` you need:

* a working Go (>= 1.6) installation
* `make`
* [Glide](https://glide.sh) (>= 0.11)


### Building

```shell
$ git clone https://github.com/j-vizcaino/goteleinfo $GOPATH/src/goteleinfo
$ cd $GOPATH/src/goteleinfo
$ make
```

## Development

### Adding a new exporter

In order to add a new exporter:

* the object must implement the `Exporter` interface defined in `teleinfo/exporter.go`
* the module must register using `exporters.Register()`

Once you're done, write some tests and feel free to open a PR.

