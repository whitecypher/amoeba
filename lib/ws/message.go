package ws

import "bytes"

type Message interface {
	Bytes() []byte
}

type Messages []Message

func (m Messages) Bytes() []byte {
	r := [][]byte{}
	for _, i := range m {
		r = append(r, i.Bytes())
	}
	return bytes.Join(r, []byte("\n"))

}
