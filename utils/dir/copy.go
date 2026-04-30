package dir

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyDir(src string, dst string) error {
	// 校验源目录是否存在
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源目录信息失败: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("源路径不是目录: %s", src)
	}

	// 创建目标目录（继承源目录权限）
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 读取源目录所有内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("读取源目录失败: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归拷贝子目录
			if err = CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 拷贝文件
			if err = CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile 拷贝单个文件（保留权限）
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dstFile.Close()

	// 分块拷贝（大文件不会爆内存）
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("拷贝文件内容失败: %w", err)
	}

	// 获取源文件权限并设置到目标文件
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}
