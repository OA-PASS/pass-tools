# PASS tools

[![Build Status](https://travis-ci.com/OA-PASS/pass-tools.svg?branch=master)](https://travis-ci.com/OA-PASS/pass-tools)

Contains miscellaneous PASS CLI tools

## Usage

Pre-compiled binaries are present in [the releases section](https://github.com/oa-pass/pass-tools/releases/).  It's possible to download and extract the binaries to your `PATH` as follows:

For Mac OS:

    $ base=https://github.com/oa-pass/pass-tools/releases/download/v0.0.1 &&
      curl -L $base/pass-tools-$(uname -s)-$(uname -m) >/usr/local/bin/pass-tools &&
      chmod +x /usr/local/bin/pass-tools

For Linux:

    $ base=https://github.com/oa-pass/pass-tools/releases/download/v0.0.1 &&
      curl -L $base/pass-tools-$(uname -s)-$(uname -m) >/tmp/pass-tools &&
      sudo install /tmp/pass-tools /usr/local/bin/pass-tools

For Windows, using Git Bash:

    $ base=https://github.com/oa-pass/pass-tools/releases/download/v0.0.1 &&
      mkdir -p "$HOME/bin" &&
      curl -L $base/pass-tools-Windows-x86_64.exe > "$HOME/bin/pass-tools.exe" &&
      chmod +x "$HOME/bin/pass-tools.exe"

If you have `go` installed and wish to build `pass-tools`, you can simply install the `pass-tools` executable via

    go get github.com/oa-pass/pass-tools/cmd/pass-tools

 This will install the binary to your `${GOPATH/bin}`.  If you have that in your `$PATH`, this is particularly convenient for building and running cli commands.

Otherwise (e.g. for development) you can [build it](#building) from a local codebase

For help with commands, use

    pass-tools -h

This will show sub-commands, flags, arguments, etc.

For help with any sub-command, use the `-h` flag.  For example

    pass-tools migrate -h
    pass-tools assign pi -h

Fedora and elasticsearch parameters are described in the `pass-tools -h` help:

    NAME:
       pass-utils - PASS utilities

    USAGE:
       pass-tools [global options] command [command options] [arguments...]

    VERSION:
       0.0.0

    COMMANDS:
         assign   Assign ownership of a PASS resource to a user
         migrate  Migrate PASS data from an old format/schema/context to a new one
         help, h  Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --fedora value, --pass.fedora.baseurl value               Fedora baseURL (default: "http://localhost:8080/fcrepo/rest/") [$PASS_FEDORA_BASEURL]
       --es value, --pass.elasticsearch.url value                Elasticsearch URL (default: "http://localhost:9200/pass/_search") [$PASS_ELASTICSEARCH_URL]
       --pass.fedora.user value, --username value, -u value      Username for basic auth to Fedora (default: "fedoraAdmin") [$PASS_FEDORA_USER]
       --pass.fedora.password value, --password value, -p value  Password for basic auth to Fedora (default: "moo") [$PASS_FEDORA_PASSWORD]
       --help, -h                                                show help
       --version, -v                                             print the version

Notice that Fedora and elasticsearch connection parameters are global options, have default values, and can be defined by environment variables as described in the help.

For example, to specify a different Fedora user and password, do:

    pass-tools -u myUser -h myPass assign pi [args]

or, using environment variables

    export PASS_FEDORA_USER=myUser
    export PASS_FEDORA_PASSWORD=myPass
    pass-tools assign pi [args]
    pass-tools assign pi [diffentArgs]

### Assigning a PI (and associated submissions) to users

Use the `assign pi` subcommands to assign a grant PI to a new person   This is mainly intended for massaging data on demo and test PASS instances.
The `-s` flag optionally assigns any submissions that were associated with the grant (and have the grant's former PI as the submitter) to the new user.

The first argument is expected to be an ID (URL, or locatorID) of the person to 
whom grants are being assigned, followed by any number of grant IDs (URLs, or localKeys) of grants to assign to the given user.

For example:

    pass-tools assign pi -s johnshopkins.edu:jhed:newsubmitter1 johnshopkins.edu:grant:1234 johnshopkins.edu:grant:5678

Or, with global options

    pass-tools -u myUser -p myPass assign pi  -s johnshopkins.edu:jhed:newsubmitter1 johnshopkins.edu:grant:1234 johnshopkins.edu:grant:5678

### Migrating metadata blobs

The `migrate metadata` subcommands migrate submission metadata blobs from the original format (not defined by any schema), to blobs governed by the schema https://oa-pass.github.io/metadata-schemas/jhu/global.json

    pass-tools migrate metadata

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
