#!/bin/sh

. /usr/bin/setpwd.sh

echo "pgbench -i -s $1 -h $PG_HOST -p $PG_PORT -U $PG_USER $PG_DB"
pgbench -i -s $1 -h $PG_HOST -p $PG_PORT -U $PG_USER $PG_DB
