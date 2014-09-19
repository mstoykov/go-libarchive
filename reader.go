package archive

/*
#cgo pkg-config: libarchive
#include <archive.h>
#include <archive_entry.h>
#include <stdlib.h>
#include "reader.h"
*/
import "C"
import (
	"io"
	"unsafe"
)

// Reader represents libarchive archive
type Reader struct {
	archive *C.struct_archive
	reader  io.Reader // the io.Reader from which we Read
	buffer  []byte    // buffer for the raw reading
}

// NewReader returns new Archive by calling archive_read_open
func NewReader(reader io.Reader) (r *Reader, err error) {
	r = new(Reader)
	r.buffer = make([]byte, 1024)
	r.archive = C.archive_read_new()
	C.archive_read_support_filter_all(r.archive)
	C.archive_read_support_format_all(r.archive)

	r.reader = reader

	e := C.go_libarchive_open(r.archive, unsafe.Pointer(r))

	err = codeToError(int(e))
	return
}

//export myopen
func myopen(archive *C.struct_archive, client_data unsafe.Pointer) C.int {
	// actually write something
	return ARCHIVE_OK
}

//export myclose
func myclose(archive *C.struct_archive, client_data unsafe.Pointer) C.int {
	// actually write something
	return ARCHIVE_OK
}

//export myread
func myread(archive *C.struct_archive, client_data unsafe.Pointer, block unsafe.Pointer) C.size_t {
	reader := (*Reader)(client_data)
	read, err := reader.reader.Read(reader.buffer)
	if err != nil && err != io.EOF {
		// set error
		read = -1
	}

	*(*uintptr)(block) = uintptr(unsafe.Pointer(&reader.buffer[0]))

	return C.size_t(read)
}

// Next calls archive_read_next_header and returns an
// interpretation of the ArchiveEntry which is a wrapper around
// libarchive's archive_entry, or Err.
//
// ErrArchiveEOF is returned when there
// is no more to be read from the archive
func (r *Reader) Next() (ArchiveEntry, error) {
	e := new(entryImpl)

	err := codeToError(int(C.archive_read_next_header(r.archive, &e.entry)))

	if err != nil {
		e = nil
	}

	return e, err
}

// Read calls archive_read_data which reads the current archive_entry.
// It acts as io.Reader.Read in any other aspect
func (r *Reader) Read(b []byte) (n int, err error) {
	n = int(C.archive_read_data(r.archive, unsafe.Pointer(&b[0]), C.size_t(cap(b))))
	err = codeToError(n)
	if n == 0 {
		err = io.EOF
	}
	return
}

// Free frees the resources the underlying libarchive archive is using
// calling archive_read_free
func (r *Reader) Free() error {
	if C.archive_read_free(r.archive) == ARCHIVE_FATAL {
		return ErrArchiveFatal
	}
	return nil
}

// Close closes the underlying libarchive archive
// calling archive read_cloe
func (r *Reader) Close() error {
	if C.archive_read_close(r.archive) == ARCHIVE_FATAL {
		return ErrArchiveFatal
	}
	return nil
}
