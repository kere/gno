package httpd

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

// IUpload interface
type IUpload interface {
	Auth(ctx *fasthttp.RequestCtx) error
	Success(name, folder string) error
	StoreDir() string
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
		if err != nil {
			ctx.WriteString(err.Error())
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			ctx.WriteString(err.Error())
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			ctx.WriteString(err.Error())
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		defer file.Close()

		src, err := ioutil.ReadAll(file)
		if err != nil {
			ctx.WriteString(err.Error())
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		h := md5.New()
		h.Write(src)

		ext := filepath.Ext(util.Bytes2Str(name))
		storeName := fmt.Sprintf("%x%s", h.Sum(nil), ext)
		folder := up.StoreDir()
		_, err = os.Stat(folder)
		if os.IsNotExist(err) {
			os.Mkdir(folder, os.ModeDir)
		}

		newFile := filepath.Join(folder, storeName)
		_, err = os.Stat(newFile)
		if os.IsExist(err) {
			ctx.WriteString(newFile)
			return
		}

		nf, err := os.OpenFile(newFile, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			ctx.WriteString(err.Error())
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
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

		ctx.WriteString(newFile)
		up.Success(storeName, folder)

	})
}
