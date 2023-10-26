package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
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
