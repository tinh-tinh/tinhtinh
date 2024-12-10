package core

func Composition() *DynamicController {
	return &DynamicController{}
}

func (c *DynamicController) Composition(composer Controller) Controller {
	c.middlewares = append(c.middlewares, composer.getMiddlewares()...)
	c.metadata = append(c.metadata, composer.getMetadata()...)
	c.Dtos = append(c.Dtos, composer.GetDtos()...)
	return c
}
