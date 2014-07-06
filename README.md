# Welcome to Bob!

[![Build Status](https://travis-ci.org/modcloth/docker-builder.svg?branch=master)](https://travis-ci.org/modcloth/docker-builder)
[![GoDoc](https://godoc.org/github.com/modcloth/docker-builder?status.png)](https://godoc.org/github.com/modcloth/docker-builder)

Bob, the `docker-builder`, builds, tags, and pushes multiple Docker images, all
from a friendly `Bobfile` config file.

## Other Topics

* [Motivation for Docker Builder](docs/why.md)
* [Docker Builder Server](docs/advanced-usage.md)

## Quick Start

Steps to quick start:

```bash
# install docker-builder
go get github.com/modcloth/docker-builder

# create a Bobfile
docker-builder init .

# build the aforementioned Bobfile
docker-builder build Bobfile
```

## More Reasonably-Paced Start
0. Install `docker-builder`
0. Explore the [Writing a Bobfile](docs/writing-a-bobfile.md) doc
0. Run`docker-builder --help`
0. Run `docker-builder help build`

### First, Install `docker-builder`

#### Easiest

```bash
go get github.com/modcloth/docker-builder
```

#### From pre-compiled binaries

```bash
# on Mac OS X
curl -sL https://github.com/modcloth/docker-builder/releases/download/v0.6.0/docker-builder-v0.6.0-darwin-amd64.tar.gz | \
  tar -xzf - -C /usr/local --strip-components=1

# on Linux, note: you may need sudo
curl -sL https://github.com/modcloth/docker-builder/releases/download/v0.6.0/docker-builder-v0.6.0-linux-amd64.tar.gz | \
  sudo tar -xzf - -C /usr/local --strip-components=1
```

These commands will place `docker-builder` at
`/usr/local/bin/docker-builder`, so to use `docker-builder`, make sure
and check that `/usr/local/bin` is in your `$PATH` or change the `-C`
option.

**NOTE:** You may see some junk output when running `tar -xzf`.  This
has something to do with the archives being build on Mac OSX.  The
output is harmless and safe to ignore.

#### From source

To build from source, run `make build`.  You may have to install some
things first, such as `go`.
