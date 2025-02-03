package main

import (
	"io/fs"
	"os"
	"time"
)

type mockDirEntryImpl struct {
	name     string
	isDir    bool
	fileInfo mockFileInfo
}

type mockFileInfo struct {
}

func mockDirEntry(name string, isDir bool) os.DirEntry {
	return &mockDirEntryImpl{name: name, isDir: isDir, fileInfo: mockFileInfo{}}
}
func (m *mockDirEntryImpl) Name() string {
	return m.name
}
func (m *mockDirEntryImpl) IsDir() bool {
	return m.isDir
}
func (m *mockDirEntryImpl) Type() os.FileMode {
	if m.isDir {
		return os.ModeDir
	}
	return 0
}
func (m *mockDirEntryImpl) Info() (fs.FileInfo, error) {
	return &m.fileInfo, nil
}

func (m *mockFileInfo) Name() string {
	return ""
}
func (m *mockFileInfo) Size() int64 {
	return 0
}
func (m *mockFileInfo) Mode() fs.FileMode {
	return 0
}
func (m *mockFileInfo) ModTime() time.Time {
	return time.Now()
}
func (m *mockFileInfo) IsDir() bool {
	return false
}
func (m *mockFileInfo) Sys() any {
	return nil
}
