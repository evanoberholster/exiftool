package mknote

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/tiff"
)

// Errors
var (
	ErrNikonMkNote = errors.New("err makernote is not a Nikon makernote")
)

// NikonMkNoteHeader parses the Nikon Makernote from reader and returns byteOrder and error
// TODO: Exhaustatively Test and refactor
func NikonMkNoteHeader(reader io.Reader) (byteOrder binary.ByteOrder, err error) {
	// Nikon Makernotes header is 18 bytes. Move Reader up necessary bytes
	mknoteHeader := make([]byte, 18)
	if n, err := reader.Read(mknoteHeader); n < 18 || err != nil {
		return nil, ErrNikonMkNote
	}
	// Nikon makernote header starts with "Nikon" with the first 5 bytes
	if isNikonMkNoteHeaderBytes(mknoteHeader[:5]) {
		byteOrder := tiff.BinaryOrder(mknoteHeader[10:14])
		if byteOrder != nil {
			return byteOrder, nil
		}
		//if isTiffBigEndian(mknoteHeader[10:14]) {
		//	byteOrder = binary.BigEndian
		//	return byteOrder, nil
		//} else if isTiffLittleEndian(mknoteHeader[10:14]) {
		//	byteOrder = binary.LittleEndian
		//	return byteOrder, nil
		//}
	}

	return nil, ErrNikonMkNote
}

// Nikon Makernote Header
// isNikonMkNoteHeaderBytes represents "Nikon" the first 5 bytes of the
func isNikonMkNoteHeaderBytes(buf []byte) bool {
	return buf[0] == 'N' &&
		buf[1] == 'i' &&
		buf[2] == 'k' &&
		buf[3] == 'o' &&
		buf[4] == 'n'
	//nikonMkNoteHeaderBytes = []byte{0x4e, 0x69, 0x6b, 0x6f, 0x6e}
}