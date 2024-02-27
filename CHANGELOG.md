# Changelog

This section details changes between revisions of this utility.

- [Changelog](#changelog)
  - [0.2.0](#020)
  - [0.1.2](#012)
  - [0.1.1](#011)
  - [0.1.0](#010)


## 0.2.0

The increment to the minor version here represents an organizational re-arrangement in addition to expanded functionality. This repository will maintain an identity solely as the SAGA tool while other information has been migrated to [another repository](https://gitlab.com/ucfdracolab/saga-data).

The following changes are present in v0.2.0

- Bumped go version to latest (1.22.x)
- Re-organized and renamed project to match publishing
- Added checkpoint system to save populations over time
- Added `-resume` argument to restart from a checkpoint file
- Added `-chkfreq` and `-chkpath` arguments for configuring checkpoint settings
- Added `-config` verb to run a given config file
  - See `//input/config/test.json` for available arguments and an example
  - Tags are subject to addition until the first major release
- Added system for declaring GAs as data
- Added system for declaring jobs as data

## 0.1.2

- Updated project organization to more clearly segregate documentation/data from code/logic.
- Updated go version to latest (1.21.5).
- Reduced default output during genetic evolution. Prior behavior can be obtained with additional runtime data using the `-verbose` CLI flag.
- Replaced `-minimize` flag with `-evolve` flag.
- Replaced `-population` flag with `-pop` flag.
- Updated reporting of argument validation and defaults.
- Updated in-line documentation to conform to standard degree of explanation.

## 0.1.1

- Finalize project migration to public namespace. Currently documentation and interaction are sub-optimal, so we'll do one big pass for the 0.1.2 publication before we do the major overhaul on configuration and execution that'll be numbered 0.2.0

## 0.1.0

- Added initial project version