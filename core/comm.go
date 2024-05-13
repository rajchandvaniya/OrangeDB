package core

import "syscall"

type FDComm struct {
	Fd int
}

func (f FDComm) Read(buffer []byte) (int, error) {
	return syscall.Read(f.Fd, buffer)
}

func (f FDComm) Write(buffer []byte) (int, error) {
	return syscall.Write(f.Fd, buffer)
}
