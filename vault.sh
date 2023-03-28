#!/usr/bin/env bash

set -e

export GITHUB_OWNER=$(vault read -field=GITHUB_OWNER pulumi/gh-own)
export GITHUB_TOKEN=$(vault read -field=GITHUB_TOKEN pulumi/gh-tkn-dt)
export GITHUB_TOKEN=$(vault read -field=GITHUB_TOKEN pulumi/gh-tkn-dt-fg)
export PULUMI_ACCESS_TOKEN=$(vault read -field=PULUMI_ACCESS_TOKEN pulumi/tkn-dt)
export PULUMI_ORG_NAME=$(vault read -field=PULUMI_ORG_NAME pulumi/org-1)
