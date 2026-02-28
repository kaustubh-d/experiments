#!/usr/bin/env python3

import os
import sys
import shutil
import subprocess
from datetime import datetime
from pathlib import Path

# Configuration
BACKUP_DIR = os.getenv('BACKUP_DIR', '/tmp/backups')
LOG_FILE = os.path.join(BACKUP_DIR, 'backup.log')
CHANGELOG_FILE = os.path.join(BACKUP_DIR, 'CHANGELOG')
MAX_BACKUPS = 5


def get_timestamp(fmt='%Y%m%d_%H%M%S'):
  """Generate timestamp in specified format"""
  return datetime.now().strftime(fmt)


def log(message):
  """Log message with timestamp"""
  timestamp = get_timestamp('%Y-%m-%d %H:%M:%S')
  log_message = f"[{timestamp}] {message}"
  print(log_message)
  with open(LOG_FILE, 'a') as f:
    f.write(log_message + '\n')


def create_backup_folder(backup_dir):
  """Create new folder with timestamp"""
  timestamp = get_timestamp()
  new_folder = os.path.join(backup_dir, f'backup_{timestamp}')

  try:
    os.makedirs(new_folder, exist_ok=True)
    log(f"Created backup folder: {new_folder}")
    return new_folder
  except Exception as e:
    log(f"ERROR: Failed to create backup folder at {backup_dir}: {e}")
    return None


def copy_files_to_backup(source_dir, backup_folder, copy_script):
  """Copy files from source to backup folder"""
  if not os.path.isfile(copy_script):
    log(f"ERROR: Copy script not found: {copy_script}")
    return False

  try:
    subprocess.run(['bash', copy_script, source_dir, backup_folder], check=True)
    log(f"Successfully copied files to {backup_folder}")
    return True
  except subprocess.CalledProcessError as e:
    log(f"ERROR: Failed to copy files to {backup_folder}: {e}")
    return False


def cleanup_old_backups(backup_dir, max_backups):
  """Retain only the latest N backups"""
  backup_folders = sorted(
    [d for d in Path(backup_dir).iterdir() if d.is_dir() and d.name.startswith('backup_')],
    key=lambda x: x.stat().st_mtime
  )

  if len(backup_folders) <= max_backups:
    log(f"Backup count ({len(backup_folders)}) within limit ({max_backups}). No cleanup needed.")
    return []

  to_delete = backup_folders[:len(backup_folders) - max_backups]
  deleted_folders = []

  for old_folder in to_delete:
    try:
      shutil.rmtree(old_folder)
      log(f"Deleted old backup folder: {old_folder}")
      deleted_folders.append(str(old_folder))
    except Exception as e:
      log(f"ERROR: Failed to delete {old_folder}: {e}")

  return deleted_folders


def update_changelog(backup_dir, new_folder, deleted_folders):
  """Generate changelog"""
  timestamp = get_timestamp('%d-%b-%Y')

  try:
    with open(CHANGELOG_FILE, 'a') as f:
      f.write(f"[CREATED] {timestamp} - {os.path.basename(new_folder)}\n")
      for deleted_folder in deleted_folders:
        f.write(f"[DELETED] {timestamp} - {os.path.basename(deleted_folder)}\n")
    log(f"Updated changelog: {CHANGELOG_FILE}")
    return True
  except Exception as e:
    log(f"ERROR: Failed to update changelog: {e}")
    return False


def main():
  """Main execution"""
  if len(sys.argv) < 4:
    print("Usage: python backup-proc.py <source_dir> <backup_dir> <copy_script>")
    sys.exit(1)

  source_dir = sys.argv[1]
  backup_dir = sys.argv[2]
  copy_script = sys.argv[3]

  os.makedirs(backup_dir, exist_ok=True)

  log("========== Backup Process Started ==========")

  # Step 1
  new_folder = create_backup_folder(backup_dir)
  if not new_folder:
    sys.exit(1)

  # Step 2
  if not copy_files_to_backup(source_dir, new_folder, copy_script):
    sys.exit(1)

  # Step 3
  deleted = cleanup_old_backups(backup_dir, MAX_BACKUPS)

  # Step 4
  if not update_changelog(backup_dir, new_folder, deleted):
    sys.exit(1)

  log("========== Backup Process Completed Successfully ==========")


if __name__ == '__main__':
  main()