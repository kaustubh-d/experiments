@echo off
setlocal enabledelayedexpansion

:CopyFilesFromList
setlocal
@REM File containing list of files to copy
set "listFile=%~1"
@REM Source directory where the files are located
set "sourceDir=%~2"
@REM Destination directory where the files will be copied
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

@REM Read each file name from the list and copy it to the destination directory
for /f "usebackq delims=" %%A in ("!listFile!") do (
  @REM Generate the full path of the source file
  set "filePath=%%A"

  @REM Extract the file name from the file path
  for %%F in ("!filePath!") do set "fileName=%%~nxF"

  @REM Extract the directory path from the file path
  for %%D in ("!filePath!") do set "relDir=%%~dpD"

  set "absSrcFilePath=!sourceDir!\!filePath!"
  if exist "!absSrcFilePath!" (
    @REM Construct the source file path by combining the source directory with the relative path
    set "srcAbsDirPath=!sourceDir!!relDir!"

    @REM Construct the destination path by combining the destination directory with the relative path
    set "destAbsDirPath=!destDir!!relDir!"

    @REM Copy fileName from srcAbsDirPath to destAbsDirPath, creating directories as needed
    robocopy "!srcAbsDirPath!" "!destAbsDirPath!" "!fileName!" /NJH /NJS /NS /NC /NFL /NDL >nul
    if !errorlevel! leq 1 (
      echo Copied: !filePath! to !destAbsDirPath!
    ) else (
      echo Warning: Failed to copy: !filePath! to !destAbsDirPath!
    )
  ) else (
    echo Warning: File not found: !absSrcFilePath!
  )
)

endlocal
exit /b 0