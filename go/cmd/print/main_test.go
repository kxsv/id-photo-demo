package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func Test_jfif_dpi(t *testing.T) {
	photo := image.NewRGBA(image.Rect(0, 0, 16, 24))
	photo.Set(0, 0, color.White)
	var encoded bytes.Buffer
	if err := jpeg.Encode(&encoded, photo, nil); err != nil {
		t.Fatal(err)
	}
	original := append([]byte(nil), encoded.Bytes()...)
	with_dpi := patchJfifDpi(original, 300)
	if len(with_dpi) != len(original)+18 {
		t.Fatal("JPEG without JFIF must gain one APP0 segment")
	}
	if !bytes.Equal(original, encoded.Bytes()) {
		t.Fatal("patch must not modify its input")
	}
	assert_dpi := func(data []byte, expected uint16) {
		t.Helper()
		if !bytes.Equal(data[2:11], []byte{0xFF, 0xE0, 0, 16, 'J', 'F', 'I', 'F', 0}) {
			t.Fatal("missing JFIF signature")
		}
		if data[11] != 1 || data[12] != 1 || data[13] != 1 {
			t.Fatal("incorrect JFIF version or density units")
		}
		if binary.BigEndian.Uint16(data[14:16]) != expected || binary.BigEndian.Uint16(data[16:18]) != expected {
			t.Fatal("incorrect horizontal or vertical DPI")
		}
		decoded, err := jpeg.Decode(bytes.NewReader(data))
		if err != nil || decoded.Bounds() != photo.Bounds() {
			t.Fatalf("JPEG changed or became unreadable: %v", err)
		}
	}
	assert_dpi(with_dpi, 300)
	updated := patchJfifDpi(with_dpi, 600)
	assert_dpi(updated, 600)
	if len(updated) != len(with_dpi) {
		t.Fatal("updating JFIF must not insert a duplicate segment")
	}
	assert_dpi(with_dpi, 300)
}

func Test_jfif_invalid_input(t *testing.T) {
	for _, data := range [][]byte{
		nil,
		{0xFF, 0xD8},
		[]byte("not a JPEG"),
		{0xFF, 0xD8, 0xFF, 0xE0, 0, 1},
		{0xFF, 0xD8, 0xFF, 0xE0, 0, 16, 'J'},
	} {
		if !bytes.Equal(data, patchJfifDpi(data, 300)) {
			t.Fatalf("invalid input was modified: %x", data)
		}
	}
}
