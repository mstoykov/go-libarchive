#include <archive.h>
#include <archive_entry.h>
ssize_t go_libarchive_reader(struct archive *a, void *client_data, const void **block);
ssize_t go_libarchive_open(struct archive *a, void *client_data);


