package miniox

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"github.com/minio/minio-go/v7"
	"github.com/nfnt/resize"
	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/helper/idx"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type Client struct {
	conn        *minio.Client
	Bucket      string        // 桶名
	Location    string        // 地点/服务
	Timeout     time.Duration // 超时
	ContentType string
}
type Options struct {
	Expires    int64 // 过期时间-单位分钟(默认30分钟)
	Width      int   // 指定宽
	Height     int   // 指定高
	IsDownload bool  // 是否下载。图片默认False.返回预览链接地址。文件或者下载图片可以为Ture为下载连接
}

// Upload 上传/更新文件或者图片
func (c *Client) Upload(ctx context.Context, f io.Reader, name string, size int64) error {

	if c.ContentType == "" {
		c.ContentType = consts.FileType
	}
	// 将文件上传到Minio服务器
	_, err := c.conn.PutObject(ctx, c.Bucket, name, f, size, minio.PutObjectOptions{
		ContentType: c.ContentType,
	})

	return errors.WithStack(err)
}

// GetUrl 获取地址
func (c *Client) GetUrl(ctx context.Context, name string, ops Options) string {

	params := url.Values{}
	expires := consts.DefaultExpireTime

	if !ops.IsDownload {
		if ops.Width != 0 && ops.Height != 0 {
			params.Set("resize", fmt.Sprintf("%dx%d", ops.Width, ops.Height))
		}
		params.Set("response-content-disposition", "inline")
		params.Set("response-content-type", "image/jpeg")
	}

	if ops.Expires != 0 {
		expires = time.Duration(ops.Expires) * time.Minute
	}
	res, err := c.conn.PresignedGetObject(ctx, c.Bucket, name, expires, params)
	if err != nil {
		logx.Error(err)
		return ""
	}

	return res.String()
}

// GetUrls 批量获取地址
func (c *Client) GetUrls(ctx context.Context, names []string, ops Options) map[string]string {

	params := url.Values{}
	if !ops.IsDownload {
		if ops.Width != 0 && ops.Height != 0 {
			params.Set("resize", fmt.Sprintf("%dx%d", ops.Width, ops.Height))
		}
		params.Set("response-content-disposition", "inline")
		params.Set("response-content-type", "image/jpeg")
	}
	expires := consts.DefaultExpireTime
	if ops.Expires != 0 {
		expires = time.Duration(ops.Expires) * time.Minute
	}

	var data = make(map[string]string)
	for _, name := range names {
		res, err := c.conn.PresignedGetObject(ctx, c.Bucket, name, expires, params)
		if err != nil {
			logx.Error(err)
			continue
		}
		data[name] = res.String()
	}

	return data
}

// DelObject 删除对象
func (c *Client) DelObject(ctx context.Context, names ...string) error {

	for _, name := range names {
		err := c.conn.RemoveObject(ctx, c.Bucket, name, minio.RemoveObjectOptions{})
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// Deprecated: 建议使用 UploadWebp
// UploadByRequest 上传图片
// 最大支持 2048 KB，即 2 M
func (c *Client) UploadByRequest(ctx context.Context, r *http.Request, prefix string) (string, error) {

	err := r.ParseMultipartForm(consts.FileMaxSize)
	if err != nil {
		return "", err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}

	defer file.Close()
	name := idx.GenUUID()
	if prefix != "" {
		name = prefix + "/" + name
	}
	name += ".webp"

	// 读取响应体内容
	imageData, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// 获取图片格式
	mimeType := http.DetectContentType(imageData)
	var img image.Image
	switch mimeType {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(bytes.NewReader(imageData))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(imageData))
	case "image/gif":
		img, err = gif.Decode(bytes.NewReader(imageData))
	default:
		return "", fmt.Errorf("unsupported image format: %s", mimeType)
	}

	if err != nil {
		return "", err
	}

	// 压缩图片
	resizedImg := resize.Resize(0, 0, img, resize.Lanczos3)

	// 创建缓冲区
	var buf bytes.Buffer
	// 将压缩后的图片转换为 WebP 格式并保存到缓冲区
	err = webp.Encode(&buf, resizedImg, &webp.Options{Quality: 90})
	if err != nil {
		var buf bytes.Buffer
		// 编码图像为JPEG格式，并将其写入到字节缓冲区中
		err = jpeg.Encode(&buf, img, nil)
		if err != nil {
			return "", err
		}

		err = c.Upload(ctx, bytes.NewReader(buf.Bytes()), name, int64(buf.Len()))
		if err != nil {
			return "", errors.WithStack(err)
		}
	} else {
		err = c.Upload(ctx, bytes.NewReader(buf.Bytes()), name, int64(buf.Len()))
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	return name, nil
}

// UploadWebp 上次文件
func (c *Client) UploadWebp(ctx context.Context, r *http.Request, prefix string) (string, error) {
	return c.UploadByRequest(ctx, r, prefix)
}

// UploadFile 上次文件
func (c *Client) UploadFile(ctx context.Context, r *http.Request, prefix string) (string, error) {

	err := r.ParseMultipartForm(consts.FileMaxSize)
	if err != nil {
		return "", err
	}

	file, head, err := r.FormFile("file")
	if err != nil {
		return "", err
	}

	defer file.Close()
	name := idx.GenUUID()
	tem := strings.Split(head.Filename, ".")
	if len(tem) > 1 {
		name += "." + tem[len(tem)-1]
	}

	if prefix != "" {
		name = prefix + "/" + name
	}

	c.ContentType = consts.FileType
	err = c.Upload(ctx, file, name, head.Size)
	if err != nil {
		return "", errors.WithStack(err)
	}

	c.ContentType = consts.WebpType
	return name, nil
}
