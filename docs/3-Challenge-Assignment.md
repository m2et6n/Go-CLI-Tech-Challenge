# Part 3: Tech Challenge Assignment

## Overview

For this assignment, you will be creating a File System Command Line Tool (CLI). This CLI will be responsible for creating, listing, updating (renaming or moving), and deleting files within a file system.

## Required CLI commands
| Command | Arguments | Example Command | Explanation |
| --- | --- | --- | --- |
| list [directory] | directory: path of directory | go run main.go list meetings/agendas/ | List all files and folders in the specified directory. |
| create [file-path] | file-path: name of the file, including the relative path | go run main.go create meetings/agendas/meeting_agenda_1.txt | Create the specified file. |
| move [source-path] [destination-path] | source-path: the source path of the file to move <br> destination-path: the destination path for the file | go run main.go move meetings/agendas/meeting_agenda_1.txt meetings/icebreakers/icebreaker_1.txt | Move a file from source location to destination location. Also used for renaming a file. |
| delete [file-path] | file-path: name of the file, including the relative path | go run main.go delete meetings/agendas/meeting_agenda_1.txt | Delete the file at the specified location. |

## Project Requirements Checklist

- [ ] Your program should be able to run successfully.
- [ ] For each CLI tool, your program should be able to run as intended and accept the appropriate command line arguments.
- [ ] Your program should handle input validation and display an appropriate error message.
- [ ] Your project should include unit tests with 80% test coverage.

## Submission Steps

1. Navigate to your Github repository
2. Make sure that all your work is checked into a separate branch from the starting branch
3. Open a pull request from your working branch to the starting branch
4. Add the following reviewers so they will receive a notification that your solution is ready for review:
   golang-tech-challenge-reviewers