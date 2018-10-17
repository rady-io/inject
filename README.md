# Welcome to Rady-framework!

[![Coverage Status](https://coveralls.io/repos/github/rady-io/inject/badge.svg)](https://coveralls.io/github/rady-io/inject)
[![Go Report Card](https://goreportcard.com/badge/github.com/rady-io/inject)](https://goreportcard.com/report/github.com/rady-io/inject)
[![Build Status](https://travis-ci.org/rady-io/inject.svg?branch=master)](https://travis-ci.org/rady-io/inject)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/rady-io/inject/blob/master/LICENSE)
[![Documentation](https://godoc.org/github.com/rady-io/inject?status.svg)](https://godoc.org/github.com/rady-io/inject)

#### Example and Docs are under updating

## What can rady do now?
- Dependency injection (Include components and value in config file).
- Structured route registration (annotation route. router, controller and middleware, can be embedded in other router).
- Middleware registration.
- Initialize components in factory function.
- Entities registration.
- Config file hot-reload (Include factories' recall).
- Some [wrappers](https://github.com/Hexilee/rady-middleware) (cors, jwt, logger) for echo-middleware.
- DI test
- Env-dependent config file

## Todos
- Complete lifetime system for beans
- Gorm integration (In project [rorm](https://github.com/Hexilee/rorm)).
- Integration with [htest](https://github.com/Hexilee/htest)
- Editor plugin (Goland and vscode):
    - Tag indecator.
    - Route inspection.
    - Injection inspection.
    - Config file injection inspection (Can jump between config and code).

- AOP
- More middleware wrappers
- Cache
- Dashbord
- Cloud


