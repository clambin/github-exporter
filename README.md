# github-exporter
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/clambin/github-exporter?color=green&label=Release&style=plastic)
![Build)](https://github.com/clambin/github-exporter/workflows/Build/badge.svg)
![Codecov](https://img.shields.io/codecov/c/gh/clambin/github-exporter?style=plastic)
![Go Report Card](https://goreportcard.com/badge/github.com/clambin/github-exporter)

Prometheus exporter for GitHub repositories

## Installation
Binaries are available on the [release](https://github.com/clambin/github-exporter/releases) page. Container images are available on [ghcr.io](https://ghcr.io/clambin/github-exporter).

## Usage
### Command-line options
The following command-line arguments can be passed:

```
Usage:
  github-exporter [flags]

Flags:
      --config string   Configuration file
      --debug           Log debug messages
  -h, --help            help for github-exporter
  -v, --version         version for github-exporter
```

By default, github-monitor looks for the configuration file (`config.yaml`) in the following locations:
- `/etc/github-exporter`
- `$HOME/.github-exporter`
- `.`

### Configuration
The configuration file contains the following options:

```
# set debug to true to log debug messages
debug: false
# listener address for the /metrics endpoint
addr: :9090
# this section lists all repos to be monitored. Repos can be specified in either the `repo` section as an individual
# repo, or in the `user` section, which will monitor all repos for that user.
# notes: 
#   - github-exporter will not remove any duplicate repos
#   - organizations are currently not supported
repos:
  user:
    - clambin
  repo:
    - clambin/github-exporter
  # set archived to true to report metrics for archived repos. By default these are not reported on.
  archived: false
git:
  # token contains your github token to access the GitHub API.
  token: <your-token>
  # cache specifies how long to cache GitHub information.  
  cache: 1h
```

Any value in the configuration file may be overriden by setting an environment variable with a prefix `GITHUB_EXPORTER_`.
E.g. to override the git token, set the following variable:

```
export GITHUB_EXPORTER_GIT.TOKEN="your-token"
```

## Prometheus metrics

| metric | type |  labels | help |
| --- | --- |  --- | --- |
| github_exporter_api_inflight_current | GAUGE | |current in flight requests |
| github_exporter_api_inflight_max | GAUGE | |maximum in flight requests |
| github_exporter_forks | GAUGE | archived, repo|Total number of forks |
| github_exporter_http_request_duration_seconds | SUMMARY | code, method, path|http request duration in seconds |
| github_exporter_http_requests_total | COUNTER | code, method, path|total number of http requests |
| github_exporter_issues | GAUGE | archived, repo|Total number of open issues |
| github_exporter_pulls | GAUGE | archived, repo|Total number of open pull requests |
| github_exporter_stars | GAUGE | archived, repo|Total number of stars |

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
