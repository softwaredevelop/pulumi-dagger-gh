#!/usr/bin/env bash

set -e

find "$HOME"/.cache/go-build -type f -not -newermt "$(date -R -d '1 day ago')" -delete &&
    go clean -cache -modcache &&
    export GOCACHE=/tmp/gocache
