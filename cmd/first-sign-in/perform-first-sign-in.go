package firstsignin

import (
	"github.com/spf13/cobra"
	"gtech.dev/spark/cognito"
)

var FirstSignInCmd = &cobra.Command{
	Use:   "first-sign-in",
	Short: "Performs a first sign in",
	Run: func(cmd *cobra.Command, args []string) {
		cognitoClient := cognito.New()

		cognitoClient.PerformFirstSignIn()
	},
}
