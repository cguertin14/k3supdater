# k3supdater

[![Go Report Card](https://goreportcard.com/badge/github.com/cguertin14/k3supdater)](https://goreportcard.com/report/github.com/cguertin14/k3supdater)
[![codecov](https://codecov.io/gh/cguertin14/k3supdater/branch/main/graph/badge.svg?token=BUUUB7F5HX)](https://codecov.io/gh/cguertin14/k3supdater)
[![Docker Repository on Quay](https://quay.io/repository/cguertin14/k3supdater/status "Docker Repository on Quay")](https://quay.io/repository/cguertin14/k3supdater)


Updater (similar to Renovate Bot) for k3s ansible playbook versions.


## Github Access Token

First, you need to create a Github access token with write access to the repository which you want the bot to push to. Then, you'll need to set the `GITHUB_ACCESS_TOKEN` environment variable on your machine. 

## Usage

To use `k3supdater`, you'll need to [download the appropriate binary for your machine](https://github.com/cguertin14/k3supdater/releases) or compile it locally. It can also be used via docker, using the `quay.io/cguertin14/k3supdater` image.

Then, you'll be able to run update commands like so:
```bash
$ k3supdater update --repo-owner cguertin14 --repo-name k3s-ansible-ha
```

## Kubernetes Manifests

If you want to use `k3supdater` inside your kubernetes cluster, make sure to check out the [k8s manifests](./manifests/README.md) we have defined for this project.