# Part 1: Technical Setup

## Prerequisites

If you have not already installed golang on your machine, you can follow the official golang instructions for installing golang [here](https://go.dev/doc/install) or by using brew with the command `brew install go`.

## Installing an IDE and Tools

It is recommended to use [Visual Studio Code](https://code.visualstudio.com/) as your IDE while working through this tech challenge. This is because Visual Studio has a rich set of extensions that will aid you as you progress through the exercise. However, you are free to use your preferred IDE.

After installing Visual Studio Code, you should head to the Extensions tab and install the following extensions:

- [Go](https://marketplace.visualstudio.com/items?itemName=golang.Go)

\*Note that you may have to restart VSCode after installing some extensions for them to take effect.

## Repository Setup

Create your own repository using the `Use this template` in the top right of the webpage. Make yourself the owner, and give the repository the same name as the template (Go-CLI-Tech-Challenge). Once your repository is cloned locally, you are all set to move on to the next step.

## Create Go Module for the Tech Challenge

Open up a terminal and navigate to the root directory for this tech challenge. From there, you will need to run the following command, swapping `[name]` with your name.
```
go mod init github.com/[name]/cli
```

The `go mod init` command is used to initialize your Go Module. As defined by the official Golang documentation:
> A module is a collection of Go packages stored in a file tree with a go.mod file at its root. The go.mod file defines the module’s module path, which is also the import path used for the root directory, and its dependency requirements, which are the other modules needed for a successful build. Each dependency requirement is written as a module path and a specific semantic version.

Next, you will need to install the list of dependencies that will be used for this challenge. In your terminal, run the following command.
```
go get github.com/stretchr/testify/assert
```

For more information on Go modules and managing dependencies, you can go to [Using Go Modules](https://go.dev/blog/using-go-modules).

## Next Steps

So far, you have created and cloned down your own repository for completing the tech challenge, along with initializing your go application. In part 2, we will walk through how to create one of the command line tools for the challenge: The `list` command. In this walkthrough, we will discuss about idiomatic go practices and patterns that you should utilize when writing effective go code, such as creating interfaces to decouple code and keeping a clean `main()` function. When you are ready, click [here](2-Detailed-Walkthrough.md) to proceed to part 2. You also have the option to skip part 2 if you are familiar with writing a CLI tool in Golang. In that case, click [here](3-Challenge-Assignment.md) to proceed to part 3.