package cmd

import (
	"github.com/auribuo/novasearch/api"
	"github.com/auribuo/novasearch/data"
	"github.com/auribuo/novasearch/log"
	"github.com/spf13/cobra"
	"time"
)

var host string
var port int
var logLevel string
var debug bool

var rootCmd = &cobra.Command{
	Use:   "novasearch",
	Short: "A library to calculate the optimal path to follow to find the most supernovae",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(log.Setup(logLevel))
		initStart := time.Now()
		cobra.CheckErr(data.Init(data.DefaultFetcher))
		log.Logger.Infof("repo initialized in %dms. starting api...", time.Since(initStart).Milliseconds())
		api.Start(host, port, debug)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&host, "host", "H", "localhost", "The host to listen on")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port to listen on")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "The log level to use")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug mode")
	err := rootCmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"debug", "info", "warn", "error"}, cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)
}
