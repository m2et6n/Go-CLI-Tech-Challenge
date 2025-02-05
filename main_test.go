package main

import (
	"bytes"
	"fmt"
	"os"
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
			name: "list happy path",
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
			name: "list error reading directory",
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
			name: "missing required list arguments",
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
			name: "create happy path",
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
			name: "create error reading directory",
			fields: fields{
				fw: mockFW{
					write: func(filePath string) (*os.File, error) {
						return nil, fmt.Errorf("fake-error")
					},
				},
				args: []string{"main.go", "create", "meetings"},
			},
			want:    "",
			wantErr: "Error creating file: fake-error",
		},
		{
			name: "missing required list arguments",
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
