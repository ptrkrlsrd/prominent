// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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

	prominent "github.com/ptrkrlsrd/prominent/pkg"
	"github.com/spf13/cobra"
)

var port string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		prominent.Serve(fmt.Sprintf(":%s", port))
	},
}

func init() {
	serveCmd.Flags().StringVarP(&port, "port", "p", ":3000", "Port to use")
	rootCmd.AddCommand(serveCmd)
}
