package file

import (
	"fmt"
	"os"
)

func ClearFile(filePath string) error {
	// 方法2: 重新打开文件并截断
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	file.Close()

	fmt.Println("原文件已清空")
	return nil
}
