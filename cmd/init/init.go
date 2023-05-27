package init

import (
	"os"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gtech.dev/spark/config"
)

func getClientIdFromAllOptions() string {
	var clientId string
	clientIdFromCli := viper.GetString(FlagKey.ClientId)
	if clientIdFromCli == "" {
		clientIdPrompt := promptui.Prompt{
			Label: "Client ID",
		}
		clientId, _ = clientIdPrompt.Run()
	} else {
		clientId = clientIdFromCli
	}
	if clientId == "" {
		color.Red("Client ID cannot be empty")
		os.Exit(1)
	}
	return clientId
}

func getRegionFromAllOptions() string {
	var region string
	regionFromCli := viper.GetString(FlagKey.Region)
	if regionFromCli == "" {
		regionPrompt := promptui.Prompt{
			Label:    "Region",
			Validate: config.IsValidAwsRegion,
		}
		region, _ = regionPrompt.Run()
	} else {
		region = regionFromCli
	}
	if region == "" {
		color.Red("Region cannot be empty")
		os.Exit(1)
	} else if err := config.IsValidAwsRegion(region); err != nil {
		color.Red(err.Error())
	}
	return region
}

func getPoolIdFromAllOptions() string {
	var poolId string
	poolIdFromCli := viper.GetString(FlagKey.Region)
	if poolIdFromCli == "" {
		poolIdPrompt := promptui.Prompt{
			Label: "Pool ID",
		}
		poolId, _ = poolIdPrompt.Run()
	} else {
		poolId = poolIdFromCli
	}
	if poolId == "" {
		color.Red("Pool ID cannot be empty")
		os.Exit(1)
	}
	return poolId
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialized the spark CLI's configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag(FlagKey.ClientId, cmd.Flags().Lookup(FlagKey.ClientId))
		viper.BindPFlag(FlagKey.PoolId, cmd.Flags().Lookup(FlagKey.PoolId))
		viper.BindPFlag(FlagKey.Overwrite, cmd.Flags().Lookup(FlagKey.Overwrite))
		viper.BindPFlag(FlagKey.Region, cmd.Flags().Lookup(FlagKey.Region))
	},
	Run: func(cmd *cobra.Command, args []string) {
		overwrite := viper.GetBool(FlagKey.Overwrite)
		if configuration, _ := config.CognitoIsInitialized(); configuration == nil {
			clientId := getClientIdFromAllOptions()
			poolId := getPoolIdFromAllOptions()
			region := getRegionFromAllOptions()

			if err := config.UpdateCognitoConfig(config.CognitoConfig{
				ClientId: clientId,
				Region:   region,
				PoolId:   poolId,
			}); err != nil {
				color.Red("Error updating config: %v", err.Error())
				os.Exit(1)
			}
		} else {
			var clientId string
			var region string
			var poolId string
			if overwrite {
				clientId = getClientIdFromAllOptions()
				region = getRegionFromAllOptions()
				poolId = getPoolIdFromAllOptions()
			} else {
				if configuration.ClientId != "" && !overwrite {
					color.Cyan("Client ID already configured")
					clientId = configuration.ClientId
				} else if configuration.ClientId == "" {
					clientId = getClientIdFromAllOptions()
				}

				if configuration.Region != "" && !overwrite {
					color.Cyan("Region already configured")
					region = configuration.Region
				} else if configuration.Region == "" {
					region = getRegionFromAllOptions()
				}

				if configuration.PoolId != "" && !overwrite {
					color.Cyan("Pool ID already configured")
					poolId = configuration.PoolId
				} else if configuration.PoolId == "" {
					poolId = getPoolIdFromAllOptions()
				}
			}

			if err := config.UpdateCognitoConfig(config.CognitoConfig{
				ClientId:    clientId,
				Region:      region,
				PoolId:      poolId,
				AccessToken: configuration.AccessToken,
				IdToken:     configuration.IdToken,
				Session:     configuration.Session,
				Expires:     configuration.Expires,
			}); err != nil {
				color.Red("Error updating config: %v", err.Error())
				os.Exit(1)
			}
		}
	},
}

func init() {
	InitCmd.Flags().StringP(FlagKey.ClientId, string(FlagKey.ClientId[0]), "", "Client ID that Spark CLI will authenticate to Cognito with")
	InitCmd.Flags().BoolP(FlagKey.Overwrite, string(FlagKey.Overwrite[0]), false, "Overwrite current configuration")
	InitCmd.Flags().StringP(FlagKey.Region, string(FlagKey.PoolId[0]), "", "Region where the Cognito pool is hosted")
}
