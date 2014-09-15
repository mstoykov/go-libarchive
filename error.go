package archive

// #include <archive.h>
import "C"
import (
	"errors"
	"io"
)


const (
	ARCHIVE_EOF   = C.ARCHIVE_EOF
	ARCHIVE_OK    = C.ARCHIVE_OK
	ARCHIVE_RETRY = C.ARCHIVE_RETRY
	ARCHIVE_WARN  = C.ARCHIVE_WARN
	ARCHIVE_FATAL = C.ARCHIVE_FATAL
)

var (
	ErrArchiveEOF   = io.EOF 
	ErrArchiveRetry = errors.New("libarchive: RETRY [operation failed but can be retried]")
	ErrArchiveWarn  = errors.New("libarchive: WARN [success but non-critical error]")
	ErrArchiveFatal = errors.New("libarchive: FATAL [critical error, archive closing]")
)

func codeToError(e int) error {
	switch e {
	case ARCHIVE_EOF:
		return  ErrArchiveEOF
	case ARCHIVE_FATAL:
		return  ErrArchiveFatal
	case ARCHIVE_RETRY:
		return  ErrArchiveRetry
	case ARCHIVE_WARN:
		return  ErrArchiveWarn
	}
	return nil
}
