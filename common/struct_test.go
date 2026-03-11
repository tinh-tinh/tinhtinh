package common_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common"
)

func Test_Name(t *testing.T) {
	type Person struct{}
	require.Equal(t, "Person", common.GetStructName(Person{}))
	require.Equal(t, "Person", common.GetStructName(&Person{}))
}

func Test_Partial(t *testing.T) {
	type LargeStruct struct {
		ID             int     `json:"id,omitempty"`
		Name           string  `json:"name,omitempty"`
		Email          string  `json:"email,omitempty"`
		Age            int     `json:"age,omitempty"`
		Address        string  `json:"address,omitempty"`
		PhoneNumber    string  `json:"phone_number,omitempty"`
		IsActive       bool    `json:"is_active,omitempty"`
		Balance        float64 `json:"balance,omitempty"`
		Score          int     `json:"score,omitempty"`
		Department     string  `json:"department,omitempty"`
		Position       string  `json:"position,omitempty"`
		YearsEmployed  int     `json:"years_employed,omitempty"`
		Salary         float64 `json:"salary,omitempty"`
		VacationDays   int     `json:"vacation_days,omitempty"`
		EmployeeID     string  `json:"employee_id,omitempty"`
		SecurityLevel  int     `json:"security_level,omitempty"`
		LastLoginDate  string  `json:"last_login_date,omitempty"`
		PreferredShift string  `json:"preferred_shift,omitempty"`
		EmergencyPhone string  `json:"emergency_phone,omitempty"`
		BadgeNumber    string  `json:"badge_number,omitempty"`
	}

	largeInput := LargeStruct{
		ID:             0,
		Name:           "",
		Email:          "",
		Age:            0,
		Address:        "",
		PhoneNumber:    "",
		IsActive:       false,
		Balance:        0,
		Score:          0,
		Department:     "",
		Position:       "",
		YearsEmployed:  0,
		Salary:         0,
		VacationDays:   0,
		EmployeeID:     "",
		SecurityLevel:  0,
		LastLoginDate:  "",
		PreferredShift: "",
		EmergencyPhone: "",
		BadgeNumber:    "",
	}

	largeBefore, err := json.Marshal(&largeInput)
	require.Nil(t, err)
	require.Equal(t, `{}`, string(largeBefore))

	largePointer := common.PartialStruct(largeInput)
	largeAfter, _ := json.Marshal(&largePointer)
	require.Equal(t, `{"id":0,"name":"","email":"","age":0,"address":"","phone_number":"","is_active":false,"balance":0,"score":0,"department":"","position":"","years_employed":0,"salary":0,"vacation_days":0,"employee_id":"","security_level":0,"last_login_date":"","preferred_shift":"","emergency_phone":"","badge_number":""}`, string(largeAfter))
}

func Test_Pick(t *testing.T) {
	type TestStruct struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Email    string  `json:"email"`
		Age      int     `json:"age"`
		Balance  float64 `json:"balance"`
		IsActive bool    `json:"is_active"`
	}

	input := TestStruct{
		ID:       123,
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		Balance:  1000.50,
		IsActive: true,
	}

	// Test picking a subset of fields
	fields := []string{"Name", "Email", "Balance"}
	result := common.PickStruct(input, fields)

	// Marshal both original and picked structs to compare
	originalJSON, err := json.Marshal(input)
	require.NoError(t, err)
	require.Equal(t, `{"id":123,"name":"John Doe","email":"john@example.com","age":30,"balance":1000.5,"is_active":true}`, string(originalJSON))

	pickedJSON, err := json.Marshal(result)
	require.NoError(t, err)
	require.Equal(t, `{"name":"John Doe","email":"john@example.com","balance":1000.5}`, string(pickedJSON))

	// Test picking no fields
	emptyResult := common.PickStruct(input, []string{})
	emptyJSON, err := json.Marshal(emptyResult)
	require.NoError(t, err)
	require.Equal(t, `{}`, string(emptyJSON))

	// Test picking non-existent fields (should be ignored)
	invalidResult := common.PickStruct(input, []string{"NonExistent", "Name"})
	invalidJSON, err := json.Marshal(invalidResult)
	require.NoError(t, err)
	require.Equal(t, `{"name":"John Doe"}`, string(invalidJSON))
}

