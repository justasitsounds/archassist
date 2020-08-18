#!/bin/sh

for filename in "${PWD}/*.puml"; do
    rf=$(eval echo "$filename")
    ff="${rf/.puml/.png}"
    cat $filename | docker run --rm -i -v $AWSICONDIST:/include -e ALLOW_PLANTUML_INCLUDE=true think/plantuml -tpng > $ff
done
