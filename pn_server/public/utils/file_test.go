package utils

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestFileToBase64(t *testing.T) {
	file, err := os.Open("/root/workspace/work2/dc3/data/g2y.jpg")
	if err != nil {
		t.Errorf("open file error: %v", err)
	}
	defer file.Close()
	base64Str, err := FileToBase64(file)
	if err != nil {
		t.Errorf("FileToBase64 error: %v", err)
	}
	t.Logf("[len(base64Str)]%v\n", len(base64Str))

	reader, err := Base64ToReader(base64Str)
	if err != nil {
		t.Errorf("Base64ToReader error: %v", err)
	}

	// write to file
	f, err := os.Create("/root/workspace/work2/dc3/data/g2y2.jpg")
	if err != nil {
		t.Errorf("create file error: %v", err)
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	if err != nil {
		t.Errorf("copy file error: %v", err)
	}
	f.Sync()

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println(base64Str[:100])
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}
