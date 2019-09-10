#!/bin/sh
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -eo pipefail
IFS=$'\n\t'

printenv

echo "Settings:"
echo ""
echo "lintUnchangedFiles: ${INPUT_LINTUNCHANGEDFILES}"
echo "lintDirectory: ${INPUT_LINTDIRECTORY}"
echo "fileGlob: ${INPUT_FILEGLOB}"
echo ""

cd $GITHUB_WORKSPACE
vale --glob="${INPUT_FILEGLOB}" "${INPUT_LINTDIRECTORY}" --output=JSON
