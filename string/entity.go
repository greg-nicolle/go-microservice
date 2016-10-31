package stringModule

import (
  "errors"
  "strings"
)

// StringService provides operations on strings.
type StringService interface {
  uppercase(string) (string, error)
  count(string) (int, error)
}

type stringService struct{}

func (stringService) uppercase(s string) (string, error) {
  if s == "" {
    return "", ErrEmpty
  }
  return strings.ToUpper(s), nil
}

func (stringService) count(s string) (int, error) {
  return len(s), nil
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// ServiceMiddleware is a chainable behavior modifier for StringService.
type ServiceMiddleware func(StringService) StringService
