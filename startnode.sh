#!/bin/bash

rm -rf $HOME/.statesetd/

cd $HOME

statesetd init --chain-id=testing testing --home=$HOME/.statesetd
statesetd keys add validator --keyring-backend=test --home=$HOME/.statesetd
statesetd add-genesis-account $(statesetd keys show validator -a --keyring-backend=test --home=$HOME/.osmosisd) 1000000000stake,1000000000valtoken --home=$HOME/.statesetd
statesetd gentx validator 500000000stake --keyring-backend=test --home=$HOME/.statesetd --chain-id=testing
statesetd collect-gentxs --home=$HOME/.statesetd

statesetd start --home=$HOME/.statesetd