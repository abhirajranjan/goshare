package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"sync"

	"github.com/pkg/profile"
)

type file struct {
	filename string
	data     io.Reader
	size     int
	div      int
}

type part struct {
	file   file
	number int
	size   int
}

type Multipart struct {
	wg *sync.WaitGroup
	mu sync.Locker
	multipart.Writer
}

func (m *Multipart) writePart(prt part) {
	m.mu.Lock()
	defer m.mu.Unlock()
	defer m.wg.Done()

	wr, err := m.CreateFormFile(strconv.FormatInt(int64(prt.number), 10), prt.file.filename)
	if err != nil {
		log.Println(err)
		return
	}

	reader := io.LimitReader(prt.file.data, int64(prt.size))
	written, err := io.Copy(wr, reader)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("written %d data of file %s %d\n", written, prt.file.filename, prt.size)

}

func (m *Multipart) makeMultiPartReq(files []file) {
	for _, f := range files {
		for i := 0; i < f.size; i += f.div {
			fmt.Println("started ")
			m.wg.Add(1)
			go m.writePart(part{
				file:   f,
				number: i,
				size:   f.div,
			})
		}
	}

}

func newMultipart(w io.Writer) *Multipart {
	return &Multipart{
		wg:     new(sync.WaitGroup),
		mu:     new(sync.Mutex),
		Writer: *multipart.NewWriter(w),
	}
}

func GenerateData(nums int, size int, div int) []file {
	files := make([]file, nums)
	for i := range nums {
		files = append(files, file{
			filename: fmt.Sprintf("%d", i),
			data:     rand.Reader,
			div:      div,
			size:     size,
		})
	}
	return files
}

func main() {
	defer profile.Start(profile.MemProfile).Stop()

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	reader, writer := io.Pipe()
	mp := newMultipart(writer)
	mp.makeMultiPartReq(GenerateData(2, 200*1024*1024*8, 2*1024*1024*8))

	go func() {
		mp.wg.Wait()
		mp.Close()
		writer.Close()
	}()

	resp, err := http.Post("http://localhost:3000", fmt.Sprintf("multipart/form-data; boundary=%s", mp.Boundary()), reader)
	if err != nil {
		panic(err)
	}

	content, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode, content)
}
