go build -ldflags "-s -w" tds.go
@ECHO OFF
if %ERRORLEVEL%==0 tds run
if not %ERRORLEVEL%==0 echo "failed..." exit