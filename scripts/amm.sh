frogchain query amm list-pool
frogchain query amm show-pool [id]
frogchain query amm get-pool-param [id]
frogchain query amm get-pool-share-token [id]
frogchain query amm get-pool-assets [id]
frogchain query amm get-swap-exact-tokens-for-tokens [pool-id] [amount-in] [path]
frogchain query amm get-swap-tokens-for-exact-tokens [pool-id] [amount-out] [path]

frogchain tx amm create-pool [pool-param] [pool-assets] [asset-amounts]
frogchain tx amm add-liquidity [pool-id] [desired-amounts] [min-amounts]
frogchain tx amm remove-liquidity [pool-id] [desired-amount] [min-amounts]
frogchain tx amm swap-exact-tokens-for-tokens [pool-id] [amount-in] [amount-out-min] [path] [to] [deadline]
frogchain tx amm swap-tokens-for-exact-tokens [pool-id] [amount-out] [path] [to] [deadline]