package utils

import (
	"regexp"
	"strings"
)

func matchCompatibleSubtitles(text string, subsRe *regexp.Regexp) []string {
	splitRe := regexp.MustCompile("(?i)compativel")

	// Split text by "compativel" word
	textParts := splitRe.Split(text, -1)

	subtitles := subsRe.FindAllString(textParts[0], -1)

	// For each part of the text, find subtitles and return only the last
	for i := 1; i < len(textParts)-1; i++ {
		matches := subsRe.FindAllString(textParts[i], -1)

		if len(matches) > 1 {
			subtitles = append(subtitles, matches[len(matches)-1])
		}
	}

	return subtitles
}

func FindSubtitlesInText(text, title string) []string {
	text = RemoveAccents(text)

	titleRe := regexp.MustCompile(`(?mi)\w+`)
	title = strings.Join(titleRe.FindAllString(title, -1), "[. ]")

	var subsRe = regexp.MustCompile(`(?mi)^(` + title + `.+)`)

	if strings.Contains(strings.ToLower(text), "compativel") {
		matches := matchCompatibleSubtitles(text, subsRe)

		return matches
	}

	return subsRe.FindAllString(text, -1)
}
