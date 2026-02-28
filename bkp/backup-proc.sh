#!/bin/bash

# Configuration
LOG_FILE="${BACKUP_DIR}/backup.log"
CHANGELOG_FILE="${BACKUP_DIR}/CHANGELOG"
MAX_BACKUPS=5

# Helper function to generate timestamp
get_timestamp() {
  local format="${1:-%Y%m%d_%H%M%S}"
  date "+$format"
}

# Logging function
log() {
  echo "[$(get_timestamp '%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Step 1: Create new folder with timestamp
create_backup_folder() {
  local backup_dir="$1"
  local timestamp=$(get_timestamp)
  local new_folder="${backup_dir}/backup_${timestamp}"

  if mkdir -p "$new_folder"; then
    log "Created backup folder: $new_folder"
    echo "$new_folder"
    return 0
  else
    log "ERROR: Failed to create backup folder at $backup_dir"
    return 1
  fi
}

# Step 2: Copy files from source
copy_files_to_backup() {
  local source_dir="$1"
  local backup_folder="$2"
  local copy_script="$3"

  if [ ! -f "$copy_script" ]; then
    log "ERROR: Copy script not found: $copy_script"
    return 1
  fi

  if bash "$copy_script" "$source_dir" "$backup_folder"; then
    log "Successfully copied files to $backup_folder"
    return 0
  else
    log "ERROR: Failed to copy files to $backup_folder"
    return 1
  fi
}

# Step 3: Retain only latest 5 backups
cleanup_old_backups() {
  local backup_dir="$1"
  local max_backups="$2"

  local backup_count=$(find "$backup_dir" -maxdepth 1 -type d -name 'backup_*' | wc -l)

  if [ "$backup_count" -le "$max_backups" ]; then
    log "Backup count ($backup_count) within limit ($max_backups). No cleanup needed."
    return 0
  fi

  local to_delete=$((backup_count - max_backups))
  find "$backup_dir" -maxdepth 1 -type d -name 'backup_*' -printf '%T@ %p\n' | \
  sort -n | head -n "$to_delete" | cut -d' ' -f2- | while read old_folder; do
    if rm -rf "$old_folder"; then
      log "Deleted old backup folder: $old_folder"
      echo "$old_folder"
    else
      log "ERROR: Failed to delete $old_folder"
    fi
  done

  return 0
}

# Step 4: Generate changelog
update_changelog() {
  local backup_dir="$1"
  local new_folder="$2"
  local deleted_folders="$3"

  local timestamp=$(get_timestamp '%d-%b-%Y')

  # Log created backup
  echo "[CREATED] $timestamp - $(basename "$new_folder")" >> "$CHANGELOG_FILE"

  # Log deleted backups
  if [ -n "$deleted_folders" ]; then
    echo "$deleted_folders" | while read deleted_folder; do
      echo "[DELETED] $timestamp - $(basename "$deleted_folder")" >> "$CHANGELOG_FILE"
    done
  fi

  log "Updated changelog: $CHANGELOG_FILE"
  return 0
}

# Main execution
main() {
  local source_dir="${1:?Source directory required}"
  local backup_dir="${2:?Backup directory required}"
  local copy_script="${3:?Copy script path required}"

  log "========== Backup Process Started =========="

  # Step 1
  new_folder=$(create_backup_folder "$backup_dir") || exit 1

  # Step 2
  copy_files_to_backup "$source_dir" "$new_folder" "$copy_script" || exit 1

  # Step 3
  deleted=$(cleanup_old_backups "$backup_dir" "$MAX_BACKUPS")

  # Step 4
  update_changelog "$backup_dir" "$new_folder" "$deleted" || exit 1

  log "========== Backup Process Completed Successfully =========="
}

main "$@"