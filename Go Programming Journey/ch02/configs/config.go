// Code generated for package configs by go-bindata DO NOT EDIT. (@generated)
// sources:
// configs/config.yml
package configs

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

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _configsConfigYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x54\x5d\x4f\x1b\x47\x14\x7d\x47\xe2\x3f\x5c\xc9\x2f\xed\x43\xcc\x1a\x07\x92\xce\x53\x80\x90\x04\x84\x5b\x0b\x6f\xea\xf6\x29\x1a\xbc\xc3\xb2\xd5\xee\xce\x32\x33\x4b\x4c\x9e\x48\x05\x6d\x42\x63\x40\x2a\x60\x44\x40\x02\x95\x16\x57\x55\x1c\xaa\x22\x40\x06\xca\x9f\xf1\xec\x9a\x7f\x51\xcd\xae\xb1\x37\x04\xf5\xcd\x3e\x7b\xef\xb9\x1f\xe7\xdc\x29\x10\x36\x47\x18\x82\x14\x04\x3b\x15\xb9\xbc\x7f\xbd\x54\x09\x2f\xeb\xbd\x3d\x00\x93\xbe\x9b\xa3\x06\x41\x60\x90\x29\xdf\x84\x14\xb4\xea\xff\x86\x97\x75\x30\x2d\x17\xc2\xed\xc5\xd6\xd5\x5a\x6b\xff\x5d\x50\xdb\x97\x17\xab\x2a\xfc\x99\x10\x5e\x9e\x32\x81\xe0\xa1\xa6\x69\xd0\x21\x0c\xff\xfa\x28\x57\x7f\x93\xab\xa7\x11\x27\xc1\x86\x6e\x39\x84\xfa\x02\xc1\xa0\xa6\xa0\x22\xb3\x04\xf9\x04\x1b\xf2\x3c\xd5\x8f\x6c\xac\x87\xeb\xb5\x6e\x3f\x8f\xc9\x34\xf6\x6d\x91\xc7\x26\x29\x58\xaf\x08\x82\x4c\x94\x9f\xc3\xe5\x24\x14\x61\x13\xd4\x2c\xe0\x39\x92\xc7\x62\x06\x01\x17\x94\x61\x93\xf4\xd9\xd4\xe4\x90\x82\xeb\xf3\xad\x56\xfd\x20\x26\x0f\xaa\xbf\xcb\xab\xaa\xfc\xb0\x25\x7f\xac\x35\x2f\x6f\x0a\x4d\x50\xf3\x89\x65\x93\xaf\xb1\x43\x10\x60\xcf\xbb\x3b\x69\xad\x12\x1e\x1e\x25\xc2\x47\xcb\x02\x41\xda\xa6\xe6\x9d\xe1\xc1\xe6\xcf\xcd\xf3\x13\xb9\xb6\x12\x5e\x2c\xc8\xb5\x8a\xca\x4b\x41\x70\x7a\x2e\x97\xf7\x9a\x67\xcb\xcd\x8b\xbd\xf0\xfd\x99\x5c\xfa\xa7\x3b\xed\x73\xcf\xa6\xd8\xf8\x7c\x0a\x3f\xc2\xd5\x20\x71\x5e\x4c\x1c\x6e\x2f\x06\x3b\x0b\xe1\xf9\x9b\xe6\xd5\xae\xfc\xb0\x15\xbe\xaf\xcb\xcb\x8d\x04\x4d\x24\xf2\x73\x66\x23\x98\x11\xc2\x43\x7d\x7d\x99\xfe\x07\x69\x2d\xad\xa5\x33\x48\xa9\xd5\xc7\x05\x16\x56\xe9\x16\xa9\xea\x76\xbd\xd6\x6c\xac\xc8\xbf\x37\xc2\x83\x86\xaa\x11\xe1\xb1\xae\x72\xe7\x48\xee\x2e\x74\x6b\x8c\x39\xd8\x24\x39\x5c\x8e\x85\x18\x80\x5b\x64\xc1\xdb\x05\xb9\xf4\xba\x55\x3f\x8b\x5b\x95\x07\x87\xe1\x9f\x8d\xeb\xea\xf1\x17\xb9\xe1\x2f\x6f\x91\x0c\xd9\x36\x7d\x39\x5a\x16\x1c\xfd\x1f\x47\x62\xa1\x2a\x1f\xe0\x1e\xa4\x7f\xf0\xcc\xc4\x6f\xd2\xfd\xe3\xb9\x66\xbc\xf2\xd8\xc3\xad\x93\xa5\xa0\x7a\x12\x54\x4f\xae\xab\xc7\x09\x6b\x8d\x50\x57\x90\xb2\xe8\x98\x51\x19\xec\x31\x16\x78\x0a\x73\x12\x5d\xc8\xc6\x51\x50\xa9\xcb\xc6\xaf\x09\x53\x0e\xeb\xf3\x1e\x41\xe0\xcc\xf3\x59\x3b\x9a\x83\x13\xe6\x46\xce\x61\x94\x0a\x48\x66\xb5\x8e\xff\x68\x9f\x41\x1e\x73\xfe\x92\x32\xe3\x8e\x20\xf9\xf1\xa7\x70\xef\x75\x74\x50\x94\xab\x1e\x3a\x42\x65\xb3\xda\x60\x5c\x32\x76\x66\x69\x46\xeb\xff\x24\xb3\xe3\x48\x1d\x4f\xd9\x24\xcf\xc8\xb4\x55\x46\x30\x65\x53\xf3\x85\x9a\x7c\xbf\x16\x47\xc8\xb7\x95\xf6\xce\x46\x66\x30\xe3\x44\x20\xf0\xc5\xf4\xc3\xb8\x2f\xc6\xa3\x5b\x44\xa0\x33\x9f\xb4\xcf\x6b\xcc\xb0\xc9\x08\x75\x5d\x9e\x38\xb9\x6f\x3c\xe2\xb6\xb1\xac\xd6\xdb\xd3\xdb\x93\x82\xf1\xa2\x0e\xf2\xcd\xae\x3c\xfc\x45\xbe\xdb\xbc\xd9\xd0\x78\x51\x47\x2a\xa5\x40\x4a\x4c\x55\xc2\x86\x63\xb9\x0a\x18\xe3\xdc\x57\xef\x8e\x6a\xef\x1e\x27\x6c\xce\x2a\x45\xf5\x46\xcb\x9e\xc5\x08\x82\x07\xfd\x5a\x9b\x77\xd4\xc1\x96\xfd\x39\x73\x04\xa3\xee\xa2\xb8\x23\xbc\x74\x66\x30\x9b\x2e\x51\x27\x9a\x25\x7a\x8b\xee\x0f\x0e\xdc\xa8\x12\x6f\xcd\xa0\x73\xe4\xc5\xab\xf9\xd2\xa3\x64\x68\x47\x0e\xfd\x7b\x7d\x7c\xf2\xdb\xe2\xe4\x77\x4f\x87\x9e\x3d\xd5\x87\x8b\x6a\x6f\x2b\xa7\x72\x75\x33\xdc\x5e\x84\x42\x4e\xcf\x43\x57\x9f\x31\x5e\x28\x4c\x20\x10\xed\x45\x3d\x61\xd4\xb9\x9b\x5e\xa7\xe8\xc6\x88\x99\x8c\x96\x1d\xf8\xaa\x5f\xbb\xaf\x3d\x9a\x9d\x55\xdf\xff\x0b\x00\x00\xff\xff\xaf\x12\xf6\xc8\x7f\x05\x00\x00")

func configsConfigYmlBytes() ([]byte, error) {
	return bindataRead(
		_configsConfigYml,
		"configs/config.yml",
	)
}

func configsConfigYml() (*asset, error) {
	bytes, err := configsConfigYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "configs/config.yml", size: 1407, mode: os.FileMode(438), modTime: time.Unix(1652715613, 0)}
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
	"configs/config.yml": configsConfigYml,
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
	"configs": &bintree{nil, map[string]*bintree{
		"config.yml": &bintree{configsConfigYml, map[string]*bintree{}},
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
