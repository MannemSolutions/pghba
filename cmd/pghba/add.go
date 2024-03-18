package main

import (
	"os/user"

	"github.com/mannemsolutions/pghba/pkg/hba"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

func addCommand() *cobra.Command {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("current user couldn't be detected")
	}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add one or more rules to pg_hba.conf",
		Long: `Use the add command to add or replace rules in the pg_hba.conf file.
           Note that existing rules with the same identifiers will be replaced.
           Also not that the new rule will automatically be added in the correct location (based on its scope).`,
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetInt("verbose") > 0 {
				atom.SetLevel(zapcore.DebugLevel)
				log.Info("Debug logging enabled")
			}
			fileName := viper.GetString("hbaFile")
			hbaFile := hba.NewFile(fileName)
			err := hbaFile.Read()
			if err != nil {
				log.Fatalf("error reading file %s: %e", fileName, err)
			}
			rules := argsToRules("add")
			if hbaFile.AddRules(*rules, true) {
				if err = hbaFile.Save(false); err != nil {
					log.Fatalf("error saving file %s: %e", fileName, err)
				}
				log.Info("Rules successfully inserted")
			}
		},
	}
	addCmd.PersistentFlags().StringP("connection_type", "t", "",
		`connection type that the rule applies to. 
Defaults might be derived from PGHOST where a PGHOST starting wth / would default to local and anything else would default to host.
See --help and documentation for options on globbing, and regular expressions.`)
	log.Warn("we need to be smart about converting PGHOST into a connection_type in the code")
	bindArgument("add", "connection_type", addCmd, []string{""}, "local")

	addCmd.PersistentFlags().StringP("database", "d", "",
		`database(s) that the rule applies to. 
Defaults to env var PGDATABASE, or PGUSER, or the linux user itself.
See --help and documentation for options on globbing, and regular expressions.`)
	bindArgument("add", "database", addCmd, []string{"PGDATABASE", "PGUSER"}, currentUser.Username)

	addCmd.PersistentFlags().StringP("user", "U", "",
		`user(s) that the rule applies to. 
Defaults to env var PGUSER, or the linux user itself.
See --help and documentation for options on globbing, and regular expressions.`)
	bindArgument("add", "user", addCmd, []string{"PGUSER"}, currentUser.Username)

	addCmd.PersistentFlags().StringP("source", "s", "",
		`source(s) that the rule applies to. 
Defaults might be derived from PGHOST (like localhost, 127.0.0.1, or ::1).
See --help and documentation for options on globbing, and regular expressions.`)
	log.Warn("we need to be smart about converting PGHOST into a source in the code")
	bindArgument("add", "source", addCmd, []string{"PGHOST"}, "localhost")

	addCmd.PersistentFlags().StringP("mask", "m", "",
		`source mask that the rule applies to. 
Usually left empty. For IP addresses it defaults to a CIDR for one IP.`)
	bindArgument("add", "mask", addCmd, []string{"PGHBAMASK"}, "")

	addCmd.PersistentFlags().StringP("authMethod", "a", "",
		`Authentication method that the new rule should have. 
Defaults to scram-sha256.`)
	bindArgument("add", "authMethod", addCmd, []string{"PGHBAMETHOD"}, "scram-sha-256")

	addCmd.PersistentFlags().StringP("authOptions", "o", "",
		`Authentication options that the new rule should have.`)
	bindArgument("add", "authOptions", addCmd, []string{"PGHBAOPTIONS"}, "")

	addCmd.PersistentFlags().StringP("line", "l", "",
		`line number to add before. 0 means before first, -1 means at the end, 
and 'auto' will prepend before the rule with a larger span.`)
	bindArgument("add", "line", addCmd, []string{"PGHBALINE"}, "0")

	return addCmd
}
