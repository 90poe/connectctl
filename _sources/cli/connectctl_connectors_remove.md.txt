## connectctl connectors remove

Remove a connector

### Synopsis


Removes a named connector from a cluster.
It can remove one or more connectors in a single execution.


```
connectctl connectors remove [flags]
```

### Options

```
  -h, --help        help for add
  -c, --clusterURL  the url of the kafka connect cluster to remove the connectors from
  -n, --connectors  the names of the connectors to remove. Multiple connector names 
                    can be specified either by comma separating conn1,conn2
                    or by repeating the flag --n conn1 --n conn2.
```
### Options inherited from parent commands
## connectctl connectors remove

Remove a connector

### Synopsis


Removes a named connector from a cluster.
It can remove one or more connectors in a single execution.


```
connectctl connectors remove [flags]
```

### Options

```
  -h, --help        help for add
  -c, --clusterURL  the url of the kafka connect cluster to remove the connectors from
  -n, --connectors  the names of the connectors to remove. Multiple connector names 
                    can be specified either by comma separating conn1,conn2
                    or by repeating the flag --n conn1 --n conn2.
```
### Options inherited from parent commands

```
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl connectors](connectctl_connectors.md)	 - Manage connectors
```
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl connectors](connectctl_connectors.md)	 - Manage connectors