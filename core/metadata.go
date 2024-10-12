package core

import "slices"

type Metadata struct {
	Key   string
	Value interface{}
}

func SetMetadata(key string, value interface{}) *Metadata {
	return &Metadata{
		Key:   key,
		Value: value,
	}
}

func (controller *DynamicController) Metadata(meta ...*Metadata) *DynamicController {
	controller.metadata = append(controller.metadata, meta...)
	return controller
}

func (ctx *Ctx) GetMetadata(key string) interface{} {
	metaIdx := slices.IndexFunc(ctx.metadata, func(meta *Metadata) bool {
		return meta.Key == key
	})
	if metaIdx != -1 {
		return ctx.metadata[metaIdx].Value
	}
	return nil
}

func (ctx *Ctx) SetMetadata(meta ...*Metadata) *Ctx {
	if len(meta) == 0 {
		ctx.metadata = []*Metadata{}
		return ctx
	}
	ctx.metadata = meta
	return ctx
}
