# goteleinfo: monitor your house electrical current

[![Build Status](https://travis-ci.org/j-vizcaino/goteleinfo.svg)](https://travis-ci.org/j-vizcaino/goteleinfo)

`goteleinfo` is a small library capable of decoding Teleinfo frames sent by an electrical meter.

Teleinfo is the protocol used by ERDF (*Electricité Réseau Distribution France*) electrical meters. Using this
protocol, the electrical meter sends status frames periodically (~2 per second) that gives
information about your electrical current consumption and status.

The Teleinfo protocol is described [in this document](https://www.enedis.fr/sites/default/files/Enedis-NOI-CPT_02E.pdf)

In order to convert read the frames, the [Micro Teleinfo](https://www.tindie.com/products/Hallard/micro-teleinfo-v11/) or [PiTinfo](https://www.tindie.com/products/Hallard/pitinfo/) is needed.

## JSON frames

A simple CLI is provided in `cmd/teleinfo-info` that reads frames and prints them, JSON encoded.
This program can serve as an example of what the library provides.
