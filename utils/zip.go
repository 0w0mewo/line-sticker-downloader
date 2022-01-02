package utils

import (
	"archive/zip"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type ZipFile struct {
	path   string
	zw     *zip.Writer
	logger *logrus.Logger
}

func NewZip(path string) *ZipFile {
	var p string
	// check if the given path to zip file has extension of zip
	// add the extension if not
	if filepath.Ext(path) != ".zip" {
		p = path + ".zip"
	} else {
		p = path
	}

	// ensure the zip file path is abs
	p, err := filepath.Abs(p)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	// fd for zip
	zfd, err := os.Create(p)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	return &ZipFile{
		path:   p,
		zw:     zip.NewWriter(zfd),
		logger: logger,
	}

}

func (z *ZipFile) Zip(path string) error {
	return filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// ensure path is abs path
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		// file to be zipped
		fd, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fd.Close()

		// convert to zip file-rooted relative path
		f, err := filepath.Rel(filepath.Dir(z.path), path)
		if err != nil {
			return err
		}

		// zip it
		zipped, err := z.zw.Create(f)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipped, fd)

		return err

	})
}

func (z *ZipFile) Close() error {
	return z.zw.Close()

}
