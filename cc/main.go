package main

import (
	"strings"
)

func JustifyText(wordSequence []string, targetLineLength int) (res []string) {
	if len(wordSequence) == 0 {
		return res
	}
	curLineStartIdx, curLineEndIdx, curLineSpaceLeft := 0, 0, targetLineLength-len(wordSequence[0])
	for i := 1; i < len(wordSequence); i++ {
		word := wordSequence[i]
		if curLineSpaceLeft-1-len(word) >= 0 {
			curLineEndIdx = i
			curLineSpaceLeft -= 1 + len(word)
		} else {
			res = append(res, generateLine(wordSequence, curLineStartIdx, curLineEndIdx, curLineSpaceLeft, targetLineLength))
			curLineStartIdx, curLineEndIdx, curLineSpaceLeft = i, i, targetLineLength-len(word)
		}
	}
	res = append(res, generateLine(wordSequence, curLineStartIdx, curLineEndIdx, curLineSpaceLeft, targetLineLength))
	return
}

func generateLine(wordSequence []string, curLineStartIdx int, curLineEndIdx int, curLineSpaceLeft int, targetLineLength int) (res string) {
	first_space, intermediate_spaces, last_space := 1, 1, 0
	if curLineEndIdx == len(wordSequence)-1 || curLineStartIdx == curLineEndIdx {
		last_space = curLineSpaceLeft
	} else {
		first_space += curLineSpaceLeft/(curLineEndIdx-curLineStartIdx) + curLineSpaceLeft%(curLineEndIdx-curLineStartIdx)
		intermediate_spaces += curLineSpaceLeft / (curLineEndIdx - curLineStartIdx)
	}

	//Add first word
	res = ""
	res = res + wordSequence[curLineStartIdx]
	//Add spaces and rest of the words
	for i, cur_space := curLineStartIdx+1, first_space; i <= curLineEndIdx; i, cur_space = i+1, intermediate_spaces {
		res = res + strings.Repeat("-", cur_space) + wordSequence[i]
	}
	//Add last space
	res = res + strings.Repeat("-", last_space)
	return
}

func reverseString(word string) (res string) {
	// byteArr := []byte(word)
	// for i, j := 0, len(byteArr)-1; i < j; i, j = i+1, j-1 {
	// 	temp := byteArr[i]
	// 	byteArr[i] = byteArr[j]
	// 	byteArr[j] = temp
	// }

	for _, c := range word {
		res = res + string(c)
	}
	return
}

func main() {
}
