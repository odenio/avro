![Logo](http://svg.wiersma.co.za/hamba/project?title=avro&tag=A%20fast%20Go%20avro%20codec)

[![Go Report Card](https://goreportcard.com/badge/github.com/odenio/avro)](https://goreportcard.com/report/github.com/odenio/avro)
[![Build Status](https://github.com/odenio/avro/actions/workflows/test.yml/badge.svg)](https://github.com/odenio/avro/actions)
[![Coverage Status](https://coveralls.io/repos/github/odenio/avro/badge.svg?branch=master)](https://coveralls.io/github/odenio/avro?branch=master)
[![GoDoc](https://godoc.org/github.com/odenio/avro?status.svg)](https://godoc.org/github.com/odenio/avro)
[![GitHub release](https://img.shields.io/github/release/odenio/avro.svg)](https://github.com/odenio/avro/releases)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/odenio/avro/master/LICENSE)

A fast Go avro codec

## Overview

Install with:

```shell
go get github.com/odenio/avro/v2
```

## Usage

```go
type SimpleRecord struct {
	A int64  `avro:"a"`
	B string `avro:"b"`
}

schema, err := avro.Parse(`{
    "type": "record",
    "name": "simple",
    "namespace": "org.hamba.avro",
    "fields" : [
        {"name": "a", "type": "long"},
        {"name": "b", "type": "string"}
    ]
}`)
if err != nil {
	log.Fatal(err)
}

in := SimpleRecord{A: 27, B: "foo"}

data, err := avro.Marshal(schema, in)
if err != nil {
	log.Fatal(err)
}

fmt.Println(data)
// Outputs: [54 6 102 111 111]

out := SimpleRecord{}
err = avro.Unmarshal(schema, data, &out)
if err != nil {
	log.Fatal(err)
}

fmt.Println(out)
// Outputs: {27 foo}
```

More examples in the [godoc](https://godoc.org/github.com/odenio/avro).

#### Types Conversions

| Avro                    | Go Struct                          | Go Interface              |
| ----------------------- | ---------------------------------- | ------------------------- |
| `null`                  | `nil`                              | `nil`                     |
| `boolean`               | `bool`                             | `bool`                    |
| `bytes`                 | `[]byte`                           | `[]byte`                  |
| `float`                 | `float32`                          | `float32`                 |
| `double`                | `float64`                          | `float64`                 |
| `long`                  | `int64`                            | `int64`                   |
| `int`                   | `int`, `int32`, `int16`, `int8`    | `int`                     |
| `string`                | `string`                           | `string`                  |
| `array`                 | `[]T`                              | `[]interface{}`           |
| `enum`                  | `string`                           | `string`                  |
| `fixed`                 | `[n]byte`                          | `[]byte`                  |
| `map`                   | `map[string]T{}`                   | `map[string]interface{}`  |
| `record`                | `struct`                           | `map[string]interface{}`  |
| `union`                 | *see below*                        | *see below*               |
| `int.date`              | `time.Time`                        | `time.Time`               |
| `int.time-millis`       | `time.Duration`                    | `time.Duration`           |
| `long.time-micros`      | `time.Duration`                    | `time.Duration`           |
| `long.timestamp-millis` | `time.Time`                        | `time.Time`               |
| `long.timestamp-micros` | `time.Time`                        | `time.Time`               |
| `bytes.decimal`         | `*big.Rat`                         | `*big.Rat`                |
| `fixed.decimal`         | `*big.Rat`                         | `*big.Rat`                |

##### Unions

The following union types are accepted: `map[string]interface{}`, `*T` and `interface{}`.

* **map[string]interface{}:** If the union value is `nil`, a `nil` map will be en/decoded. 
When a non-`nil` union value is encountered, a single key is en/decoded. The key is the avro
type name, or scheam full name in the case of a named schema (enum, fixed or record).
* ***T:** This is allowed in a "nullable" union. A nullable union is defined as a two schema union, 
with one of the types being `null` (ie. `["null", "string"]` or `["string", "null"]`), in this case 
a `*T` is allowed, with `T` matching the conversion table above.
* **interface{}:** An `interface` can be provided and the type or name resolved. Primitive types
are pre-registered, but named types, maps and slices will need to be registered with the `Register` function. In the 
case of arrays and maps the enclosed schema type or name is postfix to the type
with a `:` separator, e.g `"map:string"`. If any type cannot be resolved the map type above is used unless
`Config.UnionResolutionError` is set to `true` in which case an error is returned.

##### TextMarshaler and TextUnmarshaler

The interfaces `TextMarshaler` and `TextUnmarshaler` are supported for a `string` schema type. The object will
be tested first for implementation of these interfaces, in the case of a `string` schema, before trying regular
encoding and decoding. 

### Recursive Structs

At this moment recursive structs are not supported. It is planned for the future.

## Benchmark

Benchmark source code can be found at: [https://github.com/nrwiersma/avro-benchmarks](https://github.com/nrwiersma/avro-benchmarks)

```
BenchmarkGoAvroDecode-10       	  495176	      2413 ns/op	     418 B/op	      27 allocs/op
BenchmarkGoAvroEncode-10       	  420168	      2917 ns/op	     948 B/op	      63 allocs/op
BenchmarkGoGenAvroDecode-10    	  757150	      1552 ns/op	     728 B/op	      45 allocs/op
BenchmarkGoGenAvroEncode-10    	 1882940	       639.0 ns/op	     256 B/op	       3 allocs/op
BenchmarkHambaDecode-10        	 3138063	       383.0 ns/op	      64 B/op	       4 allocs/op
BenchmarkHambaEncode-10        	 4377513	       273.3 ns/op	     112 B/op	       1 allocs/op
BenchmarkLinkedinDecode-10     	 1000000	      1109 ns/op	    1688 B/op	      35 allocs/op
BenchmarkLinkedinEncode-10     	 2641016	       456.0 ns/op	     248 B/op	       5 allocs/op
```

Always benchmark with your own workload. The result depends heavily on the data input.

## EXPERIMENTAL - Go structs generation

Go structs can be generated for you from the schema. The types generated follow the same logic in [types conversions](#types-conversions)

Install the struct generator with:

```shell
go install github.com/odenio/avro/v2/cmd/avrogen@<version>
```

Example usage assuming there's a valid schema in `in.avsc`:

```shell
./app -pkg avro -o bla.go -tags json:snake,yaml:upper-camel in.avsc
```

Check the options and usage with `-h`:

```shell
avrogen -h
```

Or use it as a lib in internal commands, it's the `gen` package
