# Indexing Engine

Made with :heart: by<br/>
<a href="https://figment.io"><img alt="Figment" src="assets/figment-logo.svg" height="32px" align="bottom"/></a>

Indexing Engine is a library that bundles together common building blocks of a blockchain indexer in the form of separate packages.
Although the packages are intended to be used together, each of them provides different functionality and can be used in isolation to satisfy a particular need.

The library consists of the following packages:

- [`pipeline`](pipeline)
- [`worker`](worker)
- [`metrics`](metrics)
- [`health`](health)
- [`store`](store)

## Installation

To install the latest version of the library, run the following command:

```shell
go get -u github.com/figment-networks/indexing-engine
```

## Documentation

You can find the documentation for each of the packages inside its directory.

## License

The library is licensed under [Apache License 2.0](LICENSE).
