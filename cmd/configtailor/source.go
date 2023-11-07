package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/santiago-labs/configtailor/internal"

	"github.com/spf13/cobra"
)

var sourcePath string

var sourceCmd = &cobra.Command{
	Use:   "exec",
	Short: "execs to the generated config output",
	Run: func(cmd *cobra.Command, args []string) {
		env := os.Environ()
		envMap := make(map[string]string, len(env))
		for _, e := range env {
			envMap[e] = strings.Split(e, "=")[0]
		}

		fileBytes, err := ioutil.ReadFile(sourcePath)
		if err != nil {
			panic(err)
		}

		config := map[string]interface{}{}
		if err := json.Unmarshal(fileBytes, &config); err != nil {
			panic(err)
		}

		flattenedConfig := internal.Flatten(config, "", "")
		for k, v := range flattenedConfig {
			titleK := strings.ToUpper(k)
			if envMap[titleK] != "" {
				fmt.Println("overriding", titleK)
			}
			os.Setenv(titleK, v)
		}
	},
}

func init() {
	rootCmd.AddCommand(sourceCmd)
	sourceCmd.Flags().StringVar(&sourcePath, "path", "", "Root path to search for config files")
	sourceCmd.MarkFlagRequired("path")
}
