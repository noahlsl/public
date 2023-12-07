package iox

import (
	"io"
	"mime/multipart"
)

// File2Reader 将 multipart.File 转换为 io.Reader
func File2Reader(file multipart.File, size int64) io.Reader {
	// 创建一个读取器，通过读取 multipart.File 的内容来实现
	reader := make([]byte, size) // 根据实际需求调整缓冲区大小
	return &readerWrapper{file: file, reader: reader}
}

// 读取器包装器
type readerWrapper struct {
	file   multipart.File
	reader []byte
}

// 实现 io.Reader 的 Read 方法
func (rw *readerWrapper) Read(p []byte) (n int, err error) {
	return rw.file.Read(rw.reader)
}
