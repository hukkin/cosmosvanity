[![Build Status](https://travis-ci.com/hukkinj1/cosmosvanity.svg?branch=master)](https://travis-ci.com/hukkinj1/cosmosvanity)
[![codecov.io](https://codecov.io/gh/hukkinj1/cosmosvanity/branch/master/graph/badge.svg)](https://codecov.io/gh/hukkinj1/cosmosvanity)
[![GolangCI](https://golangci.com/badges/github.com/hukkinj1/cosmosvanity.svg)](https://golangci.com/r/github.com/hukkinj1/cosmosvanity)
# cosmosvanity

<!--- Don't edit the version line below manually. Let bump2version do it for you. -->
> Version 0.1.0

> CLI tool for generating [Cosmos](https://cosmos.network) vanity addresses

## Features
* Generate Cosmos bech32 vanity addresses
* Use all CPU cores
* Specify a substring that the addresses must
    * start with
    * end with
    * contain
* Set required minimum amount of letters (a-z) or digits (0-9) in the addresses
* Binaries built for Linux, macOS and Windows

## Installing
Download the latest binary release from the [_Releases_](https://github.com/hukkinj1/cosmosvanity/releases) page. Alternatively, build from source yourself.

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

Find an address consisting of letters only (e.g. cosmos1rfqkejeaxlxwtjxucnrathlzgnvgcgldzmuxxe)
```bash
./cosmosvanity --letters 38
```

Find an address with at least 26 digits (e.g. cosmos1r573c4086585u084926726x535y3k2ktxpr88l)
```bash
./cosmosvanity --digits 26
```

Generate 5 addresses (the default is 1)
```bash
./cosmosvanity -n 5
```

Combine flags introduced above
```bash
./cosmosvanity --contains 8888 --startswith a --endswith c
```
