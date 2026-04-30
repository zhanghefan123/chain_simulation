package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFileWithName 拷贝文件到指定目录并重命名
func CopyFileWithName(src, destDir, newFileName string) error {
	// 检查源文件是否存在
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("源文件不存在: %s", src)
	}

	// 创建目标目录
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 构建目标文件路径
	destPath := filepath.Join(destDir, newFileName)

	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %v", err)
	}
	defer sourceFile.Close()

	// 创建目标文件
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer destFile.Close()

	// 拷贝内容
	bytesCopied, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("拷贝内容失败: %v", err)
	}

	// 确保数据写入磁盘
	if err = destFile.Sync(); err != nil {
		return fmt.Errorf("同步文件失败: %v", err)
	}

	// 可选：复制文件权限
	var srcInfo os.FileInfo
	if srcInfo, err = os.Stat(src); err == nil {
		os.Chmod(destPath, srcInfo.Mode())
	}

	fmt.Printf("已拷贝 %d 字节到 %s\n", bytesCopied, destPath)
	return nil
}
