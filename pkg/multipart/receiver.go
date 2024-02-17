package multipart

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	_ "net/http/pprof"
	"net/textproto"
)

type subpart struct {
	*multipart.Part
}

func (p *subpart) Header() textproto.MIMEHeader {
	return p.Part.Header
}

type Part interface {
	io.Reader
	Header() textproto.MIMEHeader
	FileName() string
	FormName() string
}

func MultipartReceiverHandler(handlerFunc func(Part) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headers, err := getheaders(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

			log.Println(err)
			return
		}

		reader := multipart.NewReader(r.Body, headers.Boundary)
		var part *multipart.Part
		for {
			part, err = reader.NextPart()
			if err != nil {
				log.Println(err)
				break
			}

			if err := handlerFunc(&subpart{part}); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				log.Println(err)
				return
			}

			part.Close()
		}
	}
}

// func handlePart(part *multipart.Part) {
// 	wd, _ := os.Getwd()

// 	tempdir := path.Join(wd, part.FileName())
// 	os.Mkdir(tempdir, fs.FileMode(os.O_CREATE))
// 	file, err := os.CreateTemp(tempdir, part.FormName())
// 	if err != nil {
// 		log.Println("os.CreateTemp err", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	if _, err := file.ReadFrom(part); err != nil {
// 		log.Println("file.ReadFrom err", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// }
