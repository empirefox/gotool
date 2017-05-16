package crypt

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/brentp/xopen"
	"github.com/mcuadros/go-defaults"
)

var (
	ErrNoFiles      = errors.New("No files to encrypt")
	ErrWeakPassword = errors.New("Weak password")
)

type Xps struct {
	Password   string            `json:"password,omitempty"`
	Files      map[string]string `json:"files,omitempty"`
	XpsFile    string            `json:"xps-file,omitempty" default:"xps.tar.gz"`
	GzipLevel  int               `json:"gzip-level"         default:"-1"`
	ConfigFile string            `json:"config-file"        default:"config.json"` // in xps file: json ymal toml json5
	EquipTag   string            `json:"equip-tag"          default:"xps"`
}

func NewXps(filepath, filetype string) (*Xps, error) {
	file, err := xopen.Ropen(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	xps := new(Xps)
	err = UnmarshalFormat(content, xps, filetype)
	if err != nil {
		return nil, err
	}

	defaults.SetDefaults(xps)
	return xps, nil
}

func (xps *Xps) EncryptXhexFile(password string) error {
	if password == "" {
		password = xps.Password
	}
	if len(password) < 8 {
		return ErrWeakPassword
	}
	if len(xps.Files) == 0 {
		return ErrNoFiles
	}

	xx20, salt, err := NewAEAD([]byte(password), nil)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	gw, err := gzip.NewWriterLevel(buf, xps.GzipLevel)
	if err != nil {
		return err
	}
	tw := tar.NewWriter(gw)

	// file meta info
	// v0.1 only store salt
	hdr := &tar.Header{
		Name: "",
		Mode: 0600,
		Size: 2 + int64(len(salt)),
	}
	if err = tw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err = tw.Write([]byte{0, 1}); err != nil {
		return err
	}
	if _, err = tw.Write(salt); err != nil {
		return err
	}

	// TODO goroutine
	for vname, filename := range xps.Files {
		plaintext, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		ciphertext, err := EncryptXX20p1305(xx20, plaintext)
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name: vname,
			Mode: 0600,
			Size: int64(len(ciphertext)),
		}
		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err = tw.Write(ciphertext); err != nil {
			return err
		}

		//		fmt.Printf("%s md5:%x\n", vname, md5.Sum(ciphertext))
	}

	// Make sure to check the error on Close.
	if err = tw.Close(); err != nil {
		return err
	}
	if err = gw.Close(); err != nil {
		return err
	}

	return ioutil.WriteFile(xps.XpsFile, buf.Bytes(), os.ModePerm)
}

func (xps *Xps) DecryptXhexFile(password string) (map[string][]byte, error) {
	if password == "" {
		password = xps.Password
	}
	if len(password) < 8 {
		return nil, ErrWeakPassword
	}

	// Open the tar archive for reading.
	file, err := os.OpenFile(xps.XpsFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)

	_, err = tr.Next()
	if err != nil {
		return nil, err
	}

	versalt, err := ioutil.ReadAll(tr)
	if err != nil {
		return nil, err
	}
	fmt.Printf("xps file v%d.%d\n", versalt[0], versalt[1])

	xx20, _, err := NewAEAD([]byte(password), versalt[2:])
	if err != nil {
		return nil, err
	}

	files := make(map[string][]byte)
	// Iterate through the files in the archive.
	// TODO goroutine
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return nil, err
		}
		//		fmt.Println(hdr.Name)

		ciphertext, err := ioutil.ReadAll(tr)
		if err != nil {
			return nil, err
		}

		//		fmt.Printf("%s md5:%x\n", hdr.Name, md5.Sum(ciphertext))

		if _, ok := xps.Files[hdr.Name]; !ok {
			fmt.Println("File not in config:", hdr.Name)
		}

		files[hdr.Name], err = DecryptXX20p1305(xx20, ciphertext)
		if err != nil {
			return nil, err
		}
	}

	if len(xps.Files) != len(files) {
		for vname := range xps.Files {
			if _, ok := files[vname]; !ok {
				fmt.Println("File not in xps:", vname)
			}
		}
	}

	return files, nil
}
