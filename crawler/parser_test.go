package gcrawler

import (
	"testing"
	"io/ioutil"
)

func TestGoogleLinkParser(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/parser.html")
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	parser := GoogleLinkParser{}
	links := parser.Parse(content)

	if len(links) != 3 {
		t.Fatalf("There are not 3 parsed links as expected. Count: %b\n", len(links))
	}

	assertString(t, "http://stis.ping-pong.cz/", links[0]);
	assertString(t, "http://www.pingpong.cz/", links[1]);
	assertString(t, "http://www.pinces.cz/", links[2]);
}

func TestGoogleLinkParserEmpty(t *testing.T) {
	parser := GoogleLinkParser{}
	links := parser.Parse("")

	if len(links) != 0 {
		t.Fatalf("There are not 0 parsed links as expected. Count: %b\n", len(links))
	}
}

func TestLibraryParser(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/parser.html")
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	parser := LibraryParser{}
	links := parser.Parse(content)

	if len(links) != 5 {
		t.Fatalf("There are not 5 parsed links as expected. Count: %b\n", len(links))
	}

	assertString(t, "sdk.js#xfbml=1&amp;version=v2.0", links[0]);
	assertString(t, "jquery-1.12.4.js", links[1]);
	assertString(t, "jquery-1.12.6.js", links[2]);
	assertString(t, "jquery-1.12.4.js", links[3]);
	assertString(t, "jquery-1.12.5.js", links[4]);
}

func TestLibraryParserEmpty(t *testing.T) {
	parser := LibraryParser{}
	links := parser.Parse("")

	if len(links) != 0 {
		t.Fatalf("There are not 0 parsed links as expected. Count: %b\n", len(links))
	}
}

func assertString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Expected string: `%s`, Actual string: `%s`", expected, actual)
	}
}