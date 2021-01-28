package environment

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// NewFromYAML loads configuration written on YAML file into the "envOutput" parameter
func NewFromYAML(envConfigPath string, envName string, envOutput interface{}) error {
	configFilePath := fmt.Sprintf("%s/env.%s.yaml", envConfigPath, envName)
	environmentConfigPath, _ := filepath.Abs(configFilePath)
	environmentConfig, err := ioutil.ReadFile(filepath.Clean(environmentConfigPath))
	if err != nil {
		return err
	}

	insensitiveUnmarshal := caseInsensitiveUnmarshal(environmentConfig)
	err = yaml.Unmarshal(insensitiveUnmarshal, envOutput)
	if err != nil {
		return err
	}

	// no errors and environment specific configuration was loaded correctly
	return nil
}

// caseInsensitiveUnmarshal transforms all keys of yaml to lowercase
// since yaml.Unmarshal doesn't support this feature yet
func caseInsensitiveUnmarshal(in []byte) []byte {
	// var lines []string
	lines := make([]string, 0, len(in))
	for _, line := range strings.Split(string(in), "\n") {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			line = fmt.Sprintf("%s:%s", strings.ToLower(parts[0]), parts[1])
		}
		lines = append(lines, line)
	}
	return []byte(strings.Join(lines, "\n"))
}
