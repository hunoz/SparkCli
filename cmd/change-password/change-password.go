package changepassword

import (
	"github.com/spf13/cobra"
	"gtech.dev/spark/cognito"
)

var ChangePasswordCmd = &cobra.Command{
	Use:   "change-password",
	Short: "Change your password. Do not use if needing to reset password",
	Run: func(cmd *cobra.Command, args []string) {
		cognitoClient := cognito.New()

		cognitoClient.ChangePassword()
	},
}
