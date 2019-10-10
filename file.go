package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const DIRECTORY_SEPARATOR = string(os.PathSeparator)

var RobotLogDir = GetCurrentDirectory() + "/logs"

func FileGetContents(filename string) ([]byte, error) {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

func FilePutContents(filename string, fileData []byte) error {
	return ioutil.WriteFile(filename, fileData, os.ModePerm)
}

func GetSize(filename string) (int64, error) {
	file, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return file.Size(), nil
}

func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + DIRECTORY_SEPARATOR + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		panic(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func GetParentDirectory(dirctory string) string {
	return SubStr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func Wlog(dirName, prefix, content string) {
	dir := RobotLogDir + DIRECTORY_SEPARATOR + dirName
	err := IsNotExistMkDir(dir)
	if err != nil {
		return
	}
	fileName := dir + DIRECTORY_SEPARATOR + GetNow().Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	defer logFile.Close()
	if err != nil {
		panic("open file error")
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		content += fmt.Sprintf(" %s:%d", file, line)
	}
	l := log.New(logFile, "["+prefix+"] ", log.LstdFlags) //|log.Llongfile
	l.Println(content + "\n")
}

//完整名称
func ParseFileName(fullFilename string) (string, string, string) {
	var filenameWithSuffix, fileSuffix, filenameOnly string
	filenameWithSuffix = path.Base(fullFilename)
	fileSuffix = path.Ext(filenameWithSuffix)
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	return filenameWithSuffix, filenameOnly, fileSuffix
}

func FileCopy(source, dest string) (bool, error) {
	fd1, err := os.Open(source)
	if err != nil {
		return false, err
	}
	defer fd1.Close()

	//os.O_TRUNC 覆盖写入
	//
	fd2, err := os.OpenFile(dest, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return false, err
	}
	defer fd2.Close()
	_, e := io.Copy(fd2, fd1)
	if e != nil {
		return false, e
	}
	return true, nil
}

func GetAllFile(pathName string, s []string) ([]string, error) {
	if !strings.HasSuffix(pathName, DIRECTORY_SEPARATOR) {
		pathName += DIRECTORY_SEPARATOR
	}
	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathName + DIRECTORY_SEPARATOR + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				return s, err
			}
		} else {
			fullName := pathName + DIRECTORY_SEPARATOR + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}
