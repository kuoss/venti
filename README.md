# venti
Metrics and logs data visualizer

[![release](https://github.com/kuoss/venti/actions/workflows/release.yml/badge.svg)](https://github.com/kuoss/venti/actions)
[![pull-request](https://github.com/kuoss/venti/actions/workflows/pull-request.yml/badge.svg)](https://github.com/kuoss/venti/actions)
[![Coverage Status](https://coveralls.io/repos/github/kuoss/venti/badge.svg?branch=main)](https://coveralls.io/github/kuoss/venti?branch=main)
[![codecov](https://codecov.io/gh/kuoss/venti/branch/main/graph/badge.svg?token=EXPE6OS8HJ)](https://codecov.io/gh/kuoss/venti)
[![GitHub license](https://img.shields.io/github/license/kuoss/venti.svg)](https://github.com/kuoss/venti/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/kuoss/venti.svg)](https://github.com/kuoss/venti/stargazers)
[![contribuiton welcome](https://img.shields.io/badge/contributions-welcome-orange.svg)](https://github.com/kuoss/venti/blob/main/CONTRIBUTING.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/kuoss/venti)](https://goreportcard.com/report/github.com/kuoss/venti)

# users

```
$ sudo apt-get install apache2-utils
$ htpasswd -nbBC 12 testuser topsecret
testuser:$2y$12$LbMigrXXYhtJcQ0kwI5Wue1uYzF20idYdWECtl3P79Ack.GhwnDOO
```
```
# /app/etc/users.yml
users:
...
- username: testuser
  hash: $2y$12$LbMigrXXYhtJcQ0kwI5Wue1uYzF20idYdWECtl3P79Ack.GhwnDOO
  isAdmin: false
```

## Contributors

<a href="https://github.com/kuoss/venti/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=kuoss/venti" />
</a>
