package zepto

import (
	"io"
	"mime/multipart"
	"os"
)

type File struct {
	multipart.File
	*multipart.FileHeader
}

func (c *Context) PostFile(key string) (File, error) {
	file, header, err := c.Request.FormFile(key)
	if err != nil {
		return File{}, err
	}

	return File{File: file, FileHeader: header}, nil
}

func (c *Context) SaveFile(file File, path string) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, file.File)
	if err != nil {
		return err
	}

	return nil
}
