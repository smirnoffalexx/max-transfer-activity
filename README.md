# max-transfer-activity
There is a tracker of maximum transfer activity (send or receive ERC20 tokens) on Ethereum mainnet in 100 latest blocks.

For running service add GET_BLOCK_URL to .env file as in example.env and then run in the terminal:
```
go build
./max-transfer-activity
```

Top 5 addresses with their activity are presented in terminal logs.
