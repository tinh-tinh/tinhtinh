package storage

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"net/http"
)

func HandleFile(r *http.Request, opt UploadFileOption) (*File, error) {
	// Validate limit
	if err := validateLimit(opt.Limit, r); err != nil {
		return nil, err
	}

	if len(r.MultipartForm.File) == 0 {
		return nil, errors.New("no file uploaded")
	}

	file, fileHeader, err := r.FormFile(opt.FieldName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := validateFilterFile(r, fileHeader, opt); err != nil {
		return nil, err
	}

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

	mimeType, err := detectAndValidateContentType(rs, fileHeader.Filename)
	if err != nil {
		return nil, err
	}
	mediaType, params, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return nil, err
	}
	encode := params["charset"]

	uploadedFile, err := storeFile(opt.FieldName, fileHeader, r, opt)
	if err != nil {
		return nil, err
	}

	uploadedFile.MimeType = mediaType
	uploadedFile.Encoding = encode

	return uploadedFile, nil
}
