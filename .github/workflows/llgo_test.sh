#!/bin/bash
set -e

root_dir=$(pwd)

for dir in "$@"; do
    echo "Testing $dir"
    cd "$root_dir/$dir"
    llgo test .
done

echo "All tests passed!"
exit 0
