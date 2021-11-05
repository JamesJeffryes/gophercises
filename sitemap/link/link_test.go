package link

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestEx1(t *testing.T) {
	expected := []Link{
		{"/other-page", "A link to another page"},
	}
	parsed := ParseFile("../ex1.html")
	td.Cmp(t, parsed, expected)
}

func TestEx2(t *testing.T) {
	expected := []Link{
		{"https://www.twitter.com/joncalhoun", "Check me out on twitter"},
		{"https://github.com/gophercises", "Gophercises is on Github!"},
	}
	parsed := ParseFile("../ex2.html")
	td.Cmp(t, parsed, expected)
}

func TestEx3(t *testing.T) {
	expected := []Link{
		{"#", "Login"},
		{"/lost", "Lost? Need help?"},
		{"https://twitter.com/marcusolsson", "@marcusolsson"},
	}
	parsed := ParseFile("../ex3.html")
	td.Cmp(t, parsed, expected)
}

func TestEx4(t *testing.T) {
	expected := []Link{
		{"/dog-cat", "dog cat"},
	}
	parsed := ParseFile("../ex4.html")
	td.Cmp(t, parsed, expected)
}
