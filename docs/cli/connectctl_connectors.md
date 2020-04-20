## connectctl connectors

perform a connector operation against a Kafka Connect cluster

### Synopsis


Perform operations against a Kafka Connect cluster that relate to connectors.
Operations are always against a specific cluster and URL must be supplied.


```
connectctl connectors <command...> [flags]
```

### Options

None, all options are at the subcommand level

### Options inherited from parent commands

```
  -h, --help                    Help for connectctl
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl](connectctl.md)	 - connectctl: work with Kafka Connect easily
* [connectctl connectors status](connectctl_connectors_status.md)     - Connectors status
* [connectctl connectors add](connectctl_connectors_add.md)     - Add connectors
* [connectctl connectors remove](connectctl_connectors_remove.md)     - Remove connectors
* [connectctl connectors list](connectctl_connectors_list.md)     - List connectors
* [connectctl connectors restart](connectctl_connectors_restart.md)     - Restart connectors
* [connectctl connectors pause](connectctl_connectors_pause.md)     - Pause connectors
* [connectctl connectors resume](connectctl_connectors_resume.md)     - Resume connectors
* [connectctl connectors manage](connectctl_connectors_manage.md)     - Actively manage connectors