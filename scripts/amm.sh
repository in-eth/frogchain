frogchaind query amm list-pool
frogchaind query amm show-pool 1
frogchaind query amm get-pool-param 1
frogchaind query amm get-pool-share-token 1
frogchaind query amm get-pool-assets 1
frogchaind query amm get-swap-exact-tokens-for-tokens 1 500 'foocoin,token'
frogchaind query amm get-swap-tokens-for-exact-tokens 1 500 'foocoin,token'

frogchaind tx amm create-pool '{"SwapFee":"10","ExitFee":"10","FeeCollector":"cosmos1g3z20q5jskz3g2anvs5hnxpn7tsa7dvylv34yd"}' '30000foocoin,300token' '1,1' --from Alice
frogchaind tx amm add-liquidity 0 '3000,30' '3000,30' --from Bob
frogchaind tx amm remove-liquidity 0 "100" '1000,10' --from Bob
frogchaind tx amm swap-exact-tokens-for-tokens 0 "5000" "5" 'foocoin,token' "cosmos1rsl3u6e0l8923m9hgnaxuqfpcu05xa9emh3dl2" "2024-01-02 15:04:05 -07:00" --from Alice
frogchaind tx amm swap-tokens-for-exact-tokens 0 "100" 'foocoin,token' "cosmos1rsl3u6e0l8923m9hgnaxuqfpcu05xa9emh3dl2" "2024-01-02 15:04:05 -07:00" --from Alice
