package main

import (
	"fmt"
	"strings"

	"github.com/mannemsolutions/pghba/pkg/hba"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

func deleteCommand() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete one or more rules from pg_hba.conf",
		Long: `Use the delete command to remove rules from a pg_hba.conf file.
			   Note that the original rules will be preserved as good as possible,
						 including comments. The location in the file can be automatically added
						 before the line with a bigger span, or can be manually set.
						 Alternatively all lines can be automatically sorted.`,
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
			rules := argsToRules("delete")
			if hbaFile.DeleteRules(rules) {
				err = hbaFile.Save(false)
				if err != nil {
					log.Fatalf("error saving file %s: %e", fileName, err)
				}
				log.Info("Rule was successfully deleted and altered file is saved.")
			} else {
				log.Info("Rule was not found.")
			}
		},
	}
	defCT := fmt.Sprintf("(%s)", strings.Join(hba.AllConnTypes(), "|"))
	deleteCmd.PersistentFlags().StringP("connection_type", "t", "",
		`Connection type that the rule applies to. When not set apply to all connection types.`)
	bindArgument("delete", "connection_type", deleteCmd, []string{}, defCT)

	deleteCmd.PersistentFlags().StringP("database", "d", "",
		`Database(s) that the rule applies to. Defaults to 'all'.`)
	bindArgument("delete", "database", deleteCmd, []string{"PGDATABASE"}, "")

	deleteCmd.PersistentFlags().StringP("user", "U", "",
		`User(s) that the rule applies to. Defaults to 'all'.`)
	bindArgument("delete", "user", deleteCmd, []string{"PGUSER"}, "")

	deleteCmd.PersistentFlags().StringP("source", "s", "",
		`source(s) that the rule applies to. When not set apply to all source nets.`)
	bindArgument("delete", "source", deleteCmd, []string{"PGHOST"}, "")

	deleteCmd.PersistentFlags().StringP("mask", "m", "",
		`source mask that the rule applies to. Usually left empty. 
When not set apply to one IP.`)
	bindArgument("delete", "mask", deleteCmd, []string{"PGHBAMASK"}, "")

	deleteCmd.PersistentFlags().StringP("authMethod", "a", "",
		`Authentication method that the new rule should have. 
When not set apply to all methods.`)
	bindArgument("delete", "authMethod", deleteCmd, []string{"PGHBAMETHOD"}, "")

	return deleteCmd
}
