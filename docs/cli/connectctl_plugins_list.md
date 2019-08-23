## connectctl plugins list

List connector plugins

### Synopsis


Lists all the connector plugins installed on a given 
Kafka Connect cluster node.
The output can be formatted as JSON or a table.


```
connectctl plugins list [flags]
```

### Options

```
  -h, --help        help for add
  -c, --clusterURL  the url of the kafka connect cluster to list the plugins from
  -o, --output      specify the format of the list of plugins.  Valid options
                    are json and table. The default is json.

```
### Options inherited from parent commands

```
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl plugins](connectctl_plugins.md)	 - Manage plugins