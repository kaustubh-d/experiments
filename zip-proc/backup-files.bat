@echo off
setlocal enabledelayedexpansion

:CopyFilesFromList
setlocal
set "listFile=%~1"
set "sourceDir=%~2"
set "destDir=%~3"

if not exist "!listFile!" (
  echo Error: List file not found: !listFile!
  endlocal
  exit /b 1
)

if not exist "!sourceDir!" (
  echo Error: Source directory not found: !sourceDir!
  endlocal
  exit /b 1
)

if not exist "!destDir!" (
  echo Creating destination directory: !destDir!
  mkdir "!destDir!"
)

for /f "usebackq delims=" %%A in ("!listFile!") do (
  set "filePath=!sourceDir!\%%A"
  if exist "!filePath!" (
    copy "!filePath!" "!destDir!\" >nul
    echo Copied: %%A
  ) else (
    echo Warning: File not found: !filePath!
  )
)

endlocal
exit /b 0