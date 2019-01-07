rm tds
go build -gccgoflags "-s -w" tds.go flow.go click.go import.go
./tds run