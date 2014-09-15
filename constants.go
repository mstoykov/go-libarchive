package archive

//#cgo pkg-config: libarchive
//#include<archive_entry.h>
import "C"

var (
	FileTypeRegFile = C.AE_IFREG
	FileTypeSymLink = C.AE_IFLNK
	FileTypeSocket  = C.AE_IFSOCK
	FileTypeCharDev = C.AE_IFCHR
	FileTypeBlkDev  = C.AE_IFBLK
	FileTypeDir     = C.AE_IFDIR
	FileTypeFIFO    = C.AE_IFIFO
)
