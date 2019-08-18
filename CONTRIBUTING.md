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

#### 2. Fork and clone the repo

Make a fork of this repository and clone it by running:

```bash
git clone git@github.com:<yourusername>/connectlctl.git
```

It is not recommended to clone under your `GOPATH` (if you define one). Otherwise, you will need to set
`GO111MODULE=on` explicitly.

#### 3. Run the tests and build

Make sure you can run the tests and build the binary.

```bash
make install-deps
make test
make build
```

#### 4. Write your feature

- Find an [issue](https://github.com/90poe/connectctl/issues) to work on or
  create your own. If you are a new contributor take a look at issues marked
  with [good first issue](https://github.com/90poe/connectctl/labels/good%20first%20issue).

- Then create a topic branch from where you want to base your work (usually branched from master):

    ```bash
    git checkout -b <feature-name>
    ```

- Write your feature. Make commits of logical units and make sure your
  commit messages are in the [proper format](#format-of-the-commit-message).

- If needed, update the documentation, either in the [README](README.md) or in the [docs](docs/) folder.

- Make sure the tests are running successfully.

#### 5. Submit a pull request

Push your changes to your fork and submit a pull request to the original repository. If your PR is a work in progress
then make sure you prefix the title with `WIP: `. This lets everyone know that this is still being worked on. Once its
ready remove the `WIP: ` title prefix and where possible squash your commits.

```bash
git push <username> <feature-name>
```

## Acceptance policy

These things will make a PR more likely to be accepted:

- a well-described requirement
- tests for new code
- tests for old code!
- new code and tests follow the conventions in old code and tests
- a good commit message (see below)

In general, we will merge a PR once a maintainer has reviewed and approved it.
Trivial changes (e.g., corrections to spelling) may get waved through.
For substantial changes, more people may become involved, and you might get asked to resubmit the PR or divide the
changes into more than one PR.

### Format of the Commit Message

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

