package main

import (
	"flag"
	"user_wallet/pkg/wallet"
)

var userSharding = flag.Int("userSharding", 0, "index of user")

func main() {
	flag.Parse()
	us := int32(*userSharding)
	wallet.Run(us)
}
