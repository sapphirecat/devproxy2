# Change Log

All notable changes to this project will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/),
with the clarification that "updating dependencies" is a feature update.
Unsupported Go versions MAY be removed in such an update, at the
developer's sole discretion.
 
## [2.1.0] - 2021-12-19

Updated dependencies and minimally tested (it builds and starts up) on
Go 1.13 and Go 1.18beta1.
 
### Changed
 
- Updated go-toml from 1.7.0 to 2.0.0-beta.4
 
## [2.0.0] - 2020-11-21
  
Initial devproxy2 release, superseding the archived
[1.x](https://github.com/sapphirecat/devproxy)

### Changed

1. Configuration is done by file
2. A hostname matcher never receives `host:port` format
3. Actions are a specific, declarative redirection, not a Go function
4. Removed command-line options: `-target`, `-listen`, and `-port`
5. Added `debug_rule` configuration allows per-rule control of verbosity
6. Using Go Modules, supported on Go 1.13+
7. Restructured as an application, not a library
