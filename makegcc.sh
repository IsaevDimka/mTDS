go build -compiler gccgo -gccflags "-s -w" tds.go flow.go click.go
tds run