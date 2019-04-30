# PASS tools

[![Build Status](https://travis-ci.com/OA-PASS/pass-tools.svg?branch=master)](https://travis-ci.com/OA-PASS/pass-tools)

Contains miscellaneous PASS CLI tools

## Usage

If you have go installed, you can simply install the `pass-tools` executable via

    go get github.com/oa-pass/pass-tools/cmd/pass-tools

 This will install the binary to your `${GOPATH/bin}`.  If you have that in your `$PATH`, this is particularly convenient for building and running cli commands.

Otherwise (e.g. for development) you can [build it](#building) from a local codebase

For help with commands, use

    pass-tools help

## Building

Building pass tools requires go 1.12 or later.

First, clone

    git clone https://github.com/OA-PASS/pass-tools.git

Then, you can build the executable (which will be placed at the root of the pass-tools directory) via

    go build ./cmd/pass-tools

Otherwise, you can install it to `${GOPATH/bin}` via

    go install ./cmd/pass-tools

## Testing

To run unit tests, do

    go test ./...

For integration tests, you need to have Fedora and Elasticsearch running.  Use the provided `docker-compose` file to do that

   docker-compose up -d

Then, run with integration tests

    go test -tags=integration ./...
