package miniox

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"gitlab.galaxy123.cloud/base/public/helper/idx"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.galaxy123.cloud/base/public/constants/consts"
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

// UploadByRequest 上次文件
// suffix 不传就默认去文件后缀
// 最大支持 1024 KB，即 1 M
func (c *Client) UploadByRequest(r *http.Request, prefix string, suffix ...string) (string, error) {

	err := r.ParseMultipartForm(consts.FileMaxSize)
	if err != nil {
		return "", err
	}

	file, fileWriter, err := r.FormFile("file")
	if err != nil {
		return "", err
	}

	defer file.Close()
	name := idx.GenUUID()
	if prefix != "" {
		name = prefix + "/" + name
	}
	// 从文件名中提取后缀名
	fileExt := filepath.Ext(fileWriter.Filename)
	if len(suffix) != 0 {
		fileExt = suffix[0]
	}
	name += fileExt

	err = c.Upload(r.Context(), file, name, fileWriter.Size)
	if err != nil {
		return "", err
	}

	return name, nil
}
