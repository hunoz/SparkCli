package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gtech.dev/spark/cognito"
	"gtech.dev/spark/config"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate to Cognito and cache credentials",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag(FlagKey.Force, cmd.Flags().Lookup(FlagKey.Force))
	},
	Run: func(cmd *cobra.Command, args []string) {
		var configuration *config.CognitoConfig
		config.CheckIfCognitoIsInitialized()
		if config, e := config.GetCognitoConfig(); e != nil {
			if strings.Contains(e.Error(), "Invalid region") {
				fmt.Println("Spark has not been initialized. Please run 'spark init' to initialize Spark.")
			} else {
				fmt.Printf("Error getting config: %s\n", e.Error())
			}
			os.Exit(1)
		} else {
			configuration = config
		}
		passwordValidator := cognito.CheckIfValidPassword
		usernamePrompt := promptui.Prompt{
			Label: "Username",
		}
		username, err := usernamePrompt.Run()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		passwordPrompt := promptui.Prompt{
			Label:    "Password",
			Mask:     '*',
			Validate: passwordValidator,
		}
		password, err := passwordPrompt.Run()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		cognitoClient := cognito.New(configuration)

		cognitoClient.InitiateAuth(username, password, viper.GetBool(FlagKey.Force))
	},
}

func init() {
	AuthCmd.Flags().Bool(FlagKey.Force, false, "Force update session, even if expiration time is > 6 hours")
}
