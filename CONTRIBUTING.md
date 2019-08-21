# Contributing to connectctl

Firstly, thanks for considering contributing to *connectctl*. To make it a really 
great tool we need your help.

*connectctl* is [Apache 2.0 licenced](LICENSE) and accepts contributions via GitHub
pull requests. There are many ways to contribute, from writing tutorials or blog posts, 
improving the documentation, submitting bug reports and feature requests or writing code.


## Certificate of Origin

By contributing to this project you agree to the [Developer Certificate of
Origin](https://developercertificate.org/). This was created by the Linux
Foundation and is a simple statement that you, as a contributor, have the legal 
right to make the contribution. 

To signify that you agree to the DCO you must signoff all commits:

```bash
git commit --signoff
```

## Getting Started

- Fork the repository on GitHub
- Read the [README](README.md) for getting started as a user and learn how/where to ask for help
- If you want to contribute as a developer, continue reading this document for further instructions
- Play with the project, submit bugs, submit pull requests!

### Contribution workflow

#### 1. Set up your Go environment

This project is written in Go. To be able to contribute you will need:

1. A working Go installation of Go >= 1.12. You can check the
[official installation guide](https://golang.org/doc/install).

2. Make sure that `$(go env GOPATH)/bin` is in your shell's `PATH`. You can do so by
   running `export PATH="$(go env GOPATH)/bin:$PATH"`

3. Fork this repository and clone it by running:

```bash
git clone git@github.com:<yourusername>/connectlctl.git
```

> As the project uses modules its recommeneded that you NOT clone under the `GOPATH`.

#### 2. Test and build

Make sure you can run the tests and build the binary.

```bash
make install-deps
make test
make build
```

#### 3. Find a feature to work on

- Look at the existing [issues](https://github.com/90poe/connectctl/issues) to see if there is anything
you would like to work on. If don't see anything then feel free to create your own feature request.

- If you are a new contributor then take a look at the issues marked 
with [good first issue](https://github.com/90poe/connectctl/labels/good%20first%20issue).

- Make your code changes within a feature branch:

    ```bash
    git checkout -b <feature-name>
    ```

- Try to commit changes in logical units with a commit message in this [format](#format-of-the-commit-message). Remember
to signoff your commits.

- Don't forget to update the docs if relevent. The [README](README.md) or the [docs](docs/) folder is where docs usually live.

- Make sure the tests pass and that there are no linting problems.

#### 4. Create a pull request

Push your changes to your fork and then create a pull request to origin. Where possible use the PR template.

You can mark a PR as wotk in progress by prefixing the title of your PR with `WIP: `.


### Commit Message Format

We follow a rough convention for commit messages that is designed to answer two
questions: what changed and why. The subject line should feature the what and
the body of the commit should describe the why.

```text
Added restart connectors

Added a new subcommand for restarting connectors. It can be called in 2 ways.
To restart all connecors in a cluster:

connectctl connectors restart

And to restart only specific connectors:

connectctl connectors restart connector1,connectors2


Issue #2
```

The format can be described more formally as follows:

```text
<short title for what changed>
<BLANK LINE>
<why this change was made and what changed>
<BLANK LINE>
<footer>
```

The first line is the subject and should be no longer than 70 characters, the
second line is always blank, and other lines should be wrapped at 80 characters.
This allows the message to be easier to read on GitHub as well as in various git tools.

