# tinyman-go-sdk

# Overview
This is a Golang SDK providing access to the [Tinyman AMM](https://docs.tinyman.org/) on the Algorand blockchain. It currently supports V1.1 Tinyman.

# Installation
```command
go get github.com/synycboom/tinyman-go-sdk@v0.1.0
```

# Package overview
`v1` package provides a Tinyman client which is a main entry point for this SDK.
`v1/constants` contains constants for using with the SDK.
`v1/contracts` provides a getter function to retrieve the pool logic signature account.
`v1/pools` provides a liquidity pool utilities that you'll use to interact with it.
`v1/prepare` contains functions that prepare transaction groups to interact with the Tinyman contracts.

`utils` provides utilities like converting numbers, getting states, etc.

`types` contains data types used in the SDK.

`examples` are the example codes.

# Usage

## Boostrapping
Bootstrap a liquidity pool [/example/bootstrap](/example/bootstrap)

## Minting
Add assets to an existing pool in exchange for the liquidity pool asset [/example/mint](/example/mint).

## Burning
Exchange the liquidity pool asset for the pool assets [/example/burn](/example/burn).

## Swapping
Swap one asset for another in an existing pool [/example/swap](/example/swap).

## Redeeming
Redeem excess amounts from previous transactions [/example/redeem](/example/redeem).

## Running example
To run the examples, create a new /example/.env file by following the variables in /example/.env.example
Then setup /.vscode/launch.json, and use it to run the examples

# License
tinyman-go-sdk is licensed under a MIT license except for the exceptions listed below. See the LICENSE file for details.

## Exceptions
`v1/contracts/asc-v1_1.json` is currently unlicensed. It may be used by this SDK but may not be used in any other way or be distributed separately without the express permission of Tinyman.

# Disclaimer
Nothing in the repo constitutes professional and/or financial advice. Use this SDK at your own risk.
