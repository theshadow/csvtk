// Copyright © 2018 Xander Guzman <xander.guzman@xanderguzman.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version defines the version of the tool
var Version = "dev-build"

// BuildDate defines the date this build was built
var BuildDate = "dev-build"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of csvtk",
	Long:  "Renders the application version and build date.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (build: %s)\n", Version, BuildDate)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
