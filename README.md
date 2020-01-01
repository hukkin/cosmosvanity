[![Build Status](https://travis-ci.com/hukkinj1/cosmosvanity.svg?branch=master)](https://travis-ci.com/hukkinj1/cosmosvanity)
[![codecov.io](https://codecov.io/gh/hukkinj1/cosmosvanity/branch/master/graph/badge.svg)](https://codecov.io/gh/hukkinj1/cosmosvanity)
[![GolangCI](https://golangci.com/badges/github.com/hukkinj1/cosmosvanity.svg)](https://golangci.com/r/github.com/hukkinj1/cosmosvanity)
# cosmosvanity

<!--- Don't edit the version line below manually. Let bump2version do it for you. -->
> Version 0.0.1

> [Cosmos](https://cosmos.network) vanity address generator. Command line user interface.

## Features
* Generate Cosmos bech32 vanity addresses
* Use all CPU cores
* Specify a substring that the address must
    * start with (`--startswith` flag)
    * end with (`--endswith` flag)
    * contain (`--contains` flag)
* Supports Linux, macOS and Windows

## Installing
Download the latest binary release from the [_Releases_ page](https://github.com/hukkinj1/cosmosvanity/releases). Alternatively, build from source yourself.

## Usage examples
Find an address that starts with "00000" (e.g. cosmos100000v3fpv4qg2a9ea6sj70gykxpt63wgjen2p)
```bash
./cosmosvanity --startswith 00000
```

Find an address that ends with "8888" (e.g. cosmos134dck5uddzjure8pyprmmqat96k3jlypn28888)
```bash
./cosmosvanity --endswith 8888
```

Find an address containing the substring "gener" (e.g. cosmos1z39wgener7azgh22s5a3pyswtnjkx2w0hvn3rv)
```bash
./cosmosvanity --contains gener
```

Combine flags shown above
```bash
./cosmosvanity --contains 8888 --startswith a --endswith c
```
