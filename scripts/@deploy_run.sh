#!/bin/bash

CONFIG=config.toml
NEW_BINARY=droplist-api
BINARY=/tmp/droplist
LOG=/home/tenlab/droplist/log

run() {
	echo "Killing existing process $BINARY"
	pkill -f .*$BINARY.*
	echo "Starting $BINARY.."
	until nohup $BINARY >> $LOG 2>&1; do
		echo "[${date}] Crashed with exit code $?.  Respawning.." >> $LOG
		sleep 3
	done
}

rm -f $BINARY
cp $NEW_BINARY $BINARY
cp $CONFIG /tmp/$CONFIG

run &
