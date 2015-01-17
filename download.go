package simplegrab

import (
	"encoding/base32"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 将数据(url)下载到本地(filepath)
// 注意：需要保证url链接的网页(或文件)无需设置相关cookie和header也可下载
func download(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	io.Copy(file, resp.Body)
	return nil
}

// 根据url和后缀(ext)得到在本地存储的filename
func GetFilenameFromUrl(url, ext string) string {
	ext = "." + strings.TrimLeft(ext, ".")
	return base32.StdEncoding.EncodeToString([]byte(url)) + ext
}

// 根据filename得到相应的url
func GetUrlFromFilename(filename string) string {
	// 去掉文件后缀
	filename = filename[:len(filename)-len(path.Ext(filename))]
	data, _ := base32.StdEncoding.DecodeString(filename)
	return string(data)
}

// 下载url，保存到本地目录dir下，文件名为GetFilenameFromUrl(url, ext)，返回filepath
// 注意：需要保证url链接的网页(或文件)无需设置相关cookie和header也可下载
func Download(url, dir, ext string) (string, error) {
	dir = strings.TrimRight(dir, "/") + "/"
	filepath := dir + GetFilenameFromUrl(url, ext)
	// 若本地已经存在就不再下载
	if _, err := os.Stat(filepath); err != nil {
		if err = download(url, filepath); err != nil {
			return "", err
		}
	}
	return filepath, nil
}

// 下载网页(url)，并返回它的Document结构
func GetDocument(url, dir, ext string) (*goquery.Document, error) {
	filepath, err := Download(url, dir, ext)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
