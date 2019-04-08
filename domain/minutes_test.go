package domain_test

import (
	"github.com/leanovate/mite-go/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseMinutes(t *testing.T) {
	expected := domain.NewMinutes(23)
	actual, err := domain.ParseMinutes("23m")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	actual, err = domain.ParseMinutes("22m33s")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	_, err = domain.ParseMinutes("1970-01-01")
	assert.NotNil(t, err)
}

func TestMinutes_Value(t *testing.T) {
	expected := 23
	actual := domain.NewMinutes(23).Value()

	assert.Equal(t, expected, actual)
}

func TestMinutes_String(t *testing.T) {
	expected := "23m"
	actual := domain.NewMinutes(23).String()

	assert.Equal(t, expected, actual)
}
