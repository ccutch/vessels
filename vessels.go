package vessels

import (
	"io"
	"io/ioutil"
)

type Vessel io.ReadWriteCloser

// Read to a Vessel we use the ioutil.ReadAll method
func Read(v Vessel) []byte {
	bs, err := ioutil.ReadAll(v)
	if err != nil {
		panic(err)
	}
	return bs
}

// Write to a Vessel we use the Vessel.Write method
func Write(v Vessel, bs []byte) {
	_, err := v.Write(bs)
	if err != nil {
		panic(err)
	}
}

// Close a Vessel we use the Vessel.Close method
func Close(v Vessel) {
	err := v.Close()
	if err != nil {
		panic(err)
	}
}
