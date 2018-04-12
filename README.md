```go
package main

import (
	"os"
	"fmt"
	"tgzlib/tgzlib"
)

func main()  {

	//Write file to tgz
	file, err := os.OpenFile("../test.tgz", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	w, err := tgzlib.NewWriter(file, tgzlib.DefaultCompressLevel)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := w.Write("../tgzlib"); err != nil {
		fmt.Println(err)
		return
	}
	if err := w.Close(); err != nil {
		fmt.Println(err)
		return
	}
	if err := file.Close(); err != nil {
		fmt.Println(err)
		return
	}

	//write file to buffer
	/*

	bw := tgzlib.NewDefaultWriter()
	if err != nil {
		fmt.Println(err)
	}
	if err := bw.Write("../tgzlib"); err != nil {
		fmt.Println(err)
	}
	w.Close()
	fileConent := bw.Body()

	*/

	//Read tgz file to buffer
	f1, err := os.OpenFile("../test.tgz", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	r, err := tgzlib.NewReader(f1)
	if err != nil {
		fmt.Println(err)
		return
	}
	bf, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range bf {
		//filename
		fmt.Println(v.Name)
		//fileContent
		fmt.Println(string(v.Data))
	}
}
```