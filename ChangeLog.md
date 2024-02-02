ChangeLog
=========

All noticeable changes in the project  are documented in this file.

Format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project uses [semantic versions](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.22.1] 2024-02-01

Return git's stderr contents along with errors

## [1.22.0] 2024-01-10

Add more Pause-related functions

## [1.21.0] 2023-10-07

Add DiffUpstream and RemoteUpdate to git

## [1.20.1] 2023-09-19

Git interface issues

## Fixed
* `PopStash` implementation

## Added
* `Unstage`

## [1.20.0] 2023-09-09

Misc. changes

### Modified

* Renamed the `ing2` package to `chars`
* Added check for level to onerror
* Misc. git-related adjustments
* Added Taskfile

## [1.19.1] 2023-09-02

Use slo for git instead of stdout

## [1.19.0] 2023-09-01

Complete the git.Handler interface

## [1.18.0] 2023-08-30

Add git Init, Fetch and Status

## [1.17.0] 2023-08-30

Rework the git interface

## [1.16.2] 2023-08-29

Remove redundant caller attributes

## [1.16.1] 2023-08-26

Rework onerror implementation

### Fixed

* `onerror.Wit()` mismatch in relation to `slog.With()`
* missing `Fatal()` signature in `onerror.Recorder` interface

## [1.16.0] 2023-08-24

Additions from other packages

### Added

* onerror package from onerror
* cron package from retrace
* rsynccb package from retrace

### Modified

* Updated README

## [1.15.0] 2023-08-23

Add the git package from bumpy-ride

## [1.14.1] 2023-08-13

Upgrade Go version

## [1.14.0] 2023-05-13

Add FindExec and FindLibExec

## [1.13.0] 2023-03-31

Add the CreatetheseDirs function

## [1.12.0] 2022-10-21

Add the pointers package

## [1.11.0] 2022-08-26

Enhancements

### Added
* ing2.GetRandomLetters
    * Reimplements `ing2.GetRandomString`, using `crypto/rand`

### Modified
* Upgrade Go version
* Use an atomic pointer for `ing2.randomseed`

## [1.10.0] 2022-03-24

Misc. changes

### Modified

* Renamed package: stringing -> ing2
* Updated go version

### Removed

* Deleted the slice package
    * Superseeded by `golang.org/x/exp/slice`

## [1.9.0] 2022-02-03

Add PathToURLUnchecked

## [1.8.2] 2022-01-15

Make path to URL absolute (reprise)

## [1.8.1] 2022-01-15

Make path to URL absolute

## [1.8.0] 2022-01-12

Added and removed

### Added

* Common string utilities

### Removed

* Moved logger and onerror to their own repos

## [1.7.0] 2021-12-11

Remove env.ParseArgs

## [1.6.0] 2021-12-03

Update Go version and bumpy-ride

## [1.5.0] 2021-09-26

Remove obsolete args in env.SetDirs

## [1.4.0] 2021-09-20

Subpackages

### Modified

* Moved files to its own packages
* Simplified naming

## [1.3.1] 2021-07-05

Improvements

### Added

* Integer validators

### Modified

* httpstatus helpers names

## [1.3.0] 2021-06-27

Improvements

### Modified

* HTTP status error handler

## [1.2.2] 2021-06-26

Tag release

## [1.2.1] 2021-06-26

Quick fix

### Fixed

* HTTP helpers for nil response

## [1.2.0] 2021-06-26

New Helpers

### Added

* HTTP Status Helpers

## [1.1.1] 2021-06-24

Quick fix

### Fixed
* Handling of percentage when parsing URL

## [1.1.0] 2021-04-17

Added some utilities

### Added
* Parsing of global/environment arguments
* Router/logger-related middleware

## [1.0.0] 2021-04-16

Initial release.
