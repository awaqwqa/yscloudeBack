package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

func Md5Encrypt(input string) string {
	// 创建一个新的MD5哈希实例
	hasher := md5.New()
	// 将字符串写入到哈希实例中
	hasher.Write([]byte(input))
	// 计算哈希值
	hashBytes := hasher.Sum(nil)
	// 将字节哈希转换为16进制字符串
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}
func stringToNumericID(input string) (uint16, error) {
	// 创建一个新的SHA-1哈希实例
	hasher := sha1.New()
	// 将字符串写入到哈希实例中
	hasher.Write([]byte(input))
	// 计算哈希值
	hashBytes := hasher.Sum(nil)
	// 从哈希值的前两个字节中提取16位数字ID
	numericID := binary.BigEndian.Uint16(hashBytes[:2])
	return numericID, nil
}
func GenerateRandomKey() (string, error) {
	randomBytes := make([]byte, 32) // 生成32个字节的随机数
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	_, err = hash.Write(randomBytes)
	if err != nil {
		return "", err
	}

	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes)[32:], nil
}

// 生成file对应hash值
func HashFileSHA256(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建一个新的哈希对象
	hasher := sha256.New()
	// 使用文件内容更新哈希对象
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// 计算并返回最终的哈希值
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash), nil
}
func GenerateUniqueIntID() int64 {
	rand.Seed(time.Now().UnixNano())
	min := int64(10000000) // 最小值（8位数的最小值）
	max := int64(99999999) // 最大值（8位数的最大值）
	return rand.Int63n(max-min+1) + min
}
