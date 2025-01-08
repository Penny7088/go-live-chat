package strutil

import (
	"fmt"
	"math/big"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenValidateCode 生成数字验证码
func GenValidateCode(length int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	for i := 0; i < length; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[newRand.Intn(10)])
	}
	return sb.String()
}

// Random 生成随机字符串
func Random(length int) string {
	var result []byte
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyz")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}

// MtSubstr 字符串截取
func MtSubstr(value string, start, end int) string {

	if start > end {
		return ""
	}

	str := []rune(value)

	if length := len(str); end > length {
		end = length
	}

	return string(str[start:end])
}

func BoolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}

// FileSuffix 获取文件后缀名
func FileSuffix(filename string) string {
	return strings.TrimPrefix(path.Ext(filename), ".")
}

func NewMsgId() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func NewUuid() string {
	return uuid.New().String()
}

// StringToInt 将字符串转换为 int
func StringToInt(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("无法将字符串转换为整数: %v", err)
	}
	return num, nil
}

// StringToUint 将字符串转换为 uint
func StringToUint(s string) (uint, error) {
	num, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("无法将字符串转换为无符号整数: %v", err)
	}
	return uint(num), nil
}

// StringToBigInt 将字符串转换为 big.Int
func StringToBigInt(s string) (*big.Int, error) {
	bi := new(big.Int)
	_, success := bi.SetString(s, 10) // 设置为十进制
	if !success {
		return nil, fmt.Errorf("无法将字符串转换为大整数")
	}
	return bi, nil
}

// StringToBigUint 将字符串转换为 big.Uint
func StringToBigUint(s string) (*big.Int, error) {
	bi := new(big.Int)
	_, success := bi.SetString(s, 10) // 设置为十进制
	if !success {
		return nil, fmt.Errorf("无法将字符串转换为大无符号整数")
	}
	return bi, nil
}
