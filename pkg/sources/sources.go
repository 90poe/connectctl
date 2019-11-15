package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/90poe/connectctl/pkg/client/connect"
)

// Files returns the aggregrated connectors loaded from a set of filepaths or an error
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

// Directory returns the aggregrated connectors loaded from a directory and its children or an error
// Note - Files need to end with .json.
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

// EnvVarValue returns the connectors loaded from an environmental variable or an error
func EnvVarValue(env string) func() ([]connect.Connector, error) {
	return func() ([]connect.Connector, error) {
		value, ok := os.LookupEnv(env)

		if !ok {
			return nil, fmt.Errorf("error resolving env var : %s", env)
		}

		value = strings.TrimSpace(value)
		return processBytes([]byte(value))
	}
}

// StdIn returns the connectors piped via stdin or an error
func StdIn(in io.Reader) func() ([]connect.Connector, error) {
	return func() ([]connect.Connector, error) {
		data, err := ioutil.ReadAll(in)

		if err != nil {
			return nil, errors.Wrap(err, "error reading from StdIn")
		}
		return processBytes(data)
	}
}

func processBytes(data []byte) ([]connect.Connector, error) {
	if bytes.HasPrefix(data, []byte("[")) { // REVIEW : is there a better test for a JSON array?
		c, err := newConnectorsFromBytes(data)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling connectors from bytes")
		}
		return c, nil
	}

	c, err := newConnectorFromBytes(data)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling connector from bytes")
	}

	return []connect.Connector{c}, nil
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
