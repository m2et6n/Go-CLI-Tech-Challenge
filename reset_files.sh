#!/bin/bash

# Function to create a directory
create_directory() {
    local dir_path="$1"
    if [ -d "$dir_path" ]; then
        echo "Directory '$dir_path' already exists. Recreating '$dir_path'."
        rm -rf "$dir_path"
        mkdir -p "$dir_path"
        echo "Directory '$dir_path' created."
    else
        mkdir -p "$dir_path"
        echo "Directory '$dir_path' created."
    fi
}

# Function to create a file
create_file() {
    local file_path="$1"
    if [ -f "$file_path" ]; then
        echo "File '$file_path' already exists."
    else
        touch "$file_path"
        echo "File '$file_path' created."
    fi
}

echo "Resetting files..."

create_directory "meetings"
create_directory "meetings/agendas"
create_file "meetings/agendas/meeting_agenda_1.txt"
create_directory "meetings/attachments"
create_directory "meetings/icebreakers"
create_file "meetings/icebreakers/icebreaker1.txt"
create_file "meetings/icebreakers/icebreaker2.txt"
create_file "meetings/icebreakers/icebreaker3.txt"
create_file "meetings/icebreakers/icebreaker4.txt"
create_directory "meetings/notes"
create_file "meetings/notes/note1.txt"

echo "Files reset!"