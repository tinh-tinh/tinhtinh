package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func ValidateLimit(limit *UploadFileLimit, r *http.Request) error {
	if limit == nil {
		return nil
	}

	if limit.FileSize > 0 {
		err := r.ParseMultipartForm(limit.FileSize << 20)
		if err != nil {
			return err
		}
	}

	if limit.Fields > 0 {
		numFields := len(r.MultipartForm.File)
		if numFields > limit.Fields {
			errStr := fmt.Sprintf("number of fields exceeds limit %d", limit.Fields)
			return errors.New(errStr)
		}
	}

	return nil
}

func ValidateFilterFile(r *http.Request, file *multipart.FileHeader, opt UploadFileOption) error {
	if opt.FileFilter != nil && !opt.FileFilter(r, file) {
		return errors.New("file filter failed")
	}

	return nil
}

func DetectAndValidateContentType(file io.ReadSeeker, filename string) (string, error) {
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

func ReadFile(file multipart.File) (io.ReadSeeker, error) {
	var rs io.ReadSeeker
	if seeker, ok := file.(io.ReadSeeker); ok {
		rs = seeker
	} else {
		// fallback: copy to a buffer; consider limiting size and using a temp file for large uploads
		buf, err := io.ReadAll(io.LimitReader(file, 10<<20)) // cap to 10MB
		if err != nil {
			return nil, err
		}
		rs = bytes.NewReader(buf)
	}
	return rs, nil
}
