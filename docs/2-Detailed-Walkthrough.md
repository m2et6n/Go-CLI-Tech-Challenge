<!-- omit in toc -->
# Part 2: Detailed Walkthrough

<!-- omit in toc -->
## Table of Contents
- [Project Structure](#project-structure)
- [Create the `list` CLI Tool](#create-the-list-cli-tool)
	- [Create `run()` Function](#create-run-function)
	- [Create `list()` Function](#create-list-function)
	- [De-Couple from the `os` Package](#de-couple-from-the-os-package)
		- [De-Coupling By Introducing Interfaces](#de-coupling-by-introducing-interfaces)
		- [De-Coupling By Passing In Operating System Fundamentals as Arguments](#de-coupling-by-passing-in-operating-system-fundamentals-as-arguments)
- [Add Testing](#add-testing)
	- [Mock out the `readFileSystem` Interface](#mock-out-the-readfilesystem-interface)
	- [Write Table-Driven Tests for `list()`](#write-table-driven-tests-for-list)
	- [Implement the Remaining Tests for the `list` Command](#implement-the-remaining-tests-for-the-list-command)
- [Implement the Tech Challenge](#implement-the-tech-challenge)
## Project Structure

By default, you should see the following file structure in your root directory

```
.
├── meetings/
│   └── agendas/
│       └── meeting_agenda_1.txt
│   └── attachments/
│   └── icebreakers/
│       └── icebreaker1.txt
│       └── icebreaker2.txt
│       └── icebreaker3.txt
│       └── icebreaker4.txt
│   └── notes/
│       └── note1.txt
├── .gitignore
├── file_mock_test.go
├── main_test.go
├── main.go
├── README.md
└── reset_files.sh
```

The `meetings/` folder will serve as the directory for us to test our CLI tools on, and is pre-populated with files and folders. Over the course of implementing this challenge, you will likely add, remove, and rename several of the provided files and folders. If at any point you would like to reset this folder to it's default state, you can do so by running the provided `reset_files.sh` bash script. You can run this by executing the following command in your terminal at the root directory for the project: `./reset_files.sh`

The `file_mock_test.go` file contains helper functions we will be utilizing when we get to unit testing our application. For now, you may ignore this file.

The `main_test.go` file will hold all our unit tests for the application.

Lastly, the `main.go` file serves as the entrypoint into our application, and will hold all of our code for the project. While it is quite common to break out the functionality of the application into several files/folders, our project will be fairly small and contain very little code. Therefore, we will be storing all of our application's logic in this `main.go` file. For more information on Go project structuring, please visit [Organizing a Go Module](https://go.dev/doc/modules/layout) from the Go team.

For now, lets open the `main.go` file and begin building our CLI tool.

## Create the `list` CLI Tool

For this walkthrough, we cover how to create the first of our CLI tools: The `list` command. The `list` command will be used to list all the files and subdirectories found within the specified directory. As such, `list` will take only one argument: The relative path to the directory we wish to display. We will also need validation to ensure this argument is provided when running the `list` command, and display an appropriate error when it is not specified or we are unable to list the given path (either because it does not exist or is not a valid directory).

### Create `run()` Function

Currently in our `main.go` file, we only have one empty method defined: The `main()` function. `main()` serves as the entrypoint into our application, and it does not require any inputs nor returns any values. While we could very well implement the entire application within this one function, it would not be idiomatic Go code. Instead, we will utilize one of the first effective practices for writing Go code: **Keeping a clean `main()`**. To do this, we will define a `run()` function that will serve as the location for our initial application logic. Defining a `run()` method allows us to separate our initialization logic from the main execution flow while also promoting readability and flexibility. Additionally, we are able to specify operating system fundamentals as arguments and return errors to our `main()` function for more graceful error handling. We will cover more in-depth on the type of operating system arguments we want to pass into `run()` later, but for now we will keep it simple. Update your `main.go` file to add the `run()` function.

```
func run() error {
	// TODO: Add initial application logic here

	return nil
}
```

Additionally, you will need to update `main()` so that it can call `run()` and handle any error returned.

```
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
```

> Note: You will need to update your list of imports to include 'fmt' and 'os'.

Next, we will want to fetch the command-line arguments to determine if they called the `list` command. Additionally, we will want to store the path to the directory the user wishes to list. Add the following logic at the beginning of the `run()` method, right below the `TODO` message:

```
args := os.Args

if len(args) < 2 {
    return fmt.Errorf("invalid command.\nUsage: go run main.go [command] [args]")
}
command := args[1]

switch command {
case "list":
    // TODO: implement the list() function and call here
default:
    return fmt.Errorf("unknown command:%s\nUsage: go run main.go [list|create|delete|move] [args]", command)
}
```

In the above code, we use `os.Args` to fetch the arguments that were passed in to the command line. With Golang's `os` package, the first item in the arguments array will always be the executable file that was ran, followed by any additional arguments listed. Therefore, we check to make sure that the length of the array is more than two, as we will store that second argument as the command to run. If we don't receive at least two arguments, then we will return an appropriate error message for `main()` to handle.

After storing the second argument as the command to run, we then want to check if the command matches the syntax we expect. We do this using a `switch` statement, specifying all the cases as our expected commands to run. In our case, we utilize the `default` path for if the user doe not enter an appropriate command, to which we return an appropriate error message for `main()` to handle. At this point, we are now able to create the `list()` method where we will handle our logic for listing the contents of the specified directory.

### Create `list()` Function

The `list()` function will hold all our application logic for listing the contents of the directory provided in the command-line arguments. Therefore, `list()` will take in as an argument the string array `args`. Additionally, as it is idiomatic to propagate all our error messages to the highest level for error handling, we will have `list()` return an error if any occurs, or `nil` if none occur. Create the `list()` function so it matches below.

```
func list(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("invalid arguments.\nUsage: go run main.go list [directory]")
	}

	dir := args[2]
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory: %s", err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Fprintf(os.Stdout, file.Name()+"/\n")
		} else {
			fmt.Fprintf(os.Stdout, file.Name()+"\n")
		}
	}

	return nil
}
```

Let's highlight a few important steps that are happening in the above code:
- We are first checking to see if the user provided a directory to read from by checking if there exists a third item in the `args` array and returning an appropriate error message if not.
- To list all the files and subdirectories within the given directory, we are utilizing the `os.ReadDir()` function. This function returns an array of type [fs.DirEntry](https://pkg.go.dev/io/fs#DirEntry) along with an error if something happened while reading that directory.
- As we want to differentiate a subdirectory from a file, we then iterate through all the objects in the returned array to format our output. We do this by utilizing the `DirEntry.IsDir()` method (which returns a boolean indicating if the object is a directory or not) and `DirEntry.Name` (which returns the name of the file/subdirectory).
- When printing the contents of the directory, we use `fmt.Fprintf()`. This is because we can specify which writer we want to write to. Currently, we are writing to the standard output, but we will change this later.

Now that we have defined out `list()` function, we can go back to our `run()` method and call `list()` within our `switch` statement.

```
...
switch command {
case "list":
    return list(args)
default:
    return fmt.Errorf("unknown command:%s\nUsage: go run main.go [list|create|delete|move] [args]", command)
}
...
```

We can now test out the `list` command. In your terminal, navigate to the root directory of the application, and run `go run main.go list meetings`. You should be able to see the list of subdirectories under `meetings/`. If it does not work, go back and make sure that you do not have any syntax errors in your code.

### De-Couple from the `os` Package

At this point, you have one of the four commands working for the command line tool. However, our application is tightly coupled to the `os` package for reading in our given directory. This will also be an issue later when we want to test the `list` command without actually reading the given directory. Additionally, we are currently printing everything to the standard output, when in reality we would want the ability to switch out where our outputs are being written to. To solve the first problem, we will utilize another pattern for writing effective Go code: **De-coupling by introducing interfaces**. Additionally, we will utilize Golang's ability to implement multiple interfaces by defining interfaces that serve only one purpose (i.e. read, write, etc) and implementing them with one type. For the second problem, we will circle back to the practice of **keeping a clean main()**.

#### De-Coupling By Introducing Interfaces

To begin, we will want to define an interface that will be responsible for one purpose: Reading the given directory. Before, we were utilizing the `os.ReadDir()` method, which takes in as an input the path to the directory to read and return an array of directory entries, or an error if any occurs. Therefore, we will want to define our interface with a method that matches that signature. In `main.go`, define the following `readFileSystem` interface as follows.

```
type readFileSystem interface {
	ReadDir(name string) ([]os.DirEntry, error)
}
```

Next, we will define a type that implements the `readFileSystem` interface. In Golang, types implicitly implement interfaces by simply implementing the interface's methods. Below the interface definition, define the `FS` struct and implement the `ReadDir()` method. In it, we will simply call `os.ReadDir()` and return it's contents.

```
type FS struct{}

func (f FS) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}
```

We can now go through our application and change our logic to where we are initializing this FS object and using it versus the `os` package. Working backwards, first update the `list()` function to pass in a `readFileSystem` object as a parameter, and call it in replace of `os.ReadDir()`.

```
func list(rfs readFileSystem, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("invalid arguments.\nUsage: go run main.go list [directory]")
	}

	dir := args[2]
	files, err := rfs.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory: %s", err)
	}
	for _, file := range files {
		if file.IsDir() {
			fmt.Fprintf(os.Stdout, file.Name()+"/\n")
		} else {
			fmt.Fprintf(os.Stdout, file.Name()+"\n")
		}
	}

	return nil
}
```

Next, update the `run()` function to pass in a `readFileSystem` object as a parameter, and update the call to `list()` to include the object.

```
func run(fs readFileSystem) error {
	args := os.Args

	if len(args) < 2 {
		return fmt.Errorf("invalid command.\nUsage: go run main.go [command] [args]")
	}
	command := args[1]

	switch command {
	case "list":
		return list(fs, args)
	default:
		return fmt.Errorf("unknown command:%s\nUsage: go run main.go [list|create|delete|move] [args]", command)
	}
}
```

Lastly, update `main()` to initialize a `FS` object and pass it in as a parameter to the `run()` function.

```
func main() {
	fs := FS{}
	if err := run(fs); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
```

There we go! We have now de-coupled ourselves from the `os` package when reading a directory by utilizing the power of interfaces! However, there is still an issue with the way we are passing in our file system through the application. Take a moment to see if you can spot the issue.

Currently, we are only passing in one interface object into our `run()` function because we have only written out one command-line tool. However, if we were to repeat this with our remaining 3 commands, that would be 4 interfaces we are passing in! That is not very ideal. One option would be to simply initialize our interfaces in `run()` so that we don't need to pass them in as parameters. This is not a good idea, as we want to initialize all our interface instances at the top level (`main()`) and pass them in with our operating system fundamentals (more on this later). Another option would be to combine all 4 interfaces into one interface and pass just the one in. While this would simplify things, it also is not very idiomatic, as instances of this larger interface would be able to utilize all the methods. In general, this is not a good practice, as we want to encapsulate an instance to its required behavior (the instance in `list()` should only be able to read a directory, not create or delete files). If only we were able to use inheritance with interfaces...

Unfortunately, Golang does not support inheritance. However, Golang does support something called embedding. With embedding, interfaces can embed (or take) other interfaces, essentially giving them the same required methods as the embedded interface. To embed one interface into another, we simply need to list the embedded interface within the embedding interface's interface definition. For us, we can accomplish this by adding the following interface in `main.go`.

```
type fileSystem interface {
	readFileSystem
}
```

Then, we can update the signature for `run()` to pass in the `fileSystem` instance in replace of the `readFileSystem` instance.

```
func run(fs fileSystem) error {
	// TODO: Add initial application logic here
	args := os.Args
    ...
}
```

And we are done! We have introduced a new interface that can embed all of our file system interfaces, and then pass it in as a requirement to the `run()` method. Let's take a second to understand how exactly this works.
- In `main()`, we initialized an object of type `FS`. Since this type implements the `ReadDir()` method, `FS` implements the `readFileSystem` interface. Additionally, as the `fileSystem` interface only embeds the `readFileSystem` interface and contains no other methods, `FS` also is implementing the `fileSystem` interface.
- We are then passing in the `fs` object as a parameter to the `run()` method. Because this parameter is of type `fileSystem` and `fs` implements `fileSystem`, the `fs` instance is also of type `fileSystem`. This allows us to pass in the `fs` instance as a parameter to all our different flows we will have (i.e. `list()`, `create()`, etc).
- For each method that handles one command (i.e. `list()`), we will specify a parameter to be of type the more-specific interface we want to use here. In our above case, we had `list()` require as a parameter a `readFileSystem` instance. That means when we pass in the `fs` object to `list()` from the `run()` function, it is being casted to the type `readFileSystem` within `list()`. This means that we are not able to utilize any other methods defined in the `fileSystem` interface.

By extracting our call to `os.ReadDir()` using interfaces, we have been able to de-couple our code from the `os` package when we are reading the contents of the given directory, allowing for flexibility and better testability. Additionally, we used the concept of embedding to pass in our file system instance while also restricting it's control based on where we are within our application. As you begin to make other single-use interfaces for the remaining commands, make sure that you embed them within the `fileSystem` interface.

#### De-Coupling By Passing In Operating System Fundamentals as Arguments

While we have de-coupled ourselves from the `os` package in terms of reading directories, we are still tightly coupled in the way that we fetch our command-line arguments and output our messages. Wouldn't it be nice to be able to swap out where we output our messages to? We can accomplish this by circling back to the practice of **keeping a clean main()**. To recap, we want to create a `run()` function to separate our initialization logic from our main execution flow for the application. Another aspect of this practice that we will now dive into is having the `run()` function take in operating system fundamentals as arguments. In doing so, we are able to begin our main execution flow by calling `run()` but having the ability to change out how we get our command-line arguments and where we write our outputs to. Mat Ryer does a very good job at explaining this in his article [How I write HTTP services in Go after 13 years](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/#func-main-only-calls-run). He also does a great job at listing the type of OS fundamentals to pass into the `run()` function. For our purposes, we will only require three new arguments:

| Argument | Type | Description |
| --- | --- | --- |
| ctx | `context.Context` | The context for the application. They are useful for passing metadata and control signals between goroutines. |
| args | `[]string` | The arguments passed into the application. We will be grabbing these from `main()` using `os.Args` |
| stdout | `io.Writer` | The writer we want to output to. |

First, update the `main()` function to initialize a context instance and pass it into `run()` along with `os.Args` and `os.Stdout` as our values for `args` and `stdout`.

```
func main() {
	ctx := context.Background()
	fs := FS{}
	if err := run(ctx, fs, os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
```
> Note: You will need to add `context` to the list of imports.

Next, update the function signature for `run()` to require the new arguments. You will also need to signal the context and defer it's cancellation at the beginning of the function. The beginning of `run()` should look like the following.

```
func run(
	ctx context.Context,
	fs fileSystem,
	args []string,
	stdout io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if len(args) < 2 {
		return fmt.Errorf("invalid command.\nUsage: go run main.go [command] [args]")
	}

    ...
}
```
> Note: You will need to add `io` and `os/signal` to the list of imports.

Now, update the call to `list()` within the switch statement to take `stdout` as one of its arguments.

```
switch command {
case "list":
    return list(fs, stdout, args)
default:
    return fmt.Errorf("unknown command:%s\nUsage: go run main.go [list|create|delete|move] [args]", command)
}
```

Finally, update the function signature for `list()` to include the `stdout` parameter, and replace all instances of os.Stdout in the `list()` function with `stdout`. This will instruct the program to write our output messages to the variable writer object, meaning we can swap out where we write our messages out to. Your `list()` function should look similar to below.

```
func list(rfs readFileSystem, stdout io.Writer, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("invalid arguments.\nUsage: go run main.go list [directory]")
	}

	dir := args[2]
	files, err := rfs.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory: %s", err)
	}
	for _, file := range files {
		if file.IsDir() {
			fmt.Fprintf(stdout, file.Name()+"/\n")
		} else {
			fmt.Fprintf(stdout, file.Name()+"\n")
		}
	}

	return nil
}
```

We have now been able to de-couple ourselves from the `os` package. This enables flexibility in how we read in our command-line arguments, handle file/directory read/writes, and where we output our print statements to. Additionally, we have improved the testability of our code. At this point, you have now completed the `list` command for our CLI tool, utilizing practices for writing effective Go code. Take a moment to try out your new CLI tool. Next, we will cover unit testing our application by writing unit tests for our `list()` function.

## Add Testing

When it comes to unit testing in Golang, there are a few different methods. In this walkthrough, we will be covering one of the popular patterns for unit testing our application through the use of **table-driven testing**. Instead of writing separate unit tests for each case, table-driven testing allows us to organize and run multiple test cases in a clear and concise manner. Additionally, we will be using the practice of **mocking** to isolate and test our application code without relying on dependencies. For this section, we will walk through how to write table-driven tests for the `list()` function.

### Mock out the `readFileSystem` Interface

Looking at the `list()` function, we require three arguments to be passed in: (1) A `readFileSystem` instance that handles reading our desired directory, (2) The writer we want to write our print statements to, and (3) The arguments that were passed into the application in order to fetch our desired directory to list. While we could pass in an instance of our `FS` type for testing, that would not be very ideal. In reality, we want to simulate the behavior that this `readFileSystem` instance should do and control what is returned from our instance (either a list of directory entries or an error). To do this, we can utilize the practice of **mocking** our `readFileSystem` interface. We will define a new struct that implements the `readFileSystem` interface, but allows us to pass in what we want the `ReadDir()` function to return. In `main_test.go`, define the `mockRFS` struct with the `read` functional property and have it implement the `readFileSystem` interface.

```
type mockRFS struct {
	read func(name string) ([]os.DirEntry, error)
}

func (m mockRFS) ReadDir(name string) ([]os.DirEntry, error) {
	return m.read(name)
}
```
> Note: You will need to import the `os` package.

As shown above, we have defined a new struct called `mockRFS` that implements the `readFileSystem` interface. The `mockRFS` has a property called `read`, which is a function that has the same signature as our `ReadDir()` function. By defining this property, we can simulate what `ReadDir()` function will return when testing by providing this logic when initializing the `mockRFS` instance in our tests. That way, we can pass in this `mockRFS` struct into our `list()` function. We will see this more in practice when we write our table-driven tests later.

### Write Table-Driven Tests for `list()`

We are now able to write our table-driven tests for `list()`. Table-driven tests can look different depending on the size of the application, but all follow the same pattern:

1. We first define a struct for our test cases. This struct typically contains the inputs and fields required to run the test, along with the expected output from running the test.
2. We then create a slice of test cases, where each slice contains all the fields defined in our struct necessary to run each test.
3. Next, we then iterate over each test in our slice and run our test case. For each iteration, we run our test case and assert that the expected behavior is reached. It is here that we see the power of table-driven testing as we are often able to run these test cases in parallel and use subtesting to improve output readability.

For more information on table-driven testing, you can check out the article [Table-driven unit tests](https://yourbasic.org/golang/table-driven-unit-test/).

In `main_test.go`, write the following function signature for `TestList()`.

```
func TestList(t *testing.T) {
    // TODO: Add table-driven test
}
```
> Note: You will need to add `testing` to the list of imports.

Next, define the struct for our test cases. You can define the struct in however way you want, but it is recommended to follow the syntax shown below. Note: You can also combine this step with initializing our slice of test cases by using the below syntax.

```
func TestList(t *testing.T) {
	// TODO: Write table-driven test
	t.Parallel()
	type fields struct {
		rfs  readFileSystem
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr string
	}{}
}
```
> Note: Because we are not performing any sequential or blocking actions like reading from a database, we can call `t.Parallel()` at the beginning of our table-driven test so it can be run in parallel with other tests we make later on.

In our struct definition above, we defined tests to be a slice of type struct that contains the 4 following properties:

| Name | Type | Description |
| --- | --- | --- |
| name | string | the name of the specific test. This helps with printing and debugging. |
| fields | fields (struct) | the fields required to set up and call `list()`. The `fields` struct will contain our `readFileSystem` instance along with the arguments we want to pass into the function. |
| want | string | if we expect `list()` to function correctly, `want` will contain the expected printed statements outputted. |
| wantErr | string | if we expect `list()` to return an error, `wantErr` will contains the expected error message. |

Next, we will be creating our first test case for `list()`. For this test case, we will want to be testing the happy path for `list()`, which will print out the list of files and subdirectories within the directory we provide. To do this, copy the test case below and add it into the tests slice.

```
tests := []struct {
    ...
}{
    {
        name: "happy path",
        fields: fields{
            rfs: mockRFS{
                read: func(name string) ([]os.DirEntry, error) {
                    return []os.DirEntry{
                        mockDirEntry("attachments", true),
                        mockDirEntry("agendas", true),
                        mockDirEntry("sample_file.txt", false),
                    }, nil
                },
            },
            args: []string{"main.go", "list", "meetings"},
        },
        want:    "attachments/\nagendas/\nsample_file.txt\n",
        wantErr: "",
    },
}
```

Let's break down what is happening here. We defined an instance of our custom struct type. For this instance, we named the instance "happy path".
For our `fields` property, we created a new instance of the `fields` struct. As `fields` has two properties of its own, we created new instances for the `args` and `rfs` property for `fields`. Because we want to simulate a successful retrieval of the entries for the provided directory, we initialized `read` with a `mockRFS` instance. For this instance, we set `read` to a function that will return a slice of 3 directory entries: An `attachments/` subdirectory, an `agendas/` subdirectory, and a `sample_file.txt` file. To be able to initialize these three `os.DirEntry` objects, we are utilizing the `mockDirEntry()` function that has been provided to you in the `file_mock_test.go` file. We will not go into depth on how this function works, but just know that you can use it to make instances of the `os.DirEntry` interface and specify if it should be a directory or file. By setting `rfs` to a `mockRFS` instance instead of `FS`, we are able to control our expected output when `list()` calls `ReadDir()` and isolate that dependency when testing our `list()` function. `fields` also contains the `args` property, which references the arguments that were initially passed into the application. Because we expect the `args` property to contain at least three entries, we set `args` to a new string slice containing three strings with our chosen directory to be called "meetings".
Now that we have specified our `fields` and `name` property, we now move on to `want` and `wantErr`. Because we expect the `list()` function to work successfully, we set the `wantErr` message to be an empty string. For `want`, we expect `list()` to output a formatted print statement of the two subdirectories and file that we return from `ReadDir()`. Therefore, we set `want` to be a string containing the three entries in the order we specified them in the returned slice, formatting it to include the `/` for all subdirectories and expecting a new-line character after each entry.

That was a lot to unpack. Make sure you spend some time looking over how we initialized our test case so you can replicate this for future test cases.

Now that we have initialized our test case, we can move on to the next step, which is iterating over our test cases and asserting our expected behavior. Add the following code to the end of `TestList()` below the initialization of `tests`.

```
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        var bytes bytes.Buffer
        err := list(tt.fields.rfs, &bytes, tt.fields.args)
        if tt.wantErr == "" {
            assert.NoError(t, err)
        } else {
            assert.Equal(t, tt.wantErr, err.Error())
        }
        got := bytes.String()
        assert.Equal(t, tt.want, got)
    })
}
```
> Note: You will need to add `bytes` and `"github.com/stretchr/testify/assert"` to the list of imports.

Let's go over what is happening here. First, we are iterating through each test in the slice, calling this test `tt`. To improve the output readability of our testing, it is here that we are calling `t.Run()` specifying the name of the test we are running. Since these tests also do not contain any sequential or blocking logic (like reading to a database), we are able to call `t.Parallel()` to run each test case in parallel, improving performance. Next, we want to grab the printed statements that `list()` will output upon a successful run. To do this, we create a buffer and pass it in as our writer for `list()` to write to. It is here that you may recall we were using `fmt.Fprintf()` to write our print statements to the writer we passed into our function. Without using this, we would have not been able to capture the output from our `list()` function, which would have made it difficult to assert it worked as intended. After defining our buffer, we then call our `list()` function, providing the fields we initialized within our test struct and the newly created buffer, and store the returned error in `err`. We then check to see if we were expecting an error to be returned by `list()` by checking the contents of the test's `wantErr` property. If we expected no error (indicated by an empty string), we then assert there was no error. Otherwise, we assert that the returned error message matches the error message we expected. Finally, we grab the output written to buffer and assert that it matches our expected print statements.

We have now completed our first test case for `list()`, which covers the happy path. To run this unit test, open up your terminal and run the following command from the root directory of your repository: `go test -v`. You should see that the test passed.
> Note: You use the `-v` flag to get verbose output that lists all of the tests and their results.

Let's write one more test case for `list()`, covering the exception when we are not able to read the provided directory. When `list()` is not able to read the provided directory, it returns a formatted error containing the following pattern: "error reading directory: {error}". Therefore, we can append to our slice of test cases the following test.

```
{
    name: "error reading directory",
    fields: fields{
        rfs: mockRFS{
            read: func(name string) ([]os.DirEntry, error) {
                return nil, fmt.Errorf("fake-error")
            },
        },
        args: []string{"main.go", "list", "meetings"},
    },
    want:    "",
    wantErr: "error reading directory: fake-error",
},
```
> Note: You will need to add `fmt` to the list of imports

In this test case, we define the `read` property to be a function that simply returns an error message with the message "fake-error". As a result, the error message we expect `list()` to return will be "error reading directory: fake-error". Run the table-driven test again by running the command `go test -v` in your terminal. You should see that both the happy path and the exception test case were ran and passed.

### Implement the Remaining Tests for the `list` Command

Congratulations! At this point, you have been able to write table-driven tests to unit test part of the `list` command. Furthermore, you were able to use the practice of mocking to isolate the application code you are testing from it's dependencies. At this point, you have been given all the tools and knowledge to write out the remaining tests for the `list` command.

> Note: As a requirement for completing the tech challenge, you will be required to have at least 80% code coverage in your unit tests. To check the code coverage at any point, you can run the following commands. These commands will generate the code coverage into a coverage profile file, and then open a web browser displaying a visual representation of your code coverage, highlighting lines that were executed during testing.
> ```
> go test -coverprofile=coverage.out ./...
> go tool cover -html=coverage.out
> ```

## Implement the Tech Challenge

Congratulations! At this point, you have implemented one of the four command-line tools for the tech challenge, utilizing effective practices for writing idiomatic Go code. Additionally, you wrote table-driven tests to unit test the various flows of your application. At this point, you should be able to complete the remaining three CLI tools to finish the Tech Challenge. When you are ready, you can move on to [Part 3](3-Challenge-Assignment.md), where we go over the challenge assignment and list the requirements for completing the challenge.