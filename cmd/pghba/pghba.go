package pghba

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/user"
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
		log.Fatal(err)
		os.Exit(1)
	}
}


func init() {
	initLogger()

	cobra.OnInitialize(initConfig)
	pgData := os.Getenv("PGDATA")

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("current user couldn't be detected")
	}

	rootCmd.PersistentFlags().CountP("verbose", "v",
		`Be more verbose in the output.`)
	bindArgument("", "verbose", rootCmd, []string{"PGHBAVERBOSE"}, "0")

	rootCmd.PersistentFlags().StringP("cfgFile", "c", "", "config file (default is $HOME/.pghba.yaml)")
	bindArgument("", "cfgFile", rootCmd, []string{"PGHBACFG"}, filepath.Join(currentUser.HomeDir, ".pghba.yaml"))
	viper.AddConfigPath(viper.GetString("cfgFile"))
	err = viper.ReadInConfig() // Find and read the config file
	if err == nil {            // Handle errors reading the config file
		fmt.Printf("pghba is reading config from this config file: %s", viper.ConfigFileUsed())
	}

	rootCmd.PersistentFlags().StringP("hbaFile", "f", "", "pg_hba.conf file (default is $PGDATA/pg_hba.conf)")
	bindArgument("", "hbaFile", rootCmd, []string{}, filepath.Join(pgData, "pg_hba.conf"))

	initAdd()
	initDelete()
	initVersion()
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
