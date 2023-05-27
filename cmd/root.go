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
	"gtech.dev/spark/cmd/update"
)

var RootCmd = &cobra.Command{
	Use:   "spark",
	Short: "Utilities for interacting with Cognito",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := cmd.Flags().GetBool("version")
		if err != nil {
			return
		}

		if version {
			fmt.Println(update.CmdVersion)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	RootCmd.Flags().BoolP("version", "v", false, "Current version of Spark")
	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(changepassword.ChangePasswordCmd)
	RootCmd.AddCommand(registertotp.RegisterTotpCmd)
	RootCmd.AddCommand(resetpassword.ResetPasswordCmd)
	RootCmd.AddCommand(firstsignin.FirstSignInCmd)
	RootCmd.AddCommand(cmdInit.InitCmd)
	RootCmd.AddCommand(update.UpdateCmd)
}
