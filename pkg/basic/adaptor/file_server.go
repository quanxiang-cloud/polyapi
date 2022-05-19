package adaptor

import (
	"context"
)

// FileServerOper is the interface for file server proxy
type FileServerOper interface {
	// UploadFile upload a file to file server
	UploadFile(c context.Context, filename string, content []byte) (string, error)
}

// SetFileServerOper set the instance of file server oper
func SetFileServerOper(f FileServerOper) FileServerOper {
	i := getInst()
	old := i.fileServerOper
	i.fileServerOper = f
	return old
}

// GetFileServerOper get the instance of file server oper
func GetFileServerOper() FileServerOper {
	i := getInst()
	return i.fileServerOper
}
