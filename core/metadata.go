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

func (controller *DynamicController) Metadata(meta *Metadata) *DynamicController {
	controller.metadata = append(controller.metadata, meta)
	return controller
}

func (controller *DynamicController) GetMetadata(key string) interface{} {
	metaIdx := slices.IndexFunc(controller.metadata, func(meta *Metadata) bool {
		return meta.Key == key
	})
	if metaIdx != -1 {
		return controller.metadata[metaIdx].Value
	}
	return nil
}
