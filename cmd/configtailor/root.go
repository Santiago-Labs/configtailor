package configtailor

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "configtailor",
	Short: "configtailor - a simple CLI to combine hierarchical configs with environments",
	Long: `configtailor is a CLI tool to help manage hierarchical configs with different cells and environments in JSON.
   
For example, if you have a config directory structure like this:
	config/
	├── region
	│   ├── cell
	│   │   └── config.json
	│   └── config.json

Where config/region/cell/config.json contains:
	{
		"region": "$region",
		"cell": "$cell"
	}

When Compiling for:
	- region: us-west-2, us-east-1
	- cell: us000, us001 

We will generate the following directory structure:
	generated
	├── us-west-2
	│   ├── us000 
	│   │   └── config.json
	│   ├── us001
	│   │   └── config.json
	├── us-east-1
	│   ├── us000
	│   │   └── config.json
	│   ├── us001
	│   │   └── config.json

With the following contents for "us-west-2/us000/config.json":
	{
		"region": "us-west-2",
		"cell": "us000"
	}
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "Please pass in a command. See more with -h\n")
		os.Exit(1)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
