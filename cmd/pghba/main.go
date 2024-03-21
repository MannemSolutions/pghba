package main

// cobra and viper are used to create a uniform interface on CLI and configuration file.
import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/mannemsolutions/pghba/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

// requireSubcommand returns an error if no sub command is provided
// This was copied from skopeo, which copied it from podman: `github.com/containers/podman/cmd/podman/validate/args.go
// Some small style changes to match skopeo were applied, but try to apply any
// bugfixes there first.
func requireSubcommand(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		suggestions := cmd.SuggestionsFor(args[0])
		if len(suggestions) == 0 {
			return fmt.Errorf("unrecognized command `%[1]s %[2]s`\nTry '%[1]s --help' for more information", cmd.CommandPath(), args[0])
		}
		return fmt.Errorf("unrecognized command `%[1]s %[2]s`\n\nDid you mean this?\n\t%[3]s\n\nTry '%[1]s --help' for more information", cmd.CommandPath(), args[0], strings.Join(suggestions, "\n\t"))
	}
	return fmt.Errorf("missing command '%[1]s COMMAND'\nTry '%[1]s --help' for more information", cmd.CommandPath())
}

// This function returns either a validly formed command for main() to run, or
// an error. Initializes a cobra command structure using the settings from the
// configuration file. Override the default location with -c,--cfgFile).
// Override the target pg_hba.conf file with -f, --hbaFile
func createApp() *cobra.Command {

	cobra.OnInitialize(initConfig)
	pgData := os.Getenv("PGDATA")

	rootCmd := &cobra.Command{
		Use:   "pghba",
		Short: "A tool to manage PostgreSQL pg_hba.conf files",
		Long: `The pghba tool delivers admins with a tool to manage pg_hba files with a bit more sophistication.
Complete documentation is available at https://github.com/mannemsolutions/pghba/`,
		RunE:              requireSubcommand,
		CompletionOptions: cobra.CompletionOptions{},
		TraverseChildren:  true,
		Version:           internal.AppVersion,
		//SilenceErrors: true,
		//SilenceUsage: true,
	}
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

	rootCmd.AddCommand(
		addCommand(),
		deleteCommand(),
	)
	return rootCmd
}

// Execute the fully formed pghba command and handle any errors.
func main() {
	initLogger()
	rootCmd := createApp()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Info("finished")
}

// Read settings as key value pairs from the ".pghba" config file in the home directory.
// This is (obscurely) referenced from the "createApp" function above.
// TODO would this be clearer if moved above createApp?
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
