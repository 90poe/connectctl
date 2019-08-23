## connectctl connectors list

List the connectors

### Synopsis


Lists all the connectors in a given Kafka Connect cluster.
The output includes the connecor status and the format can be
specified.


```
connectctl connectors list [flags]
```

### Options

```
  -h, --help        help for add
  -c, --clusterURL  the url of the kafka connect cluster to remove the connectors from
  -o, --output      specify the format of the list of connectors.  Valid options
                    are json and table. The default is json.

```
### Options inherited from parent commands

```
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl connectors](connectctl_connectors.md)	 - Manage connectors