package core

func Execution[S any](key any, ctx Ctx) *S {
	data, ok := ctx.Get(key).(*S)
	if !ok {
		return nil
	}
	return data
}
