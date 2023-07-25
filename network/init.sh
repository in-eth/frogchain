#!/bin/bash

BINARY=frogchaind
CHAIN_DIR=./data
CHAINID_1=frogchain
CHAINID_2=localosmosis

VAL_MNEMONIC_1="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion"
WALLET_MNEMONIC_1="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"
WALLET_MNEMONIC_2="veteran try aware erosion drink dance decade comic dawn museum release episode original list ability owner size tuition surface ceiling depth seminar capable only"
RLY_MNEMONIC_1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"

P2PPORT_1=16656
P2PPORT_2=16657
RPCPORT_1=16657
RPCPORT_2=26657
RESTPORT_1=1316
RESTPORT_2=1317
ROSETTA_1=8080
ROSETTA_2=9080
pprof_laddr_1=6061

# Stop if it is already running 
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    killall $BINARY
fi

echo "Removing previous data..."
rm -rf $CHAIN_DIR/$CHAINID_1 &> /dev/null

# Add directories for both chains, exit if an error occurs
if ! mkdir -p $CHAIN_DIR/$CHAINID_1 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

echo "Initializing $CHAINID_1..."

$BINARY init test --home $CHAIN_DIR/$CHAINID_1 --chain-id=$CHAINID_1

echo "Adding genesis accounts..."
echo $VAL_MNEMONIC_1 | $BINARY keys add val1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_1 | $BINARY keys add wallet1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $WALLET_MNEMONIC_2 | $BINARY keys add wallet2 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test
echo $RLY_MNEMONIC_1 | $BINARY keys add rly1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend=test 

$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID_1 keys show val1 --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR/$CHAINID_1
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID_1 keys show wallet1 --keyring-backend test -a) 100000000000frog  --home $CHAIN_DIR/$CHAINID_1
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID_1 keys show wallet2 --keyring-backend test -a) 100000000000frog  --home $CHAIN_DIR/$CHAINID_1
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID_1 keys show rly1 --keyring-backend test -a) 100000000000frog  --home $CHAIN_DIR/$CHAINID_1

echo "Creating and collecting gentx..."
$BINARY gentx val1 7000000000stake --home $CHAIN_DIR/$CHAINID_1 --chain-id $CHAINID_1 --keyring-backend test
$BINARY collect-gentxs --home $CHAIN_DIR/$CHAINID_1

echo "Creating and collecting gentx..."
export alice=$($BINARY keys show alice -a --home $CHAIN_DIR/$CHAINID_1) && echo $alice;

echo "Changing defaults and ports in app.toml and config.toml files..."
sed -i -e 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's#"localhost:6060"#"localhost:'"$pprof_laddr_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's/1317/'"$RESTPORT_1"'/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's/":8080"/":'"$ROSETTA_1"'"/g' $CHAIN_DIR/$CHAINID_1/config/app.toml

# Update host chain genesis to allow x/bank/MsgSend ICA tx execution
# cat $CHAIN_DIR/$CHAINID_1/config/genesis.json | jq '(.app_state.interchainaccounts.host_genesis_state.params.allow_messages) |= ["/cosmos.bank.v1beta1.MsgSend"]' > $CHAIN_DIR/$CHAINID_1/config/genesis.json.tmp
# mv $CHAIN_DIR/$CHAINID_1/config/genesis.json.tmp $CHAIN_DIR/$CHAINID_1/config/genesis.json
