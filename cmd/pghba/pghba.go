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

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("current user couldn't be detected")
	}

	rootCmd.PersistentFlags().StringP("cfgFile", "c", "", "config file (default is $HOME/.pghba.yaml)")
	viper.BindPFlag("cfgFile", rootCmd.PersistentFlags().Lookup("cfgFile"))
	viper.SetDefault("cfgFile", filepath.Join(currentUser.HomeDir, ".pghba.yaml"))
	viper.BindEnv("cfgFile", "PGHBACFG")
	viper.AddConfigPath(viper.GetString("cfgFile"))
	err = viper.ReadInConfig() // Find and read the config file
	if err == nil { // Handle errors reading the config file
		fmt.Printf("pghba is reading config from this config file: %s", viper.ConfigFileUsed())
	}

	rootCmd.PersistentFlags().StringP("hbaFile", "f", "", "pg_hba.conf file (default is $PGDATA/pg_hba.conf)")
	viper.BindPFlag("hbaFile", rootCmd.PersistentFlags().Lookup("hbaFile"))
	viper.SetDefault("hbaFile", filepath.Join(pgData, "pg_hba.conf"))

	rootCmd.PersistentFlags().StringP("connection_type", "t", "",
		`connection type that the rule applies to. 
Defaults might be derived from PGHOST where a PGHOST starting wth / would default to local and anything else would default to host.
See --help and documentation for options on globbing, and regular expressions.`)
	viper.BindPFlag("connection_type", rootCmd.PersistentFlags().Lookup("connection_type"))
	log.Warn("we need to be smart about converting PGHOST into a connection_type in the code")
	viper.BindEnv("connection_type", "PGHOST")
	viper.SetDefault("connection_type", "local")

	rootCmd.PersistentFlags().StringP("database", "d", "",
		`database(s) that the rule applies to. 
Defaults to env var PGDATABASE, or PGUSER, or the linux user itself.
See --help and documentation for options on globbing, and regular expressions.`)
	viper.BindPFlag("pgdatabase", rootCmd.PersistentFlags().Lookup("database"))
	viper.BindEnv("pgdatabase", "PGDATABASE", "PGUSER")
	viper.SetDefault("pgdatabase", currentUser.Username)

	rootCmd.PersistentFlags().StringP("user", "U", "",
		`user(s) that the rule applies to. 
Defaults to env var PGUSER, or the linux user itself.
See --help and documentation for options on globbing, and regular expressions.`)
	viper.BindPFlag("pguser", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindEnv("pguser")
	viper.SetDefault("pguser", currentUser)

	rootCmd.PersistentFlags().StringP("source", "s", "",
		`source(s) that the rule applies to. 
Defaults might be derived from PGHOST (like localhost, 127.0.0.1, or ::1).
See --help and documentation for options on globbing, and regular expressions.`)
	viper.BindPFlag("source", rootCmd.PersistentFlags().Lookup("source"))
	// We need to be smart about this in the code
	log.Warn("we need to be smart about converting PGHOST into a source in the code")
	viper.BindEnv("source", "PGHOST")
	viper.SetDefault("source", "localhost")

	rootCmd.PersistentFlags().StringP("mask", "m", "",
		`source mask that the rule applies to. 
Usually left empty. For IP addresses it defaults to a CIDR for one IP.`)
	viper.BindPFlag("mask", rootCmd.PersistentFlags().Lookup("mask"))
	viper.BindEnv("mask", "PGHBAMASK")
	viper.SetDefault("mask", "")

	rootCmd.PersistentFlags().StringP("authmethod", "a", "",
		`Authentication method that the new rule should have. 
Defaults to scram-sha256.`)
	viper.BindPFlag("authmethod", rootCmd.PersistentFlags().Lookup("authmethod"))
	viper.BindEnv("authmethod", "PGHBAMETHOD")
	viper.SetDefault("authmethod", "scram-sha256")

	rootCmd.PersistentFlags().StringP("authoptions", "o", "",
		`Authentication options that the new rule should have. 
Defaults to having no options.`)
	viper.BindPFlag("authoptions", rootCmd.PersistentFlags().Lookup("authoptions"))
	viper.BindEnv("authoptions", "PGHBAOPTIONS")
	viper.SetDefault("authoptions", "")

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