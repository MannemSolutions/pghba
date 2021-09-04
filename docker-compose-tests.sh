#!/bin/bash
set -x
set -e
for TESTFILE in tests/*; do
  cat "$TESTFILE" | docker-compose run -i pghba
done
echo "All is as expected"
