package archive

// #include <archive.h>
import "C"
import (
	"errors"
	"fmt"
	"io"
)

const (
	ARCHIVE_EOF    = C.ARCHIVE_EOF
	ARCHIVE_OK     = C.ARCHIVE_OK
	ARCHIVE_RETRY  = C.ARCHIVE_RETRY
	ARCHIVE_WARN   = C.ARCHIVE_WARN
	ARCHIVE_FAILED = C.ARCHIVE_FAILED
	ARCHIVE_FATAL  = C.ARCHIVE_FATAL
)

var (
	ErrArchiveEOF   = io.EOF
	ErrArchiveFatal = errors.New("libarchive: FATAL [critical error, archive closing]")
)

func codeToError(archive *C.struct_archive, e int) error {
	switch e {
	case ARCHIVE_EOF:
		return ErrArchiveEOF
	case ARCHIVE_FATAL:
		return errors.New(fmt.Sprintf("libarchive: FATAL [%s]", errorString(archive)))
	case ARCHIVE_FAILED:
		return errors.New(fmt.Sprintf("libarchive: FAILED [%s]", errorString(archive)))
	case ARCHIVE_RETRY:
		return errors.New(fmt.Sprintf("libarchive: RETRY [%s]", errorString(archive)))
	case ARCHIVE_WARN:
		return errors.New(fmt.Sprintf("libarchive: WARN [%s]", errorString(archive)))
	}
	return nil
}

func errorString(archive *C.struct_archive) string {
	return C.GoString(C.archive_error_string(archive))
}
