package main

import (
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"io"
	"os"
	"strings"
)

var OperatorName string = "my-operator"

var ChartDir string = "./charts/"
var OutDir string = "./charts/my-operator/templates"
var Image string = "registry.example.com/org/my-operator:1.2.3"

func init() {
	// parse flags for operator name, chart dir and image
	chartDir := flag.String("chart-dir", ChartDir, "The directory where the chart is located")
	operatorName := flag.String("operator-name", OperatorName, "The name of the operator")
	image := flag.String("image", Image, "The image of the operator")

	help := false
	flag.BoolVar(&help, "help", false, "Prints this help message")
	flag.BoolVar(&help, "h", false, "Prints this help message")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if chartDir != nil {
		ChartDir = *chartDir
	}
	if operatorName != nil {
		OperatorName = *operatorName
	}
	if image != nil {
		Image = *image
	}

	// make sure the last char of the chart dir is a slash
	if !strings.HasSuffix(ChartDir, "/") {
		ChartDir = fmt.Sprintf("%s/", ChartDir)
	}

	OutDir = fmt.Sprintf("%s%s/templates", ChartDir, OperatorName)
}

func main() {
	// Read the Kubernetes manifest from standard input
	manifestBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	manifest := string(manifestBytes)

	// Split the manifest into individual resources
	resources := strings.Split(manifest, "---\n")

	// Loop through each resource and create a file for it
	for _, resource := range resources {
		// Skip over empty lines and '---' delimiters
		if len(strings.TrimSpace(resource)) == 0 {
			continue
		}

		// Unmarshal the resource into a map
		var resourceMap map[string]interface{}
		if err := yaml.Unmarshal([]byte(resource), &resourceMap); err != nil {
			fmt.Fprintf(os.Stderr, "error unmarshaling resource: %v\n", err)
			os.Exit(1)
		}

		// Extract the resource name and kind
		name := resourceMap["metadata"].(map[string]interface{})["name"].(string)
		name = strings.ReplaceAll(name, OperatorName+"-", "")
		kind := resourceMap["kind"].(string)

		// Create a file with the name "<NAME>-<KIND>.yaml" and write the resource to it
		filename := fmt.Sprintf("%s/%s-%s.yaml", OutDir, strings.ToLower(name), strings.ToLower(kind))
		resource = ReplaceDeploymentImageWithHelmValues(resource)
		err := os.WriteFile(filename, []byte(resource), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func ReplaceDeploymentImageWithHelmValues(resource string) string {
	return strings.ReplaceAll(resource, "image: "+Image, "image: {{ .Values.image.repository }}:{{ .Values.image.tag }}")
}
