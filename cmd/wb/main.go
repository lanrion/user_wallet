package main

import (
	"flag"
	"user_wallet/pkg/wb"
)

var userSharding = flag.Int("userSharding", 0, "index of user")

func main() {
	flag.Parse()
	ss := wb.NewWbService(int32(*userSharding))

	ss.Consume()
}
