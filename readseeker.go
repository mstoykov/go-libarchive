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
	"errors"
	"io"
	"unsafe"
)

// ReedSeeker represents libarchive archive
type ReedSeeker struct {
	archive *C.struct_archive
	reader  io.ReadSeeker // the io.ReedSeeker from which we Read
	buffer  []byte        // buffer for the raw reading
	index   int64         // current reading index
}

// NewReader returns new Archive by calling archive_read_open
func NewReadSeeker(reader io.ReadSeeker) (r *ReedSeeker, err error) {
	r = new(ReedSeeker)
	r.buffer = make([]byte, 1024)
	r.archive = C.archive_read_new()
	C.archive_read_support_filter_all(r.archive)
	C.archive_read_support_format_all(r.archive)

	seek_callback := (*C.archive_seek_callback)(C.go_libarchive_seek)
	C.archive_read_set_seek_callback(r.archive, seek_callback)

	r.reader = reader

	e := C.go_libarchive_open(r.archive, unsafe.Pointer(r))

	err = codeToError(r.archive, int(e))
	return
}

//export myseek
func myseek(archive *C.struct_archive, client_data unsafe.Pointer, request C.int64_t, whence C.int) C.int64_t {
	reader := (*ReedSeeker)(client_data)
	offset, err := reader.reader.Seek(int64(request), int(whence))
	if err != nil {
		return C.int64_t(0)
	}
	return C.int64_t(offset)
}

// Next calls archive_read_next_header and returns an
// interpretation of the ArchiveEntry which is a wrapper around
// libarchive's archive_entry, or Err.
//
// ErrArchiveEOF is returned when there
// is no more to be read from the archive
func (r *ReedSeeker) Next() (ArchiveEntry, error) {
	e := new(entryImpl)

	errno := int(C.archive_read_next_header(r.archive, &e.entry))
	err := codeToError(r.archive, errno)

	if err != nil {
		e = nil
	}

	return e, err
}

// Read calls archive_read_data which reads the current archive_entry.
// It acts as io.ReedSeeker.Read in any other aspect
func (r *ReedSeeker) Read(b []byte) (n int, err error) {
	n = int(C.archive_read_data(r.archive, unsafe.Pointer(&b[0]), C.size_t(cap(b))))
	if n == 0 {
		err = ErrArchiveEOF
	} else if 0 > n { // err
		err = codeToError(r.archive, ARCHIVE_FAILED)
		n = 0
	}
	r.index += int64(n)
	return
}

// Seek sets the offset for the next Read to offset
func (r *ReedSeeker) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case 0:
		abs = offset
	case 1:
		abs = int64(r.index) + offset
	case 2:
		abs = int64(len(r.buffer)) + offset
	default:
		return 0, errors.New("libarchive: SEEK [invalid whence]")
	}
	if abs < 0 {
		return 0, errors.New("libarchive: SEEK [negative position]")
	}
	r.index = abs
	return abs, nil
}

// Size returns compressed size of the current archive entry
func (r *ReedSeeker) Size() int {
	return int(C.archive_filter_bytes(r.archive, C.int(0)))
}

// Free frees the resources the underlying libarchive archive is using
// calling archive_read_free
func (r *ReedSeeker) Free() error {
	if C.archive_read_free(r.archive) == ARCHIVE_FATAL {
		return ErrArchiveFatal
	}
	return nil
}

// Close closes the underlying libarchive archive
// calling archive read_cloe
func (r *ReedSeeker) Close() error {
	if C.archive_read_close(r.archive) == ARCHIVE_FATAL {
		return ErrArchiveFatal
	}
	return nil
}
