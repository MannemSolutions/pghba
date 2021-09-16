package pghba

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
	Use:   "pghba",
	Short: "A tool to manage PostgreSQL pg_hba.conf files",
	Long: `The pghba tool delivers admins with a tool to manage pg_hba files with a bit more sophistication.
                Complete documentation is available at https://github.com/mannemsolutions/pghba/`,
}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	pgData := os.Getenv("PGDATA")
	if pgData == "" {
		pgData = "./"
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "cfgFile", "", "config file (default is $HOME/.pghba.yaml)")
	rootCmd.PersistentFlags().StringP("hbaFile", "f", "", "pg_hba.conf file (default is $PGDATA/pg_hba.conf)")
	viper.BindPFlag("hbaFile", rootCmd.PersistentFlags().Lookup("hbaFile"))
	viper.SetDefault("hbaFile", filepath.Join(pgData, "pg_hba.conf"))

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pghba")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}