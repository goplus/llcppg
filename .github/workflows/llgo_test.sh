#!/bin/bash
set -e

root_dir=$(pwd)

for dir in "$@"; do
    echo "Testing $dir"
    cd "$root_dir/$dir"
    output=$(llgo test . 2>&1)
    echo "$output"
done

echo "All tests passed!"
exit 0
