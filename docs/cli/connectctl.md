## connecctl

connectctl: work with Kafka Connect easily

### Synopsis

*connectctl* is a cli that makes working with kafka connect easier. It can be used to manage connectors and plugins and to also also actively manage/reconcile the state of a cluster.

The operations you can perform are split into 2 subcommands:
    connectors  Manage Kafka Connect connectors
    plugins     Manage Kafka connect connector plugins

Example usage:

	$ connectctl connectors add  \
		-c http://connect:8083 
	$ connectctl connectors list -c http://connect:8083 

### Options

```
  -h, --help                    Help for connectctl
  -l, --loglevel loglevel       Specify the loglevel for the program (default info)
      --logfile                 Specify a file to output logs to
```

### SEE ALSO

* [connectctl connectors](connectctl_connectors.md)	 - Perform connector operations against a Kafka Connect cluster
* [connectctl plugins](connectctl_plugins.md)	 - Perform connector plugin operations against a Kafka Connect cluster
* [connectctl version](connectctl_version.md)	 - Display version information