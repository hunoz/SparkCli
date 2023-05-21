package resetpassword

import (
	"github.com/spf13/cobra"
	"gtech.dev/spark/cognito"
)

var ResetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset your password",
	Run: func(cmd *cobra.Command, args []string) {
		cognitoClient := cognito.New()

		cognitoClient.ResetPassword()
	},
}
