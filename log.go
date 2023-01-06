package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

type FLog struct {
	f *os.File
	*logrus.Logger
}

func NewLog(f string) (*FLog, error) {
	const flag = os.O_CREATE | os.O_APPEND | os.O_RDWR
	fd, err := os.OpenFile(f, flag, 0644)
	if err != nil {
		return nil, err
	}
	l := logrus.New()
	l.SetOutput(fd)
	flog := &FLog{
		f:      fd,
		Logger: l,
	}
	return flog, nil
}

func (l *FLog) Close() {
	if l.f != nil {
		err := l.f.Close()
		if err != nil {
			println(err)
		}
	}
}
