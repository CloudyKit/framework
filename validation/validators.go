// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func Sub(runner func(At)) Tester {
	return func(c *Validator) {
		cc := *c
		value := c.Value
		prefix := c.prefix + c.Name
	restart:
		switch value.Kind() {
		case reflect.Array, reflect.Slice:
			length := value.Len()
			for i := 0; i < length; i++ {
				c.target = value.Index(i)
				c.prefix = prefix + fmt.Sprintf("[%d]", i)
				runner(c.Test)
			}
		case reflect.Struct, reflect.Map:
			c.target = value
			c.prefix = prefix
			runner(c.Test)
		case reflect.Ptr, reflect.Interface:
			value = value.Elem()
			goto restart
		}

		*c = cc
	}
}

func IsZero(v reflect.Value) bool {

	if kind := v.Kind(); kind == reflect.Struct {
		size := v.NumField()
		for i := 0; i < size; i++ {
			if IsZero(v.Field(i)) == false {
				return false
			}
		}
		return true
	} else if kind == reflect.Array {
		size := v.Len()
		for i := 0; i < size; i++ {
			if IsZero(v.Index(i)) == false {
				return false
			}
		}
		return true
	} else if kind == reflect.Bool {
		return v.Bool() == false
	} else if kind == reflect.String {
		return v.String() == ""
	} else if kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
		return v.Uint() == 0
	} else if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 {
		return v.Int() == 0
	} else if kind == reflect.Float32 || kind == reflect.Float64 {
		return v.Float() == 0
	} else if kind == reflect.Slice {
		return v.Len() == 0
	} else if kind == reflect.Invalid {
		return true
	}

	return v.IsNil()
}

func BeforeNow(msg string) Tester {
	return func(c *Validator) {
		timeNow := time.Now()
		if c.Value.Interface().(time.Time).After(timeNow) {
			c.Err(msg)
		}
	}
}

func AfterNow(msg string) Tester {
	return func(c *Validator) {
		timeNow := time.Now()
		if c.Value.Interface().(time.Time).Before(timeNow) {
			c.Err(msg)
		}
	}
}

func NoEmpty(msg string) Tester {
	return func(c *Validator) {
		if IsZero(c.Value) {
			c.Err(msg)
		}
	}
}

func Empty(msg string) Tester {
	return func(c *Validator) {
		if IsZero(c.Value) == false {
			c.Err(msg)
		}
	}
}

func OneOf(msg string, list ...interface{}) Tester {
	return func(c *Validator) {
		for i := 0; i < len(list); i++ {
			if reflect.DeepEqual(list[i], c.Value.Interface()) {
				return
			}
		}
		c.Err(msg)
		return
	}
}

func StringContains(msg string, item string) Tester {
	return func(c *Validator) {
		if strings.Contains(c.Value.String(), item) == false {
			c.Err(msg)
		}
	}
}

func SliceContains(msg string, item interface{}) Tester {
	return func(c *Validator) {
		size := c.Value.Len()
		for i := 0; i < size; i++ {
			if reflect.DeepEqual(c.Value.Index(i).Interface(), item) {
				return
			}
		}
		c.Err(msg)
		return
	}
}

func SameAs(msg string, FieldName string) Tester {
	return func(c *Validator) {
		if reflect.DeepEqual(c.Field(FieldName).Interface(), c.Value.Interface()) == false {
			c.Err(msg)
		}
	}
}

func MinLength(msg string, length int) Tester {
	return func(c *Validator) {
		if c.Value.Len() < length {
			c.Err(msg)
		}
	}
}

func MaxLength(msg string, length int) Tester {
	return func(c *Validator) {
		if c.Value.Len() > length {
			c.Err(msg)
		}
	}
}

func MinUint(msg string, i uint64) Tester {
	return func(c *Validator) {
		if c.Value.Uint() < i {
			c.Err(msg)
		}
	}
}

func MaxUint(msg string, i uint64) Tester {
	return func(c *Validator) {
		if c.Value.Uint() > i {
			c.Err(msg)
		}
	}
}

func MinInt(msg string, i int64) Tester {
	return func(c *Validator) {
		if c.Value.Int() < i {
			c.Err(msg)
		}
	}
}

func MaxInt(msg string, i int64) Tester {
	return func(c *Validator) {
		if c.Value.Int() > i {
			c.Err(msg)
		}
	}
}

func MinFloat(msg string, i float64) Tester {
	return func(c *Validator) {
		if c.Value.Float() < i {
			c.Err(msg)
		}
	}
}

func MaxFloat(msg string, i float64) Tester {
	return func(c *Validator) {
		if c.Value.Float() > i {
			c.Err(msg)
		}
	}
}

var Email = NewRegexValidator("^(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])$")

func NewRegexValidator(pattern string) func(msg string) func(*Validator) {
	regExp := regexp.MustCompile(pattern)
	return func(msg string) func(*Validator) {
		return func(c *Validator) {
			str := fmt.Sprint(c.Value.Interface())
			if !regExp.MatchString(str) {
				c.Err(msg)
			}
			return
		}
	}
}
