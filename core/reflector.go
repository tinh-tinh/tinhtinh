package core

func Reflector[M any](key string, ctx Ctx) M {
	data, ok := ctx.GetMetadata(key).(M)
	if !ok {
		return *new(M)
	}
	return data
}
