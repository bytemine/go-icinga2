# About

`icinga2` is a [Go](https://golang.org) client package for the [Icinga2](https://www.icinga.org/products/icinga-2/)-API. Currently only
[event streams](http://docs.icinga.org/icinga2/snapshot/doc/module/icinga2/chapter/icinga2-api#icinga2-api-event-streams) are supported.
Eventually other features of the Icinga2-API will be supported, PR are welcome of course!

# Design

This package provides functionality to open an event stream. The returned [io.Reader](https://golang.org/pkg/io/#Reader) can be used directly in conjunction with a
[json.Decoder](https://golang.org/pkg/encoding/json/#Decoder) decoding into a matching type from the [event](./event) package. This only works if you are requesting only one type of events.
If multiple types are requested, the returned io.Reader can be passed on to an `event.Mux` which multiplexes the different event types 
into several readers.

# Usage

## Package name

The package name is just `icinga2` while the import path is `github.com/go-icinga2` to prevent confusion. While this is a bit ugly, I think this is common enough to be ok.
In effect you'll use `import "github.com/bytemine/go-icinga2"` and after that `icinga2.NewClient()` etc.

## Example

A usage example can be found in `client_test.go`, showing how multiple event streams can be consumed.
