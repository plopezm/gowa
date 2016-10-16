package gowa

import (
	"runtime"
	"path"
	"fmt"
)

func GowaStart(dbtype string, dbpath string, pageSize uint32) *GowaManager{
	GM = new(GowaManager)
	GM.Init(dbtype, dbpath, pageSize)
	return GM
}

func GowaGetTemplatePath() string{
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(path.Join(path.Dir(filename), "template"))
	return path.Join(path.Dir(filename), "template")
}
