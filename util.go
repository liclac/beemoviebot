package main

import (
	"strings"
)

func MakeChunks(s string, max int) []string {
	delimiter := "\n\n"
	paragraphs := strings.Split(s, delimiter)

	var chunks []string
	var chunk string
	for _, p := range paragraphs {
		if len(chunk)+len(p)+len(delimiter) > max {
			chunks = append(chunks, chunk)
			chunk = ""
		}
		chunk += p + delimiter
	}
	chunks = append(chunks, chunk)

	return chunks
}
