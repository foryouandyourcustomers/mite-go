package date_test

import (
	"github.com/leanovate/mite-go/date"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseMinutes(t *testing.T) {
	expected := date.NewMinutes(23)
	actual, err := date.ParseMinutes("23m")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	actual, err = date.ParseMinutes("23m11s")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	_, err = date.ParseMinutes("1970-01-01")
	assert.NotNil(t, err)
}

func TestMinutes_Value(t *testing.T) {
	expected := 23
	actual := date.NewMinutes(23).Value()

	assert.Equal(t, expected, actual)
}

func TestMinutes_String(t *testing.T) {
	expected := "23m0s"
	actual := date.NewMinutes(23).String()

	assert.Equal(t, expected, actual)
}
