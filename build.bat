@echo off

set OUTPUT_FOLDER=bin
set BINARY_NAME=KC-Checker

echo Compiling for Windows...
set GOOS=windows
set GOARCH=amd64
go build -o %OUTPUT_FOLDER%\%BINARY_NAME%.exe
if %ERRORLEVEL% neq 0 (
    echo Compilation failed for Windows
    exit /b %ERRORLEVEL%
)
echo Compilation for Windows successful

echo Compiling for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o %OUTPUT_FOLDER%\%BINARY_NAME%
if %ERRORLEVEL% neq 0 (
    echo Compilation failed for Linux
    exit /b %ERRORLEVEL%
)
echo Compilation for Linux successful

echo Compilation completed successfully for both Windows and Linux
exit /b 0
