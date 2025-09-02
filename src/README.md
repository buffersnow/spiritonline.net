# bFXServer

### bFX => buffersnow.com Fast eXchange

bFX is the heart of SpiritOnline and the internal codename of the backend services. It is composed of a series of microservices, each handling one individual service or group of services.

This directory contains the following:

- [**cmd**](./cmd): Basic main functions of each microservice.
- [**internal**](./internal): Private implementations of a particular service.
- [**pkg**](./pkg): Common code and shareable/reusable libraries.
- [**tests**](./tests): Core-Services integration tests, automated UIA tests (where possible)
- [**tools**](./tools): Helper tooling, mainly for project and test generation
