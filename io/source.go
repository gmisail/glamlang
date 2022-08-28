package io

import (
	"fmt"
	"os"
	"strings"
)

type SourceFile struct {
	contents string
}

func CreateSource(fileName string) *SourceFile {
	fileData, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Errorf(err.Error())

		return nil
	}

	return &SourceFile{contents: string(fileData)}
}

func (s *SourceFile) GetLine(absolute int) string {
	start := absolute
	end := absolute

	// move cursor back until we find the beginning of file or newline
	for {
		char := s.contents[start]

		if char == '\n' || start == 0 {
			break
		}

		start--
	}

	// move cursor forward until we find the end of file or newline
	for {
		char := s.contents[end]

		if char == '\n' || end == len(s.contents)-1 {
			break
		}

		end++
	}

	return strings.TrimSpace(s.contents[start:end])
}
