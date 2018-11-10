# goteleinfo: monitor your house electrical current

[![Build badge]][Build] [![GoDoc badge]][GoDoc] [![GoReport badge]][GoReport]

[Build badge]: https://travis-ci.org/j-vizcaino/goteleinfo.svg
[Build]: https://travis-ci.org/j-vizcaino/goteleinfo
[GoDoc badge]: https://godoc.org/github.com/j-vizcaino/goteleinfo?status.png
[GoDoc]: https://godoc.org/github.com/j-vizcaino/goteleinfo
[GoReport badge]: https://goreportcard.com/badge/github.com/j-vizcaino/goteleinfo
[GoReport]: https://goreportcard.com/report/github.com/j-vizcaino/goteleinfo

`goteleinfo` is a small library capable of decoding Teleinfo frames sent by an electrical meter.

Teleinfo is the protocol used by ERDF (*Electricité Réseau Distribution France*) electrical meters. Using this
protocol, the electrical meter sends status frames periodically (~2 per second) that gives
information about your electrical current consumption and status.

The Teleinfo protocol is described [in this document](https://www.enedis.fr/sites/default/files/Enedis-NOI-CPT_02E.pdf)

In order to convert read the frames, the [Micro Teleinfo](https://www.tindie.com/products/Hallard/micro-teleinfo-v11/) or [PiTinfo](https://www.tindie.com/products/Hallard/pitinfo/) is needed.

## Metrics

The library provides counters to track frame reading and decoding errors. Those are stored internally as Prometheus metrics.

* `teleinfo_frames_read_total`: total number of raw Teleinfo frames read successfully.
* `teleinfo_frames_read_errors_total`: total number of read errors. This metric is tagged by `error_type`.
* `teleinfo_frames_decoded_total`: total number of decoded Teleinfo frames.
* `teleinfo_frames_decode_errors_total`: total number of frame decoding errors. This is metric is tagged by `error_type`

### Difference between reading and decoding

Bytes coming from the serial port are accumulated by the reader, looking for _frame start_ and _frame end_ markers to extract raw Teleinfo frames.
This is called "reading".

Decoding refers to verifying checksum and extracting fields from frame (eg. `OPTARIF`, `HPHC` fields, etc...).

## Example: serving Teleinfo as JSON frames

A simple HTTP server example is provided in [`cmd/teleinfo-info`](https://github.com/j-vizcaino/goteleinfo/blob/master/cmd/teleinfo-json/main.go) that provides the last *n* frames, JSON encoded.
Frames are served in `/frames`, while Prometheus metrics can be read from `/metrics` URL.

