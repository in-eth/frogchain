#!/bin/bash
set -e

# Load shell variables
. ./network/hermes/variables.sh

### Configure the clients and connection
echo "Initiating connection handshake..."
$HERMES_BINARY --config $CONFIG_DIR create connection --a-chain frogchain --b-chain localosmosis

sleep 2

# echo "Creating transfer channel..."
$HERMES_BINARY --config $CONFIG_DIR create channel --a-chain frogchain --a-connection connection-0 --a-port transfer --b-port transfer

sleep 2
