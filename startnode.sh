#!/bin/bash

rm -rf $HOME/.statesetd/

cd $HOME

statesetd init --chain-id=testing testing --home=$HOME/.statesetd
statesetd keys add validator --keyring-backend=test --home=$HOME/.statesetd
statesetd add-genesis-account $(statesetd keys show validator -a --keyring-backend=test --home=$HOME/.statesetd) 1000000000stake,1000000000valtoken --home=$HOME/.statesetd
statesetd gentx validator 500000000stake --keyring-backend=test --home=$HOME/.statesetd --chain-id=testing
statesetd collect-gentxs --home=$HOME/.statesetd

cat $HOME/.statesetd/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="10s"' > $HOME/.statesetd/config/tmp_genesis.json && mv $HOME/.statesetd/config/tmp_genesis.json $HOME/.statesetd/config/genesis.json

statesetd start --home=$HOME/.statesetd