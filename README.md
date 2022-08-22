# gomodsearch

Find go module direct, indirect and transitive dependencies.

## Installation

Clone repo and then run
```bash
$ cd gomodsearch
$ make install
```

## Usage

```bash
$ gomodsearch <path> <module>[@version] ...
```
Example
```bash
$ gomodsearch . golang.org/x/mod golang.org/x/mod@v0.0.1 golang.org/x/net
```
> If the version is not declared it searches all existing versions.
## License
[MIT](https://github.com/francogeller/gomodsearch/blob/develop/LICENCE.md)