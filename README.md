# libnf-go

Go bindings for the [`libnf`](https://github.com/netx-as/libnf) C library — enabling parsing of NetFlow data in Go applications.

## Features

- **Go Bindings for libnf**: Provides Go language bindings to the `libnf` C library.
- **NetFlow Data Parsing**: Enables parsing and processing of NetFlow data within Go programs.

## Prerequisites

Before using this Go module, ensure that the following are installed on your Linux-based operating system:

- [`libnf`](https://github.com/netx-as/libnf)
- [`SWIG`](http://www.swig.org/)

## Installation

To install the package, run:

```bash
go install github.com/matejnesuta/libnf-go@latest
go get github.com/matejnesuta/libnf-go
```

## Demo examples

In order to run these, clone the repository and run:
```bash
go run main.go <example-name> [--profile]
```

Results of the profiler can be viewed like this: 
```bash
go tool pprof -http=:6061 .prof/cpu.prof
```

This package uses a Godoc for documentation. In order to generate it, run this command: 
```bash
godoc -http=:6060
```

This repository contains various automated tests. To run them from the root, use this: 
```bash
go test ./...
```
