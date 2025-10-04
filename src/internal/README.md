# bFXServer/[_internal_](.)

_NOTE: tbi stands for To Be Implemented!_

## Core

- [**router**](.): Internal services routing through a custom tcp protocol (tbi)
- [**proxy**](./proxy): Custom HTTP/1.1 reverse proxy capable of handling non-standard clients

## Service

- [**myspace**](./myspace): MySpaceIM server (uses [**gamespy**](/src/pkg/gp) presence protocol)
- [**iwmaster**](./iwmaster): InfinityWard master server for online play via custom clients

### Gamespy

- [**gsp**](./gsp): GameSpy Presence, contains GPCM (wip) and GPSP (tbi) servers
- [**wfc**](./wfc): Nintendo WiFi Connection, houses NAS-Wii, NAS-DS (WIP), DLS1 (tbi), Conntest, SAKE (tbi), RACE (tbi)
- [**qr**](./qr): GameSpy Query and Report (wip), also known as the Gamespy master server.

## Edge

_none yet_
