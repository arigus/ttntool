// Copyright © 2016 Hylke Visser
// MIT Licensed - See LICENSE file

package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/htdvisser/ttntool/cmd"
)

func main() {
	cli.Colors[log.DebugLevel] = 90
	cli.Colors[log.InfoLevel] = 32
	log.SetHandler(cli.New(os.Stdout))

	log.SetLevel(log.DebugLevel)

	cmd.Execute()
}
