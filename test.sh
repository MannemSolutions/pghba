#!/bin/bash
set -ex
rm -f pg_hba.conf
./pghba add -a md5 -t '(local|hostssl)' -d '(db_[a-e])' -s '(127.0.0.1|192.168.2.13)' -U '(postgres|test{1..5})'
NUMLINES=$(cat ./pg_hba.conf | wc -l)
test ${NUMLINES} -eq 90
UNIQUE_LINES=$(sort -u ./pg_hba.conf | wc -l)
test ${NUMLINES} -eq ${UNIQUE_LINES}
test $(grep -c db_e ./pg_hba.conf) -eq 18
test $(grep -c local ./pg_hba.conf) -eq 30
echo OK
