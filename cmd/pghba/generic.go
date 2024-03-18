package main

import (
	"fmt"

	"github.com/mannemsolutions/pghba/pkg/hba"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getArg(ns, key string) string {
	return viper.GetString(fmt.Sprintf("%s.%s", ns, key))
}

func argsToRules(namespace string) *hba.Rules {
	var rowNum int
	if namespace != "add" || viper.GetString("line") == "auto" {
		rowNum = -1
	} else {
		rowNum = viper.GetInt("line")
	}

	rules, err := hba.NewRules(
		rowNum, getArg(namespace, "connection_type"), getArg(namespace, "database"),
		getArg(namespace, "user"), getArg(namespace, "source"), getArg(namespace, "mask"),
		getArg(namespace, "authMethod"), getArg(namespace, "authOptions"))
	if err != nil {
		log.Fatalf("cannot parse arguments into rule: %e", err)
	}
	return rules
}
func bindArgument(ns string, key string, cmd *cobra.Command, envVars []string, defaultValue string) {
	var err error
	var viperKey string
	if ns == "" {
		viperKey = key
	} else {
		viperKey = fmt.Sprintf("%s.%s", ns, key)
	}
	err = viper.BindPFlag(viperKey, cmd.PersistentFlags().Lookup(key))
	if err != nil {
		log.Fatalf("error while binding argument for %s: %e", key, err)
	}
	if len(envVars) > 0 {
		envVars = append([]string{key}, envVars...)
		err = viper.BindEnv(envVars...)
		if err != nil {
			log.Fatal("error while binding env var for %s: %e", viperKey, err)
		}
	}
	viper.SetDefault(viperKey, defaultValue)
}
