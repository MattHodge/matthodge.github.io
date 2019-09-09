#!/bin/sh
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -eo pipefail
IFS=$'\n\t'

vale --version

printenv

echo "Settings:"
echo ""
echo "lintAllFiles: ${INPUT_LINTALLFILES}"
echo "lintDirectory: ${INPUT_LINTDIRECTORY}"
echo "fileGlob: $INPUT_FILEGLOB"
