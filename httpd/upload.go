package httpd

import (
	"crypto/md5"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

// IUpload interface
type IUpload interface {
	Auth(ctx *fasthttp.RequestCtx) error
	// Do(ctx *fasthttp.RequestCtx, token, ext, folder string, now time.Time) error
	Do(ctx *fasthttp.RequestCtx, fileHeader *multipart.FileHeader) error
	// StoreDir(now time.Time) string
}

func doUploadErr(ctx *fasthttp.RequestCtx, err error) bool {
	if err == nil {
		return false
	}
	ctx.WriteString(err.Error())
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	return true
}

// RegistUpload router
func (s *SiteServer) RegistUpload(rule string, up IUpload) {
	s.Router.POST(rule, func(ctx *fasthttp.RequestCtx) {
		name := ctx.FormValue("name")
		// filename := ctx.FormValue("filename") // filename to store
		size := ctx.FormValue("size")
		last := ctx.FormValue("lastModified")
		typ := ctx.FormValue("type")

		req := &ctx.Request
		apiToken := req.Header.Peek(APIFieldToken)
		pToken := req.Header.Peek(APIFieldPageToken)

		u32 := buildUploadToken(req, name, size, last, typ, pToken)

		// auth api token
		if u32 != util.Bytes2Str(apiToken) {
			ctx.SetStatusCode(fasthttp.StatusForbidden)
			return
		}

		err := up.Auth(ctx)
		if doUploadErr(ctx, err) {
			return
		}

		fileHeader, err := ctx.FormFile("file")
		if doUploadErr(ctx, err) {
			return
		}

		err = up.Do(ctx, fileHeader)
		if doUploadErr(ctx, err) {
			return
		}

		//复制文件内容
		// _, err = io.Copy(nf, file)
		// if err != nil {
		// 	ctx.WriteString(err.Error())
		// 	ctx.SetStatusCode(fasthttp.StatusBadRequest)
		// 	return
		// }
	})
}

var filepool bytebufferpool.Pool

// DoUpload upload
func DoUpload(name, storeDir string, fileHeader *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(name)
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := filepool.Get()
	defer filepool.Put(buf)

	_, err = buf.ReadFrom(file)
	if err != nil {
		return "", err
	}

	token := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))

	_, err = os.Stat(storeDir)
	if os.IsNotExist(err) {
		os.MkdirAll(storeDir, os.ModeDir)
	}

	fileName := token + ext
	newFile := filepath.Join(storeDir, fileName)
	_, err = os.Stat(newFile)
	if os.IsExist(err) {
		return fileName, nil
	}

	nf, err := os.OpenFile(newFile, os.O_CREATE|os.O_RDWR, 0666)
	defer nf.Close() //使用完需要关闭
	if err != nil {
		return fileName, err
	}

	nf.Write(buf.Bytes())
	return fileName, nil
}
