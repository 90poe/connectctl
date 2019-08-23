## connectctl connectors resume

Resume connectors

### Synopsis


Resume connectors in a specified Kafka Connect cluster.
It can resume one or more connectors in a single execution.


```
connectctl connectors resume [flags]
```

### Options

```
  -h, --help        help for add
  -c, --clusterURL  the url of the kafka connect cluster to resume connectors in
  -n, --connectors  the names of the connectors to resume. Multiple connector names 
                    can be specified either by comma separating conn1,conn2
                    or by repeating the flag --n conn1 --n conn2. If no name is
                    supplied then ALL connectors will be resumed.
```
### Options inherited from parent commands

```
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl connectors](connectctl_connectors.md)	 - Manage connectors