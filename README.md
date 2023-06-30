# github-exporter
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/clambin/github-exporter?color=green&label=Release&style=plastic)
![Build)](https://github.com/clambin/github-exporter/workflows/Build/badge.svg)
![Codecov](https://img.shields.io/codecov/c/gh/clambin/github-exporter?style=plastic)
![Go Report Card](https://goreportcard.com/badge/github.com/clambin/github-exporter)

Prometheus exporter for GitHub repositories

## Installation
Binaries are available on the [release](https://github.com/clambin/github-exporter/releases) page. Docker images are available on [ghcr.io](https://ghcr.io/clambin/github-exporter).

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

| metric                           | type    | description                                          | labels                        |
|----------------------------------|---------|------------------------------------------------------|-------------------------------|
| github_monitor_forks             | gauge   | Number of forks                                      | archived, fork, private, repo |
| github_monitor_issues            | gauge   | Number of open issues raised against the repo        | archived, fork, private, repo |
| github_monitor_pulls             | gauge   | Number of open pull requests raised against the repo | archived, fork, private, repo |
| github_monitor_stars             | gauge   | Number of stars                                      | archived, fork, private, repo |
| github_monitor_api_latency       | summary | Latency of GitHub API calls                          | application, method, path     |
| github_monitor_api_errors_total  | counter | Number of errors raised by GitHub API calls          | application, method, path     |
| github_monitor_api_in_flight     | gauge   | Number of GitHub API calls in flight                 | application                   |
| github_monitor_api_max_in_flight | gauge   | Highest number of GitHub API calls in flight         | application                   |


## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.


