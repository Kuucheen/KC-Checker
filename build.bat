@echo off

set OUTPUT_FOLDER=bin
set BINARY_NAME=KC-Checker
set VERSION=1.0.4
set DESCRIPTION=Open source proxy checker built in Go, fast and beautiful. Created by Kuchen (Kuucheen on GitHub).

echo Compiling for Windows...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-X 'main.Version=%VERSION%' -X 'main.Description=%DESCRIPTION%'" -o %OUTPUT_FOLDER%\%BINARY_NAME%.exe

if %ERRORLEVEL% neq 0 (
    echo Compilation failed for Windows
    exit /b %ERRORLEVEL%
)
echo Compilation for Windows successful

echo Compiling for Linux...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-X 'main.Version=%VERSION%' -X 'main.Description=%DESCRIPTION%'" -o %OUTPUT_FOLDER%\%BINARY_NAME%
if %ERRORLEVEL% neq 0 (
    echo Compilation failed for Linux
    exit /b %ERRORLEVEL%
)
echo Compilation for Linux successful

echo Compilation completed successfully for both Windows and Linux
exit /b 0
