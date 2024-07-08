package main

import (
	"fmt"
	"os"

	"github.com/NishantJoshi00/yamlink"
)

func main() {
	config_file, err1 := os.LookupEnv("CONFIG_FILE")

	if !err1 {
		panic("Please provide a config file")
	}

	if len(os.Args) < 2 {
		panic("Please provide a query")
	}

	if len(os.Args) >= 3 {
		panic("Please provide only one query")
	}

	query := os.Args[1]

	loaded_config, err := yamlink.ReadFile(config_file)

	if err != nil {
		panic(err)
	}

	resolve, err := yamlink.PathLookup(query, loaded_config)

	if err != nil {
		panic(err)
	}

	fmt.Fprint(os.Stdout, resolve)
}
