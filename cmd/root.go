package cmd

import (
	"time"

	"github.com/auribuo/novasearch/api"
	"github.com/auribuo/novasearch/data"
	"github.com/auribuo/novasearch/log"
	"github.com/spf13/cobra"
)

var appConfig struct {
	Host    string
	Port    int
	Logging struct {
		Level   string
		NoColor bool
	}
	DevMode bool
}

var rootCmd = &cobra.Command{
	Use:   "novasearch",
	Short: "A library to calculate the optimal path to follow to find the most supernovae",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(log.Setup(appConfig.Logging.Level, appConfig.Logging.NoColor))
		if appConfig.DevMode {
			log.SetLogLevel(log.DebugLevel)
			log.Default.Warn("Development mode enabled")
		}
		initStart := time.Now()
		cobra.CheckErr(data.NewDefault().Init())
		log.Default.Infof("repo initialized in %dms. starting api...", time.Since(initStart).Milliseconds())
		api.Start(appConfig.Host, appConfig.Port, appConfig.DevMode)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&appConfig.Host, "host", "H", "localhost", "The host to listen on")
	rootCmd.Flags().IntVarP(&appConfig.Port, "port", "p", 8080, "The port to listen on")
	rootCmd.Flags().StringVarP(&appConfig.Logging.Level, "log-level", "l", "info", "The log level to use")
	rootCmd.PersistentFlags().BoolVar(&appConfig.Logging.NoColor, "no-color", false, "Disable color output")
	rootCmd.PersistentFlags().BoolVar(&appConfig.DevMode, "dev", false, "Enable development mode")
	err := rootCmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"debug", "info", "warn", "error"}, cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)
}
