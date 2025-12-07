package storage

import (
	"mime"
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
			defer file.Close()

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

			uploadedFile, err := storeFile(field, fileHeader, r, opt)
			if err != nil {
				return nil, err
			}

			uploadedFile.MimeType = mediaType
			uploadedFile.Encoding = encode

			uploadFiles = append(uploadFiles, uploadedFile)
		}
	}

	return uploadFiles, nil
}
