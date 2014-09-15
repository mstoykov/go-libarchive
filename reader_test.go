package archive

import (
	"bytes"
	"os"
	"testing"
)

func TestNewArchive(t *testing.T) {
	test_file, err := os.Open("./fixtures/test.tar")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

	archive, err := NewArchive(test_file)
	if err != nil {
		t.Fatalf("Error on creating Archive from a io.Reader:\n %s", err)
	}
	defer func() {
		err := archive.Free()

		if err != nil {
			t.Fatalf("Error on archive Free:\n %s", err)
		}
	}()

	defer func() {
		err := archive.Close()
		if err != nil {
			t.Fatalf("Error on archive Close:\n %s", err)
		}
	}()
	entry, err := archive.Next()
	if err != nil {
		t.Fatalf("got error on archive.Next():\n%s", err)
	}
	name := entry.PathName()
	if name != "a" {
		t.Fatalf("got %s expected %s as Name of the first entry", name, "a")
	}

	b := make([]byte, 512)
	size, err := archive.Read(b)
	if err != nil {
		t.Fatalf("got error on archive.Read():\n%s", err)
	}
	if size != 14 {
		t.Fatalf("got %d as size of the read but expected %d", size, 14)
	}

	expectedContent := []byte("Sha lalal lal\n")
	if !bytes.Equal((b[:size]), expectedContent) {
		t.Fatalf("The contents:\n [%s] are not the expectedContent:\n [%s]", b[:size], expectedContent)
	}

	_, err = archive.Next()
	if err != ErrArchiveEOF {
		t.Fatalf("Expected EOF on second archive.Next() got err :\n %s", err)
	}
}

func TestTwoReaders(t *testing.T) {
	test_file, err := os.Open("./fixtures/test.tar")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

	_,err = NewArchive(test_file)

	test_file2, err := os.Open("./fixtures/test2.tar")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

	_, err = NewArchive(test_file2)
	if err != nil {
		t.Fatalf("Error on creating Archive from a io.Reader:\n %s", err)
	}

}


