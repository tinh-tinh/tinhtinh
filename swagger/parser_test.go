package swagger

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Recursive(t *testing.T) {
	t.Parallel()

	t.Run("test case", func(t *testing.T) {
		doc := NewSpecBuilder()
		mapper := recursiveParsePath(doc)
		jsonBytes, _ := json.Marshal(mapper)
		fmt.Println(string(jsonBytes))
	})
}
