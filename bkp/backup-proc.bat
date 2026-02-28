@echo off
setlocal enabledelayedexpansion

REM Configuration
set "LOG_FILE=%BACKUP_DIR%\backup.log"
set "CHANGELOG_FILE=%BACKUP_DIR%\CHANGELOG"
set "MAX_BACKUPS=5"

REM Helper function to generate timestamp
for /f "tokens=2-4 delims=/ " %%a in ('date /t') do (set mydate=%%c%%a%%b)
for /f "tokens=1-2 delims=/:" %%a in ('time /t') do (set mytime=%%a%%b)
set "TIMESTAMP=%mydate%_%mytime%"

REM Logging function
setlocal enabledelayedexpansion
(
  echo [%date% %time%] %1
) >> "%LOG_FILE%"
endlocal

REM Step 1: Create backup folder
if not exist "%BACKUP_DIR%\backup_%TIMESTAMP%" (
  mkdir "%BACKUP_DIR%\backup_%TIMESTAMP%"
  echo Created backup folder: %BACKUP_DIR%\backup_%TIMESTAMP%
  set "NEW_FOLDER=%BACKUP_DIR%\backup_%TIMESTAMP%"
) else (
  echo ERROR: Failed to create backup folder
  exit /b 1
)

REM Step 2: Copy files
if exist "%3" (
  call "%3" "%1" "!NEW_FOLDER!"
) else (
  echo ERROR: Copy script not found: %3
  exit /b 1
)

REM Step 3: Cleanup old backups
for /f "tokens=*" %%d in ('dir /b /ad "%BACKUP_DIR%\backup_*" 2^>nul') do (
  set /a COUNT+=1
)
if !COUNT! gtr %MAX_BACKUPS% (
  set /a TO_DELETE=!COUNT!-%MAX_BACKUPS%
  for /f "tokens=*" %%d in ('dir /b /ad /o-d "%BACKUP_DIR%\backup_*" 2^>nul') do (
    if !TO_DELETE! gtr 0 (
      rmdir /s /q "%BACKUP_DIR%\%%d"
      set /a TO_DELETE-=1
    )
  )
)

REM Step 4: Update changelog
echo [CREATED] %date% - backup_%TIMESTAMP% >> "%CHANGELOG_FILE%"

endlocal