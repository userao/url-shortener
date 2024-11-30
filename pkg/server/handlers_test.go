package server

import (
	"errors"
	"reflect"
	"testing"
)

type TestData struct {
	data     any
	expected any
}

func TestValidateUrl(t *testing.T) {

	urls := []TestData{
		{"https://google.com", nil},
		{"https://lms.ascon.ru/course/view.php?id=732#section-2", nil},
		{"https://habr.com/ru/articles/541676/", nil},
		{"google.com", nil},
		{"http://google.com", nil},
		{"htt://google.com", errors.New("provided url not valid")},
		{"", errors.New("provided url not valid")},
		{" ", errors.New("provided url not valid")},
	}
	for _, url := range urls {
		result := validateUrl(url.data.(string))
		if reflect.TypeOf(result) != reflect.TypeOf(url.expected) {
			t.Fatalf("For URL %s want %v, got %v", url.data.(string), url.expected, result)
		}
	}
}
