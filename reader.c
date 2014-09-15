#include <archive.h>
#include <archive_entry.h>
#include "_cgo_export.h"

ssize_t readWrap(struct archive *a, void *client_data, const void **block) {
	return myread( a, client_data, block);
}

ssize_t go_libarchive_open(struct archive *a, void *client_data ) {
	return archive_read_open(a,client_data, myopen, readWrap, myclose);
}


