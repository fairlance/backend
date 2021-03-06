// Code generated by go-bindata.
// sources:
// templates/index.html
// templates/payment.html
// templates/websockettest.html
// DO NOT EDIT!

package importer

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x3c\xeb\x72\xdb\x36\xb3\xff\xf3\x14\x5b\x7e\x3e\x96\x54\x5b\xa4\x2f\x49\xda\x28\x92\x4e\xe3\xca\x49\xd3\x5b\x32\x89\xdb\x33\x9d\x4c\x26\x03\x91\x90\x84\x98\x04\x38\x00\x28\x59\x71\xf5\xee\x67\x00\x90\x12\xef\xa2\x64\x27\x93\x76\x3e\xfd\x09\xb9\xd8\x5d\xec\x0d\x8b\x05\xb8\x4e\xff\x9b\xd1\xab\x1f\xaf\xfe\x7a\x7d\x09\x33\x19\xf8\xc3\x07\x7d\xf3\xcf\x83\xfe\x0c\x23\x6f\xf8\x00\x00\xa0\x2f\x89\xf4\xf1\xf0\x39\x22\xdc\x47\xd4\xc5\xf0\x32\x08\x19\x97\x98\xf7\x1d\x33\x62\xb0\x7c\x42\xaf\x61\xc6\xf1\x64\x60\xcd\xa4\x0c\x45\xcf\x71\x02\x74\xe3\x7a\xd4\x1e\x33\x26\x85\xe4\x28\x54\x2f\x2e\x0b\x9c\x35\xc0\x39\xb7\xcf\xed\xef\x1c\x57\x88\x0d\xcc\x0e\x08\xb5\x5d\x21\x2c\xe0\xd8\x1f\x58\x42\x2e\x7d\x2c\x66\x18\x4b\x0b\x08\x95\x78\xca\x89\x5c\x0e\x2c\x31\x43\xe7\xdf\x3f\xec\x5e\xfc\xf9\x17\x21\x6f\x5f\x3e\xc7\xbf\x9c\x7a\x2f\x82\x9f\xdf\x3c\xbb\x5e\xba\xd1\x4f\xcf\x7e\x7a\x33\x3d\x3f\x7b\x15\xfc\xe1\x2e\x16\xdf\x31\x7a\xfe\xe6\x2f\x6f\xfa\xf0\x4f\x74\xf4\x3a\x78\x7b\x25\x3e\x39\xbf\x3c\xfe\x7e\x3e\xf6\x2e\x3f\xce\x1e\x46\x96\x96\x5d\xfd\x5c\xce\x84\x60\x9c\x4c\x09\x1d\x58\x88\x32\xba\x0c\x58\x24\xac\x58\x39\x2d\xc5\x70\x8d\xac\x7e\xca\x50\x70\x9b\x01\xa9\x1f\x9b\x63\x3e\xf1\xd9\xa2\xbb\xec\x81\x70\x39\xf3\xfd\xa7\x19\x9c\x55\xe6\xcd\xf9\x76\xcc\xbc\x65\x09\x9b\x10\x79\x1e\xa1\xd3\xae\x64\x61\x0f\xce\x4e\xc2\x9b\xa7\x79\x8c\x2c\xd7\x6f\x9d\xcc\xbb\x8d\x7c\xcc\x65\x77\x42\x6e\xb0\x57\xc6\x9d\x09\x22\x09\xa3\x3d\x8d\xf0\x14\x0a\x08\x63\x26\x25\x0b\x7a\xf0\x48\xcd\x5c\x18\xe5\x64\x3a\x93\x3d\x38\x2f\x1d\x5c\x10\x4f\xce\xd4\xe0\xff\x14\xc7\x3e\x75\x09\xf5\xf0\x4d\xef\xc9\x93\x27\x4f\x4a\x67\xe5\x1e\xe6\x4a\xe7\xae\x8f\x27\xb2\xcb\x91\x47\x22\xd1\x2b\x9d\x26\x85\xab\xc5\xa9\x44\xde\x18\xbc\xef\xc4\x6e\xec\x3b\x26\xba\x1f\xf4\x95\xf5\x63\x17\x7b\x64\x0e\xc4\x1b\x58\x28\x0c\xad\x8d\xa7\x35\x78\xde\x25\x93\x81\x15\x88\x29\x7c\x33\x18\x40\xab\x65\xc1\xbc\x3b\x26\xd4\xeb\xb9\x3e\x12\x62\x60\xdd\x42\xcb\xd8\x5b\x44\xae\x8b\x85\x68\xf5\x40\x21\x2b\x5c\x76\xdd\x3a\x4e\x46\x17\x88\x53\x42\xa7\xa9\xd1\x37\x11\x55\x10\xdb\xb6\x37\x58\x1e\xa2\x53\xcc\x5b\x3d\x68\x9b\x09\x35\x0f\x38\x3c\x84\xe4\x35\x45\xd4\x81\x95\x95\x51\x36\x16\x48\x73\x82\x54\x0c\x58\xc3\x83\xdb\x40\x4c\x57\x7d\xc7\x23\xf3\x94\x76\x14\xcd\x13\x1a\x8a\xe6\x63\xc4\xc1\xfc\xd3\xf5\xf0\x04\x45\xbe\x4c\x5e\x85\x44\x92\xb8\xca\xda\x56\x76\x15\x68\xfb\xc4\x1c\x5c\x46\x25\x22\x14\xf3\x1c\x4e\x1e\x2f\xe6\xa9\x7c\x50\x8a\xab\xf1\x51\x0e\x7b\xcc\x11\xf5\xac\x38\xc1\x38\x28\x24\x0e\x89\xd3\x50\x05\x87\xe4\x57\xcc\x5b\x95\xe8\x7d\x07\x95\x48\x9e\x35\xd9\x1a\x1c\xf9\x29\x09\x13\x3b\x51\x34\xaf\x52\xc8\x27\x25\x51\xe3\x4a\x32\xc7\xad\x1e\x48\x34\xd6\x01\xe1\x8d\x5b\xb0\xaa\x51\xa8\x8f\x62\x0b\xfc\x47\xc5\x20\xa3\x3d\xd7\x27\xee\xf5\xc0\xd2\xf4\x9a\xdc\x1a\x8e\x2e\x4a\xf5\x30\xba\xf8\xe4\x2e\xe2\x09\x8c\xb8\x3b\xbb\x93\x88\x31\x0b\x6b\xf8\x56\x3f\xec\x27\x6a\x83\xd9\x33\x21\xe2\x84\x68\x19\x60\x2a\xad\xe1\x6b\xf3\xf0\x85\xa6\x5d\xe0\xb1\x60\xee\x35\x96\x12\x0b\x69\x0d\xff\x2f\x79\x85\x2b\x2c\x76\x95\xa1\xef\x44\x7e\x6e\xed\xe5\x16\xb3\x43\xd1\x3c\x97\xb9\xb6\xac\xcc\x34\x0e\x0f\x6f\xd8\x62\xcb\xca\x75\x99\xdf\x0d\xbc\xee\xb9\x15\xa7\xc4\x54\xd8\x56\x85\x7d\x8a\x7a\x2c\x69\x77\xca\x59\x14\x76\xe7\x98\x4b\xe2\x22\x1f\x14\x68\xec\x33\xf7\xba\x2e\xa0\xc6\x91\x94\x8c\xa6\xb8\x68\xb2\x24\x43\xa9\x67\x11\x64\x63\x2d\x0a\x3d\x24\x71\x5b\x07\x30\xa3\x03\x8e\x3f\x4c\x31\xc5\x1c\x49\xfc\x41\x79\xe2\x83\x87\x24\x6a\x75\x2c\x90\xcb\x10\x0f\x2c\xc3\xdf\x1a\xbe\x88\x71\x40\xe1\x80\xc2\xe9\x3b\x66\xec\x73\x0a\x67\x82\xe5\x03\xf2\xfd\xa2\x44\x26\x63\xc1\x33\xdf\xff\x12\x92\x78\xd8\xc7\x12\x2b\x49\x3e\x4c\x38\x0b\x3e\x98\xa5\xfa\x01\xd3\x29\xa1\xb8\x28\xdc\x48\xa3\x2b\xe1\xe0\x39\x67\x01\x98\x05\x0d\x97\x1a\xfd\x9e\xe4\x35\x20\xbd\x23\xee\x24\xba\x37\xae\x94\x17\xf9\x3e\x28\x1c\x50\x79\xb2\x4e\xc8\xaa\xc4\x5f\x01\x2e\x2e\x93\x27\x77\x5c\x26\x10\x46\xbe\x6f\xaa\x9b\x06\xcb\x23\xa3\x6c\x31\x97\x2b\xae\x9b\x02\x44\x21\x6b\x91\x3e\xb2\xb1\x68\xad\xac\x0a\x2f\xe4\x32\xb8\x26\x8a\x69\xac\xe1\xcf\x6c\x2c\x9a\xbb\x79\x3f\xf1\x26\x1c\x63\xbd\x79\xf3\x9d\xa5\x4c\x93\x5a\xc3\xe7\x9b\xb7\xbd\xbc\xae\x87\x66\xa7\xc3\x83\x5b\xc5\x7f\xd5\x77\x66\xa7\xdb\x3d\xe9\xfa\x18\xf1\x09\xb9\xa9\xf3\x9e\x08\x11\x1d\x5e\x31\x89\x7c\x78\x49\x61\x74\xd1\x33\xa0\xb5\xaa\xc8\x9b\x62\x55\xbd\x49\x85\xf2\x92\x8e\x2e\x56\x7d\x47\xd3\xc4\xff\x34\xe1\xfc\x92\x66\x17\xe7\x96\x39\x0c\xae\x41\xdd\x69\xb6\x84\x61\xa3\xb8\x5d\x53\x4a\x1c\x84\xbe\x4a\xbb\x66\xb1\x84\x1c\xcf\x5f\xa3\x29\xfe\x15\x8d\xb1\xdf\xee\x24\x25\xf7\x76\x46\xd0\x2c\xa9\xdc\x88\xf2\x4c\xc2\x26\x13\x81\xe5\xa0\x75\x04\x6d\xf3\x08\x5d\xf0\x49\x40\x64\x07\x8e\xa0\x75\xa8\x1f\x07\x2d\x38\x8a\x81\xca\x5e\x39\x51\x57\xdb\x57\xc3\x5a\x4e\x27\xd1\x7a\x3b\xee\xc1\xad\x1b\x71\x8e\xa9\x4c\x4f\xb5\xb3\x59\x29\xbe\x91\x5f\x87\x59\x8f\xb6\x99\x35\x27\xea\x67\x30\x6b\x5d\x34\xd7\x25\x80\xd4\xda\x96\x68\xec\xe3\x2e\xc7\x22\x64\x54\x90\x39\xae\x5b\xe3\x1a\x37\x43\x08\x86\x5c\x48\x4e\x42\xec\xc5\x6f\x2e\xa3\x1e\xa6\x42\x9d\xd6\xea\xa5\x97\x9b\xcb\x99\x7a\x3c\xde\xd0\xbd\x72\x36\x24\x5e\xdf\x91\xb3\xe6\xf8\x14\x05\x78\x37\x0a\x42\x27\x6c\x37\x0a\x13\x43\xde\x33\xb9\x1b\x99\xa9\x09\x9a\xd1\xf4\x9d\x6d\x46\x52\x7c\xb6\x9a\xbb\x2f\x37\x37\x09\xf5\x78\x9b\x55\x39\x61\x7c\x60\x61\x2a\x89\x5c\x02\xa1\xa0\x9f\x08\x16\x4d\x97\x64\x53\xe7\x1a\x64\x6f\x78\x70\x6b\xa6\x58\xda\xc4\x83\x55\xdf\x91\x0d\x22\xa8\x9c\x5e\xb9\x1e\x56\x9b\xf7\x09\xe1\x42\xfe\xae\x81\xb0\x81\xfa\x28\x01\xde\x65\x2e\x11\x05\x01\xe2\xcb\xf4\x74\x38\x40\xc4\xdf\x83\x6b\x52\x74\xe8\x4b\xcb\xc4\xf2\xf6\x3a\xc6\x54\xe6\x01\x49\x02\xfc\x96\x50\x17\xb7\xf3\xc3\x1d\x58\x01\x9a\xb2\xdd\x75\x69\x8c\x0c\x95\xd5\xd2\xd6\x92\x5e\x25\xf8\xb6\xc7\xdc\x48\x9d\x6d\x75\x6a\xa7\x91\xef\xc3\xe1\x21\x78\xcc\x7d\x39\x82\x6f\x06\x1b\xdf\x77\xe0\xef\xbf\x61\x8d\x1b\xa3\x5a\x3b\x49\xa9\x7e\x35\xc5\xf8\x14\xcb\x43\x3d\xaf\xce\xeb\x9b\xa0\x53\xf9\x7e\xb3\x1d\xac\x77\x83\xb2\x4d\x60\x67\x71\x32\xf5\xc7\xd4\x5f\x86\x33\xe2\x32\x0a\xeb\xa7\xae\x3b\xc3\x73\xce\x68\xd7\x63\x0b\x6a\x6d\x2d\x67\x4a\xa7\x68\xbc\x0d\x65\xa8\xf6\xf7\x68\x8d\x43\x07\x29\x87\x66\xf7\xdd\x8d\x67\x8d\x63\x3f\xa3\x29\xa3\xf0\x1f\x60\xc8\x2d\xe7\xee\xb2\x40\xfd\x2c\xf1\x17\x85\x3e\x43\xde\xda\x60\x60\xa6\xff\x07\xda\x8d\xe3\x80\xcd\xf1\x17\xb3\x1b\xe6\x48\xe0\x8d\xd9\xcc\xec\x9f\xd9\x6c\xcd\x93\xfc\xf6\xf2\x61\x8d\x99\xab\xc7\x9b\x2e\xee\x5d\x36\x9b\x1d\x2a\x02\x88\xf7\x44\x97\xf9\xca\xb0\x03\xeb\xd1\x3e\xce\xcb\x14\xb6\x51\x77\xa2\xce\x7d\xfa\xfb\xcf\x1e\xcc\x60\x53\xd9\x36\x2a\xb9\xaa\x79\x34\x2b\xc5\xaa\xe9\x79\x52\x9c\xb5\xe7\xc8\x3f\x86\x6b\xbc\xec\xa8\x02\x2d\xf1\xd8\x9e\xba\x6d\xf8\x7b\x30\x47\x3e\x99\xd2\x81\xa5\xbf\xa7\xf4\xc7\xc3\x83\xdb\x6b\xbc\x54\xe7\x9c\xe1\x6e\xf5\x45\x23\xfe\x07\xb7\x73\xe4\xef\x58\x2e\x15\x98\x36\x8e\xf3\x12\xca\xfd\xfc\xd1\x77\x74\x74\xed\xba\xb5\xec\x54\x9e\x35\x5f\xbd\x3b\x9c\xd7\x9b\x1f\x42\xb7\x58\xa6\xd6\x02\x77\xbb\x9d\x14\x41\xf7\xf4\x2c\x7f\x3d\xb9\xfe\x34\xb3\xf5\xf0\xcb\x4b\xbf\x14\xd4\x4f\xb6\xc5\x18\x29\x9a\x05\x6e\x54\xb7\x34\x17\x28\x43\x35\x61\x3c\x48\xc8\xd4\x73\x77\xc6\x38\xf9\xc4\xa8\x44\xbb\x14\x4b\x25\xd7\xbd\x8f\x77\x4c\x0c\x69\x16\x5a\x10\x7d\xe9\xbb\x4f\x1a\xde\xc3\xde\x4d\x78\xad\x2f\xa2\xef\x90\xf1\xf6\xba\xfe\x35\x81\x68\x07\x58\xce\x98\xa7\x43\x73\x8a\x65\xd3\xfb\xdf\x1c\xb1\xa1\xb5\x86\x2f\x2e\xaf\xf6\xab\x9f\x3e\x83\x36\x21\x13\xfb\xab\xa3\x89\xad\xe1\xeb\x57\x6f\xbf\x22\x85\xa2\x3b\xe8\x13\x69\x75\xfe\xf8\x7a\xb4\x31\x1f\x93\xf6\x56\x28\x26\xb7\x86\xa3\xcb\x5f\x2f\xaf\x2e\xef\xa6\x56\xcd\x6d\xe4\x3d\x93\xed\x43\xf2\xb5\x67\x30\x42\xc3\x48\xde\x3d\x87\xa5\xcf\x28\x29\x96\x5d\xe4\x79\x8c\x5a\x43\xfd\xcd\x7f\x9f\x83\x70\x66\x12\xcd\x38\x0e\x5d\x89\x6f\xa4\x95\xb1\xac\xcb\xa8\xe4\xcc\x8f\x43\xef\x1a\x2f\xa3\xd0\xc6\x54\x62\xae\x4e\x12\xe6\x83\x8d\x1a\x0b\x98\xa7\x9b\xe6\x4c\x4c\x46\xdc\xb7\x20\xf4\x91\x8b\x67\xcc\xf7\x14\xae\x86\x38\xff\xb2\x50\xdc\x15\xbd\x6c\xeb\x8e\x4b\xa1\x8a\x5c\xfd\xaf\xd8\xd9\x55\x50\x21\x8e\x51\x79\x5c\x71\xb6\x10\x03\xeb\xbb\x62\x10\xa9\x1a\x35\x17\x45\x1a\xa4\x8e\x2a\x31\xc7\x7f\x61\x50\xec\x63\xe9\xaf\xd3\xed\x3b\x5c\xff\xe4\x6f\x12\xe3\xbc\x32\xbc\xbc\xc1\x6e\x24\x1b\xf4\x74\x54\x0a\xf1\x35\xb9\xbb\xef\x28\xe7\x34\x39\xc7\x6d\x65\xd9\x04\x25\x7f\xb0\xc9\x66\x1a\x8e\x16\x23\x24\x91\x6e\xec\x6c\x92\x66\xfa\x21\xc7\xc3\x83\xdb\x2c\xf5\xaa\xef\x28\xf0\x5d\x44\xad\xfb\xea\xd9\xf8\xa4\x59\x68\x48\x5b\xbf\xc6\x8f\x7d\xc7\x9c\x78\xfb\xc2\xe5\x24\x94\x20\xb8\xbb\xe9\x16\x8f\x68\x78\x3d\xd5\xfd\xe1\xf3\x08\xff\x70\x66\x9f\xd8\xe7\x8e\x47\x84\x54\xaf\xba\x29\xfc\xa3\xd0\xf7\x80\x9a\x74\x2b\x0f\x74\x43\x98\xf8\xe1\xc4\x3e\x3d\xb3\x4f\x0c\x1b\x0d\x29\x61\x94\x70\x32\x92\xce\x11\x07\x14\x86\x30\x00\x8a\x17\xf0\x67\x84\xdb\x9b\xf6\x69\x0f\xeb\xef\x13\x98\x8b\x1e\xbc\x6b\x1d\xdc\xb6\x8e\xa1\xb5\x6a\xbd\x3f\x5e\x23\x60\xbf\x07\xad\xff\xa0\x30\x6c\x6d\x60\x1e\x92\xa8\x97\xeb\xc1\x96\x68\xdc\xd3\xed\x40\xc7\x59\xf0\x32\xc4\xbd\xb8\xbf\x26\x3b\x12\x88\x69\x0f\x5a\x39\x60\xf2\xc5\xb0\x07\xef\xde\xe7\x18\x15\x1b\x39\x7a\x70\x52\x8a\x33\xba\x28\x8c\xe8\x5b\xc8\x32\xa8\xbe\xfe\xea\xc1\xed\x2a\x3b\x62\xbe\xe4\x14\x08\xb4\xad\x0a\x50\x13\xb9\x79\x83\x80\xee\xe8\xf6\x96\x45\x25\xb5\xf6\x7a\x3b\xee\x99\xa3\x5c\x71\x38\xe2\x7e\x39\x5d\xbc\x3c\xd4\x60\x66\x6c\xd3\x3d\x91\xd2\x64\x81\x64\x51\x2e\xe3\x91\x49\x44\xf5\xdd\x37\xb4\x3b\x25\x72\x93\x09\xb4\xe5\x8c\x08\x3b\xbe\x47\xb1\xbc\xb1\x55\x86\xa7\xf9\x29\xbc\xf8\x42\xbd\x53\xc0\xc8\x76\x75\xe4\xac\xac\x63\xa6\x5e\x12\xcd\x3d\x89\x0a\x18\xc0\xbb\xf7\x5f\x4e\xd8\x32\x93\x1a\xbf\x89\xbc\x51\x0d\xcb\xb4\x32\x2f\x2e\xaf\x5e\x23\x8e\x02\x51\xa9\x95\xee\x98\x07\x6b\xd3\xfb\x6e\x15\xff\x24\x40\x2d\xdd\x50\xb3\x51\xc7\x30\xbd\xe1\xb5\xe0\xc8\xd0\xeb\xae\xb0\x23\x68\x1d\x4a\x34\x4e\x41\xd1\xb8\xc8\x46\x99\x68\x2d\x90\x4a\xcc\x11\xf5\xf0\x84\x50\xec\x55\x19\x2a\x43\xf1\xee\xe4\xbd\xce\xe6\x87\xad\x0e\x6c\xd8\x0c\xc0\x3a\xb4\xe0\x68\x03\x29\xce\x0b\xfa\xcf\x3e\x62\xec\xf8\xa1\x96\xa0\xd8\x04\xb4\x49\x5e\x4a\xbd\x22\x85\x49\x81\x53\x2c\xdb\x3e\x73\x91\xb2\xbd\x6d\xfe\xe4\x45\x99\x26\xdb\xba\xfc\x51\x30\xfa\xbf\xca\x50\x46\x92\xa2\xff\xd5\xcf\x96\x33\x4c\xdb\x1b\x3f\xc6\xed\x30\xb8\xca\x50\x5a\x86\x30\x8c\xbd\x99\x60\xdb\x2a\x45\xda\xbf\x61\x21\xd0\x14\x97\xdb\x25\x21\x4c\x05\x77\x96\xfa\x32\x1e\xa8\x27\x37\x3d\xe8\x59\xca\xab\xb2\x18\xc8\x10\x99\x86\xc2\x1c\xd5\x32\xdc\x22\x6a\xfc\x7d\x3b\x4f\xf8\x4a\x83\xeb\x49\x75\xe6\x2c\x50\xfe\xaa\xa0\xf5\x84\xf1\xf7\xa3\x1c\xe1\x48\x41\xab\x09\x55\xf4\xae\x8d\x63\xba\x54\xeb\xfc\xb7\x36\x4b\x71\x8f\x29\xce\x7c\x61\x5f\x15\xd1\xaa\x45\xc9\xf1\x1e\x5d\xd4\x70\x1c\x5d\x6c\xe7\x93\xfa\x18\x5e\x60\x33\x8a\xc7\xaa\xb9\x54\x37\xd9\x29\xde\xba\xb9\xf3\x37\x31\x6d\x77\xca\x39\xac\x2a\x96\x8c\xab\xf6\x99\xd4\x9a\xc1\x9c\x33\x5e\x67\x70\x97\x51\xc1\x7c\x6c\xfb\x6c\x1a\x23\x37\x58\x5b\xd6\xa5\xc2\x2c\x49\x92\x77\x50\x21\xb7\x21\x65\xba\x1f\x1b\x6f\x92\xf1\xb2\x18\x9a\x14\x65\xc7\xad\x80\xe5\xda\x73\x2c\x23\x4e\xb3\x84\xdd\x14\x1d\x1c\xc1\x69\x07\x8e\xc0\xea\x5a\x49\x46\x37\x58\x25\xa9\xb2\x58\x1c\x18\xe6\x96\x55\xa7\x63\xa6\x15\x71\x57\x1d\x8f\xd2\xa2\xf6\xe3\x1d\x27\x09\xde\x5d\x74\x3e\xaa\xd4\x39\x87\xd6\x4e\xe1\x7d\x0b\x67\x9d\xb2\x5d\x7b\x3f\x43\xe4\x5b\x50\xb7\xd9\xa2\x5c\x8d\x6a\xc9\x53\xc1\x50\x27\x46\x72\x26\x6d\x54\x09\xed\x5c\x33\x54\x94\x41\xba\x06\x5c\x57\x11\xc5\x32\x13\x92\xea\x2c\xa9\x29\xca\x51\x4c\xff\x44\x6f\xfd\x9d\xad\x1c\x4b\x17\xb3\x9a\xd1\xe6\x12\xaf\x1c\x33\xa9\x8a\xd3\xc8\x06\x56\x8e\x6f\xea\xeb\x34\xb6\x82\x14\x71\x8b\x69\xcf\x94\x0d\x21\x13\xe5\x75\x83\x2a\x15\x5a\xc7\x5f\x57\xa1\x90\x3b\x5e\x0f\xe0\xe7\xb7\xaf\x7e\xb7\x85\xe4\x84\x4e\xc9\x64\xd9\xce\xb2\x7c\x63\xd0\x8e\x75\x23\xc6\x31\x3c\xac\xc8\x87\xf0\xdf\xd4\x1f\xcb\x15\x53\x6f\x96\x61\xe9\x2a\xdc\x56\x8e\x2e\x08\xf5\xd8\xc2\x16\x58\x5e\x91\x00\xb3\x48\xb6\xeb\x57\x75\x56\xcf\x56\xab\xa4\x24\x3e\x86\xd3\x93\x93\x93\x9c\x52\xab\xfc\x11\x94\x13\xe4\x93\x4f\x99\x23\x08\x1b\x7f\xac\x52\x41\x48\xae\x0f\x53\xc5\xe9\x26\x8c\x43\x5b\xa7\x10\x20\x14\x14\x8b\xca\xf3\x01\x1b\x7f\xb4\x67\x48\xbc\x5a\xd0\xd7\x9c\x85\x98\xcb\x65\x3b\xec\xc0\xe1\xa1\xa2\x7a\x17\xbe\x6f\x74\xce\xd0\xc2\x4b\x6e\x87\x91\x98\xb5\x31\x75\x99\x87\xff\x78\xf3\xf2\x47\x16\x84\x8c\x62\x2a\x15\xc3\x23\xb0\x06\x96\x6e\xca\x2a\x8c\x9a\x89\x3a\x55\x0e\xaf\x4a\xe1\x6a\xc2\x8f\x8c\xd0\xb6\x75\x68\xd5\xdb\x75\xdd\x3c\x9b\xb6\xab\x3a\xea\x95\xe9\xa3\xe0\xf1\x25\xcb\x48\x1d\x30\x35\x5e\xb9\xf9\xb1\xcb\xa8\xa7\xf2\xf3\x6f\x48\xce\xec\x89\xcf\x18\x6f\xb7\xd7\x84\x1d\xe8\x82\x99\xc4\x49\x7c\x5f\xca\x86\x50\x89\xf9\x1c\xf9\x59\x3e\x09\x73\x07\xce\x4f\x1f\x9d\x3f\xae\xa0\x57\x0e\x5c\xd3\x0f\xd5\x36\x56\xbb\x79\xaf\x51\x8f\xc0\x82\x25\x46\x5c\x94\xac\xc6\xa2\xc1\xb7\x49\x78\xf6\xe8\xc9\x59\x31\xb8\xef\x2a\x5f\xc0\xa8\x9c\xdd\x8f\x80\xdf\x3f\x7e\x78\xef\xe2\x79\x68\x79\x3f\xc2\x29\xe7\xde\xb3\x6c\x33\x16\xdd\x93\x6b\x1f\xdf\xbb\x57\x09\x8d\x24\x6e\x26\x5c\x4c\x5f\x14\x4c\xa7\x93\x64\xfd\x59\x55\xff\xad\x82\x79\x8a\x37\x0c\x95\x9e\x93\x2b\xa3\xa7\x0f\xd2\x97\xac\x8e\xfe\x9f\x45\xfe\x3f\x00\x00\xff\xff\x46\x62\x37\xe6\x70\x44\x00\x00")

func templatesIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesIndexHtml,
		"templates/index.html",
	)
}

func templatesIndexHtml() (*asset, error) {
	bytes, err := templatesIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/index.html", size: 17520, mode: os.FileMode(420), modTime: time.Unix(1499600894, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesPaymentHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x57\xdf\x6f\xdb\x36\x10\x7e\xf7\x5f\xc1\x11\x79\xdb\x64\xae\x75\xb0\x16\x9e\x24\xa0\x5b\xd3\x24\x2b\xd0\x18\x71\x5a\xa0\x4f\x01\x25\x9e\x2d\x36\xe2\x8f\x92\x27\x27\x8e\x91\xff\x7d\xa0\x64\x3b\x71\x62\xd9\xf2\x16\x3d\x04\x20\xef\xbb\xe3\xf1\xbe\xe3\x77\x4e\x5c\xa0\x2a\xd3\x5e\x2f\x2e\x80\x8b\xb4\x47\x08\x21\x31\x4a\x2c\x21\x1d\xf1\xb9\x02\x8d\x31\x6b\x96\x8d\xa9\x94\xfa\x86\x14\x0e\x26\x09\x2d\x10\xad\x1f\x32\xa6\xf8\x5d\x2e\x74\x3f\x33\x06\x3d\x3a\x6e\xc3\x22\x37\x8a\xad\x37\xd8\xa0\x3f\xe8\xbf\x63\xb9\xf7\x8f\x7b\x7d\x25\x75\x3f\xf7\x9e\x12\x07\x65\x42\x3d\xce\x4b\xf0\x05\x00\x52\x22\x35\xc2\xd4\x49\x9c\x27\xd4\x17\x7c\xf0\xfe\x38\xfa\xeb\xdb\x77\x29\xc7\xe7\x9f\xe0\xf3\x1b\x71\xaa\xfe\xb9\xfc\x70\x33\xcf\xab\xb3\x0f\x67\x97\xd3\xc1\xdb\x0b\xf5\x35\xbf\xbd\x7d\x67\xf4\xe0\xf2\xbb\x98\x1e\x7f\xe3\xbf\x8e\xd4\xf8\xca\xdf\xb3\xcf\x7f\xbc\x9f\x65\xe2\xe4\x47\x71\x5c\xd1\x3a\xf7\xf0\xe5\xce\x78\x6f\x9c\x9c\x4a\x9d\x50\xae\x8d\x9e\x2b\x53\x79\x9a\xf6\x62\xd6\xdc\xbf\x17\x67\x46\xcc\x97\x97\x15\x72\x46\xf2\x92\x7b\x9f\xd0\xdc\x68\xe4\x52\x83\x8b\x26\x65\x25\x05\x4d\xd7\x21\x9f\xa2\x9c\xb9\x7d\x62\x79\x19\xa3\x8c\xee\x7c\xf4\xe6\xed\x33\x4c\x8d\xe3\xcb\xaa\x32\x6e\x25\x93\xca\x1a\x87\xe0\x68\x7a\x66\x14\xc4\x8c\x3f\x0b\xca\x84\x9c\x3d\xc9\xe0\xd9\xf2\x97\x28\x7a\xad\xac\x26\xc6\x29\xa2\x00\x0b\x23\x12\x3a\xba\x18\x5f\x51\xc2\x73\x94\x46\x27\xd4\x36\xed\xb1\xc5\xab\x69\x14\x9e\x41\x99\x8e\x9c\xf9\x01\x39\x92\x73\x31\x24\xb1\xd4\xb6\x42\xe2\xe0\x67\x25\x1d\x08\x82\x73\x0b\x09\xd5\x95\xca\xc0\x51\xa2\xb9\x82\x84\xda\x06\x7f\x2d\x05\x65\x69\xcc\x9a\x20\xbb\x0e\x18\x23\xc7\xca\x0f\x49\xec\xa1\x0c\x07\xad\xa3\x37\xf1\x7c\x6d\x6e\xc9\xb1\x0e\x63\x6c\xb8\x4e\x3a\x91\x9a\x97\xf2\x5e\xea\xe9\x35\x82\x53\x3e\x66\x4b\xc3\x5e\x4f\x0b\x5a\x04\xb7\x49\xa5\xc5\x01\x6e\x52\x5f\x5b\x67\xa6\x0e\xfc\x7f\x39\x4b\x6a\xe9\x0b\x10\xdd\x3d\x85\xd1\xb0\x1b\x1d\xb3\xa6\x84\x7b\xca\xde\x90\xd8\x70\x57\x48\x21\x40\xaf\xb8\x9b\x02\x5e\x5b\xee\xb8\xf2\x94\xcc\x78\x59\x41\x42\x17\x0b\xd2\x3f\x3d\xb9\x1a\xd5\xbb\xe4\xe1\x81\x12\xd6\x12\x36\xab\x10\x8d\x5e\xc6\xf5\x55\xa6\x24\xd2\x55\x87\x66\xa8\x49\x86\x3a\xb2\x4e\x2a\xee\xe6\x34\xfd\x6a\x05\x47\x88\x59\xe3\xb4\xa5\x6d\x59\xe8\xdb\x0e\x6f\x86\x44\xd1\xab\x3f\xe4\x8d\x27\x73\x7a\xd2\xfd\xc5\x6c\x2d\x41\x53\x5a\x33\x99\xf8\x20\x8d\xeb\xb2\x5a\x07\x33\xd2\xbf\xa8\xb7\x49\x7f\x8c\x60\x43\x71\xdb\xea\x35\x72\x30\x93\xa6\xf2\x24\xf0\x11\xc0\xe4\xe1\xa1\xbd\x7a\x07\x27\xa3\xe1\x0e\xbb\x27\xf3\x25\xa0\x3b\x25\xd2\x99\xc6\xd7\xa6\x10\x79\x56\xc2\x0a\xd9\x2c\xea\xbf\x91\x47\x27\x2d\x88\x36\x06\xf1\x71\x86\x6e\xb7\xbb\x76\xe3\x32\x40\x1a\x33\x2c\xf6\xa3\xce\x3f\x76\xc3\x7d\xe1\x61\x7c\x74\x41\x36\x42\xda\x0d\xfb\x77\x29\x41\xe3\xbe\x1c\x62\xd6\x76\xdf\xe0\xd7\x5a\xa9\x18\x1f\x27\xf0\xb6\x6f\xb1\x20\x47\x8d\xcc\x90\x61\xb2\x29\x2f\xc1\xe6\xb8\x9e\x02\x39\xba\x81\xf9\x6f\xe4\xc8\xd6\x90\xe5\x10\x0a\x88\xff\xc3\xcd\x0e\x66\xd7\xa0\x0e\xf3\x92\xd4\xbf\x75\x12\xaa\xb8\x9b\x4a\x3d\x24\xbf\xff\xb9\x63\x3c\x6d\x04\xdf\x3a\x3f\x57\x1a\xfc\xa8\xb8\x47\xb6\x7f\xfe\xb1\x16\xdb\x17\x33\xb5\x4d\x7e\x5b\x8e\xea\xae\xf2\x2b\x46\x76\x48\xfc\xf3\x6f\xb1\x90\x13\x02\x3f\x43\xba\x4d\xef\x11\xba\x31\x4a\xe9\x0e\xb6\x36\x72\x3d\x68\x74\x7c\xaa\xea\x23\xc8\x25\xe4\x20\x67\x61\x82\xee\x92\xc1\x97\x39\x83\x16\x1d\xf2\xda\xa6\x5c\x2f\x31\xfb\x3a\x2a\xb4\xdc\x13\x3e\x0f\x71\x08\xef\xfe\x40\x97\x25\x09\x87\x39\xad\x94\x60\xaf\x5b\xbb\x18\x90\x3d\x75\x8d\x59\x8b\x20\xc4\xac\x16\xe5\x2e\xf3\xa1\x17\xb3\x26\x46\xf8\x9d\x1f\xfe\xdf\xf9\x37\x00\x00\xff\xff\xef\x07\xbe\x74\xf6\x0c\x00\x00")

func templatesPaymentHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesPaymentHtml,
		"templates/payment.html",
	)
}

func templatesPaymentHtml() (*asset, error) {
	bytes, err := templatesPaymentHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/payment.html", size: 3318, mode: os.FileMode(420), modTime: time.Unix(1500040958, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesWebsockettestHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x5a\x5f\x6f\xdb\xb0\x11\x7f\xf7\xa7\xe0\xd4\x00\x72\xd0\x44\x72\xd2\x6d\x28\x54\x49\x2b\xb6\xf5\xa1\xc3\xda\x06\x48\xb7\x3c\x14\xc1\x40\x4b\xb4\xc5\x9a\x12\x05\x92\x76\x1a\x64\xfe\xee\x03\x49\xeb\x2f\x29\x4b\x4e\x93\xf9\x21\x10\xc9\xe3\xdd\xef\x8e\x77\x3f\xfe\x41\xc2\x4c\xe4\x24\x9e\xcd\xc2\x0c\xc1\x34\x9e\x01\x00\x40\xc8\xc5\x23\x41\xfa\x5b\xfe\xc4\x92\xa6\x8f\x40\xb0\xa0\x10\xd9\x65\x92\x61\x92\xce\x69\x9a\x9e\x83\xa7\x5a\x42\xfe\x96\x30\xd9\xac\x19\xdd\x16\xe9\x65\x42\x09\x65\x01\x78\x93\x24\xc9\x87\x5a\x66\xaf\x75\xfb\x07\xe5\xb3\xd0\xd7\x16\x67\xa1\x54\x7f\xb0\x9c\xe2\x1d\xc0\x69\xe4\xc0\xb2\x74\x1a\x00\xb2\x3b\xee\x18\x0b\xcb\x38\x84\x20\x63\x68\x15\x39\x3f\xe1\x0e\xf2\x84\xe1\x52\x04\x19\xe6\x82\xb2\x47\x4f\x62\x99\x9f\x3b\xf1\x5f\x61\xb2\x09\x7d\x18\x83\xff\x82\x2d\x23\x41\x88\x8b\x72\x2b\xc0\xee\x32\xa7\x29\x22\x91\xb3\x84\x1c\xfd\x8b\x11\x27\x0e\x7d\x35\xd2\xb3\xe1\x97\x2d\x08\x7e\x07\x43\x98\xb1\xa3\xf0\xb2\xeb\xf8\x9f\x74\x8d\x8b\xd0\xcf\xae\x0d\xe4\x25\x14\x59\x70\xf6\x04\xcb\x92\xe0\x04\x0a\x4c\x8b\x1b\x28\xb2\x7d\xc7\x9e\x12\x25\x70\x89\x48\xfc\x29\x87\x98\x04\x9d\x11\x35\xda\x73\x06\x49\x31\x27\x36\xe5\x38\x22\x28\x69\x09\xea\x36\x4a\x3f\x0d\x4c\x50\x93\x68\x29\x71\x81\xdd\xe5\x8a\xb2\x83\x6e\x80\x0b\xa0\x3e\xb8\x13\x9f\x3d\xa9\xaf\x7d\xe8\x6b\x41\x8b\x59\x5f\xdb\xe9\xc7\x54\xfb\x64\xf3\xf3\x06\x72\xfe\x40\x59\x3a\xee\x6a\x79\x90\x74\xfc\x13\xdc\xad\xd4\x4f\xf3\xb8\x32\x21\x9d\xae\xbe\x95\xdf\x55\xe3\x85\x5c\x5f\x6e\x85\x50\x56\x69\x11\x24\x04\x27\x9b\xc8\x21\x32\x71\x9c\x2a\x7f\xb4\x80\x91\x43\xdf\xe9\x06\x15\x41\x28\xd0\x2f\x01\x19\x82\x20\xa1\x84\x47\xce\xd5\x62\xe1\x00\x46\x1f\x78\xe4\xbc\x93\x60\x85\x94\xda\x87\x7e\x25\x16\xdb\x72\x7a\x24\x8f\xbf\x20\xce\xe1\x1a\x17\xeb\x90\x97\x50\x22\xc5\xab\xc8\xc9\xab\xce\x3b\x0e\xfe\x10\x81\x62\x4b\x88\x13\x83\x4b\xf0\xed\xe6\xd3\xd7\xd0\x97\x82\xb1\x35\xf5\x8d\x58\x1d\x6a\xa1\xd6\xa7\x2a\xc1\x8c\xa8\x25\x4c\xb4\x44\x45\x8d\xed\xee\xd6\x89\xbf\x95\x68\x20\x5e\x43\x2a\x12\x42\x39\xea\xe8\xf8\x9b\xec\x19\x08\xfa\x40\x79\xde\x30\xfa\x53\x26\x1c\x4e\x2b\x7e\x11\x8f\x25\x8a\x9c\x62\x9b\x2f\x11\x73\x5a\x59\xab\x05\xff\x83\x53\x07\xf8\xb1\x3d\x1d\x8c\x15\x50\xbd\xbd\x55\x7e\xbf\x68\x69\xad\x23\xa7\xfd\x40\x92\xcc\xea\xe5\x9e\x14\x05\x8e\x8a\xf4\x3b\xed\x84\xe1\x16\x15\xe9\x50\x14\xcc\x24\x11\x70\xd9\xde\x2e\x9a\xfe\x66\x4b\xb1\x8c\xc5\x5b\x8e\xd8\xe7\x34\xf4\x45\x76\x5c\xe6\xfb\x63\x89\xc6\xa5\x0a\x98\x8f\x48\xc9\xa8\x8c\x48\xe0\x1c\x71\x01\xf3\xf2\xb8\xd8\x61\x21\x87\xb0\xcb\x5e\xab\xe3\xa1\x68\x76\x3a\x73\x8c\x55\xd4\x93\xf3\xb5\x64\x9d\xfe\xc2\xf2\x01\xe2\xd2\xb3\xd3\xf8\xec\x29\xe7\x6b\x6f\xc5\x68\xee\x61\xc9\x4e\x62\x20\xf4\x16\x79\x99\xb0\xa7\xcd\xa8\x22\x3e\x79\x56\x0a\x05\x9c\x2c\x5c\xaf\xc3\xe4\x19\xf5\x92\x1c\x99\x11\xfa\x82\x59\x57\xcb\x5c\x95\xd0\xef\xe5\x74\x7f\xfb\xb7\x32\xe5\x57\x2a\xf0\xea\xb0\xa3\xb7\xc9\xb2\x68\xf5\xbf\x10\x5f\xb6\x55\x9e\x46\x99\x6d\x90\xbf\xc1\x9a\x7d\x35\x27\x11\xe7\x33\x09\x43\x66\xde\x48\x01\xd3\x91\xf1\x51\x2a\x91\x79\xfa\x22\x24\xc1\x10\x7c\x5d\x7e\x68\xa7\xc0\x44\x8a\xf8\xfc\xf7\xa0\xcb\x12\x17\x40\x92\x6b\x70\x3a\x15\x0c\x0e\x2a\x81\xb2\x82\x2a\xa8\x62\x32\x59\xd2\xd4\xd1\xe6\x05\xed\x18\x16\xb4\x32\x69\x49\xf6\x56\xc4\xa6\xf1\xc6\x8b\x82\x9f\xef\x20\xd9\xa2\x0b\xb0\x41\x8f\xe7\x95\x1b\x32\x3d\xe4\xc1\x6a\x83\x1e\xf7\x01\x38\x7b\x52\x22\x2f\x04\xfe\x64\xd2\x93\x29\xf6\xfa\x7c\x77\xf8\x0c\x7d\x3d\x2d\xd4\x57\x2d\xc0\x59\x12\x39\x99\x10\x25\x0f\x7c\x7f\x5b\x94\x9b\xb5\x97\xd0\xdc\xdf\x6d\xd1\xc7\x6b\x6f\xe1\xbd\xf3\x53\xcc\x85\x6c\x7a\x39\x2e\xbc\x9f\x5c\x9e\x4b\xf4\xd4\x51\x1d\xf0\x17\xa6\xfc\xe3\xc2\xbb\xba\xf6\x16\x5a\x8d\xea\x19\x56\xa4\x81\xee\x20\x03\xb0\x2c\x41\x04\x0a\xf4\x00\xfe\xbd\x45\xf3\xe6\x82\x9a\x22\x82\x73\x2c\x10\xe3\x01\xf8\xe1\x9e\x3d\xb9\x17\xc0\xdd\xbb\xf7\x17\xb5\x00\x22\x01\x70\xdf\xc0\xb2\x74\x9b\x3e\xb9\xd8\x81\x71\xcb\x55\x57\xc6\x00\xb8\x84\x26\x90\x78\x2b\x88\x19\x81\x45\x82\x3c\x4c\x83\xf7\x8b\xf7\x8b\xd6\x7c\xf9\xf3\x7d\xa0\xce\xf3\x9d\x4e\x7d\x8b\x92\x48\x56\x98\x71\xf1\x51\x36\xa5\xe7\x12\x15\x47\x09\x2d\xd2\x4e\x57\x42\x30\x2a\x5a\x52\xf7\x5d\x13\xf5\xf5\x44\x2a\x94\xb7\x1c\x39\xe7\xea\xfa\xdd\x1f\xff\xf4\xe7\xbe\x68\xe7\xea\x17\x00\xc3\xbc\x55\xb8\xbe\x97\x81\x83\x76\xd3\x97\x31\x55\xe5\x51\x15\xea\x96\x12\x00\xb7\xd7\xdd\xbb\x1f\x07\xc0\xf5\x61\x89\xfd\x56\xb7\x19\xeb\xfa\xd0\xd4\x19\xe8\xdc\x2e\x2a\x3d\x75\x67\x4f\x8b\x71\xee\x0a\xc0\x8f\xfb\x01\x91\x3b\x1e\xa8\xbd\xfc\xb8\x06\xd3\xb5\xe6\x2a\x10\x80\x2b\xc3\x89\x36\xb3\x77\xc6\xfa\xbb\x7e\xe5\x4a\xbb\xbf\x67\xc9\xb6\x4b\x98\x0e\x75\x4f\x28\x36\x9f\x2c\x7a\xba\x6e\xed\x9b\xcf\x07\x28\x92\xac\x5f\x38\xbd\xd4\x5b\x6d\x8b\x44\xdd\xb7\xe7\x68\x87\x0a\xd1\x7f\x4c\x92\x3f\x91\x61\xee\xe9\x97\x87\x48\x37\x3a\x3a\x3e\x74\x26\xec\xc7\x72\x77\xaa\xc1\xfa\xe2\xdf\xb3\x59\x69\x1a\x34\xdb\xfa\x4c\x68\x5e\x6e\x05\x4a\xcd\x18\x14\x69\xfb\xe0\xa4\x88\xa4\x01\x66\xc3\xc4\x90\xd8\xb2\x42\x43\x39\x70\x0f\x78\xab\x9b\xfd\x6c\x00\x6f\x81\xeb\x4b\x13\x6e\x17\xa2\x0d\x61\x8e\x44\x46\x25\x63\x74\x2d\x2a\xaa\x9a\x14\x2b\x49\xb6\x1c\x91\xd5\x21\x4c\x1f\x0c\x01\x4d\xd9\x25\xe5\x62\xee\x4a\x8a\x0f\x7c\xdf\xad\x90\xf7\x1c\xe9\x15\xba\xf2\x43\x21\x71\x2f\x2c\x96\x41\x43\x3a\x4d\x86\x5c\x58\xc5\x1a\xda\xe9\x2c\xad\x29\xbc\x3f\xb7\xce\xf7\x44\x86\x8a\x79\x13\x0d\x86\x78\x49\x0b\x8e\x6c\x01\xa9\x7e\x78\xd5\xc8\xa9\xe3\x82\x97\xd0\x14\x81\x28\x02\xd7\x8b\xc5\xb1\x89\x40\xa7\xed\xca\x53\x74\x08\x22\xd0\xd5\xa2\xfe\xa8\x21\x33\xd6\xb5\x1b\x83\x23\x09\x2d\x38\x25\xc8\x23\x74\x3d\xd7\x0f\x4c\x81\x73\x51\x5b\x38\xb7\xab\x1c\x8a\x4a\x22\xeb\xbb\x15\x16\xc4\x18\x65\xc7\x5c\x53\xd1\xef\x3f\xb6\xcc\x07\xac\x0e\xc3\xd5\x76\x26\x61\xed\xd1\x41\xef\xad\x68\x52\x86\xcb\x85\x54\xc0\x2d\x6f\x5d\x43\xce\xea\x6a\x35\x11\x9a\x0b\x33\x5a\x3f\x86\x69\x7d\xac\xb9\x43\xcb\x5b\x9a\x6c\x90\x98\xbb\x0f\xfc\x48\x4d\x75\xb6\x3c\xf0\x16\x38\xbe\x53\x0d\x35\x9b\x8f\xea\x7f\xe0\x7f\x51\x69\x15\xd5\x02\xaa\x69\x89\x73\x1f\x92\x47\x0b\x19\x58\x10\x4d\x09\x27\xe8\x2f\xab\xa9\x48\xad\xb1\x9c\xaf\x92\xdd\x02\x60\x3f\x09\x93\x6e\xa1\xe9\xb0\x54\xd5\x19\xdb\x3e\x88\x06\x06\xbc\x84\x16\x09\x14\xf3\x7f\xdc\x7e\xfb\xea\x95\x90\x71\x34\x6f\xa1\xb6\xd7\xcc\x11\xcf\x0f\x70\x5f\xc8\x79\x55\x23\xcf\x75\x5d\x67\xd9\x96\x10\x7b\x8d\x41\x82\x98\x98\x3b\x2b\x88\x09\x4a\x9d\x81\x42\x3c\xe2\xa9\xc2\x56\xfb\x39\xee\xa2\xb1\xa7\xf7\xde\x29\xa7\x6f\xea\x6d\x1c\x52\x4d\xaf\xb0\x0f\x2b\x3b\x9a\xf2\x5f\xea\xcc\x72\xdd\xa3\x50\xfb\x5c\xf7\x6a\x84\x63\x78\xa7\x2c\x0f\x51\xab\x8d\x53\x06\x57\xdb\xea\xba\x9c\xf2\xe3\x7e\x8c\xe0\x2c\xec\xdb\x7d\x2f\x3a\x2d\x1e\xf6\xf7\xb3\xff\x17\x07\xf7\xac\x9f\x46\xc3\x96\x33\xda\xc9\x8c\xdb\x05\xf0\x7b\xa4\x6b\xd5\xf5\x6c\xea\x31\xb4\x3d\x8f\x7a\x6d\x77\x94\x8a\x7d\x6d\x63\xbf\x4d\xc0\x43\xb8\xc7\x03\x31\x31\x0e\xcf\x61\x61\x33\xcf\x5e\x90\x88\xed\x08\x9f\xcf\xc5\x96\x37\xe0\xd7\xac\x69\x5b\x94\xc7\x99\xee\x84\x88\x1a\xf2\xcf\xe5\x3b\x3d\xb0\x3f\x9f\xcd\x9a\x57\xaa\x59\xe8\xab\x7f\xb2\xf8\x5f\x00\x00\x00\xff\xff\xba\xd3\x45\x89\x6b\x21\x00\x00")

func templatesWebsockettestHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesWebsockettestHtml,
		"templates/websockettest.html",
	)
}

func templatesWebsockettestHtml() (*asset, error) {
	bytes, err := templatesWebsockettestHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/websockettest.html", size: 8555, mode: os.FileMode(420), modTime: time.Unix(1499600894, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/index.html": templatesIndexHtml,
	"templates/payment.html": templatesPaymentHtml,
	"templates/websockettest.html": templatesWebsockettestHtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"index.html": &bintree{templatesIndexHtml, map[string]*bintree{}},
		"payment.html": &bintree{templatesPaymentHtml, map[string]*bintree{}},
		"websockettest.html": &bintree{templatesWebsockettestHtml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

