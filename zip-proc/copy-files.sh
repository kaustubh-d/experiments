#!/bin/bash

# File containing list of files to copy
listFile="$1"
# Source directory where the files are located
sourceDir="$2"
# Destination directory where the files will be copied
destDir="$3"

if [[ ! -f "$listFile" ]]; then
  echo "Error: List file not found: $listFile"
  exit 1
fi

if [[ ! -d "$sourceDir" ]]; then
  echo "Error: Source directory not found: $sourceDir"
  exit 1
fi

if [[ ! -d "$destDir" ]]; then
  echo "Creating destination directory: $destDir"
  mkdir -p "$destDir"
fi

# Read each file name from the list and copy it to the destination directory
success_log=./success.log
failure_log=./failures.log
# --info=NAME1 for linux
rsync -R --files-from="$listFile" --itemize --ignore-errors \
  "$sourceDir/./" "$destDir/" > "${success_log}" 2> "${failure_log}"

EXIT_STATUS=$?

if [ $EXIT_STATUS -eq 0 ]; then
    echo "Successfully copied all files."
    echo "Success log: ${success_log}" && cat "${success_log}"
elif [ $EXIT_STATUS -eq 23 ]; then
    echo "Warning: Some files failed to copy"
    echo "Success log: ${success_log}" && cat "${success_log}"
    echo "Failure log: ${failure_log}" && cat "${failure_log}"
else
    echo "Error: Rsync failed with code $EXIT_STATUS"
    echo "Failure log: ${failure_log}" && cat "${failure_log}"
fi

# For any failed copy operations, exit with a non-zero status code
exit $EXIT_STATUS