package post

import "github.com/tinh-tinh/tinhtinh/core"

const PostService core.Provide = "PostService"

func service(module *core.DynamicModule) *core.DynamicProvider {
	postSv := core.NewProvider(module)

	postSv.Set(PostService, nil)
	return postSv
}
