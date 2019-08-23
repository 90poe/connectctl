## connectctl connectors restart

Restart connectors

### Synopsis


Restart connectors in a specified Kafka Connect cluster.
It can restart one or more connectors in a single execution.


```
connectctl connectors restart [flags]
```

### Options

```
  -h, --help        help for add
  -c, --clusterURL  the url of the kafka connect cluster to restart connectors in
  -n, --connectors  the names of the connectors to restart. Multiple connector names 
                    can be specified either by comma separating conn1,conn2
                    or by repeating the flag --n conn1 --n conn2. If no name is
                    supplied then ALL connectors will be restarted.
```
### Options inherited from parent commands

```
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl connectors](connectctl_connectors.md)	 - Manage connectors