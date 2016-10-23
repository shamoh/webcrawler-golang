package gcrawler

import (
	"regexp"
	"strings"
)

var (
	// Patterns dedicated to Google Search page parsing.
	googleHrPattern = regexp.MustCompile(`<h3 class="r">(?P<url_part>.+?)</h3>`)
	googleUrlPattern = regexp.MustCompile(`url\?q=(?P<url_part>.+?)&amp;`)

	// Concrete content page parsing patterns.
	scriptPattern = regexp.MustCompile(`<script (?P<script_part>.+?)(</script>|/>)`)
	srcPattern = regexp.MustCompile(`src="(?P<src_part>.+?)"`)
)

type (
	Parser interface {
		// Get the information from the content source according to final implementation.
		Parse(rawContent string) []string
	}

	GoogleLinkParser struct {}

	LibraryParser struct {}
)
// Two-level parser implementation. Using two `Pattern` parses {@code content} string. The first
// pattern `firstLevelPattern` loops the `content` and `secondLevelPattern` finds one concrete string
// in the previous result.
func twoLevelParser(rawContent string, firstLevelPattern, secondLevelPattern *regexp.Regexp) []string {
	content := strings.Replace(rawContent, "\n", "", -1)
	firstResults := firstLevelPattern.FindAllStringSubmatch(content, -1)

	results := make([]string, 0)
	for _, result := range firstResults {
		secondResults := secondLevelPattern.FindAllStringSubmatch(result[1], 1)

		if len(secondResults) > 0 {
			results = append(results, secondResults[0][1])
		}
	}
	return results
}

// Implementation of the string parser which takes a `content` and parses google search results and `a` tag and its
// `href` attribute.
func (GoogleLinkParser) Parse(content string) []string {
	return twoLevelParser(content, googleHrPattern, googleUrlPattern)
}

// Implementation of the string parser which takes a `content` and parses `script` tag and its `src` attribute.
func (LibraryParser) Parse(content string) []string {
	libraryPath := twoLevelParser(content, scriptPattern, srcPattern)
	libraries := make([]string, len(libraryPath))
	for index, library := range libraryPath {
		// Take the name of the library as a last part of the library path after the slash.
		strippedPath := library[strings.LastIndexAny(library, "/") + 1:]

		// Take the name of the library before the question mark (strip the path parameters)
		var libraryName string
		questIndex := strings.Index(strippedPath, "?")
		if questIndex >= 0 {
			libraryName = strippedPath[0:questIndex]
		} else {
			libraryName = strippedPath
		}

		libraries[index] = libraryName
	}

	return libraries
}