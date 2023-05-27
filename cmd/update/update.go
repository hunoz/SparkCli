package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fatih/color"
)

var CmdVersion = "v1.0.3"

type Release struct {
	Url       string `json:"url,omitempty"`
	AssetsUrl string `json:"assets_url,omitempty"`
	UploadUrl string `json:"upload_url,omitempty"`
	TagName   string `json:"tag_name,omitempty"`
}

func CmdIsLatestVersion() (*string, bool) {
	currentVersion := strings.Split(CmdVersion, "v")[1]
	response, err := http.Get("https://api.github.com/repos/hunoz/SparkCli/releases/latest")
	if err != nil {
		color.Red("Error fetching latest release: %v", err.Error())
		os.Exit(1)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	release := Release{}
	err = json.Unmarshal(body, &release)
	if err != nil {
		color.Red("Error reading latest release: %v", err.Error())
		os.Exit(1)
	}

	latestVersion := strings.Split(release.TagName, "v")[1]

	if latestVersion == currentVersion || latestVersion < currentVersion {
		return &latestVersion, true
	}

	return &latestVersion, false
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update CLI if an update is available",
	Run: func(cmd *cobra.Command, args []string) {
		latestVersion, isLatestVersion := CmdIsLatestVersion()
		if isLatestVersion {
			color.Green("Spark is already running at the latest version")
			os.Exit(0)
		}

		path, _ := os.Executable()
		out, err := os.Create(path)
		if err != nil {
			color.Red("Error creating file at %v: %v", path, err.Error())
			os.Exit(1)
		}

		// Here we need to get the asset for the current architecture
		assetFilename := fmt.Sprintf("spark-%v-%v", runtime.GOOS, runtime.GOARCH)
		downloadUrl := fmt.Sprintf("https://github.com/hunoz/SparkCli/releases/download/%v/%v", latestVersion, assetFilename)

		response, err := http.Get(downloadUrl)
		if err != nil {
			color.Red("Error downloading latest version: %v", err.Error())
			os.Exit(1)
		}

		defer response.Body.Close()

		_, err = io.Copy(out, response.Body)
		if err != nil {
			color.Red("Error updating to latest version: %v", err.Error())
		}

		color.Green("Successfully updated Spark to version %v", latestVersion)
	},
}
