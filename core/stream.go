package core

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/tinh-tinh/tinhtinh/v2/common"
)

type StreamableFileOptions struct {
	// FilePath is the path to the file to be streamed.
	FilePath string `json:"file_path" yaml:"file_path"`
	// Download indicates whether the file should be downloaded or displayed in the browser.
	Download bool `json:"download" yaml:"download"`
}

func (ctx *DefaultCtx) StreamableFile(filePath string, opts ...StreamableFileOptions) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file %s: %w", filePath, err)
	}

	// Detect MIME type based on file extension
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream" // Default fallback
	}
	ctx.w.Header().Set("Content-Type", mimeType)

	var option StreamableFileOptions
	if len(opts) > 0 {
		option = common.MergeStruct(opts...)
	}

	if option.FilePath == "" {
		option.FilePath = filePath
	}

	// Set Content-Disposition: attachment for download, otherwise inline
	dispositionType := "inline"
	if option.Download {
		dispositionType = "attachment"
	}

	ctx.w.Header().Set("Content-Disposition", dispositionType+"; filename=\""+filepath.Base(option.FilePath)+"\"")

	_, err = io.Copy(ctx.w, file)
	if err != nil {
		return fmt.Errorf("failed to stream file %s: %w", filePath, err)
	}

	return nil
}
