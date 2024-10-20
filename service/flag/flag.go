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
	// PersistentFlags 的意思是他的子命令也可以用
	// Flags 是只有当前命你可以用
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
