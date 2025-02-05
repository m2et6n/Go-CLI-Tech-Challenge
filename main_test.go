package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Tests for list command
*/
type mockRFS struct {
	read func(name string) ([]os.DirEntry, error)
}

func (m mockRFS) ReadDir(name string) ([]os.DirEntry, error) {
	return m.read(name)
}

func TestList(t *testing.T) {
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
	}{
		{
			name: "list - happy path",
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
		{
			name: "list - error reading directory",
			fields: fields{
				rfs: mockRFS{
					read: func(name string) ([]os.DirEntry, error) {
						return nil, fmt.Errorf("fake-error")
					},
				},
				args: []string{"main.go", "list", "meetings"},
			},
			want:    "",
			wantErr: "Error reading directory: fake-error",
		},
		{
			name: "list - missing required arguments",
			fields: fields{
				rfs: mockRFS{
					read: func(name string) ([]os.DirEntry, error) {
						return nil, fmt.Errorf("fake-error")
					},
				},
				args: []string{"main.go", "list"},
			},
			want:    "",
			wantErr: "Invalid arguments.\nUsage: go run main.go list [directory]",
		},
	}

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
}

/*
Tests for create command
*/
type mockFW struct {
	write func(filePath string) (*os.File, error)
}

func (m mockFW) WriteFile(filePath string) (*os.File, error) {
	return m.write(filePath)
}

