package main

import (
	"github.com/auribuo/novasearch/cmd"
	"github.com/spf13/cobra"
)

// @title           Novasearch API
// @version         1.0
// @description     Api specification for Novasearch API

// @contact.name   auribuo
// @contact.url    https://github.com/auribuo/novasearch/issues

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api
func main() {
	cobra.CheckErr(cmd.Execute())
}
