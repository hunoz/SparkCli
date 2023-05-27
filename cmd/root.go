package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gtech.dev/spark/cmd/auth"
	changepassword "gtech.dev/spark/cmd/change-password"
	firstsignin "gtech.dev/spark/cmd/first-sign-in"
	cmdInit "gtech.dev/spark/cmd/init"
	registertotp "gtech.dev/spark/cmd/register-totp"
	resetpassword "gtech.dev/spark/cmd/reset-password"
)

var cmdVersion = "1.0.1"

var RootCmd = &cobra.Command{
	Use:   "spark",
	Short: "Utilities for interacting with Cognito",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := cmd.Flags().GetBool("version")
		if err != nil {
			return
		}

		if version {
			fmt.Println(cmdVersion)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	RootCmd.Flags().Bool("version", false, "Current version of Spark")
	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(changepassword.ChangePasswordCmd)
	RootCmd.AddCommand(registertotp.RegisterTotpCmd)
	RootCmd.AddCommand(resetpassword.ResetPasswordCmd)
	RootCmd.AddCommand(firstsignin.FirstSignInCmd)
	RootCmd.AddCommand(cmdInit.InitCmd)
}
