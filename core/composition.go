package core

func Composition() *DynamicController {
	return &DynamicController{}
}

func (c *DynamicController) Composition(composer *DynamicController) *DynamicController {
	c.middlewares = append(c.middlewares, composer.middlewares...)
	c.metadata = append(c.metadata, composer.metadata...)
	c.Dtos = append(c.Dtos, composer.Dtos...)
	return c
}
