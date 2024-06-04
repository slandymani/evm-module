#!/bin/bash

KEY="mykey"
CHAINID="evmm_1-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t evmm-datadir.XXXXX)

echo "create and add new keys"
./evmmd keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init Evmos with moniker=$MONIKER and chain-id=$CHAINID"
./evmmd init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./evmmd add-genesis-account \
"$(./evmmd keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000aphoton,1000000000000000000stake \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./evmmd gentx $KEY 1000000000000000000stake --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./evmmd collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./evmmd validate-genesis --home $DATA_DIR

echo "starting evmm node $i in background ..."
./evmmd start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started evmm node"
tail -f /dev/null