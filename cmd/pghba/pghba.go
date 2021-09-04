package pghba

import (
	"fmt"
	"github.com/mannemsolutions/pghba/internal"
	"github.com/mannemsolutions/pghba/pkg/hba"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	Cmd = &cobra.Command{
		Use:     "pghba",
		Short:   "Manage pg_hba.conf files the right way",
		Version: internal.AppVersion,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := internal.AssertRequiredSettingsSet()
			tracelog.ErrorLogger.FatalOnError(err)
			if viper.IsSet(internal.PgWalSize) {
				postgres.SetWalSize(viper.GetUint64(internal.PgWalSize))
			}
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the PgCmd.
func Execute() {
	configureCommand()
	if err := Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func configureCommand() {
	cobra.OnInitialize(internal.InitConfig, internal.Configure)

	Cmd.PersistentFlags().StringVar(&internal.CfgFile, "config", "", "config file (default is $HOME/.walg.json)")
	Cmd.PersistentFlags().BoolVarP(&internal.Turbo, "turbo", "", false, "Ignore all kinds of throttling defined in config")
	Cmd.InitDefaultVersionFlag()
	internal.AddConfigFlags(Cmd)

	// Storage tools
	Cmd.AddCommand(st.StorageToolsCmd)
}
