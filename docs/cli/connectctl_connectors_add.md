## connectctl connectors add

Add a connector

### Synopsis


Creates a new connector based on a definition ina  cluster.
It can create one or more connectors in a single execution.


```
connectctl connectors add [flags]
```

### Options

```
  -h, --help   help for add
  -c, --clusterURL  the url of the kafka connect cluster to create the connectors in
  -f, --files       the json file containing the connector definition. Multiple files can be specified
                    either by comma separating file1.json,file2.json or by repeating the flag.
  -d, --directory   a director that contains json files with the connector definitions to add
```

NOTE: the -d and -f options are mutually exclusive. 

### Options inherited from parent commands

```
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl connectors](connectctl_connectors.md)	 - Manage connectors