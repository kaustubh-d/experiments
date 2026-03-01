#!/bin/bash

# Script to generate checksum for a zip file
# Usage: ./generate-checksum.sh <zip_file> <output_folder>

if [[ $# -ne 2 ]]; then
  echo "Usage: $0 <zip_file> <output_folder>"
  exit 1
fi

ZIP_FILE="$1"
OUTPUT_FOLDER="$2"

# Verify zip file exists
if [[ ! -f "$ZIP_FILE" ]]; then
  echo "Error: Zip file '$ZIP_FILE' not found."
  exit 1
fi

# Verify output folder exists
if [[ ! -d "$OUTPUT_FOLDER" ]]; then
  echo "Error: Output folder '$OUTPUT_FOLDER' not found."
  exit 1
fi

# Function to get the appropriate checksum command for the current platform
# Parameters:
#   $1 - checksum type (sha256, sha1, md5)
# Output:
#   Returns the checksum command string
get_checksum_cmd() {
  local checksum_type="$1"
  local checksum_cmd

  case "$checksum_type" in
    sha256)
      if command -v sha256sum &> /dev/null; then
        checksum_cmd="sha256sum"
      elif command -v shasum &> /dev/null; then
        checksum_cmd="shasum -a 256"
      elif command -v sha256 &> /dev/null; then
        checksum_cmd="sha256"
      fi
      ;;
    sha1)
      if command -v sha1sum &> /dev/null; then
        checksum_cmd="sha1sum"
      elif command -v shasum &> /dev/null; then
        checksum_cmd="shasum -a 1"
      fi
      ;;
    md5)
      if command -v md5sum &> /dev/null; then
        checksum_cmd="md5sum"
      elif command -v md5 &> /dev/null; then
        checksum_cmd="md5"
      fi
      ;;
    *) checksum_cmd="shasum -a 256" ;;
  esac

  echo "$checksum_cmd"
}

# Function to generate checksums for files in the zip
# Parameters:
#   $1 - zip file path
#   $2 - output folder path
#   $3 - checksum type (optional, default: sha256, supported: sha256, sha1, md5)
# Output:
# checksums.txt file in the output folder with format: <file_name> <checksum>
# file_list.txt file in the output folder with the list of files in the zip (excluding directories)
generate_zip_meta() {
  local zip_file="$1"
  local output_folder="$2"
  # supported checksum types: sha256, sha1, md5 (default: sha256)
  local checksum_type="$3"

  local checksum_cmd="$(get_checksum_cmd "$checksum_type")"
  if [[ -z "$checksum_cmd" ]]; then
    echo "Error: No suitable checksum command found for type '$checksum_type'."
    exit 1
  fi

  local checksum_file="$output_folder/checksums.txt"
  local file_list="$output_folder/file_list.txt"

  if [[ ! -f "$zip_file" ]]; then
    echo "Error: Zip file '$zip_file' not found."
    return 1
  fi

  if [[ ! -d "$output_folder" ]]; then
    echo "Error: Output folder '$output_folder' not found."
    return 1
  fi

  # Generate the list of files in the zip (excluding directories)
  zipinfo -1 "$ZIP_FILE" | cut -d'/' -f2- | grep -v "^$" | grep -v '/$' > "$file_list"
  zip_folder_name=$(basename "$zip_file" .zip)

  # Iterate over each file in the list, extract file content from the zip,
  # compute its checksum, and save to the output file
  cat "$file_list" | while read -r file_name; do
    if [[ -n "$file_name" && "$file_name" != */ ]]; then
      checksum=$(unzip -p "$zip_file" "$zip_folder_name/$file_name" | $checksum_cmd | awk '{print $1}')
      echo "$file_name $checksum" >> "$checksum_file"
    fi
  done
}

generate_zip_meta "$ZIP_FILE" "$OUTPUT_FOLDER" "sha256"

echo "File list saved to $OUTPUT_FOLDER/file_list.txt"
echo "Checksums saved to $OUTPUT_FOLDER/checksums.txt"