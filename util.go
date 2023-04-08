package BaiduPanApi

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const chunkSize = 16 * 1024 * 1024 // 16MB

type chunkInfo struct {
	blockList []string
}

type FileChecksums struct {
	totalSize int64
	checksums []string
}

// VerifyFileExists 检查文件是否存在
func VerifyFileExists(filepath string) error {
	fileStat, err := os.Stat(filepath)
	if err != nil {
		return fmt.Errorf("failed to verify file exists: %w", err)
	} else if fileStat.IsDir() {
		return fmt.Errorf("specified path is a directory")
	}

	return nil
}

// ChecksumFile 将大文件按块计算MD5校验和
func ChecksumFile(filepath string) (FileChecksums, error) {
	err := VerifyFileExists(filepath)
	if err != nil {
		return FileChecksums{}, err
	}

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return FileChecksums{}, err
	}

	var checksums []string
	numChunks := (fileInfo.Size() + chunkSize - 1) / chunkSize
	file, err := os.Open(filepath)
	if err != nil {
		return FileChecksums{}, err
	}
	defer file.Close()

	buf := make([]byte, chunkSize)
	for i := int64(0); i < numChunks; i++ {
		n, err := file.ReadAt(buf, i*chunkSize)
		switch {
		case err == io.EOF || (err != nil && err.Error() == fmt.Sprintf("short read of %d bytes", chunkSize)):
			data := buf[:n]
			md5sum := md5.Sum(data)
			checksums = append(checksums, hex.EncodeToString(md5sum[:]))
		case err == nil:
			data := buf
			md5sum := md5.Sum(data)
			checksums = append(checksums, hex.EncodeToString(md5sum[:]))
		default:
			return FileChecksums{}, err
		}
	}
	return FileChecksums{totalSize: fileInfo.Size(), checksums: checksums}, nil
}
