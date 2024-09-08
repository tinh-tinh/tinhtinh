package post

import "github.com/tinh-tinh/tinhtinh/core"

const PostService core.Provide = "PostService"

func service(module *core.DynamicModule) *core.DynamicProvider {
	postSv := module.NewProvider(nil)

	return postSv
}
