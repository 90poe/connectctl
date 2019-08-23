## connectctl connectors manage

Actively manage connectors in a cluster

### Synopsis


This command will actively manage connectors in a cluster by
creating, updating, removing and restarting connectors.
The command can be run once or it can run continously where it
will sync the desired state and actually state on a periodic basis.


```
connectctl connectors manage [flags]
```



### Options

```
  -h, --help   help     for add
  -c, --clusterURL      the url of the kafka connect cluster to avtively manage
  -f, --files           the json file(s) containing the connector definition. Multiple files can be specified
                        either by comma separating file1.json,file2.json or by repeating the flag.
  -d, --directory       a director that contains json files with the connector definitions to add
  -s, --sync-period     how often to check the current state of the connectors in the lcuster specified
                        by -c and the desired stats of the connectors as specified by -f or -d.
                        The default is 5 minutes.
      --allow-purge     if specified then any connectors that are found in the cluster that aren't
                        in the desired state (as spcified by -f or -d) will be deleted from the cluster.
                        The default is false.
      --auto-restart    if specified then connector tasks will be restarted if they are in a FAILED state
                        The default is false.
      --once            if specified the command will run once and then exit. The default is false.

```

NOTE: the -d and -f options are mutually exclusive. If you don't specify --once then the command will
run continuosly and will try and synchronise the state of the cluster according the duration 
specific by the -s option.

### Options inherited from parent commands

```
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl connectors](connectctl_connectors.md)	 - Manage connectors