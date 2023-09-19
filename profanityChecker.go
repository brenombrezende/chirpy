package main

import (
	"strings"
)

func profanityChecker(s string) string {

	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	separated := strings.Split(s, " ")
	for i, word := range separated {
		for _, curse := range profanity {
			if strings.ToLower(word) == curse {
				separated[i] = "****"
			}
		}
	}
	return strings.Join(separated, " ")
}