func Test_Omit(t *testing.T) {
	type TestStruct struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Email    string  `json:"email"`
		Age      int     `json:"age"`
		Balance  float64 `json:"balance"`
		IsActive bool    `json:"is_active"`
	}

	input := TestStruct{
		ID:       123,
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		Balance:  1000.50,
		IsActive: true,
	}

	// Test omitting a subset of fields
	fields := []string{"Age", "Balance", "IsActive"}
	result := common.OmitStruct(input, fields)

	// Marshal both original and omitted structs to compare
	originalJSON, err := json.Marshal(input)
	require.NoError(t, err)
	require.Equal(t, `{"id":123,"name":"John Doe","email":"john@example.com","age":30,"balance":1000.5,"is_active":true}`, string(originalJSON))

	omittedJSON, err := json.Marshal(result)
	require.NoError(t, err)
	require.Equal(t, `{"id":123,"name":"John Doe","email":"john@example.com"}`, string(omittedJSON))

	// Test omitting no fields (should return same as original)
	fullResult := common.OmitStruct(input, []string{})
	fullJSON, err := json.Marshal(fullResult)
	require.NoError(t, err)
	require.Equal(t, `{"id":123,"name":"John Doe","email":"john@example.com","age":30,"balance":1000.5,"is_active":true}`, string(fullJSON))

	// Test omitting all fields
	allFields := []string{"ID", "Name", "Email", "Age", "Balance", "IsActive"}
	emptyResult := common.OmitStruct(input, allFields)
	emptyJSON, err := json.Marshal(emptyResult)
	require.NoError(t, err)
	require.Equal(t, `{}`, string(emptyJSON))

	// Test omitting non-existent fields (should be ignored)
	invalidResult := common.OmitStruct(input, []string{"NonExistent", "Age"})
	invalidJSON, err := json.Marshal(invalidResult)
	require.NoError(t, err)
	require.Equal(t, `{"id":123,"name":"John Doe","email":"john@example.com","balance":1000.5,"is_active":true}`, string(invalidJSON))
}

func Test_AssertType(t *testing.T) {
	type Person struct {
		Name    string
		Age     int
		Address string
		Phone   string
	}

	type PartialPerson struct {
		Name    *string
		Age     *int
		Address *string
		Phone   *string
	}

	type PickPerson struct {
		Name string
		Age  int
	}

	type OmitPerson struct {
		Name    string
		Age     int
		Address string
	}

	person := Person{
		Name:    "John",
		Age:     30,
		Address: "123 Main St",
		Phone:   "555-1234",
	}

	// Test common.PartialStruct type assertion
	partial := common.PartialStruct(person)
	partialJSON, err := json.Marshal(partial)
	require.Nil(t, err)

	var partialPerson PartialPerson
	err = json.Unmarshal([]byte(partialJSON), &partialPerson)
	require.Nil(t, err)
	require.Equal(t, "John", *partialPerson.Name)
	require.Equal(t, 30, *partialPerson.Age)
	require.Equal(t, "123 Main St", *partialPerson.Address)
	require.Equal(t, "555-1234", *partialPerson.Phone)

	// Test common.PickStruct type assertion
	picked := common.PickStruct(person, []string{"Name", "Age"})
	pickedJSON, err := json.Marshal(picked)
	require.Nil(t, err)

	var pickPerson PickPerson
	err = json.Unmarshal([]byte(pickedJSON), &pickPerson)
	require.Nil(t, err)
	require.Equal(t, "John", pickPerson.Name)
	require.Equal(t, 30, pickPerson.Age)

	// Test common.OmitStruct type assertion
	omitted := common.OmitStruct(person, []string{"Phone"})
	omittedJSON, err := json.Marshal(omitted)
	require.Nil(t, err)

	var omitPerson OmitPerson
	err = json.Unmarshal([]byte(omittedJSON), &omitPerson)
	require.Nil(t, err)
	require.Equal(t, "John", omitPerson.Name)
	require.Equal(t, 30, omitPerson.Age)
	require.Equal(t, "123 Main St", omitPerson.Address)
}

func Test_MergeStruct(t *testing.T) {
	type Person struct {
		Name    string
		Age     int
		Address string
		Phone   string
	}
	abc := common.MergeStruct(Person{
		Name: "abc",
	}, Person{
		Age: 12,
	}, Person{
		Address: "44454545",
	})

	assert.Equal(t, "abc", abc.Name)
	assert.Equal(t, "44454545", abc.Address)
	assert.Equal(t, 12, abc.Age)
}
