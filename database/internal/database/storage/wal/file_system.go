package wal

import (
	"fmt"
	"os"
	"sort"
)

func SegmentUpperBound(directory string, lastSegmentName string) (string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return "", fmt.Errorf("failed to scan WAL directory: %w", err)
	}

	filenames := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filenames = append(filenames, file.Name())
	}

	sort.Strings(filenames)
	idx := upperBound(filenames, lastSegmentName)
	if idx < len(filenames) {
		return filenames[idx], nil
	} else {
		return "", nil
	}
}

func SegmentLast(directory string) (string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return "", fmt.Errorf("failed to scan WAL directory: %w", err)
	}

	filename := ""
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		if file.IsDir() {
			continue
		}

		filename = file.Name()
		break
	}

	return filename, nil
}

func upperBound(array []string, target string) int {
	low, high := 0, len(array)-1

	for low <= high {
		mid := (low + high) / 2
		if array[mid] > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return low
}
