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
	"sync"
	"unsafe"
)

var readers = make(map[uint64]*Reader)

var readersIndex uint64

var readersLock sync.Mutex

// Reader represents libarchive archive
type Reader struct {
	archive *C.struct_archive
	reader  io.ReadSeeker // the io.Reader from which we Read
	buffer  []byte        // buffer for the raw reading
	offset  int64         // current reading offset
	index   *uint64       // lookup index
}

// NewReader returns new Archive by calling archive_read_open
func NewReader(reader io.ReadSeeker) (r *Reader, err error) {
	r = new(Reader)

	readersLock.Lock()
	r.index = new(uint64)
	*r.index = readersIndex
	readers[readersIndex] = r
	readersIndex++
	readersLock.Unlock()

	r.buffer = make([]byte, 1024)
	r.archive = C.archive_read_new()
	C.archive_read_support_filter_all(r.archive)
	C.archive_read_support_format_all(r.archive)

	seekCb := (*C.archive_seek_callback)(C.go_libarchive_seek)
	C.archive_read_set_seek_callback(r.archive, seekCb)

	r.reader = reader

	e := C.go_libarchive_open(r.archive, unsafe.Pointer(r.index))

	err = codeToError(r.archive, int(e))
	return
}

//export myopen
func myopen(_ *C.struct_archive, _ unsafe.Pointer) C.int {
	// actually write something
	return C.ARCHIVE_OK
}

//export myclose
func myclose(_ *C.struct_archive, _ unsafe.Pointer) C.int {
	// actually write something
	return C.ARCHIVE_OK
}

//export myread
func myread(_ *C.struct_archive, clientData unsafe.Pointer, block unsafe.Pointer) C.ssize_t {
	readersLock.Lock()
	index := *(*uint64)(clientData)
	reader := readers[index]
	readersLock.Unlock()

	read, err := reader.reader.Read(reader.buffer)
	if err != nil && err != ErrArchiveEOF {
		// set error
		read = -1
	}

	*(*uintptr)(block) = uintptr(unsafe.Pointer(&reader.buffer[0]))

	return C.ssize_t(read)
}

//export myseek
func myseek(_ *C.struct_archive, clientData unsafe.Pointer, request C.int64_t, whence C.int) C.int64_t {
	readersLock.Lock()
	index := *(*uint64)(clientData)
	reader := readers[index]
	readersLock.Unlock()

	offset, err := reader.reader.Seek(int64(request), int(whence))
	if err != nil {
		return C.int64_t(0)
	}
	return C.int64_t(offset)
}

// Next calls archive_read_next_header and returns an
// interpretation of the Entry which is a wrapper around
// libarchive's archive_entry, or Err.
//
// ErrArchiveEOF is returned when there
// is no more to be read from the archive
func (r *Reader) Next() (Entry, error) {
	e := new(entryImpl)

	errno := int(C.archive_read_next_header(r.archive, &e.entry))
	err := codeToError(r.archive, errno)

	if err != nil {
		e = nil
	}

	return e, err
}

// Read calls archive_read_data which reads the current archive_entry.
// It acts as io.Reader.Read in any other aspect
func (r *Reader) Read(b []byte) (n int, err error) {
	n = int(C.archive_read_data(r.archive, unsafe.Pointer(&b[0]), C.size_t(cap(b))))
	if n == 0 {
		err = ErrArchiveEOF
	} else if 0 > n { // err
		err = codeToError(r.archive, C.ARCHIVE_FAILED)
		n = 0
	}
	r.offset += int64(n)
	return
}

// Seek sets the offset for the next Read to offset
func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case 0:
		abs = offset
	case 1:
		abs = r.offset + offset
	case 2:
		abs = int64(len(r.buffer)) + offset
	default:
		return 0, errors.New("libarchive: SEEK [invalid whence]")
	}
	if abs < 0 {
		return 0, errors.New("libarchive: SEEK [negative position]")
	}
	r.offset = abs
	return abs, nil
}

// Size returns compressed size of the current archive entry
func (r *Reader) Size() int {
	return int(C.archive_filter_bytes(r.archive, C.int(0)))
}

// Free frees the resources the underlying libarchive archive is using
// calling archive_read_free
func (r *Reader) Free() error {
	readersLock.Lock()
	delete(readers, *r.index)
	readersLock.Unlock()

	if C.archive_read_free(r.archive) == C.ARCHIVE_FATAL {
		return ErrArchiveFatal
	}
	return nil
}

// Close closes the underlying libarchive archive
// calling archive read_cloe
func (r *Reader) Close() error {
	if C.archive_read_close(r.archive) == C.ARCHIVE_FATAL {
		return ErrArchiveFatal
	}
	return nil
}
