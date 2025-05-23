#!/bin/bash
set -e

root_dir=$(pwd)

for dir in "$@"; do
    echo "Testing $dir"
    cd "$root_dir/$dir" 
    output=$(llgo test . 2>&1)
    echo "$output"
    if echo "$output" | grep -q "exit code [^0]"; then
        echo "llgo test $(basename "$dir") failed"
        exit 1
    fi
done

echo "All tests passed!"
exit 0