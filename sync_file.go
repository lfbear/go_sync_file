package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Error in args. eg: cmd from_path to_path")
		return
	}

	fmt.Printf("Rsync start at %s", time.Now().Format("2006/01/02 15:04:05"))
	fmt.Println("")

	srcDir := os.Args[1]
	tarDir := os.Args[2]

	for {
		foo := syncFile(srcDir, tarDir)
		if foo == false {
			break
		}
		time.Sleep(10 * time.Second)
	}
}

func syncFile(srcDir string, tarDir string) bool {
	if testDir(srcDir) && testDir(tarDir) {
		x := eachAllFiles(srcDir, "")
		for k, v := range x {
			if v == "dir" { //目标目录中缺少子目录
				todir := tarDir + k
				if testDir(todir) == false {
					fmt.Printf("%s | dir not existed: %s [mkdir ok]", time.Now().Format("2006/01/02 15:04:05"), k)
					fmt.Println("")
					cmd := exec.Command("mkdir", "-p", todir)
					cmd.Run()
				}
			} else { //文件检测
				tofile := tarDir + k
				testmd5 := fileMd5sum(tofile)
				if testmd5 != v { //md5值不同
					fmt.Printf("%s | diff found: %s [rsync ok]", time.Now().Format("2006/01/02 15:04:05"), k)
					fmt.Println("")
					fromfile := srcDir + k
					cmd := exec.Command("cp", "-r", "-f", fromfile, tofile)
					cmd.Run()
				}
			}
		}
		return true
	} else {
		fmt.Println("Error: dir is not existed!")
		return false
	}
}

func testDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		//fmt.Println(err)
		return false
	}
	defer f.Close()
	return true
}

func eachAllFiles(dir string, sub string) map[string]string {
	rs := map[string]string{}
	subdir := ""
	if sub != "/" {
		subdir += sub + "/"
	}
	dirname := dir + subdir
	files, _ := ioutil.ReadDir(dirname)
	for _, f := range files {
		if f.IsDir() {
			xdir := subdir + f.Name()
			rs[xdir] = "dir"
			mergeMap(rs, eachAllFiles(dir, xdir))
		} else {
			fullname := dirname + f.Name()
			bytes, _ := ioutil.ReadFile(fullname)
			h := md5.New()
			io.WriteString(h, string(bytes))
			key := subdir + f.Name()
			rs[key] = fmt.Sprintf("%x", h.Sum(nil))
		}
	}
	return rs
}

func fileMd5sum(filename string) string {
	if testDir(filename) == false {
		return ""
	}
	bytes, _ := ioutil.ReadFile(filename)
	h := md5.New()
	io.WriteString(h, string(bytes))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func mergeMap(a map[string]string, b map[string]string) {
	for k, v := range b {
		a[k] = v
	}
}
