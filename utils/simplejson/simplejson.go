/*
this package is forked from https://github.com/bitly/go-simplejson by MIT license.
but fix https://github.com/bitly/go-simplejson/issues/42,
https://github.com/bitly/go-simplejson/issues/41
https://github.com/bitly/go-simplejson/issues/15,
https://github.com/bitly/go-simplejson/issues/55.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package simplejson

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/murphysecurity/murphysec/utils/must"
	"io"
	"log"
	"reflect"
	"strconv"
)

type JSON struct {
	data interface{}
}

func (j JSON) MarshalBinary() (data []byte, err error) {
	if s, ok := j.data.(string); ok {
		return []byte(s), nil
	}
	return j.MarshalJSON()
}

// MarshalJSON implements the json.Marshaller interface.
func (j *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(&j.data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (j *JSON) UnmarshalJSON(p []byte) error {
	dec := json.NewDecoder(bytes.NewBuffer(p))
	dec.UseNumber()
	return dec.Decode(&j.data)
}

func (j *JSON) MarshalString() string {
	if j == nil {
		return ""
	}
	if s, ok := j.data.(string); ok {
		return s
	}
	return string(must.A(j.MarshalJSON()))
}

// NewJSON returns a pointer to a new `JSON` object
// after unmarshalling `body` bytes
func NewJSON(body []byte) (*JSON, error) {
	j := new(JSON)
	err := j.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// MustJSON returns a pointer to a new `JSON` object
// after unmarshalling `body` bytes, panic if error occurred.
func MustJSON(body []byte) *JSON {
	j, err := NewJSON(body)
	if err != nil {
		panic(err)
	}
	return j
}

// IsNull returns true if the value is null.
func (j *JSON) IsNull() bool {
	return j.data == nil
}

// New returns a pointer to a new, empty `JSON` object
func New() *JSON {
	return &JSON{
		data: make(map[string]interface{}),
	}
}

// NewArray returns a pointer to a new, empty `JSON` array
func NewArray() *JSON {
	return &JSON{
		data: make([]interface{}, 0),
	}
}

// NewFrom returns a pointer to the json object converted from any other object.
func NewFrom(i interface{}) *JSON {
	if j, ok := i.(*JSON); ok {
		return j
	}
	b := func() []byte {
		if b, ok := i.([]byte); ok {
			return b
		}
		var err error
		b, err := json.Marshal(i)
		if err != nil {
			panic(err)
		}
		return b
	}()
	j, err := NewJSON(b)
	if err != nil {
		panic(err)
	}
	return j
}

// NewFromReader returns a *JSON by decoding from an io.Reader
func NewFromReader(r io.Reader) (*JSON, error) {
	j := new(JSON)
	dec := json.NewDecoder(r)
	dec.UseNumber()
	err := dec.Decode(&j.data)
	return j, err
}

// Interface returns the underlying data
func (j *JSON) Interface() interface{} {
	return j.data
}

// Set modifies `JSON` map by `key` and `value`
// Useful for changing single key/value in a `JSON` object easily.
func (j *JSON) Set(key string, val interface{}) {
	m, ok := j.CheckMap()
	if !ok {
		return
	}
	m[key] = val
}

// SetPath modifies `JSON`, recursively checking/creating map keys for the supplied path,
// and then finally writing in the value
func (j *JSON) SetPath(branch []string, val interface{}) {
	if len(branch) == 0 {
		j.data = val
		return
	}

	// in order to insert our branch, we need map[string]interface{}
	if _, ok := (j.data).(map[string]interface{}); !ok {
		// have to replace with something suitable
		j.data = make(map[string]interface{})
	}
	curr := j.data.(map[string]interface{})

	for i := 0; i < len(branch)-1; i++ {
		b := branch[i]
		// key exists?
		if _, ok := curr[b]; !ok {
			n := make(map[string]interface{})
			curr[b] = n
			curr = n
			continue
		}

		// make sure the value is the right sort of thing
		if _, ok := curr[b].(map[string]interface{}); !ok {
			// have to replace with something suitable
			n := make(map[string]interface{})
			curr[b] = n
		}

		curr = curr[b].(map[string]interface{})
	}

	// add remaining k/v
	curr[branch[len(branch)-1]] = val
}

// Del modifies `JSON` map by deleting `key` if it is present.
func (j *JSON) Del(key string) {
	m, ok := j.CheckMap()
	if !ok {
		return
	}
	delete(m, key)
}

// Get searches for the item as specified by the branch
// within a nested JSON and returns a new JSON pointer
// the pointer is always a valid JSON, allowing for chained operations
//
//	newJs := js.Get("top_level", "entries", 3, "dict")
func (j *JSON) Get(branch ...interface{}) *JSON {
	jin, ok := j.CheckGet(branch...)
	if ok {
		return jin
	}
	return &JSON{nil}
}

// getKey returns a pointer to a new `JSON` object
// for `key` in its `map` representation
// and a bool identifying success or failure
func (j *JSON) getKey(key string) (*JSON, bool) {
	m, ok := j.CheckMap()
	if ok {
		if val, ok := m[key]; ok {
			return &JSON{val}, true
		}
	}
	return nil, false
}

// getIndex returns a pointer to a new `JSON` object
// for `index` in its `array` representation
// and a bool identifying success or failure
func (j *JSON) getIndex(index int) (*JSON, bool) {
	a, ok := j.CheckArray()
	if ok {
		if len(a) > index {
			return &JSON{a[index]}, true
		}
	}
	return nil, false
}

// JSONMap guarantees the return of a `map[string]*JSON` (with optional default)
func (j *JSON) JSONMap(args ...map[string]*JSON) map[string]*JSON {
	var def map[string]*JSON

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("JSONMap() received too many arguments %d", len(args))
	}

	a, ok := j.CheckJSONMap()
	if ok {
		return a
	}

	return def
}

// JSONArray guarantees the return of a `[]*JSON` (with optional default)
func (j *JSON) JSONArray(args ...[]*JSON) []*JSON {
	var def []*JSON

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("JSONArray() received too many arguments %d", len(args))
	}

	a, ok := j.CheckJSONArray()
	if ok {
		return a
	}

	return def
}

// Map guarantees the return of a `map[string]interface{}` (with optional default)
//
// useful when you want to interate over map values in a succinct manner:
//
//	for k, v := range js.Get("dictionary").Map() {
//		fmt.Println(k, v)
//	}
func (j *JSON) Map(args ...map[string]interface{}) map[string]interface{} {
	var def map[string]interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Map() received too many arguments %d", len(args))
	}

	a, ok := j.CheckMap()
	if ok {
		return a
	}

	return def
}

// Array guarantees the return of a `[]interface{}` (with optional default)
//
// useful when you want to iterate over array values in a succinct manner:
//
//	for i, v := range js.Get("results").Array() {
//		fmt.Println(i, v)
//	}
func (j *JSON) Array(args ...[]interface{}) []interface{} {
	var def []interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Array() received too many arguments %d", len(args))
	}

	a, ok := j.CheckArray()
	if ok {
		return a
	}

	return def
}

// Bool guarantees the return of a `bool` (with optional default)
//
// useful when you explicitly want a `bool` in a single value return context:
//
//	myFunc(js.Get("param1").Bool(), js.Get("optional_param").Bool(true))
func (j *JSON) Bool(args ...bool) bool {
	var def bool

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Bool() received too many arguments %d", len(args))
	}

	b, ok := j.CheckBool()
	if ok {
		return b
	}

	return def
}

func (j *JSON) StringOrNil(args ...string) *string {
	s, ok := j.CheckString()
	if ok {
		return &s
	}

	return nil
}

// String guarantees the return of a `string` (with optional default)
//
// useful when you explicitly want a `string` in a single value return context:
//
//	myFunc(js.Get("param1").String(), js.Get("optional_param").String("my_default"))
func (j *JSON) String(args ...string) string {
	var def string

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("String() received too many arguments %d", len(args))
	}

	s, ok := j.CheckString()
	if ok {
		return s
	}

	return def
}

// Int guarantees the return of an `int` (with optional default)
//
// useful when you explicitly want an `int` in a single value return context:
//
//	myFunc(js.Get("param1").Int(), js.Get("optional_param").Int(5150))
func (j *JSON) Int(args ...int) int {
	var def int

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Int() received too many arguments %d", len(args))
	}

	i, ok := j.CheckInt()
	if ok {
		return i
	}

	return def
}

// Int64 guarantees the return of an `int64` (with optional default)
//
// useful when you explicitly want an `int64` in a single value return context:
//
//	myFunc(js.Get("param1").Int64(), js.Get("optional_param").Int64(5150))
func (j *JSON) Int64(args ...int64) int64 {
	var def int64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Int64() received too many arguments %d", len(args))
	}

	i, ok := j.CheckInt64()
	if ok {
		return i
	}

	return def
}

// Uint64 guarantees the return of an `uint64` (with optional default)
//
// useful when you explicitly want an `uint64` in a single value return context:
//
//	myFunc(js.Get("param1").Uint64(), js.Get("optional_param").Uint64(5150))
func (j *JSON) Uint64(args ...uint64) uint64 {
	var def uint64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Uint64() received too many arguments %d", len(args))
	}

	i, ok := j.CheckUint64()
	if ok {
		return i
	}

	return def
}

// Float64 guarantees the return of a `float64` (with optional default)
//
// useful when you explicitly want a `float64` in a single value return context:
//
//	myFunc(js.Get("param1").Float64(), js.Get("optional_param").Float64(5.150))
func (j *JSON) Float64(args ...float64) float64 {
	var def float64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Float64() received too many arguments %d", len(args))
	}

	f, ok := j.CheckFloat64()
	if ok {
		return f
	}

	return def
}

// CheckGet is like Get, except it also returns a bool
// indicating whenever the branch was found or not
// the JSON pointer mai be nil
//
//	newJs, ok := js.Get("top_level", "entries", 3, "dict")
func (j *JSON) CheckGet(branch ...interface{}) (*JSON, bool) {
	jin := j
	var ok bool
	for _, p := range branch {
		switch p := p.(type) {
		case string:
			jin, ok = jin.getKey(p)
		case int:
			jin, ok = jin.getIndex(p)
		default:
			ok = false
		}
		if !ok {
			return nil, false
		}
	}
	return jin, true
}

// CheckJSONMap returns a copy of a JSON map, but with values as Jsons
func (j *JSON) CheckJSONMap() (map[string]*JSON, bool) {
	m, ok := j.CheckMap()
	if !ok {
		return nil, false
	}
	jm := make(map[string]*JSON)
	for key, val := range m {
		jm[key] = &JSON{val}
	}
	return jm, true
}

// CheckJSONArray returns a copy of an array, but with each value as a JSON
func (j *JSON) CheckJSONArray() ([]*JSON, bool) {
	a, ok := j.CheckArray()
	if !ok {
		return nil, false
	}
	ja := make([]*JSON, len(a))
	for key, val := range a {
		ja[key] = &JSON{val}
	}
	return ja, true
}

// CheckMap type asserts to `map`
func (j *JSON) CheckMap() (map[string]interface{}, bool) {
	if j == nil {
		return nil, false
	}
	if m, ok := (j.data).(map[string]interface{}); ok {
		return m, true
	}
	return nil, false
}

// CheckArray type asserts to an `array`
func (j *JSON) CheckArray() ([]interface{}, bool) {
	if j == nil {
		return nil, false
	}
	if a, ok := (j.data).([]interface{}); ok {
		return a, true
	}
	return nil, false
}

// CheckBool type asserts to `bool`
func (j *JSON) CheckBool() (bool, bool) {
	if j == nil {
		return false, false
	}
	if s, ok := (j.data).(bool); ok {
		return s, true
	}
	return false, false
}

// CheckString type asserts to `string`
func (j *JSON) CheckString() (string, bool) {
	if j == nil {
		return "", false
	}
	if s, ok := (j.data).(string); ok {
		return s, true
	}
	return "", false
}

// CheckInt coerces into an int
func (j *JSON) CheckInt() (int, bool) {
	switch j.data.(type) {
	case json.Number:
		nr, err := j.data.(json.Number).Int64()
		return int(nr), err == nil
	case float32, float64:
		return int(reflect.ValueOf(j.data).Float()), true
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(j.data).Int()), true
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(j.data).Uint()), true
	}
	return 0, false
}

// CheckInt64 coerces into an int64
func (j *JSON) CheckInt64() (int64, bool) {
	switch j.data.(type) {
	case json.Number:
		nr, err := j.data.(json.Number).Int64()
		return nr, err == nil
	case float32, float64:
		return int64(reflect.ValueOf(j.data).Float()), true
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(j.data).Int(), true
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(j.data).Uint()), true
	}
	return 0, false
}

// CheckUint64 coerces into an uint64
func (j *JSON) CheckUint64() (uint64, bool) {
	switch j.data.(type) {
	case json.Number:
		nr, err := strconv.ParseUint(j.data.(json.Number).String(), 10, 64)
		return nr, err == nil
	case float32, float64:
		return uint64(reflect.ValueOf(j.data).Float()), true
	case int, int8, int16, int32, int64:
		return uint64(reflect.ValueOf(j.data).Int()), true
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(j.data).Uint(), true
	}
	return 0, false
}

// CheckFloat64 coerces into a float64
func (j *JSON) CheckFloat64() (float64, bool) {
	switch j.data.(type) {
	case json.Number:
		nr, err := j.data.(json.Number).Float64()
		return nr, err == nil
	case float32, float64:
		return reflect.ValueOf(j.data).Float(), true
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(j.data).Int()), true
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(j.data).Uint()), true
	}
	return 0, false
}

// Value returns json value, implement driver.Valuer interface.
func (j JSON) Value() (driver.Value, error) {
	if j.data == nil {
		return nil, nil
	}

	return j.MarshalJSON()
}

// Scan scans value into Json, implements sql.Scanner interface.
func (j *JSON) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return j.UnmarshalJSON(b)
}
