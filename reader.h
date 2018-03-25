#include <archive.h>
#include <archive_entry.h>

ssize_t go_libarchive_open(struct archive *a, void *client_data);

ssize_t go_libarchive_read(struct archive *a, void *client_data, const void **block);

int64_t go_libarchive_seek(struct archive * a, void *client_data, int64_t request, int whence);
