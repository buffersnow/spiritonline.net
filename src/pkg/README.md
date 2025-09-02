# bFXServer/[_pkg_](.)

_NOTE: tbi stands for To Be Implemented!_

## Core

- [**log**](./log): Simple injectable logger, supports file, console and 24 hour regular archival
- [**settings**](./settings): Comprehensive configuration package via environment variables and command-line arguments

## Network

- [**net**](./net): Easy to use injectable TCP/UDP network package
- [**web**](./web): go-fiber/v2 based wrapper package for Web with pre-configured middleware and error helpers

### Protocol

- [**com**](./com): Protocol used by [**router**](/src/cmd/router) (tbi)
- [**gp**](./gp): Shared GameSpy protocol elements

## Utility

- [**lifecycle**](./lifecycle): One file package to catch shutdown interrupts
- [**util**](./util): Basic global utility package for other packages
- [**version**](./version): Service versioning helper, embedded during compilation
- [**security**](./security): Injectable encoding, encryption and ECDSA signing utilities
