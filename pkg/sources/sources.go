package sources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/90poe/connectctl/pkg/client/connect"
)

func Files(files []string) func() ([]connect.Connector, error) {
	return func() ([]connect.Connector, error) {
		connectors := make([]connect.Connector, len(files))

		for index, file := range files {
			bytes, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, errors.Wrapf(err, "reading connector json %s", file)
			}

			c, err := newConnectorFromBytes(bytes)

			if err != nil {
				return nil, errors.Wrap(err, "unmarshalling connector from bytes")
			}

			connectors[index] = c
		}

		return connectors, nil
	}
}

func Directory(dir string) func() ([]connect.Connector, error) {
	return func() ([]connect.Connector, error) {
		var files []string

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".json" {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, errors.Wrapf(err, "list connector files in directory %s", dir)
		}

		return Files(files)()
	}
}

func EnvVarValue(env string) func() ([]connect.Connector, error) {

	return func() ([]connect.Connector, error) {
		value, ok := os.LookupEnv(env)

		if !ok {
			return nil, fmt.Errorf("error resolving env var : %s", env)
		}

		value = strings.TrimSpace(value)

		if strings.HasPrefix(value, "[") {

			c, err := newConnectorsFromBytes([]byte(value))
			if err != nil {
				return nil, errors.Wrap(err, "unmarshalling connector from bytes")
			}

			return c, nil
		}

		c, err := newConnectorFromBytes([]byte(value))
		if err != nil {
			return nil, errors.Wrap(err, "unmarshalling connector from bytes")
		}

		return []connect.Connector{c}, nil
	}
}

func newConnectorFromBytes(bytes []byte) (connect.Connector, error) {
	c := connect.Connector{}
	err := json.Unmarshal(bytes, &c)
	return c, err
}

func newConnectorsFromBytes(bytes []byte) ([]connect.Connector, error) {
	c := []connect.Connector{}
	err := json.Unmarshal(bytes, &c)
	return c, err
}
