package registertotp

import (
	"github.com/spf13/cobra"
	"gtech.dev/spark/cognito"
)

var RegisterTotpCmd = &cobra.Command{
	Use:   "register-totp",
	Short: "Register TOTP device",
	Run: func(cmd *cobra.Command, args []string) {
		cognitoClient := cognito.New()

		cognitoClient.RegisterMfaDevice()
	},
}
