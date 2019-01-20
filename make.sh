echo "Building TDS"
rm tds && go build -gccgoflags "-s -w" tds.go flow.go click.go import.go
echo "Building TDS complete ..."