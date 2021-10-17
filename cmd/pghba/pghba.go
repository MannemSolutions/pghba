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

func bindArgument(key string, envVars []string, defaultValue string) {
	var err error
	err = viper.BindPFlag(key, rootCmd.PersistentFlags().Lookup(key))
	if err != nil {
		log.Fatalf("error while binding argument for %s: %e", key, err)
	}
	if len(envVars) > 0 {
		envVars = append([]string{key}, envVars...)
		err = viper.BindEnv(envVars...)
		if err != nil {
			log.Fatal("error while binding env var for %s: %e", key, err)
		}
	}
	viper.SetDefault(key, defaultValue)
}

func init() {
	cobra.OnInitialize(initConfig)
	pgData := os.Getenv("PGDATA")

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("current user couldn't be detected")
	}

	rootCmd.PersistentFlags().StringP("cfgFile", "c", "", "config file (default is $HOME/.pghba.yaml)")
	bindArgument("cfgFile", []string{"PGHBACFG"}, filepath.Join(currentUser.HomeDir, ".pghba.yaml"))	
	viper.AddConfigPath(viper.GetString("cfgFile"))
	err = viper.ReadInConfig() // Find and read the config file
	if err == nil { // Handle errors reading the config file
		fmt.Printf("pghba is reading config from this config file: %s", viper.ConfigFileUsed())
	}

	rootCmd.PersistentFlags().StringP("hbaFile", "f", "", "pg_hba.conf file (default is $PGDATA/pg_hba.conf)")
	bindArgument("hbaFile", []string{}, filepath.Join(pgData, "pg_hba.conf"))

	rootCmd.PersistentFlags().StringP("connection_type", "t", "",
		`connection type that the rule applies to. 
Defaults might be derived from PGHOST where a PGHOST starting wth / would default to local and anything else would default to host.
See --help and documentation for options on globbing, and regular expressions.`)
	log.Warn("we need to be smart about converting PGHOST into a connection_type in the code")
	bindArgument("connection_type", []string{"PGHOST"}, "local")

	rootCmd.PersistentFlags().StringP("database", "d", "",
		`database(s) that the rule applies to. 
Defaults to env var PGDATABASE, or PGUSER, or the linux user itself.
See --help and documentation for options on globbing, and regular expressions.`)
	bindArgument("database", []string{"PGDATABASE", "PGUSER"}, currentUser.Username)

	rootCmd.PersistentFlags().StringP("user", "U", "",
		`user(s) that the rule applies to. 
Defaults to env var PGUSER, or the linux user itself.
See --help and documentation for options on globbing, and regular expressions.`)
	bindArgument("user", []string{"PGUSER"}, currentUser.Username)

	rootCmd.PersistentFlags().StringP("source", "s", "",
		`source(s) that the rule applies to. 
Defaults might be derived from PGHOST (like localhost, 127.0.0.1, or ::1).
See --help and documentation for options on globbing, and regular expressions.`)
	log.Warn("we need to be smart about converting PGHOST into a source in the code")
	bindArgument("source", []string{"PGHOST"}, "localhost")

	rootCmd.PersistentFlags().StringP("mask", "m", "",
		`source mask that the rule applies to. 
Usually left empty. For IP addresses it defaults to a CIDR for one IP.`)
	bindArgument("mask", []string{"PGHBAMASK"}, "")

	rootCmd.PersistentFlags().StringP("authMethod", "a", "",
		`Authentication method that the new rule should have. 
Defaults to scram-sha256.`)
	bindArgument("authMethod", []string{"PGHBAMETHOD"}, "scram-sha-256")

	rootCmd.PersistentFlags().StringP("authOptions", "o", "",
		`Authentication options that the new rule should have.`)
	bindArgument("authOptions", []string{"PGHBAOPTIONS"}, "")

	rootCmd.PersistentFlags().StringP("line", "l", "",
		`line number to add before. 0 means before first, -1 means at the end, 
and 'auto' will prepend before the rule with a larger span.`)
	bindArgument("line", []string{"PGHBALINE"}, "0")

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