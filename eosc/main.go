package main

import (
	// Load all contracts here, so we can always read and decode
	// transactions with those contracts.
	_ "github.com/zhongshuwen/zswchain-go/msig"
	_ "github.com/zhongshuwen/zswchain-go/system"
	_ "github.com/zhongshuwen/zswchain-go/token"

	"github.com/zhongshuwen/eosc/eosc/cmd"
)

var version = "dev"

func init() {
	cmd.Version = version
}

func main() {
	cmd.Execute()
}
