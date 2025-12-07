package storage

import (
	"errors"
	"net/http"
	"strconv"
)

func HandleFieldFiles(r *http.Request, opt UploadFileOption, fieldFiles ...FieldFile) (map[string][]*File, error) {
	result := make(map[string][]*File)

	// Validate limit
	if err := ValidateLimit(opt.Limit, r); err != nil {
		return nil, err
	}

	if opt.FieldName == "" {
		opt.FieldName = "file"
	}

	r.FormFile(opt.FieldName)
	if len(r.MultipartForm.File) == 0 {
		return nil, errors.New("no file uploaded")
	}

	fieldFileMap := make(map[string]int)
	for _, ff := range fieldFiles {
		fieldFileMap[ff.Name] = ff.MaxCount
	}

	for field, files := range r.MultipartForm.File {
		maxCount, exists := fieldFileMap[field]
		if !exists {
			continue
		}

		if maxCount > 0 && len(files) > maxCount {
			errStr := "number of files for field " + field + " exceeds limit " + strconv.Itoa(maxCount)
			return nil, errors.New(errStr)
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

			result[field] = append(result[field], uploadedFile)
		}
	}

	return result, nil
}
