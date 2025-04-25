package utils

import (
	"crypto/rand"
	"encoding/hex"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
)

var allowedTypes = map[string]bool{
	"image/jpeg":             true,
	"image/png":              true,
	"text/html":              true,
	"text/css":               true,
	"application/javascript": true,
	"text/javascript":        true,
}

func ValidateFileType(fileHeader *multipart.FileHeader) bool {
	// 获取扩展名（统一小写）
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	// 通过扩展名检测
	mimeByExt := mime.TypeByExtension(ext)

	// 提取主MIME类型（忽略参数）
	getPrimaryMIME := func(mimeStr string) string {
		if mimeStr == "" {
			return ""
		}
		primary, _, _ := mime.ParseMediaType(mimeStr)
		return primary
	}
	primaryByExt := getPrimaryMIME(mimeByExt)

	return allowedTypes[primaryByExt]
}

func GenerateUUID() string {
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		panic("UUID生成失败: " + err.Error())
	}

	// 设置版本位（Version 4）
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// 设置变体位
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return hex.EncodeToString(uuid)
}

// 提取主MIME类型（兼容空值）
func GetPrimaryMIME(mimeStr string) string {
	if mimeStr == "" {
		return ""
	}
	// 解析MIME类型并剥离参数
	primary, _, _ := mime.ParseMediaType(mimeStr)
	return primary
}
