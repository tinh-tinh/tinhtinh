package storage

import (
	"mime"
	"net/http"
)

func HandleFile(r *http.Request, opt UploadFileOption) (*File, error) {
	// Validate limit
	if err := ValidateLimit(opt.Limit, r); err != nil {
		return nil, err
	}

	if opt.FieldName == "" {
		opt.FieldName = "file"
	}

	file, fileHeader, err := r.FormFile(opt.FieldName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := ValidateFilterFile(r, fileHeader, opt); err != nil {
		return nil, err
	}

	rs, err := ReadFile(file)
	if err != nil {
		return nil, err
	}

	mimeType, err := DetectAndValidateContentType(rs, fileHeader.Filename)
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
