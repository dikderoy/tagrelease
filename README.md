TagRelease
==========

[![Build Status](https://travis-ci.org/dikderoy/tagrelease.svg?branch=master)](https://travis-ci.org/dikderoy/tagrelease)
[![Maintainability](https://api.codeclimate.com/v1/badges/c0a9e573b147851c927a/maintainability)](https://codeclimate.com/github/dikderoy/tagrelease/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/c0a9e573b147851c927a/test_coverage)](https://codeclimate.com/github/dikderoy/tagrelease/test_coverage)


`tagrelease` is a simple tool to generate release/build
identifiers from git tag names in an automated way.

`tagrelease` requires Go 1.11 or newer to compile.

Installation
------------

Release binaries are available on the
[releases](https://github.com/dikderoy/tagrelease/releases) page

Usage
-----

options:

`--help`

prints help.

---

`--beta {list of strings}`

provide list of branches recognized as trunked, separated with comma,
all releases on these branches will be marked as B (beta),
except tagged ones (default `trunk`)

---

`--rc {strings}`

provide list of branches recognized as mainstream,
all releases on these branches will be marked as RC (release candidate),
except tagged ones (default `master`)

---

`--debug`

enable debug output (to stderr)

---

`-f, --format {string}`

select output format (default `release`):

- `release` - alias for pep440
- `pep440` - output a release string according to PEP440 format:
    `{major}.{minor}.{patch}[{kind(rc|b|a)}{diff}][+{revshort}]`
    example: `0.1.0+8f70f9c`, `0.8.2a84+57a182a`
- `semver` - semver release format:
    `{major}.{minor}.{patch}`
    example: `1.2.3`
- `short` - format: `{major}.{minor}`
    example: `1.2`
- `major` - output only major version
- `minor` - output only minor version
- `patch` - output only patch version
- `revision` - full git revision (40 chars)
    example: `57a182a871e042022c22b14ad6314b0618b582f8`
- `revshort` - short git revision (7 chars)
    example: `57a182a`

---

`-i, --increment {string}`

select increment strategy (default `upstream`)

- `major` - increment to next major version (`{major+1}.0.0`)
- `minor` - increment to next minor version (`{major}.{minor+1}.0`
- `patch` - increment to next patch version `{major}.{minor}.{patch+1}`
- `upstream` - increment to next patch if present, to next minor otherwise
- `never` - do not increment anything

---

`-o, --output {string}`

select output target (default is stdout):

- `{file-path}` - if not exists, will be created, otherwise will be truncated
- `-` - for stdout

---

`-d, --workdir {string}`

select working directory to look for repository
(default is current directory - `.`)
