package file

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// 校验md5字符串
func GetMd5FromString(data string) string {
	h := md5.New()
	io.WriteString(h, data)
	sum := fmt.Sprintf("%x", h.Sum(nil))

	return sum
}

// 校验SHA256字符串
func GetSHA256FromString(data string) string {
	h := sha256.New()
	io.WriteString(h, data)
	sum := fmt.Sprintf("%x", data)

	return sum
}

// 校验文件md5
func GetMd5FromFile(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		return "", err
	}

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
	// todo JJY SUM的效果还需要再验证，其和write的关系是什么？目前的测试不符合预期
}

// 校验文件SHA256
func GetSHA256FromFile(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, nil
}

func CheckMd5Sum() {
	h := md5.New()
	io.WriteString(h, "The fog is getting thicker!")
	fmt.Printf("%x\n", h.Sum(nil))
	io.WriteString(h, "And Leon's getting laaarger!")
	fmt.Printf("%x\n", h.Sum(nil))

	h2 := md5.New()
	io.WriteString(h2, "And Leon's getting laaarger!")
	fmt.Printf("%x\n", h.Sum([]byte("The fog is getting thicker!")))
}

// http://www.codebaoku.com/it-go/it-go-211694.html