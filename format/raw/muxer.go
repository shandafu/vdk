package raw

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/deepch/vdk/codec/h265parser"

	"github.com/deepch/vdk/codec/h264parser"

	"github.com/deepch/vdk/av"
)

var startCode = []byte{0, 0, 0, 1}

type Muxer struct {
	idx int8
	w   *os.File
}

func NewMuxer(filePatch, fileName string) (*Muxer, error) {

	if _, err := os.Stat(filePatch); os.IsNotExist(err) {
		err := os.MkdirAll(filePatch, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	f2, err := os.Create(filePatch + fileName)
	if err != nil {
		return nil, err
	}

	return &Muxer{w: f2}, nil
}

func (element *Muxer) WriteHeader(streams []av.CodecData) (err error) {

	for i, stream := range streams {
		switch stream.Type() {
		case av.H264:
			_, err = element.w.Write(append(startCode, bytes.Join([][]byte{stream.(h264parser.CodecData).SPS(), stream.(h264parser.CodecData).PPS()}, startCode)...))
			element.idx = int8(i)
		case av.H265:
			_, err = element.w.Write(append(startCode, bytes.Join([][]byte{stream.(h265parser.CodecData).SPS(), stream.(h265parser.CodecData).PPS(), stream.(h265parser.CodecData).VPS()}, startCode)...))
			element.idx = int8(i)
		}
	}

	return

}

func (element *Muxer) WriteAvPacket(pkt *av.Packet) (err error) {

	if pkt.Idx == element.idx {
		// _, err = element.w.Write(startCode)
		if len(pkt.Data) < 5 {
			return
		}
		_, err = element.w.Write(Int64ToBytes(int64(len(pkt.Data) - 4)))
		if err != nil {
			return
		}
		_, err = element.w.Write(pkt.Data[4:])
	}

	return

}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func (element *Muxer) WriteRTPPacket(pkt *[]byte) (err error) {

	return

}

func (element *Muxer) Close() error {

	return element.w.Close()

}
