package datetime_test

import (
	"github.com/leanovate/mite-go/datetime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseMinutes(t *testing.T) {
	expected := datetime.NewMinutes(23)
	actual, err := datetime.ParseMinutes("23m")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	actual, err = datetime.ParseMinutes("23m11s")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	_, err = datetime.ParseMinutes("1970-01-01")
	assert.NotNil(t, err)
}

func TestMinutes_Value(t *testing.T) {
	expected := 23
	actual := datetime.NewMinutes(23).Value()

	assert.Equal(t, expected, actual)
}

func TestMinutes_String(t *testing.T) {
	expected := "23m"
	actual := datetime.NewMinutes(23).String()

	assert.Equal(t, expected, actual)
}
