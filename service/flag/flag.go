package flag

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

const (
	ConfigFile = "config-file"
)

func AddOption(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringP(ConfigFile, "c", "", "config file")

	rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		err := viper.BindPFlag(flag.Name, flag)
		if err != nil {
			panic(err)
		}
	})

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}