func TestCreate(t *testing.T) {
	t.Parallel()
	type fields struct {
		fw   fileWriter
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr string
	}{
		{
			name: "create - happy path",
			fields: fields{
				fw: mockFW{
					write: func(filePath string) (*os.File, error) {
						return &os.File{}, nil
					},
				},
				args: []string{"main.go", "create", "test_file.txt"},
			},
			want:    "File created: test_file.txt\n",
			wantErr: "",
		},
		{
			name: "create - error creating file",
			fields: fields{
				fw: mockFW{
					write: func(filePath string) (*os.File, error) {
						return nil, fmt.Errorf("fake-error")
					},
				},
				args: []string{"main.go", "create", "meetings/agendas/test_file.txt"},
			},
			want:    "",
			wantErr: "Error creating file: fake-error",
		},
		{
			name: "create - missing required arguments",
			fields: fields{
				fw: mockFW{
					write: func(fileName string) (*os.File, error) {
						return nil, fmt.Errorf("fake-error")
					},
				},
				args: []string{"main.go", "create"},
			},
			want:    "",
			wantErr: "Invalid arguments.\nUsage: go run main.go create [file-path]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var bytes bytes.Buffer
			err := writeFile(tt.fields.fw, &bytes, tt.fields.args)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.wantErr, err.Error())
			}
			got := bytes.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

/*
Tests for delete command
*/
type mockFR struct {
	delete func(filePath string) error
}

func (m mockFR) RemoveFile(filePath string) error {
	return m.delete(filePath)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	type fields struct {
		fr   fileRemover
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr string
	}{
		{
			name: "delete - happy path",
			fields: fields{
				fr: mockFR{
					delete: func(filePath string) error {
						return nil
					},
				},
				args: []string{"main.go", "delete", "test_file.txt"},
			},
			want:    "File deleted: test_file.txt\n",
			wantErr: "",
		},
		{
			name: "delete - error deleting file",
			fields: fields{
				fr: mockFR{
					delete: func(filePath string) error {
						return fmt.Errorf("fake-error")
					},
				},
				args: []string{"main.go", "delete", "meetings/agendas/test_file.txt"},
			},
			want:    "",
			wantErr: "Error deleting file: fake-error",
		},
		{
			name: "delete - missing required arguments",
			fields: fields{
				fr: mockFR{
					delete: func(fileName string) error {
						return fmt.Errorf("fake-error")
					},
				},
				args: []string{"main.go", "delete"},
			},
			want:    "",
			wantErr: "Invalid arguments.\nUsage: go run main.go delete [file-path]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var bytes bytes.Buffer
			err := deleteFile(tt.fields.fr, &bytes, tt.fields.args)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.wantErr, err.Error())
			}
			got := bytes.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

/*
Tests for move command
*/
type mockFM struct {
	move func(sourceFile string, destinationFile string) error
}

func (m mockFM) MoveFile(sourceFile string, destinationFile string) error {
	return m.move(sourceFile, destinationFile)
}

func TestMove(t *testing.T) {
	t.Parallel()
	type fields struct {
		fm   fileMover
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr string
	}{
		{
			name: "move - happy path",
			fields: fields{
				fm: mockFM{
					move: func(sourceFile string, destinationFile string) error {
						return nil
					},
				},
				args: []string{"main.go", "move", "meetings/agendas/test_file.txt", "meetings/notes/test_file.txt"},
			},
			want:    "File meetings/agendas/test_file.txt moved to: meetings/notes/test_file.txt\n",
			wantErr: "",
		},
		{
			name: "move - error moving file",
			fields: fields{
				fm: mockFM{
					move: func(sourceFile string, destinationFile string) error {
						return fmt.Errorf("fake-error")
					},
				},
				args: []string{"main.go", "move", "meetings/agendas/test_file.txt", "meetings/notes/test_file.txt"},
			},
			want:    "",
			wantErr: "Error moving file: fake-error",
		},
		{
			name: "move - missing required arguments",
			fields: fields{
				fm: mockFM{
					move: func(fsourceFile string, destinationFile string) error {
						return fmt.Errorf("fake-error")
					},
				},
				args: []string{"main.go", "move"},
			},
			want:    "",
			wantErr: "Invalid arguments.\nUsage: go run main.go move [source-path] [destination-path]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var bytes bytes.Buffer
			err := moveFile(tt.fields.fm, &bytes, tt.fields.args)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.wantErr, err.Error())
			}
			got := bytes.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

// Mock implementations
type mockFileSystem struct{}

func (m mockFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return []os.DirEntry{mockDirEntry("sample_file.txt", false)}, nil
}
func (m mockFileSystem) WriteFile(filePath string) (*os.File, error) {
	return &os.File{}, nil
}
func (m mockFileSystem) RemoveFile(filePath string) error {
	return nil
}
func (m mockFileSystem) MoveFile(sourceFile string, destinationFile string) error {
	return nil
}

// Test Writer to capture stdout
type testWriter struct {
	output strings.Builder
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	return w.output.Write(p)
}

func (w *testWriter) String() string {
	return w.output.String()
}

// Unit Test for run function
func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
		want    string
	}{
		{
			name:    "Unknown Command",
			args:    []string{"go", "unknown"},
			wantErr: "Unknown command:unknown\nUsage: go run main.go [list|create|delete|move] [args]\n",
			want:    "",
		},
		{
			name:    "Invalid Command",
			args:    []string{"go"},
			wantErr: "Invalid command.\nUsage: go run main.go [command] [args]\n",
			want:    "",
		},
		{
			name:    "List Command",
			args:    []string{"go", "list", "meetings/agendas"},
			wantErr: "",
			want:    "sample_file.txt\n",
		},
		{
			name:    "Create Command",
			args:    []string{"go", "create", "test_file.txt"},
			wantErr: "",
			want:    "File created: test_file.txt\n",
		},
		{
			name:    "Delete Command",
			args:    []string{"go", "delete", "test_file.txt"},
			wantErr: "",
			want:    "File deleted: test_file.txt\n",
		},
		{
			name:    "Move Command",
			args:    []string{"go", "move", "source.txt", "dest.txt"},
			wantErr: "",
			want:    "File source.txt moved to: dest.txt\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := mockFileSystem{}
			ctx := context.Background()

			var bytes bytes.Buffer

			err := run(ctx, fs, tt.args, &bytes)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.wantErr, err.Error())
			}
			got := bytes.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
