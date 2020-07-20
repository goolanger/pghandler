#!/bin/sh

# set defaults
. /usr/bin/setpwd.sh

export DUMP_FILE="$BACKUP_DIR/$PG_SCHEMA-$DB_ENV-$(date +"%F-%H%M%S").dump"
echo "pg_dump -v -w -Fc -h $PG_HOST -p $PG_PORT -n $PG_SCHEMA -U $PG_USER $PG_DB -f $DUMP_FILE"
pg_dump -v -w -Fc -h $PG_HOST -p $PG_PORT -n $PG_SCHEMA -U $PG_USER $PG_DB -f $DUMP_FILE
