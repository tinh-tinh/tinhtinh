package storage

import (
	"net/http"
)

func HandleFiles(r *http.Request, opt UploadFileOption) ([]*File, error) {
	uploadFiles := []*File{}

	// Validate limit
	if err := ValidateLimit(opt.Limit, r); err != nil {
		return nil, err
	}

	if opt.FieldName == "" {
		opt.FieldName = "file"
	}

	r.FormFile(opt.FieldName)
	for field, files := range r.MultipartForm.File {
		if field != opt.FieldName {
			continue
		}

		for _, fileHeader := range files {
			if err := ValidateFilterFile(r, fileHeader, opt); err != nil {
				return nil, err
			}

			file, err := fileHeader.Open()
			if err != nil {
				return nil, err
			}

			rs, err := ReadFile(file)
			file.Close()
			if err != nil {
				return nil, err
			}

			uploadedFile, err := StoreFile(field, fileHeader, r, opt)
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

			uploadFiles = append(uploadFiles, uploadedFile)
		}
	}

	return uploadFiles, nil
}
