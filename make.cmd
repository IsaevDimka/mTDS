del tds.exe
go build -ldflags "-s -w" tds.go flow.go click.go import.go
@ECHO OFF
if %ERRORLEVEL%==0 tds run
if not %ERRORLEVEL%==0 echo "failed..." exit