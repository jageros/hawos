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
	module  = flag.String("module", "", "project base package path")
	inDir   = flag.String("indir", "protos/pbdef", "proto file directory")
	outDir  = flag.String("outdir", "protos/", "generate meta go package directory")
	pbPkg   = flag.String("pbpkg", "", "proto generate pb file go package")
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

	if *module == "" {
		gomod, err := ioutil.ReadFile("go.mod")
		if err != nil {
			log.Fatalf("Input module name or exec in go mod project root directory.")
		}
		strs := strings.Split(string(gomod), "\n")
		for _, str := range strs {
			if strings.HasPrefix(str, "module") {
				ss := strings.Split(str, " ")
				*module = ss[1]
				break
			}
		}
	}

	if *pbPkg == "" {
		*pbPkg = fmt.Sprintf("%s/protos/pb", *module)
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

	sessPkg := fmt.Sprintf("%s/%s/meta/sess", *module, *outDir)
	sessPkg = strings.Replace(sessPkg, "//", "/", -1)

	err = writeToFile(metatemp.GenIMetaFile(*eumName, sessPkg, *pbPkg), fmt.Sprintf("%s/meta/imeta.go", *outDir))
	if err != nil {
		log.Fatal(err)
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
			code, msgid, err := metatemp.GenMetaFile(file.Name(), *pbPkg, sessPkg, *eumName, lines)
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
	registerFile := metatemp.GenMetaRegister(msgids, *inDir)
	err = writeToFile(registerFile, fmt.Sprintf("%s/meta/meta_register.go", *outDir))
	if err != nil {
		log.Fatal(err)
	}

	sessPath := fmt.Sprintf("%s/meta/sess", *outDir)
	sessPath = strings.Replace(sessPath, "//", "/", -1)
	_, err = os.Stat(sessPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(sessPath, fs.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = writeToFile(metatemp.ISessTemp, sessPath+"/sess.go")
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
