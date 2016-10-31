package stringModule

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestUppercaseNominalCase(t *testing.T) {
  //  Given
  entity := stringService{}
  expected := "HELLO"

  //  When
  actual, err := entity.uppercase("hello")

  //  Then
  assert.Equal(t, expected, actual)
  assert.Nil(t, err)
}

func TestUppercaseEmptyString(t *testing.T) {
  //  Given
  entity := stringService{}
  expected := ""

  //  When
  actual, err := entity.uppercase("")

  //  Then
  assert.Equal(t, expected, actual)
  assert.Equal(t, ErrEmpty, err)
}

func TestCount(t *testing.T) {
  //  Given
  entity := stringService{}
  expected := 5

  //  When
  actual, _ := entity.count("hello")

  //  Then
  assert.Equal(t, expected, actual)
}
