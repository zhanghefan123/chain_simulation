package file

import (
	"fmt"
	"strings"
)

func ModifyYml(filePath, key, value string) error {
	fmt.Printf("modifyYml %s\n", filePath)
	// 1. 进行文件内容读取
	result, err := ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read configuration yml failed: %v", err)
	}
	// 2. 进行替换
	finalString := ""
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.Contains(line, key) {
			fmt.Println("line:", line)
			keyAndValue := strings.Split(line, ":")
			finalString += fmt.Sprintf("%s\n", keyAndValue[0]+": "+value)
		} else {
			finalString += fmt.Sprintf("%s\n", line)
		}
	}
	// 3. 写入文件之中
	err = WriteStringIntoFile(filePath, finalString)
	if err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}
	return nil
}
