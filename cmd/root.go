package cmd

import (
	"fmt"

	"github.com/hunoz/spark/cmd/auth"
	changepassword "github.com/hunoz/spark/cmd/change-password"
	firstsignin "github.com/hunoz/spark/cmd/first-sign-in"
	cmdInit "github.com/hunoz/spark/cmd/init"
	"github.com/hunoz/spark/cmd/refresh"
	registertotp "github.com/hunoz/spark/cmd/register-totp"
	resetpassword "github.com/hunoz/spark/cmd/reset-password"
	"github.com/hunoz/spark/cmd/update"
	"github.com/spf13/cobra"
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
	RootCmd.AddCommand(
		auth.AuthCmd,
		changepassword.ChangePasswordCmd,
		registertotp.RegisterTotpCmd,
		resetpassword.ResetPasswordCmd,
		firstsignin.FirstSignInCmd,
		cmdInit.InitCmd,
		update.UpdateCmd,
		refresh.RefreshCmd,
	)
}
