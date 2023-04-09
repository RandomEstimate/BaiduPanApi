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
	totalSize  int64
	checksums  []string
	contentMd5 string
	sliceMd5   string
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

	// 计算整个文件的MD5
	fullMd5, err := calcFileMd5(filepath)
	if err != nil {
		fmt.Println("Error calculating full file MD5:", err)
		return FileChecksums{}, err
	}

	// 计算效验文件MD5
	headMd5, err := calcHead256KbMd5(filepath)
	if err != nil {
		fmt.Println("Error calculating head 256KB MD5:", err)
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
	return FileChecksums{totalSize: fileInfo.Size(), checksums: checksums, contentMd5: fullMd5, sliceMd5: headMd5}, nil
}

// 计算文件的MD5
func calcFileMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// 计算文件前256KB的MD5
func calcHead256KbMd5(filePath string) (string, error) {
	const size = 256 * 1024 // 256KB
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.CopyN(hash, file, size); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
