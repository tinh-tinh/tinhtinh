package core

func Execution[S any](key CtxKey, ctx Ctx) *S {
	data, ok := ctx.Get(key).(*S)
	if !ok {
		return nil
	}
	return data
}
