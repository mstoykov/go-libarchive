package archive

/*
#cgo pkg-config: libarchive
#include <archive.h>
#include <archive_entry.h>
#include <stdlib.h>
*/
import "C"

import (
	"os"
	"time"
)

// ArchiveEntry represents an libarchive archive_entry 
type ArchiveEntry interface {
	// Not done yet
	Stat() os.FileInfo
	// The name of the entry
	PathName() string
}

type entryImpl struct {
	entry *C.struct_archive_entry
}

type entryInfo struct {
	stat *C.struct_stat
}

func (h *entryImpl) Stat() os.FileInfo {
	info := &entryInfo{}
	info.stat = C.archive_entry_stat(h.entry)
	return info
}

func (h *entryImpl) PathName() string {
	name := C.archive_entry_pathname(h.entry)

	return C.GoString(name)
}

func (e *entryInfo) Name() string {
	return "" // fix
}
func (e *entryInfo) Size() int64 {
	return 0 // fix
}
func (e *entryInfo) Mode() os.FileMode {
	return os.ModeTemporary // fix
}
func (e *entryInfo) ModTime() time.Time {
	return time.Now() // fix
}
func (e *entryInfo) IsDir() bool {
	return false // fix
}
func (e *entryInfo) Sys() interface{} {
	return nil // fix
}
