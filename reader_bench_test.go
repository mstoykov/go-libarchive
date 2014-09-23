package archive_test

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"testing"

	ar "github.com/MStoykov/go-libarchive"
)

type ArchiveReader interface {
	io.Reader
	Next() (interface{}, error)
}

var buf = make([]byte, 4096) // probably fine

func BenchTestFuncLib(b *testing.B, file io.Reader) {
	reader, _ := ar.NewReader(file)
	defer reader.Free()
	defer reader.Close()
	runBenchTest(b, wrapLibArchive(reader))
}

func BenchTestFuncSTD(b *testing.B, file io.Reader) {
	reader := tar.NewReader(file)
	runBenchTest(b, wrapTarArchive(reader))
}

func wrapLibArchive(r *ar.Reader) ArchiveReader {
	return &libArchiveReaderWrapper{reader: r}
}

func wrapTarArchive(r *tar.Reader) ArchiveReader {
	return &libStdlibTarReaderWrapper{reader: r}
}

type libArchiveReaderWrapper struct {
	reader *ar.Reader
}

func (l *libArchiveReaderWrapper) Read(b []byte) (int, error) {
	return l.reader.Read(b)
}

func (l *libArchiveReaderWrapper) Next() (interface{}, error) {
	return l.reader.Next()
}

type libStdlibTarReaderWrapper struct {
	reader *tar.Reader
}

func (l *libStdlibTarReaderWrapper) Read(b []byte) (int, error) {
	return l.reader.Read(b)
}

func (l *libStdlibTarReaderWrapper) Next() (interface{}, error) {
	return l.reader.Next()
}

func runBenchTest(b *testing.B, a ArchiveReader) {
	totalBytesRead := new(int)
	defer func(t *int) {
		b.SetBytes(int64(*t))
	}(totalBytesRead)

	for {
		_, err := a.Next()
		if err != nil {
			return
		}

		var bytesread int
		for ; err == nil; bytesread, err = a.Read(buf) {
			*totalBytesRead += bytesread

		}
	}
}

func benchTemplate(b *testing.B, testFunc func(b *testing.B, file io.Reader)) {
	matches, err := filepath.Glob("./fixtures/bench[0-9]*.tar")
	if err != nil {
		b.Skipf("Error while getting benchmark fixtures:\n%s", err)
	}
	if len(matches) == 0 {
		b.Skip("No benchmark fixtures found!\nPut them in './fixtures/' with names like bench1.tar bench2.tar ...")
	}
	for _, match := range matches {
		b.StopTimer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			file, _ := os.Open(match)
			b.StartTimer()
			testFunc(b, file)
			b.StopTimer()
			file.Close()
		}
		b.StopTimer()
	}
}

func BenchmarkLibArchive(b *testing.B) {
	benchTemplate(b, BenchTestFuncLib)
}

func BenchmarkSTDArchive(b *testing.B) {
	benchTemplate(b, BenchTestFuncSTD)
}
