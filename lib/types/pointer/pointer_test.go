package pointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_To(t *testing.T) {
	t.Run("To with boolean", func(t *testing.T) {
		b := true
		assert.Equal(t, &b, To(b))
	})
	t.Run("To with string", func(t *testing.T) {
		str := "string"
		assert.Equal(t, &str, To(str))
	})
	t.Run("To with map", func(t *testing.T) {
		m := map[string]string{"string": "string"}
		assert.Equal(t, &m, To(m))
	})
	t.Run("To with struct", func(t *testing.T) {
		s := struct{ Field string }{Field: "string"}
		assert.Equal(t, &s, To(s))
	})
	t.Run("To with generic interface", func(t *testing.T) {
		var i interface{} = 10
		assert.Equal(t, &i, To(i))
	})
}

func Test_ValueOrZero(t *testing.T) {
	t.Run("Value with non nil boolean", func(t *testing.T) {
		b := true
		assert.Equal(t, b, ValueOrZero(&b))
	})
	t.Run("ValueOrZero with non nil string", func(t *testing.T) {
		str := "string"
		assert.Equal(t, str, ValueOrZero(&str))
	})
	t.Run("ValueOrZero with non nil map", func(t *testing.T) {
		m := map[string]string{"string": "string"}
		assert.Equal(t, m, ValueOrZero(&m))
	})
	t.Run("ValueOrZero with non nil struct", func(t *testing.T) {
		s := struct{ Field string }{Field: "string"}
		assert.Equal(t, s, ValueOrZero(&s))
	})
	t.Run("ValueOrZero with non nil generic interface", func(t *testing.T) {
		var i interface{} = 10
		assert.Equal(t, i, ValueOrZero(&i))
	})
	t.Run("ValueOrZero with nil boolean", func(t *testing.T) {
		var b *bool
		zeroValueOfBool := false
		assert.Equal(t, zeroValueOfBool, ValueOrZero(b))
	})
	t.Run("ValueOrZero with nil string", func(t *testing.T) {
		var str *string
		zeroValueOfString := ""
		assert.Equal(t, zeroValueOfString, ValueOrZero(str))
	})
	t.Run("ValueOrZero with nil map pointer", func(t *testing.T) {
		var m *map[string]string
		zeroValueOfMapPointer := map[string]string(nil)
		assert.Equal(t, zeroValueOfMapPointer, ValueOrZero(m))
	})
	t.Run("ValueOrZero with nil struct pointer", func(t *testing.T) {
		type anon struct {
			Field string
		}
		var s *anon
		zeroValueOfStruct := anon{}
		assert.Equal(t, zeroValueOfStruct, ValueOrZero(s))
	})
	t.Run("ValueOrZero with nil generic interface", func(t *testing.T) {
		var i *interface{}
		zeroValueOfGenericInterface := interface{}(nil)
		assert.Equal(t, zeroValueOfGenericInterface, ValueOrZero(i))
	})
}

func Test_ValueOrDefault(t *testing.T) {
	t.Run("ValueOrDefault with non nil boolean", func(t *testing.T) {
		b := true
		assert.Equal(t, b, ValueOrDefault(&b, false))
	})
	t.Run("ValueOrDefault with non nil string", func(t *testing.T) {
		str := "string"
		assert.Equal(t, str, ValueOrDefault(&str, "default"))
	})
	t.Run("ValueOrDefault with non nil map", func(t *testing.T) {
		m := map[string]string{"string": "string"}
		assert.Equal(t, m, ValueOrDefault(&m, map[string]string{"other_string": "other_string"}))
	})
	t.Run("ValueOrDefault with non nil anonymous struct", func(t *testing.T) {
		s := struct{ Field string }{Field: "string"}
		assert.Equal(t, s, ValueOrDefault(&s, struct{ Field string }{Field: "other_string"}))
	})
	t.Run("ValueOrDefault with non nil generic interface", func(t *testing.T) {
		var i interface{} = 10
		var j interface{} = "hello"
		assert.Equal(t, i, ValueOrDefault(&i, j))
	})
	t.Run("ValueOrDefault with nil boolean", func(t *testing.T) {
		var b *bool
		defaultValue := true
		assert.Equal(t, defaultValue, ValueOrDefault(b, defaultValue))
	})
	t.Run("ValueOrDefault with nil string", func(t *testing.T) {
		var str *string
		defaultValue := "string"
		assert.Equal(t, defaultValue, ValueOrDefault(str, defaultValue))
	})
	t.Run("ValueOrDefault with nil map pointer", func(t *testing.T) {
		var m *map[string]string
		defaultValue := map[string]string{"string": "string"}
		assert.Equal(t, defaultValue, ValueOrDefault(m, defaultValue))
	})
	t.Run("ValueOrDefault with nil struct pointer", func(t *testing.T) {
		type anon struct {
			Field string
		}
		var s *anon
		defaultValue := anon{Field: "string"}
		assert.Equal(t, defaultValue, ValueOrDefault(s, defaultValue))
	})
	t.Run("ValueOrDefault with nil generic interface", func(t *testing.T) {
		var i *interface{}
		defaultValue := interface{}("string")
		assert.Equal(t, defaultValue, ValueOrDefault(i, defaultValue))
	})
}
