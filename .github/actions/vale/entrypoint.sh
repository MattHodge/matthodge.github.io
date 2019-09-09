#!/bin/sh
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

vale --version

echo "Settings:"
echo ""
echo "lintAllFiles: ${INPUT_LINTALLFILES}"
echo "lintDirectory: ${INPUT_LINTDIRECTORY}"
