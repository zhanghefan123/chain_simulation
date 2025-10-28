package file

import (
	"fmt"
	"strings"
)

func ModifyYml(filePath string, mapping map[string]string) error {
	fmt.Printf("modifyYml %s\n", filePath)
	var modified = false
	// 1. 进行文件内容读取
	result, err := ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read configuration yml failed: %v", err)
	}
	// 2. 进行替换
	finalString := ""
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		// 进行 configurationSetting 之中的值的遍历
		var findCorrespondingLine = false
		for key, value := range mapping {
			if strings.Contains(line, key) {
				keyAndValue := strings.Split(line, ":")
				finalString += fmt.Sprintf("%s\n", keyAndValue[0]+": "+value)
				modified = true
				findCorrespondingLine = true
				break
			}
		}
		if !findCorrespondingLine {
			finalString += fmt.Sprintf("%s\n", line)
		}
	}
	if !modified && (len(mapping) > 0) {
		return fmt.Errorf("no key modified")
	}
	// 3. 写入文件之中
	err = WriteStringIntoFile(filePath, finalString)
	if err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}
	return nil
}
