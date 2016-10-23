package proxy

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
//  Given
  expected := []string{"Hello"}

  //  When
  actual := Split("Hello")

  //  Then
  assert.Equal(t, expected, actual)
}

func TestSplitWithTwoString(t *testing.T) {
//  Given
  expected := []string{"Hello","World"}

  //  When
  actual := Split("Hello,World")

  //  Then
  assert.Equal(t, expected, actual)
}
