package simplegrab

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 计算string的md5值，以32位字符串形式返回
func StringToMd5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

// 将数据(url)下载到本地(filepath)
// 注意：需要保证url链接的网页(或文件)无需设置相关cookie和header也可下载
func download(urlStr, filepath string) (*http.Response, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("StatusCode not 200")
	}

	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	io.Copy(file, resp.Body)
	return resp, nil
}

// 根据url和后缀(ext)得到在本地存储的filename
func GetFilenameFromUrl(urlStr, ext string) string {
	ext = "." + strings.TrimLeft(ext, ".")
	return StringToMd5(urlStr) + ext
}

// 下载url，保存到本地目录dir下，文件名为GetFilenameFromUrl(urlStr, ext)
// 注意：需要保证url链接的网页(或文件)无需设置相关cookie和header也可下载
func Download(urlStr, dir, ext string, refresh bool) (*http.Response, string, error) {
	dir = strings.TrimRight(dir, "/") + "/"
	filepath := dir + GetFilenameFromUrl(urlStr, ext)

	// 若本地已经存在就不再下载，此时http.Response为nil
	if !refresh {
		if _, err := os.Stat(filepath); err == nil {
			return nil, filepath, nil
		}
	}

	resp, err := download(urlStr, filepath)
	if err != nil {
		return nil, "", err
	}

	return resp, filepath, nil
}

// 下载网页(url)，并返回它的Document结构
// 如果网页已经存在，则不重新下载
func GetDocument(url, dir, ext string) (*goquery.Document, *http.Response, string, error) {
	resp, filepath, err := Download(url, dir, ext, false)
	if err != nil {
		return nil, nil, "", err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, resp, filepath, err
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, resp, filepath, err
	}

	return doc, resp, filepath, nil
}

// 下载网页(url)，并返回它的Document结构
// 如果网页已经存在，仍然重新下载
func GetDocumentNeedRefresh(url, dir, ext string) (*goquery.Document, *http.Response, string, error) {
	resp, filepath, err := Download(url, dir, ext, true)
	if err != nil {
		return nil, nil, "", err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, resp, filepath, err
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, resp, filepath, err
	}

	return doc, resp, filepath, nil
}

// 拷贝文件，返回写入的字节数或错误
func CopyFile(dstName, srcName string) (written int64, err error) {
	file, err := os.Stat(dstName)
	if err == nil && file.IsDir() {
		dstName = strings.TrimRight(dstName, "/") + "/" + srcName
	}

	src, err := os.Open(srcName)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	return io.Copy(dst, src)
}
