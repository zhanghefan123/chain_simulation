package breakpoint_awareness

import (
	"fmt"
	"os"
)

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func GetAlreadyCaclulatedResult(resultDirPath string) (map[string]struct{}, error) {
	finalMapping := map[string]struct{}{}
	if DirExists(resultDirPath) {
		fmt.Printf("%s exists\n", resultDirPath)
		entries, err := os.ReadDir(resultDirPath)
		if err != nil {
			return nil, fmt.Errorf("read dir failed: %v", err)
		}
		for _, entry := range entries {
			finalMapping[entry.Name()] = struct{}{}
		}
	} else {
		fmt.Printf("%s not exists\n", resultDirPath)
	}

	return finalMapping, nil
}
