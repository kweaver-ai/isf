package easylog

import (
	"fmt"
	"os"
)

// RotatingFile 自动切分文件
type RotatingFile struct {
	name        string
	maxBytes    int64
	backupCount int
	f           *os.File
	//用备份长度加速
	curBytes int64
	curFirst bool
}

// NewRotatingFile 创建RotatingFile
func NewRotatingFile(name string, maxBytes int64, backupCount int) *RotatingFile {
	return &RotatingFile{name: name, maxBytes: maxBytes, backupCount: backupCount}
}

func (rf *RotatingFile) doRollover() error {
	if rf.f != nil {
		rf.f.Close()
		rf.f = nil
	}
	if rf.backupCount > 0 {
		for i := rf.backupCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", rf.name, i)
			dfn := fmt.Sprintf("%s.%d", rf.name, i+1)
			if _, err := os.Stat(sfn); err == nil {
				os.Remove(dfn)
				os.Rename(sfn, dfn)
			}
		}

		dfn := rf.name + ".1"
		os.Remove(dfn)

		if _, err := os.Stat(rf.name); err == nil {
			os.Rename(rf.name, dfn)
		}
	}
	var err error
	rf.f, err = os.OpenFile(rf.name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	rf.curBytes = 0
	return err
}

func (rf *RotatingFile) shouldRollover(p []byte) (bool, error) {
	if rf.f == nil {
		var err error
		rf.f, err = os.OpenFile(rf.name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return false, err
		}
	}
	if rf.maxBytes > 0 {
		if !rf.curFirst {
			info, err := rf.f.Stat()
			if err != nil {
				return false, err
			}
			rf.curBytes = info.Size()
			rf.curFirst = true
		}

		if rf.curBytes+int64(len(p)) >= rf.maxBytes {
			return true, nil
		}
	}
	return false, nil
}

func (rf *RotatingFile) Write(p []byte) (n int, err error) {
	ok, err := rf.shouldRollover(p)
	if ok {
		err = rf.doRollover()
	}
	if err != nil {
		return 0, err
	}
	n, err = rf.f.Write(p)
	if err == nil {
		rf.curBytes += int64(n)
	}
	return n, err
}

// Close 关闭文件
func (rf *RotatingFile) Close() error {
	if rf.f != nil {
		return rf.f.Close()
	}
	return nil
}
