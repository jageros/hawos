/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    genmate
 * @Date:    2022/3/8 5:03 下午
 * @package: genmate
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	"flag"
	"fmt"
	"github.com/jageros/hawox"
	"github.com/jageros/hawox/tools/metactl/internal/metatemp"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	inDir   = flag.String("indir", "protos/pbdef", "proto file directory")
	outDir  = flag.String("outdir", "protos/", "generate meta go package directory")
	inPkg   = flag.String("inpkg", "", "proto generate file go package")
	eumName = flag.String("enum", "MsgID", "msg id enum name")
	version = flag.Bool("version", false, "show metactl version")
	v       = flag.Bool("v", false, "show metactl version")
)

func main() {
	flag.Parse()

	if *v || *version {
		fmt.Println(hawox.Version)
		return
	}

	if *inPkg == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("os.Getwd err=%v", err)
		}
		dirs := strings.Split(dir, "src/")
		if len(dirs) >= 2 {
			*inPkg = fmt.Sprintf("%s/protos/pb", dirs[1])
		} else {
			*inPkg = fmt.Sprintf("%s/protos/pb", dir)
		}
	}

	start := time.Now()
	files, err := ioutil.ReadDir(*inDir)
	if err != nil {
		log.Fatalf("ReadDir error: %v", err)
	}
	if len(files) <= 0 {
		log.Fatalf("Dir %s not has file", *inDir)
	}

	outPath := strings.Replace(*outDir+"/meta", "//", "/", -1)
	_, err = os.Stat(outPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(outPath, fs.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	var msgids []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".proto") {
			log.Printf("Parse %s ...", file.Name())
			start := time.Now()
			path := fmt.Sprintf("%s/%s", *inDir, file.Name())
			text, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}
			lines := strings.Split(string(text), "\n")
			code, msgid, err := metatemp.GenMetaFile(file.Name(), *inPkg, *eumName, lines)
			if err != nil {
				continue
			}
			fName := strings.Split(file.Name(), ".")[0]
			err = writeToFile(code, fmt.Sprintf("%s/meta/%s.meta.go", *outDir, fName))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Gen %s meta file done. (%s)", file.Name(), time.Now().Sub(start).String())
			msgids = append(msgids, msgid...)
		}
	}
	registerFile := metatemp.GenMetaRegister(msgids)
	err = writeToFile(registerFile, fmt.Sprintf("%s/meta/meta_register.go", *outDir))
	if err != nil {
		log.Fatal(err)
	}
	err = writeToFile(metatemp.GenIMetaFile(*eumName, *inPkg), fmt.Sprintf("%s/meta/imeta.go", *outDir))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Gen all meta file done.(%s)", time.Now().Sub(start).String())
}

func writeToFile(content, path string) error {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}
