package output

import (
	"encoding/json"
	"fmt"

	"github.com/cilginc/gochecker/internal/errors"
	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
	"github.com/goccy/go-yaml"
)

// RenderResults outputs the results in the specified format.
func RenderResults(format string, results []pkg.Result) error {
	switch format {
	case "json":
		return renderJSON(results)
	case "yaml":
		return renderYAML(results)
	case "text":
		return nil
	default:
		return errors.ErrInvalidOutputType
	}
}

// RenderPackages outputs the packages in the specified format.
func RenderPackages(format string, packages []pkg.Package) error {
	switch format {
	case "json":
		data, err := json.MarshalIndent(packages, "", "  ")
		if err != nil {
			return ui.CliError("failed to marshal JSON: %s", err)
		}
		fmt.Println(string(data))
		return nil
	case "yaml":
		data, err := yaml.Marshal(packages)
		if err != nil {
			return ui.CliError("failed to marshal YAML: %s", err)
		}
		fmt.Print(string(data))
		return nil
	case "text":
		return nil
	default:
		return errors.ErrInvalidOutputType
	}
}

// prepareResults populates ErrorMsg from Error for serialization.
func prepareResults(results []pkg.Result) []pkg.Result {
	prepared := make([]pkg.Result, len(results))
	copy(prepared, results)
	for i := range prepared {
		if prepared[i].Error != nil {
			prepared[i].ErrorMsg = prepared[i].Error.Error()
		}
	}
	return prepared
}

func renderJSON(results []pkg.Result) error {
	data, err := json.MarshalIndent(prepareResults(results), "", "  ")
	if err != nil {
		return ui.CliError("failed to marshal JSON: %s", err)
	}
	fmt.Println(string(data))
	return nil
}

func renderYAML(results []pkg.Result) error {
	data, err := yaml.Marshal(prepareResults(results))
	if err != nil {
		return ui.CliError("failed to marshal YAML: %s", err)
	}
	fmt.Print(string(data))
	return nil
}
