package httpd

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

// IUpload interface
type IUpload interface {
	Auth(ctx *fasthttp.RequestCtx) error
	Success(ctx *fasthttp.RequestCtx, token, ext, folder string, now time.Time) error
	StoreDir(now time.Time) string
}

func uploadFileName(name, size, last, typ []byte) string {
	buf := bytebufferpool.Get()
	buf.Write(name)
	buf.Write(size)
	buf.Write(last)
	buf.Write(typ)
	str := fmt.Sprintf(stri16Formart, md5.Sum(buf.Bytes()))
	bytebufferpool.Put(buf)
	return str
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

		file, err := fileHeader.Open()
		if doUploadErr(ctx, err) {
			return
		}
		defer file.Close()

		src, err := ioutil.ReadAll(file)
		if doUploadErr(ctx, err) {
			return
		}

		h := md5.New()
		h.Write(src)
		token := fmt.Sprintf("%x", h.Sum(nil))

		ext := filepath.Ext(util.Bytes2Str(name))

		now := time.Now()
		folder := up.StoreDir(now)

		_, err = os.Stat(folder)
		if os.IsNotExist(err) {
			os.MkdirAll(folder, os.ModeDir)
		}

		newFile := filepath.Join(folder, token+ext)
		_, err = os.Stat(newFile)
		if os.IsExist(err) {
			err = up.Success(ctx, token, ext, folder, now)
			doUploadErr(ctx, err)
			return
		}

		nf, err := os.OpenFile(newFile, os.O_CREATE|os.O_RDWR, 0666)
		if doUploadErr(ctx, err) {
			return
		}

		//使用完需要关闭
		nf.Write(src)
		nf.Close()

		//复制文件内容
		// _, err = io.Copy(nf, file)
		// if err != nil {
		// 	ctx.WriteString(err.Error())
		// 	ctx.SetStatusCode(fasthttp.StatusBadRequest)
		// 	return
		// }

		err = up.Success(ctx, token, ext, folder, now)
		if doUploadErr(ctx, err) {
			return
		}

	})
}
