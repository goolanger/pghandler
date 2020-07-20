#!/bin/sh

. /usr/bin/setpwd.sh

export DUMP_FILE="$BACKUP_DIR/$PG_SCHEMA-$DB_ENV-$1.dump"

if [ ! -f "$DUMP_FILE" ]; then
  echo "dump file $DUMP_FILE is missing - cannot restore database"
  exit 1
fi

echo "pg_restore -v -c -w -h $PG_HOST -p $PG_PORT -U $PG_USER -n $PG_SCHEMA -d $PG_DB $DUMP_FILE"
pg_restore -v -c -w -h $PG_HOST -p $PG_PORT -U $PG_USER -n $PG_SCHEMA -d $PG_DB $DUMP_FILE
