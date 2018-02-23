# Welcome to Rady-framework!

[![Coverage Status](https://coveralls.io/repos/github/Hexilee/rady/badge.svg)](https://coveralls.io/github/Hexilee/rady)
[![Go Report Card](https://goreportcard.com/badge/github.com/Hexilee/rady)](https://goreportcard.com/report/github.com/Hexilee/rady)
[![Build Status](https://travis-ci.org/Hexilee/rady.svg?branch=master)](https://travis-ci.org/Hexilee/rady)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Hexilee/rady/blob/master/LICENSE)
[![Documentation](https://godoc.org/github.com/Hexilee/rady?status.svg)](https://godoc.org/github.com/Hexilee/rady)

#### [Example](https://github.com/Hexilee/rady/tree/master/example) and Docs are under updating

## What can rady do now?
- Dependency injection (Include components and value in config file).
- Structured route registration (annotation route. router, controller and middleware, can be embedded in other router).
- Middleware registration.
- Initialize components in factory mode.
- Entities registration.
- Config file hot-reload (Include factories' recall).
- Gorm integration (In subproject rorm).
- Some [wrappers](https://github.com/Hexilee/rady/tree/master/middleware) (cors, jwt, logger) of echo-middleware.

## Todos
- Route Registration in Config File (Hot-reload)
- Rady-cli:
    - Generate code with route config.
    - Package manager integration (glide, godep).
    - DI test.
    
- Editor plugin (Goland and vscode):
    - Tag indecator.
    - Route inspection.
    - Injection inspection.
    - Config file injection inspection (Can jump between config and code).
    - Rady-cli integration.
- AOP
- More middleware wrappers
- Cache
- Dashbord
- Cloud


