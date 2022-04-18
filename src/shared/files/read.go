package files

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/evenq/evenq-cli/src/shared/util"
)

func Prepare(path string) (io.Reader, int64, error) {
	originalFile, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}

	stat, err := originalFile.Stat()
	if err != nil {
		return nil, 0, err
	}

	//
	if !strings.HasSuffix(path, ".gz") {
		fmt.Println("we will compress your original file for faster upload")
		s := util.Spinner("Compressing...")
		defer s.Stop()
		zippedPath := path + ".gz"

		dest, err := os.Create(zippedPath)
		if err != nil {
			return nil, 0, err
		}

		wr, err := gzip.NewWriterLevel(dest, 9)
		if err != nil {
			return nil, 0, err
		}

		_, err = io.Copy(wr, originalFile)
		if err != nil {
			return nil, 0, err
		}

		if err := wr.Close(); err != nil {
			return nil, 0, err
		}

		if err := dest.Close(); err != nil {
			return nil, 0, err
		}

		if err := originalFile.Close(); err != nil {
			return nil, 0, err
		}

		compFile, err := os.Open(zippedPath)
		if err != nil {
			return nil, 0, err
		}

		stat, err := compFile.Stat()
		if err != nil {
			return nil, 0, err
		}

		return compFile, stat.Size(), nil
	}

	return originalFile, stat.Size(), nil
}
