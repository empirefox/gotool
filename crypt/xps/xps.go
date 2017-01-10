package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/empirefox/gotool/crypt"
)

var (
	decryptdir = flag.String("d", "", "Decrypt dir to place files")
	password   = flag.String("k", "", "Password")
	configfile = flag.String("x", "xps-config.json", "Config file for xps")
)

func main() {
	flag.Parse()
	xps, err := crypt.NewXps(*configfile)
	if err != nil {
		panic(err)
	}

	if *decryptdir != "" {
		err = os.RemoveAll(*decryptdir)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(*decryptdir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		files, err := xps.DecryptXhexFile(*password)
		if err != nil {
			panic(err)
		}

		for name, content := range files {
			err = ioutil.WriteFile(filepath.Join(*decryptdir, name), content, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		fmt.Printf("Decrypted ok to dir: %s\n", *decryptdir)
	} else {
		err = xps.EncryptXhexFile(*password)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Encrypted ok to file: %s\n", xps.XpsFile)
	}
}
