#!/bin/bash

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

set -e

SSH_USER=tenlab
SSH_IP=206.189.253.106
SSH_PORT=1472
BINARY=droplist-new
DEPLOY_DIR="/home/tenlab/droplist"
CMD=api
LOG_FILE=log

echo "Building binary.."
GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o $BINARY ./../cmd/$CMD

echo "Copying binary and deploy script.."
scp -P $SSH_PORT $BINARY $SSH_USER@$SSH_IP:$DEPLOY_DIR/.
scp -P $SSH_PORT @deploy_run.sh $SSH_USER@$SSH_IP:$DEPLOY_DIR/run.sh

echo "Ready to deploy. Press any key to continue"
read _

echo "Running deploy script.."
ssh -p $SSH_PORT $SSH_USER@$SSH_IP bash -c "'
source /home/tenlab/.bash_profile
cd $DEPLOY_DIR
chmod +x run.sh
./run.sh
sleep 3
tail $LOG_FILE
exit
'"

echo "Cleaning up.."
rm $BINARY
