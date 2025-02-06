package core

import "slices"

type Metadata struct {
	Key   string
	Value interface{}
}

// SetMetadata creates a new Metadata object with the given key and value.
// The returned Metadata object can be passed to the Metadata method of a
// DynamicController to set metadata for the controller.
func SetMetadata(key string, value interface{}) *Metadata {
	return &Metadata{
		Key:   key,
		Value: value,
	}
}

// Metadata sets the given metadata for the controller. The metadata will be
// merged with any existing metadata on the controller.
func (controller *DynamicController) Metadata(meta ...*Metadata) Controller {
	controller.metadata = append(controller.metadata, meta...)
	return controller
}

// GetMetadata returns the value associated with the given key in the request context's metadata.
// If the given key is not present in the metadata, it returns nil.
func (ctx *DefaultCtx) GetMetadata(key string) interface{} {
	metaIdx := slices.IndexFunc(ctx.metadata, func(meta *Metadata) bool {
		return meta.Key == key
	})
	if metaIdx != -1 {
		return ctx.metadata[metaIdx].Value
	}
	return nil
}

// SetMetadata sets the given metadata for the request context. If the given
// metadata is empty, it will clear the request context's metadata.
func (ctx *DefaultCtx) SetMetadata(meta ...*Metadata) *DefaultCtx {
	if len(meta) == 0 {
		ctx.metadata = []*Metadata{}
		return ctx
	}
	ctx.metadata = meta
	return ctx
}
