package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
	"github.com/spf13/pflag"
)

func main() {
	var templatePath string
	var valuesPath string
	var outputPath string

	// Define command-line flags
	pflag.StringVarP(&templatePath, "template", "t", "", "Path to the template YAML file")
	pflag.StringVarP(&valuesPath, "values", "v", "", "Path to the values YAML file")
	pflag.StringVarP(&outputPath, "output", "o", "", "Path to the output file. If not provided, the output will be displayed in the terminal.")

	// Parse command-line arguments
	pflag.Parse()

	if templatePath == "" || valuesPath == "" {
		fmt.Println("Usage: go run main.go --template=template.yaml --values=values.yaml [--output=output.yaml]")
		os.Exit(1)
	}

	templateData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error reading template YAML file: %v\n", err)
		os.Exit(1)
	}

	variablesData, err := ioutil.ReadFile(valuesPath)
	if err != nil {
		fmt.Printf("Error reading variables YAML file: %v\n", err)
		os.Exit(1)
	}

	var vars map[string]string
	if err := yaml.Unmarshal(variablesData, &vars); err != nil {
		fmt.Printf("Error unmarshaling variables YAML data: %v\n", err)
		os.Exit(1)
	}

	replacedYAML := replaceVariables(string(templateData), vars)

	// Determine whether to write to a file or display in the terminal
	if outputPath != "" {
		err = ioutil.WriteFile(outputPath, []byte(replacedYAML), 0644)
		if err != nil {
			fmt.Printf("Error writing to the output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Data has been written to %s\n", outputPath)
	} else {
		fmt.Println(replacedYAML)
	}
}

func replaceVariables(input string, vars map[string]string) string {
	re := regexp.MustCompile(`\${(.*?)}`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		key := strings.Trim(match, "${}")
		if val, found := vars[key]; found {
			return val
		}
		return match // Keep the original value if not found in vars
	})
}
