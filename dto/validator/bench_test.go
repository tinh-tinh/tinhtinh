package validator_test

import (
	"testing"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

type Address struct {
	City    string `validate:"required,isAlpha"`
	ZipCode string `validate:"required,isNumber"`
}

type Profile struct {
	Email   string  `validate:"required,isEmail"`
	Age     int     `validate:"isInt"`
	Height  float64 `validate:"isFloat"`
	Address Address `validate:"nested"`
}

type Role struct {
	Code string `validate:"required,isAlphaNumeric"`
}

type User struct {
	ID      int       `validate:"required,isInt"`
	Profile *Profile  `validate:"required,nested"`
	Roles   []Role    `validate:"required,nested"`
	Joined  time.Time `validate:"isDate"`
}

type Payload struct {
	Users []User `validate:"required,nested"`
}

func buildDataset(n int) []*Payload {
	out := make([]*Payload, n)

	for i := 0; i < n; i++ {
		users := make([]User, 10)

		for j := 0; j < 10; j++ {
			users[j] = User{
				ID: j,
				Profile: &Profile{
					Email:  "invalid-email",
					Age:    70,
					Height: 170.5,
					Address: Address{
						City:    "HCM",
						ZipCode: "ABC",
					},
				},
				Roles: []Role{
					{Code: "ADMIN"},
					{Code: ""}, // ❌ sai
				},
				Joined: time.Now(),
			}
		}

		out[i] = &Payload{Users: users}
	}

	return out
}

func Benchmark_Validator_DeepNested(b *testing.B) {
	data := buildDataset(100)

	b.Run("v1_scanner", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			p := data[i%len(data)]
			_ = validator.Scanner(p)
		}
	})

	b.Run("v2_compiled", func(b *testing.B) {
		b.ReportAllocs()
		v := validator.Validator{}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			p := data[i%len(data)]
			_ = v.Validate(p)
		}
	})
}
