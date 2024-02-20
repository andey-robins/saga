# Changelog

This section details changes between revisions of this utility.

- [Changelog](#changelog)
  - [0.1.2](#012)
  - [0.1.1](#011)
  - [0.1.0](#010)


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