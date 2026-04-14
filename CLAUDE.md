# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This repository contains example applications demonstrating the usage of the [gotmc/usbtmc](https://github.com/gotmc/usbtmc) Go library for communicating with USB Test and Measurement Class (USBTMC) enabled devices.

## Multi-Module Structure

Each example is its own Go module with an independent `go.mod` (no top-level `go.mod`). Build and run commands must be executed from within the example directory or via the Justfile recipes.

- **key33220/** — Keysight/Agilent 33220A Function Generator (Go 1.21, usbtmc v0.8.0)
- **keyu2751a/** — Keysight U2751A Matrix Switch (Go 1.12, usbtmc v0.4.0)

## Development Commands

The project uses [Just](https://github.com/casey/just) as the command runner:

```bash
just check              # go fmt + go vet
just unit               # unit tests with -race -short -cover
just unit -run TestName # run a single test
just lint               # golangci-lint with .golangci.yaml config
just cover              # HTML coverage report (unit tests by default)
just tidy               # go mod tidy + verify
just update <mod>       # update a single dependency
just updateall          # update all dependencies
just outdated           # list outdated direct deps (requires go-mod-outdated)
just docs               # local pkgsite viewer
just ex1                # build and run key33220 example
```

## Key Patterns

1. **Driver Import**: Examples use the Google gousb driver via blank import: `_ "github.com/gotmc/usbtmc/driver/google"` (follows the `database/sql` registration pattern).

2. **Device Creation**: Examples connect either by VISA address string (`ctx.NewDevice(address)`) or by VID/PID (`ctx.NewDeviceByVIDPID(vid, pid)`).

3. **Debug Levels**: Both examples accept `-d` or `--debug` CLI flags to set USB debug level.
