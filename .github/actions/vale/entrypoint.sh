#!/bin/sh
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

vale --version

echo $GITHUB_WORKSPACE
ls -lah $GITHUB_WORKSPACE
