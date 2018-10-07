package archive

/*
#include <stdlib.h>
*/
import "C"

import (
	"time"
)

func (e *entryInfo) ModTime() time.Time {
	return time.Unix(int64(e.stat.st_mtimespec.tv_sec), int64(e.stat.st_mtimespec.tv_nsec))
}
