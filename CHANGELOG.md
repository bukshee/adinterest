# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Resize
- BitCopy
- Shift
- Rotate
- Left
- Right
- Mid
- Append
- SetMul
- ClearMul

## Changed
- New: in case len<=0 it no longer returns nil. It returns a Len(0) BitField.
- bumped up bitfield64 version number

## Deprecated
- Copy: this method will be renamed to Clone for added clarity.

## [1.0.0] - 2020-05-23

First release.
