package strutil

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// GenImageName 随机生成指定后缀的图片名
func GenImageName(ext string, width, height int) string {
	return fmt.Sprintf("%s_%dx%d.%s", uuid.New().String(), width, height, ext)
}

func GenFileName(ext string) string {
	return fmt.Sprintf("%s.%s", uuid.New().String(), ext)
}

func GenMediaObjectName(ext string, width, height int) string {
	var (
		mediaType = "common"
		fileName  = GenFileName(ext)
	)

	switch ext {
	case "png", "jpeg", "jpg", "gif", "webp", "svg", "ico":
		mediaType = "image"
		fileName = GenImageName(ext, width, height)
	case "mp3", "wav", "aac", "ogg", "flac":
		mediaType = "audio"
	case "mp4", "avi", "mov", "wmv", "mkv":
		mediaType = "video"
	}

	return fmt.Sprintf("media/%s/%s/%s", mediaType, time.Now().Format("200601"), fileName)
}

func GetTemplatePath(fileName string) (string, error) {
	// 获取当前工作目录
	currentDir, err := getGoModRoot()
	if err != nil {
		return "", err
	}

	// 构造相对路径
	templatePath := filepath.Join(currentDir, "template", fileName)
	return templatePath, nil
}

func getGoModRoot() (string, error) {
	// 从当前目录开始，向上查找 go.mod 文件
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// 查找 go.mod 文件
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// 如果已经到达根目录，停止查找
		if dir == filepath.Dir(dir) {
			break
		}

		// 向上移动到父目录
		dir = filepath.Dir(dir)
	}

	return "", fmt.Errorf("go.mod not found")
}
