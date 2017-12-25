#! /bin/bash

set -e

if [ ! -d .git ]; then
    echo 'Run this script from root of repository' 1>&2
    exit 1
fi

command=selfupdate-example

rm -rf release
gox -verbose ./cmd/$command
mkdir -p release
mv selfupdate-example_* release/
cd release
for bin in *; do
    mv "$bin" "$command"
    zip "${bin}.zip" "$command"
    rm "$command"
done
