# Ssh Check

Check ssh connectivity to nodes within a cluster

## What it is
This repository builds the container image used by Kuberhealthy to run the ssh-check check.

## Image
- `docker.io/kuberhealthy/ssh-check`
- `ghcr.io/kuberhealthy/ssh-check`
- Tags: `latest` and short git SHA for `main` pushes; `vX.Y.Z` for releases.

## Quick start
- Apply the example manifest: `kubectl apply -f healthcheck.yaml`
- Edit the manifest to set any required inputs for your environment.

## Build locally
- `docker build -f ./Dockerfile -t kuberhealthy/ssh-check:dev .`

## Contributing
Issues and PRs are welcome. Please keep changes focused and add a short README update when behavior changes.

## License
See `LICENSE`.
