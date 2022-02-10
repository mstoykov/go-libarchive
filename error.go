package archive

// #include <archive.h>
import "C"
import (
	"errors"
	"fmt"
	"io"
)

var (
	// ErrArchiveEOF ...
	ErrArchiveEOF   = io.EOF
	// ErrArchiveFatal ...
	ErrArchiveFatal = errors.New("libarchive: FATAL [critical error, archive closing]")
)

func codeToError(archive *C.struct_archive, e int) error {
	switch e {
	case C.ARCHIVE_EOF:
		return ErrArchiveEOF
	case C.ARCHIVE_FATAL:
		return fmt.Errorf("libarchive: FATAL [%s]", errorString(archive))
	case C.ARCHIVE_FAILED:
		return fmt.Errorf("libarchive: FAILED [%s]", errorString(archive))
	case C.ARCHIVE_RETRY:
		return fmt.Errorf("libarchive: RETRY [%s]", errorString(archive))
	case C.ARCHIVE_WARN:
		return fmt.Errorf("libarchive: WARN [%s]", errorString(archive))
	}
	return nil
}

func errorString(archive *C.struct_archive) string {
	return C.GoString(C.archive_error_string(archive))
}
