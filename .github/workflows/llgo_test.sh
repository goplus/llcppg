#!/bin/bash
set -e

root_dir=$(pwd)

for dir in "$@"; do
    echo "Testing $dir"
    llgo test "$root_dir/$dir"
done

echo "All tests passed!"
exit 0
