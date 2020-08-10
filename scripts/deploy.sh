#!/bin/bash

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

set -e

source deploy.env

BINARY=./../droplist-api
DEPLOY_DIR="/home/tenlab/droplist"
CMD=api
LOG_FILE=log

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
