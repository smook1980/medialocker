#!/usr/bin/env bash

rm rice-box.go 2&>1 /dev/null
rice embed-go
mv rice-box.go rice-box.go.pre
echo -e "// +build -dev\n" > rice-box.go.tmp
cat rice-box.go.pre >> rice-box.go.tmp
mv rice-box.go.tmp rice-box.go
rm -f rice-box.go.pre rice-box.go.tmp
