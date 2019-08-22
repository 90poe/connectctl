## connectctl
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![CircleCI](https://circleci.com/gh/90poe/connectctl/tree/master.svg?style=svg)](https://circleci.com/gh/90poe/connectctl/tree/master)
![OSS Lifecycle](https://img.shields.io/osslifecycle/90poe/connectctl)


A CLI for working with [Kafka Connect](https://docs.confluent.io/current/connect/index.html).

> This is a work in progress project. If you'd like to contribute please consider contrubiting.

### Getting started

Install the CLI for your platform from the [releases page](https://github.com/90poe/connectctl/releases).

Once installed run the cli from the terminal. You can see the available commands by running 
the following:

```bash
connectctl help
```

There are 2 top level commands (each with their own sub-commands):
- `connectctl connectors` - for operations relating to connectors in a Kafka Connect cluster
- `connectctl plugins` - for operations relating to connector plugiuns in a Kafka Connect cluster

### Contributing

We'd love you to contribute to the project. If you are interested in helping out please 
see the [contributing guide](CONTRIBUTING.md).

### Acknowledgements
The code in `pkg/client/connect` is originally from [here](https://github.com/go-kafka/connect) but has been modified for use in this utility.

### License

Copyright 2019 90poe.  This project is licensed under the Apache 2.0 License.
