
---
sidebar_position: 1
---

Ignite CLI version:		v0.27.1

# `x/amm`

## Abstract

This document specifies the amm module of the Frogchain Network.

The Automated Market Making (AMM) module is responsible for handling token swaps. It is designed to provide decentralized exchange functionality similar to Uniswap on the Ethereum network. It enables efficient and frictionless trading of digital assets within the Cosmos ecosystem.

This module is used in the Frogchain Network.

## Contents

* [Concepts](#concepts)
* [State](#pool-state)
    * [Pool Param](#Pool-Param)
    * [Pool Asset & Weight](#Pool-Asset-&-Weight)
    * [Pool Share Token](#Pool-Share-Token)
    * [Minimum Liquidity](#Minimum-Liquidity)
* [Messages](#messages)
* [Events](#events)
    * [Message Events](#message-events)
    * [Keeper Events](#keeper-events)
* [Client](#client)
    * [CLI](#cli)
    * [Query](#query)
    * [Transactions](#transactions)

## Concepts

Already mentioned before, it implements the concepts of uniswap v2.

If you want more about that:
```uniswap reference
https://docs.uniswap.org/contracts/v2/concepts/protocol-overview/how-uniswap-works
```


**Note:** The amm module is made to implement uniswap v2 platform in cosmos network, while it has some differences.

The differences are:

* `uniswap` - a pool only includes 2 coin types to swap from one coin to another coin(exactly path includes only 2 tokens). If you want to include more tokens in path, it causes more transactions generation.
* `amm` - a pool can include many coins unlike uniswap so that you can swap tokens with complex routes in one transaction.

## Pool State

Pools contain pool information for a uniquely identified pool and for swap including pool id, pool params, pool assets, and share tokens. For efficiency, since pool tokens must also be feetchd to swap, pool structs also store the balance of pool assets as `sdk.Coin`.

The `Pool` interface is defined as follows:

```go
type Pool struct {
	Id uint64
	PoolParam PoolParam // pool params for swap and exit fees
	PoolAssets   []types.Coin // assets in the pool
	AssetWeights []uint64 // represents the weight of pool asset
	ShareToken types.Coi// LP token
	MinimumLiquidity uint64 // min liquidity is the minimum amount of share token to maintain the pool
	IsActivated bool // Pool activate state
}
```

### Pool Param

PoolParam contains information for pool fees. That includes the fees and fee collector's address.
Fees are divided into two parts.

 1. Swap Fee
 : fee percent for swap tokens. Module cut fee from incoming tokens and send it to fee collector.
 1. Exit Fee
 : fee percent for remove liquidity. Amm cut fee from burning share token and send it to fee collector. It's for pool maintaince so that the pool can work regularly without that liquidity provider.

 Fee range is 0~10^8 meaning that 10^8 represents 100%.

 FeeCollector is and address represents the user collects the fees.

 The `PoolParam` interface is defined as follows:

```go
type PoolParam struct {
	// 100% of the swap fee goes to the liquidity providers — the amount of the
	// underlying token that can beredeemed by each pool token increases.
	// this represents the percent * 10^6 for the float. (0 ~ 10^8)
	SwapFee uint64 

	// When exiting a pool, the user provides the minimum amount of tokens they are
	// willing to receive as they are returning their shares of the pool. However,
	// unlike joining a pool, exiting a pool requires the user to pay the exit fee,
	// which is set as a param of the pool. The user's share tokens burnt as result.
	// Exiting the pool using a single asset is also possible.
	// this represents the percent * 10^6 for the float. (0 ~ 10^8)
	ExitFee uint64 
	
	FeeCollector string // feeCollector is an address which receive fees(for swap and exit)
}
```

### Pool Asset & Weight

Pool Asset is a `type.Coin` type that has coin `Denom` which represents the coin name and `Amount` indicates the amount of coin.

Weight is for swap tokens. It controls the coin's value so that you can change the weight when coin has not enough worth.

### Pool Share Token

Pool Share Token is a `type.Coin` type like pool asset but has a meaningful value.

When liquidity provider deposits coins for a pool, then he receives the liquidity - share token.
Share Token represents the sovereignty of the pool, meaning how many coins are owned by the provider. So if the provider withdraw the liquidity, the owned coins are withdrawed from the pool(except the fee amount).

### Minimum Liquidity

Minimum Liquidity is a value for maintaining the pool. It's cut from the liquidity of the creator so that the pool makes creator can not modify the rate of coins.

It's defined as 1000 here, but can modify if you want.

## Messages

### MsgCreatePool

Create a new pool.

```protobuf reference
https://github.com/mdanny1209/frogchain/blob/master/proto/frogchain/amm/tx.proto#L20-L26
```

The message will fail under the following conditions:

* The `PoolParam`'s `FeeCollector` is not a valid address
* The `PoolParam`'s `SwapFee` and `ExitFee` is greater than `100%=10^8`
* `PoolAssets` length is less than two
* `PoolAssets` and `AssetWeights` lengths are different
* Creater has not enough coin balance
* Minted share token amount is under `Minimum Liquidity`

### MsgAddLiquidity

Deposit coins for sovereignty of pool.

```protobuf reference
https://github.com/mdanny1209/frogchain/blob/master/proto/frogchain/amm/tx.proto#L32-L37
```

The message will fail under the following conditions:

* `poolId` is invalid, can't find pool with that id
* The `DesiredAmounts` values are under `MinAmounts`
* `PoolAssets` and `DesiredAmounts` lengths are different
* Minted share token amount is 0
* Calculated coin amount is under `MinAmounts`
* Liquidity provider has not enough balance

### MsgRemoveLiquidity

Withdraw coins from a pool.

```protobuf reference
https://github.com/mdanny1209/frogchain/blob/master/proto/frogchain/amm/tx.proto#L43-L47
```

The message will fail under the following conditions:

* `poolId` is invalid, can't find pool with that id
* Calculated coin amount is under `MinAmounts`
* Liquidity provider has not enough share token amount

### MsgSwapExactTokensForTokens

Swap coins with exact amount of coins.

```protobuf reference
https://github.com/mdanny1209/frogchain/blob/master/proto/frogchain/amm/tx.proto#L55-L63
```

The message will fail under the following conditions:

* `poolId` is invalid, can't find pool with that id
* `To` address is not valid
* `deadline` is passed
* `path` has invalid coin denom
* `path` length is less than 2
* `path` has same value
* Transaction creator has not enough coin balance

### MsgSwapTokensForExactTokens

Get exact amount of coins from coins

```protobuf reference
https://github.com/mdanny1209/frogchain/blob/master/proto/frogchain/amm/tx.proto#L70-L77
```

The message will fail under the following conditions:

* `poolId` is invalid, can't find pool with that id
* `To` address is not valid
* `deadline` is passed
* `path` has invalid coin denom
* `path` length is less than 2
* `path` has same value
* Transaction creator has not enough coin balance
* `AmountOut` is exceed the pool balance

## Events

The bank module emits the following events:

### Message Events

#### MsgCreatePool

| Type       | Attribute Key | Attribute Value    |
| ---------- | ------------- | ------------------ |
| createpool | poolparam     | {types.poolparam}  |
| createpool | poolassets    | {[]sdk.coins}      |
| createpool | assetweights  | {[]sdk.coins}      |
| message    | module        | amm                |
| message    | action        | createpool         |
| message    | sender        | {senderAddress}    |


#### MsgAddLiquidity

| Type         | Attribute Key  | Attribute Value    |
| ------------ | -------------- | ------------------ |
| addliquidity | poolid         | {uint64} 	 		 |
| addliquidity | desiredamounts | {[]uint64}         |
| addliquidity | minamounts     | {[]uint64}         |
| message      | module        	| amm                |
| message      | action        	| addliquidity       |
| message      | sender        	| {senderAddress}    |


#### MsgRemoveLiquidity

| Type            | Attribute Key | Attribute Value    |
| --------------- | ------------- | ------------------ |
| removeliquidity | poolid        | {uint64} 		   |
| removeliquidity | liquidity     | {liquidityamount}  |
| removeliquidity | minamounts    | {[]uint64}		   |
| message         | module        | amm                |
| message         | action        | removeliquidity    |
| message         | sender        | {senderAddress}    |


#### MsgSwapExactTokensForTokens

| Type     | Attribute Key | Attribute Value     |
| -------- | ------------- | ------------------- |
| swap     | poolid        | {uint64} 			 |
| swap     | amountin      | {uint64}            |
| swap     | amountoutmin  | {uint64}            |
| swap     | path          | {[]string}          |
| swap     | to        	   | {receipientaddress} |
| swap     | deadline      | {string}            |
| message  | module        | amm                 |
| message  | action        | multisend           |
| message  | sender        | {senderAddress}     |


#### MsgSwapTokensForExactTokens

| Type     | Attribute Key | Attribute Value     |
| -------- | ------------- | ------------------- |
| swap     | poolid        | {uint64} 			 |
| swap     | amountout     | {uint64}            |
| swap     | path          | {[]string}          |
| swap     | to        	   | {receipientaddress} |
| swap     | deadline      | {string}            |
| message  | module        | amm                 |
| message  | action        | multisend           |
| message  | sender        | {senderAddress}     |


### Keeper Events

In addition to message events, the amm keeper will produce events when the following methods are called (or any method which ends up calling them)

#### CreatePool

```json
{
  "type": "pool_created",
  "attributes": [
    {
      "key": "spender",
      "value": "{{sdk.AccAddress of the transaction creator creating a pool}}",
      "index": true
    },
    {
      "key": "poolid",
      "value": "{{uint64 being created}}",
      "index": true
    },
    {
      "key": "sharetoken",
      "value": "{{sdk.Coins being minted}}",
      "index": true
    }
  ]
}
```

#### AddLiquidity

```json
{
  "type": "liquidity_added",
  "attributes": [
    {
      "key": "spender",
      "value": "{{sdk.AccAddress of the transaction creator adding liquidity}}",
      "index": true
    },
    {
      "key": "poolid",
      "value": "{{uint64 for pool}}",
      "index": true
    },
    {
      "key": "sharetoken",
      "value": "{{sdk.Coins being minted}}",
      "index": true
    }
  ]
}
```

#### RemoveLiquidity

```json
{
  "type": "liquidity_removed",
  "attributes": [
    {
      "key": "spender",
      "value": "{{sdk.AccAddress of the transaction creator removing liquidity}}",
      "index": true
    },
    {
      "key": "poolid",
      "value": "{{uint64 for pool}}",
      "index": true
    },
    {
      "key": "sharetoken",
      "value": "{{sdk.Coins being burnt}}",
      "index": true
    }
  ]
}
```

#### SwapExactTokensForTokens

```json
{
  "type": "tokens_swapped",
  "attributes": [
    {
      "key": "spender",
      "value": "{{sdk.AccAddress of the transaction creator swapping}}",
      "index": true
    },
    {
      "key": "poolid",
      "value": "{{uint64 for swap}}",
      "index": true
    },
    {
      "key": "swapamountout",
      "value": "{{uint being swapped}}",
      "index": true
    }
  ]
}
```

#### SwapTokensForExactTokens

```json
{
  "type": "exact_tokens_swapped",
  "attributes": [
    {
      "key": "spender",
      "value": "{{sdk.AccAddress of the transaction creator swapping}}",
      "index": true
    },
    {
      "key": "poolid",
      "value": "{{uint64 for swap}}",
      "index": true
    },
    {
      "key": "swapamountin",
      "value": "{{uint being swapped}}",
      "index": true
    }
  ]
}
```

## Client

### CLI

A user can query and interact with the `amm` module using the CLI.

#### Query

The `query` commands allow users to query `amm` state.

```shell
frogchaind query amm --help
```

##### list-pool

The `list-pool` command allows users to query pool list.

```shell
frogchaind query amm list-pool [flags]
```

Example:

```shell
frogchaind query amm list-pool
```

Example Output:

```yml
Pool:
- assetWeights:
  - "1"
  - "1"
  id: "0"
  isActivated: true
  minimumLiquidity: "1000"
  poolAssets:
  - amount: "102000"
    denom: foocoin
  - amount: "2000"
    denom: token
  poolParam:
    exitFee: "10"
    feeCollector: cosmos1g3z20q5jskz3g2anvs5hnxpn7tsa7dvylv34yd
    swapFee: "10"
  shareToken:
    amount: "10200"
    denom: frogchain-amm-pool-0-shareToken
- assetWeights:
  - "1"
  - "1"
  id: "1"
  isActivated: true
  minimumLiquidity: "1000"
  poolAssets:
  - amount: "102000"
    denom: foocoin
  - amount: "1020"
    denom: token
  poolParam:
    exitFee: "10"
    feeCollector: cosmos1g3z20q5jskz3g2anvs5hnxpn7tsa7dvylv34yd
    swapFee: "10"
  shareToken:
    amount: "10200"
    denom: frogchain-amm-pool-1-shareToken
pagination:
  next_key: null
  total: "0"
```

##### show-pool

The `show-pool` command allows users to query a pool.

```shell
frogchaind query amm show-pool [id] [flags]
```

Example:

```shell
frogchaind query amm show-pool 1
```

Example Output:

```yml
Pool:
  assetWeights:
  - "1"
  - "1"
  id: "1"
  isActivated: true
  minimumLiquidity: "1000"
  poolAssets:
  - amount: "102000"
    denom: foocoin
  - amount: "1020"
    denom: token
  poolParam:
    exitFee: "10"
    feeCollector: cosmos1g3z20q5jskz3g2anvs5hnxpn7tsa7dvylv34yd
    swapFee: "10"
  shareToken:
    amount: "10200"
    denom: frogchain-amm-pool-1-shareToken
```

##### get-pool-param

The `get-pool-param` command allows users to query pool param.

```shell
frogchaind query amm get-pool-param [id] [flags]
```

Example:

```shell
frogchaind query amm get-pool-param 1
```

Example Output:

```yml
poolParam:
  exitFee: "10"
  feeCollector: cosmos1g3z20q5jskz3g2anvs5hnxpn7tsa7dvylv34yd
  swapFee: "10"
```

##### get-pool-share-token

The `get-pool-share-token` command allows users to query pool param.

```shell
frogchaind query amm get-pool-share-token [id] [flags]
```

Example:

```shell
frogchaind query amm get-pool-share-token 1
```

Example Output:

```yml
shareToken:
  amount: "10200"
  denom: frogchain-amm-pool-1-shareToken
```

##### get-pool-assets

The `get-pool-assets` command allows users to query pool param.

```shell
frogchaind query amm get-pool-assets [id] [flags]
```

Example:

```shell
frogchaind query amm get-pool-assets 1
```

Example Output:

```yml
assets:
- amount: "102000"
  denom: foocoin
- amount: "1020"
  denom: token
```

##### get-swap-exact-tokens-for-tokens

The `get-swap-exact-tokens-for-tokens` command allows users to query pool param.

```shell
frogchaind get-swap-exact-tokens-for-tokens [pool-id] [amount-in] [path] [flags]
```

Example:

```shell
frogchaind query amm get-swap-exact-tokens-for-tokens 1 500 'foocoin,token'
```

Example Output:

```yml
amountOut: "500"
```

##### get-swap-tokens-for-exact-tokens

The `get-swap-tokens-for-exact-tokens` command allows users to query pool param.

```shell
frogchaind get-swap-tokens-for-exact-tokens [pool-id] [amount-in] [path] [flags]
```

Example:

```shell
frogchaind query amm get-swap-tokens-for-exact-tokens 1 500 'foocoin,token'
```

Example Output:

```yml
amountIn: "98076"
```

#### Transactions

The `tx` commands allow users to interact with the `amm` module.

```shell
frogchaind tx amm --help
```

##### create-pool

The `create-pool` command allows users to create a pool.

```shell
frogchaind tx amm create-pool [pool-param] [pool-assets] [asset-weights] [flags]
```

Example:

```shell
frogchaind tx amm create-pool '{"SwapFee":10,"ExitFee":10,"FeeCollector":"cosmos1g3z20q5jskz3g2anvs5hnxpn7tsa7dvylv34yd"}' '500foocoin,100token' '1,1'
```

##### add-liquidity

The `add-liquidity` command allows users to become a liquidity provider.

```shell
frogchaind tx amm add-liquidity [pool-id] [desired-amounts] [min-amounts] [flags]
```

Example:

```shell
frogchaind tx amm add-liquidity 1 '500,100' '100,50'
```

##### remove-liquidity

The `remove-liquidity` command allows users to withdraw assets.

```shell
frogchaind tx amm remove-liquidity [pool-id] [desired-share-token-amount] [min-amounts] [flags]
```

Example:

```shell
frogchaind tx amm remove-liquidity 1 500 '100,50'
```

##### swap-exact-tokens-for-tokens

The `swap-exact-tokens-for-tokens` command allows users to withdraw assets.

```shell
frogchaind tx amm swap-exact-tokens-for-tokens [pool-id] [amount-in] [amount-out-min] [path] [to] [deadline] [flags]
```

Example:

```shell
frogchaind tx amm swap-exact-tokens-for-tokens 1 500 50 'foocoin,token' Alice '159753216578951'
```

##### swap-tokens-for-exact-tokens

The `swap-tokens-for-exact-tokens` command allows users to withdraw assets.

```shell
frogchaind tx amm swap-tokens-for-exact-tokens [pool-id] [amount-out] [path] [to] [deadline] [flags]
```

Example:

```shell
frogchaind tx amm swap-tokens-for-exact-tokens 1 500 'foocoin,token' Alice '159753216578951'
```
