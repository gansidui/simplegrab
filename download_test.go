package simplegrab

import (
	"fmt"
	"os"
	"testing"
)

func TestDownload(t *testing.T) {
	os.MkdirAll("./test", os.ModePerm)

	resp, filepath, err := Download("http://www.qq.com", "./test/", ".htm")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.Header)
	fmt.Println(filepath)

	doc, resp, filepath, err := GetDocument("http://www.baidu.com", "./test", ".html")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(doc.Find("title").Text())
	fmt.Println(resp.Header)
	fmt.Println(filepath)
}
