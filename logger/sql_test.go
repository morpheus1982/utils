package logger

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

type ExampleStruct struct {
	Name string
	Val  string
}

func (s ExampleStruct) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func format(v []byte, escaper string) string {
	return escaper + strings.ReplaceAll(string(v), escaper, "\\"+escaper) + escaper
}

func TestExplainSQL(t *testing.T) {
	type role string
	type password []byte
	var (
		tt     = time.Date(2020, time.February, 23, 11, 10, 10, 0, time.Local) //"2020-02-23 11:10:10"
		myrole = role("admin")
		pwd    = password([]byte("pass"))
		jsVal  = []byte(`{"Name":"test","Val":"test"}`)
		js     = JSON(jsVal)
		esVal  = []byte(`{"Name":"test","Val":"test"}`)
		es     = ExampleStruct{Name: "test", Val: "test"}
	)

	results := []struct {
		SQL           string
		NumericRegexp *regexp.Regexp
		Vars          []interface{}
		Result        string
	}{
		{
			SQL:    "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			Vars:   []interface{}{"jinzhu", 1, 999.99, true, []byte("12345"), tt, &tt, nil, "w@g.\"com", myrole, pwd},
			Result: `create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values ('jinzhu', 1, 999.990000, true, '12345', '2020-02-23 11:10:10', '2020-02-23 11:10:10', NULL, 'w@g."com', 'admin', 'pass')`,
		},
		{
			SQL:    "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			Vars:   []interface{}{"jinzhu?", 1, 999.99, true, []byte("12345"), tt, &tt, nil, "w@g.\"com", myrole, pwd},
			Result: `create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values ('jinzhu?', 1, 999.990000, true, '12345', '2020-02-23 11:10:10', '2020-02-23 11:10:10', NULL, 'w@g."com', 'admin', 'pass')`,
		},
		{
			SQL:    "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass, json_struct, example_struct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			Vars:   []interface{}{"jinzhu", 1, 999.99, true, []byte("12345"), tt, &tt, nil, "w@g.\"com", myrole, pwd, js, es},
			Result: fmt.Sprintf(`create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass, json_struct, example_struct) values ('jinzhu', 1, 999.990000, true, '12345', '2020-02-23 11:10:10', '2020-02-23 11:10:10', NULL, 'w@g."com', 'admin', 'pass', %v, %v)`, format(jsVal, `'`), format(esVal, `'`)),
		},
		{
			SQL:    "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass, json_struct, example_struct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			Vars:   []interface{}{"jinzhu", 1, 999.99, true, []byte("12345"), tt, &tt, nil, "w@g.\"com", myrole, pwd, &js, &es},
			Result: fmt.Sprintf(`create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass, json_struct, example_struct) values ('jinzhu', 1, 999.990000, true, '12345', '2020-02-23 11:10:10', '2020-02-23 11:10:10', NULL, 'w@g."com', 'admin', 'pass', %v, %v)`, format(jsVal, `'`), format(esVal, `'`)),
		},
	}

	for idx, r := range results {
		if result := ExplainSQL(r.SQL, r.Vars...); result != r.Result {
			t.Errorf("Explain SQL #%v expects %v, but got %v", idx, r.Result, result)
		} else {
			t.Logf("Explain SQL #%v: %v", idx, result)
		}
	}
}
