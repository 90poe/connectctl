## connectctl plugins validate

Validates connector config

### Synopsis

Validate the provided configuration values against the configuration definition. This API performs per config validation, outputs suggested values and error messages during validation.
It exits with code 1 if config is invalid.


```
connectctl plugins validate [flags]
```

### Options

```
  -c, --cluster string   the URL of the connect cluster (required)
  -h, --help             help for validate
  -i, --input string     Input data in json format (required)
  -o, --output string    specify the output format (valid options: json, table) (default "json")
  -q, --quiet            disable output logging
```
### Options inherited from parent commands

```
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl plugins](connectctl_plugins.md)	 - Manage plugins