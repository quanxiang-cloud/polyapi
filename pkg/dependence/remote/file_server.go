package apisecret

import (
	"bytes"
	"context"
	"time"

	"github.com/quanxiang-cloud/fileserver/pkg/guide"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
)

func initFileServer() error {
	g, err := guide.NewGuide(
		guide.WithTimeout(time.Second*5),
		guide.WithMaxIdleConns(3),
	)
	if err != nil {
		return err
	}
	adaptor.SetFileServerOper(&fileserver{g: g})
	return nil
}

type fileserver struct {
	g *guide.Guide
}

func (s *fileserver) UploadFile(c context.Context, fileName string, content []byte) (string, error) {
	return fileName, s.g.UploadFile(c, fileName, bytes.NewReader(content), int64(len(content)))
}
