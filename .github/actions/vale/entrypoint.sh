#!/bin/sh
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -eo pipefail
IFS=$'\n\t'

printenv

# Uses a default config path, but allows a user to override it from a file in their repository
configFilePath=/etc/vale/.vale.ini

if [ -f "${INPUT_CONFIGFILEPATH}" ]; then
    configFilePath="${INPUT_CONFIGFILEPATH}"
fi

echo "Settings:"
echo ""
echo "lintUnchangedFiles: ${INPUT_LINTUNCHANGEDFILES}"
echo "lintDirectory: ${INPUT_LINTDIRECTORY}"
echo "fileGlob: ${INPUT_FILEGLOB}"
echo "configFilePath: ${configFilePath}"
echo ""

cd $GITHUB_WORKSPACE

set -x
vale --output=JSON --config="${configFilePath}" --glob="${INPUT_FILEGLOB}" "${INPUT_LINTDIRECTORY}"
set +x
