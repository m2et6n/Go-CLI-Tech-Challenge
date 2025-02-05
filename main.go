package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	fs := FS{}
	if err := run(ctx, fs, os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, fs fileSystem, args []string, stdout io.Writer) error {

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	defer cancel()

	if len(args) < 2 {
		return fmt.Errorf("Invalid command.\nUsage: go run main.go [command] [args]\n")
	}
	command := args[1]

	switch command {
	case "list":
		return list(fs, stdout, args)
	case "create":
		return writeFile(fs, stdout, args)
	case "delete":
		return deleteFile(fs, stdout, args)
	case "move":
		return moveFile(fs, stdout, args)
	default:
		//fmt.Fprintf(stdout, "Unknown command:%s\nUsage: go run main.go [list|create|delete|move] [args]\n", command)
		return fmt.Errorf("Unknown command:%s\nUsage: go run main.go [list|create|delete|move] [args]\n", command)
	}

	//return nil
}

type FS struct{}

/*
Functions for implementing the list command
*/

// Main list() function
func list(rfs readFileSystem, stdout io.Writer, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("Invalid arguments.\nUsage: go run main.go list [directory]")
	}

	dir := args[2]
	files, err := rfs.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Error reading directory: %s", err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Fprint(stdout, file.Name()+"/\n")
		} else {
			fmt.Fprint(stdout, file.Name()+"\n")
		}
	}

	return nil
}

type readFileSystem interface {
	ReadDir(name string) ([]os.DirEntry, error)
}

func (f FS) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

/*
Functions for implementing the create command
*/

// Main create() function
func writeFile(fw fileWriter, stdout io.Writer, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("Invalid arguments.\nUsage: go run main.go create [file-path]")
	}

	fileName := args[2]
	fileHandle, err := fw.WriteFile(fileName)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}

	if fileHandle != nil {
		fmt.Fprintf(stdout, "File created: %s\n", fileName)
	}

	return nil
}

type fileWriter interface {
	WriteFile(filePath string) (*os.File, error)
}

type FW struct{}

func (f FS) WriteFile(filePath string) (*os.File, error) {
	return os.Create(filePath)
}

/*
Functions for implementing the delete command
*/

// Main delete() function
func deleteFile(fd fileRemover, stdout io.Writer, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("Invalid arguments.\nUsage: go run main.go delete [file-path]")
	}

	fileName := args[2]
	err := fd.RemoveFile(fileName)
	if err != nil {
		return fmt.Errorf("Error deleting file: %s", err)
	}

	fmt.Fprintf(stdout, "File deleted: %s\n", fileName)

	return nil
}

type fileRemover interface {
	RemoveFile(filePath string) error
}

type FR struct{}

func (f FS) RemoveFile(filePath string) error {
	return os.Remove(filePath)
}

/*
Functions for implementing the move command
*/

// Main move() function
func moveFile(fm fileMover, stdout io.Writer, args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("Invalid arguments.\nUsage: go run main.go move [source-path] [destination-path]")
	}

	sourceFile := args[2]
	destinationFile := args[3]
	err := fm.MoveFile(sourceFile, destinationFile)
	if err != nil {
		return fmt.Errorf("Error moving file: %s", err)
	}

	fmt.Fprintf(stdout, "File %s moved to: %s\n", sourceFile, destinationFile)

	return nil
}

type fileMover interface {
	MoveFile(sourceFile string, destinationFile string) error
}

type FM struct{}

func (f FS) MoveFile(sourceFile string, destinationFile string) error {
	return os.Rename(sourceFile, destinationFile)
}

type fileSystem interface {
	readFileSystem
	fileWriter
	fileRemover
	fileMover
}
