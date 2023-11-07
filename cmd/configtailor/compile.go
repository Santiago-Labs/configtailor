package cmd

import (
	"fmt"

	"github.com/santiago-labs/configtailor/configtailor"

	"github.com/spf13/cobra"
)

var rootPath string
var mappingsFlag string

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "compiles to the generated config output",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running", rootPath, mappingsFlag)
		dirs, err := configtailor.ParentDirs(rootPath)
		if err != nil {
			panic(err)
		}

		mappings, err := configtailor.Mappings(mappingsFlag)
		if err != nil {
			panic(err)
		}

		t := configtailor.ConfigTailor{
			Mappings: mappings,
			Dirs:     dirs,
			RootPath: rootPath,
		}

		if err := t.Compile(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
	compileCmd.Flags().StringVar(&rootPath, "rootpath", "", "Root path to search for config files")
	compileCmd.MarkFlagRequired("rootpath")
	compileCmd.Flags().StringVar(&mappingsFlag, "mappings", "", "The mapping of keys to values to compile for e.g. (key:value1,value2:key2:value1,value2) e.g. (cell:us1,us2,us3,us4:env:prod,dev,test)")
}
