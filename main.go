package main

import (
	"github.com/auribuo/novasearch/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.Execute())
}
