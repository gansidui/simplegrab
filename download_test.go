package simplegrab

import (
	"fmt"
	"os"
	"testing"
)

func TestDownload(t *testing.T) {
	os.MkdirAll("./test", os.ModePerm)

	filepath, err := Download("http://www.qq.com", "./test/", ".htm")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(filepath)

	doc, err := GetDocument("http://www.baidu.com", "./test", ".html")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(doc.Find("title").Text())
}
