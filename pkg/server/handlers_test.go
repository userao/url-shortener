package server

import (
	"testing"
)

type TestData struct {
	data     any
	expected any
}

func TestValidateUrl(t *testing.T) {

	urls := []TestData{
		{"https://google.com", true},
		{"https://lms.ascon.ru/course/view.php?id=732#section-2", true},
		{"https://habr.com/ru/articles/541676/", true},
		{"htt://google.com", false},
		{"", false},
		{" ", false},
	}
	for _, url := range urls {
		result := validateUrl(url.data.(string))
		if result != url.expected.(bool) {
			t.Fatalf("For URL %s want %t, got %t", url.data.(string), url.expected.(bool), result)
		}
	}
}
