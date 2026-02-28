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
hasError=0

while IFS= read -r filePath; do
  # Extract the file name from the file path
  fileName=$(basename "$filePath")
  # Extract the directory path from the file path
  relDir=$(dirname "$filePath")

  # Even if there are failures for single files, we copy the rest of the
  # files and report all failures at the end. This is to ensure we get as
  # many files copied as possible, even if some file copy had errors.
  absSrcFilePath="$sourceDir/$filePath"
  if [[ -f "$absSrcFilePath" ]]; then
    # Construct the source and destination directory paths
    srcAbsDirPath="$sourceDir/$relDir"
    destAbsDirPath="$destDir/$relDir"

    # Create destination directory if needed and copy the file
    mkdir -p "$destAbsDirPath"
    if cp "$absSrcFilePath" "$destAbsDirPath/$fileName"; then
      echo "Copied: $filePath to $destAbsDirPath"
    else
      echo "Warning: Failed to copy: $filePath to $destAbsDirPath"
      hasError=1
    fi
  else
    echo "Warning: File not found: $absSrcFilePath"
    hasError=1
  fi
done < "$listFile"

# For any failed copy operations, exit with a non-zero status code
exit $hasError