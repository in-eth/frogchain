#!/bin/bash
set -e

# Load shell variables
. ./network/hermes/variables.sh

### Sleep is needed otherwise the relayer crashes when trying to init
sleep 1
### Restore Keys
$HERMES_BINARY --config $CONFIG_DIR keys delete --chain frogchain --all
sleep 5

$HERMES_BINARY --config $CONFIG_DIR keys add --chain frogchain --mnemonic-file './network/hermes/key1'
sleep 5

$HERMES_BINARY --config $CONFIG_DIR keys delete --chain localosmosis --all
sleep 5

$HERMES_BINARY --config $CONFIG_DIR keys add --chain localosmosis --mnemonic-file './network/hermes/key2'
sleep 5
