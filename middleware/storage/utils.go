package storage

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func validateLimit(opt *UploadFileLimit, r *http.Request) error {
	if opt == nil {
		return nil
	}

	if opt.FileSize > 0 {
		err := r.ParseMultipartForm(opt.FileSize)
		if err != nil {
			return err
		}
	}

	if opt.Fields > 0 {
		numFields := len(r.MultipartForm.File)
		if numFields > opt.Fields {
			errStr := fmt.Sprintf("number of fields exceeds limit %d", opt.Fields)
			return errors.New(errStr)
		}
	}

	return nil
}

func validateFilterFile(r *http.Request, file *multipart.FileHeader, opt UploadFileOption) error {
	if opt.FileFilter != nil && !opt.FileFilter(r, file) {
		return errors.New("file filter failed")
	}

	if opt.Limit != nil && opt.Limit.FileSize > 0 {
		if file.Size > opt.Limit.FileSize {
			errStr := fmt.Sprintf("file size exceeds limit %d bytes", opt.Limit.FileSize)
			return errors.New(errStr)
		}
	}

	return nil
}

func detectAndValidateContentType(file io.ReadSeeker, filename string) (string, error) {
	const sniffLen = 512
	buf := make([]byte, sniffLen)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read for sniffing: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("reset reader: %w", err)
	}

	detected := http.DetectContentType(buf[:n])

	if detected == "application/octet-stream" && filename != "" {
		if ext := filepath.Ext(filename); ext != "" {
			if byExt := mime.TypeByExtension(ext); byExt != "" {
				detected = byExt
			}
		}
	}

	return detected, nil
}
