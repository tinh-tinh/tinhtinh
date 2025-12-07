package storage

import (
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

	if err := ValidateFilterFile(r, fileHeader, opt); err != nil {
		return nil, err
	}

	rs, err := ReadFile(file)
	file.Close()
	if err != nil {
		return nil, err
	}

	uploadedFile, err := StoreFile(opt.FieldName, fileHeader, r, opt)
	if err != nil {
		return nil, err
	}

	mimeType, err := DetectAndValidateContentType(rs, fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	err = AppendMimeExtension(uploadedFile, mimeType)
	if err != nil {
		return nil, err
	}

	return uploadedFile, nil
}
