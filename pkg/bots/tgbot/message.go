package tgbot

import (
	"unicode/utf8"
)

func SplitMessagesAgainstLengthLimitIntoMessageGroups(originalSlice []string) [][]string {
	count := 0
	tempSlice := make([]string, 0)
	batchSlice := make([][]string, 0)

	for _, s := range originalSlice {
		currentLength := utf8.RuneCountInString(s) + 20

		// If the current message itself exceeds the limit, handle it separately
		if currentLength >= MessageLengthLimit {
			// If tempSlice is not empty, save the current batch first
			if len(tempSlice) > 0 {
				batchSlice = append(batchSlice, tempSlice)
				tempSlice = make([]string, 0)
			}
			// Add the oversized message as a separate batch
			batchSlice = append(batchSlice, []string{s})
			count = 0
			continue
		}

		// If adding the current message would exceed the limit
		if count+currentLength >= MessageLengthLimit {
			// Save the current batch
			if len(tempSlice) > 0 {
				batchSlice = append(batchSlice, tempSlice)
			}
			// Start a new batch
			tempSlice = []string{s}
			count = currentLength
		} else {
			// Add to current batch
			tempSlice = append(tempSlice, s)
			count += currentLength
		}
	}

	// Handle the last batch
	if len(tempSlice) > 0 {
		batchSlice = append(batchSlice, tempSlice)
	}

	return batchSlice
}
