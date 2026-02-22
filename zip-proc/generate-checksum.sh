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

  local checksum_cmd
  case "$checksum_type" in
    sha256) checksum_cmd="shasum -a 256" ;;
    sha1) checksum_cmd="shasum -a 1" ;;
    md5) checksum_cmd="md5sum" ;;
    *) checksum_cmd="shasum -a 256" ;;
  esac

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
  zipinfo -1 "$ZIP_FILE" | grep -v '/$' > "$file_list"

  # Iterate over each file in the list, extract file content from the zip,
  # compute its checksum, and save to the output file
  cat "$file_list" | while read -r file_name; do
    if [[ -n "$file_name" && "$file_name" != */ ]]; then
      checksum=$(unzip -p "$zip_file" "$file_name" | $checksum_cmd | awk '{print $1}')
      echo "$file_name $checksum" >> "$checksum_file"
    fi
  done

  echo "Checksums saved to $output_folder/checksums.txt"
}

generate_zip_meta "$ZIP_FILE" "$OUTPUT_FOLDER" "sha256"

echo "File list saved to $OUTPUT_FOLDER/file_list.txt"
echo "Checksums saved to $OUTPUT_FOLDER/checksums.txt"