frogchaind query amm list-pool
frogchaind query amm show-pool [id]
frogchaind query amm get-pool-param [id]
frogchaind query amm get-pool-share-token [id]
frogchaind query amm get-pool-assets [id]
frogchaind query amm get-swap-exact-tokens-for-tokens [pool-id] [amount-in] [path]
frogchaind query amm get-swap-tokens-for-exact-tokens [pool-id] [amount-out] [path]

frogchaind tx amm create-pool [pool-param] [pool-assets] [asset-weights]
frogchaind tx amm add-liquidity [pool-id] [desired-amounts] [min-amounts]
frogchaind tx amm remove-liquidity [pool-id] [desired-amount] [min-amounts]
frogchaind tx amm swap-exact-tokens-for-tokens [pool-id] [amount-in] [amount-out-min] [path] [to] "2023-08-02 15:04:05Z07:00"
frogchaind tx amm swap-tokens-for-exact-tokens [pool-id] [amount-out] [path] [to] "2023-08-02 15:04:05Z07:00"