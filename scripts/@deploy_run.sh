#!/bin/bash

NEW_BINARY=shoot-new
BINARY=/tmp/shoot
LOG=/home/tenlab/shoot/log
PORT=9712

run() {
	echo "Killing existing process on port $PORT"
	pkill -f .*$BINARY.*\-httpport\=$PORT.*
	echo "Starting on port $PORT.."
	until nohup $BINARY -httpport=$PORT >> $LOG 2>&1; do
		echo "[${date}] Crashed with exit code $?.  Respawning.." >> $LOG
		sleep 3
	done
}

rm -f $BINARY
cp $NEW_BINARY $BINARY

run &
