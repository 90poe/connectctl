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

- Add yourself to [humans.txt](humans.txt) if this is your first contribution.

- Try to commit changes in logical units with a commit message in this [format](#format-of-the-commit-message). Remember
to signoff your commits.

- Don't forget to update the docs if relevent. The [README](README.md) or the [docs](docs/) folder is where docs usually live.

- Make sure the tests pass and that there are no linting problems.

#### 4. Create a pull request

Push your changes to your fork and then create a pull request to origin. Where possible use the PR template.

You can mark a PR as wotk in progress by prefixing the title of your PR with `WIP: `.


### Commit Message Format

We would like to follow the **Conventional Commits** format for commit messsages. The full specification can be 
read [here](https://www.conventionalcommits.org/en/v1.0.0-beta.3/). The format is:

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

Where `<type>` is one of the following:
* `feat` - a new feature
* `fix` - a bug fix
* `chore` - changes to the build pocess, code generation or anything that doesn't match elsewhere
* `docs` - documentation only changes
* `style` - changes that don't affect the meaning of the code (i.e. code formatting)
* `refactor` - a change that doesn't fix a feature or bug
* `test` - changes to tests only.

The `scope` can be a pkg name but is optional.
The `body` should include details of what changed and why. If there is a breaking change then the `body` should start with the 
following: `BREAKING CHANGE`.

The footer should include any related github issue numbers.

An example:

```text
feat: Added connector status command

A new command has been added to show the status of tasks within a connector. The
command will return the number of failed tasks as the exit code.

Fixes: #123
```

A tool like [Commitizen](https://github.com/commitizen/cz-cli) can be used to help with formatting commit messages.

