@echo off

if exist ..\rename-exif.exe (
    del ..\rename-exif.exe
)

go build -o ../rename-exif.exe .