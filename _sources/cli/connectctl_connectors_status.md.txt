## connectctl connectors status

Status of connectors

### Synopsis


Display status of selected connectors.
If some tasks or connectors are failing, command will exit with code 1.


```
connectctl connectors status [flags]
```

### Options

```
  -h, --help        help for add
  -c, --clusterURL  the url of the kafka connect cluster
  -n, --connectors  the names of the connectors. Multiple connector names 
                    can be specified either by comma separating conn1,conn2
                    or by repeating the flag --n conn1 --n conn2. If no name is
                    supplied status of ALL connectors will be displayed.
  -o, --output      specify the output format (valid options: json, table) (default "json")
  -q, --quiet       disable output logging
```
### Options inherited from parent commands

```
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl connectors](connectctl_connectors.md)	 - Manage connectors