package tgbot

import (
	"unicode/utf8"

	"github.com/samber/lo"
)

func SplitMessagesAgainstLengthLimitIntoMessageGroups(originalSlice []string) [][]string {
	count := 0
	tempSlice := make([]string, 0)
	batchSlice := make([][]string, 0)

	for _, s := range originalSlice {
		tempSlice = append(tempSlice, s)
		count += utf8.RuneCountInString(s) + 20

		if count >= 4096 {
			tempSlice = lo.DropRight(tempSlice, 1)     // rollback the last append
			batchSlice = append(batchSlice, tempSlice) // commit the batch

			tempSlice = make([]string, 0) // reset the temp slice
			count = 0                     // reset the count

			tempSlice = append(tempSlice, s)        // re-append the last element
			count += utf8.RuneCountInString(s) + 20 // re-calculating the count
		}
	}

	// if there are still elements in the temp slice, append them to the batch
	if len(tempSlice) > 0 {
		batchSlice = append(batchSlice, tempSlice)
	}

	return batchSlice
}
