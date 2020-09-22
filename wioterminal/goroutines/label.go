package main

import (
	"image/color"

	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

type label struct {
	buf        []uint16
	w          int16
	h          int16
	fontHeight int16
}

func NewLabel(w, h int) *label {
	return &label{
		buf:        make([]uint16, w*h),
		w:          int16(w),
		h:          int16(h),
		fontHeight: int16(tinyfont.GetGlyph(&freemono.Regular9pt7b, '0').Height),
	}
}

func (l *label) Size() (int16, int16) {
	return l.w, l.h
}

func (l *label) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || y < 0 || l.w < x || l.h < y {
		return
	}
	l.buf[y*l.w+x] = ili9341.RGBATo565(c)
}

func (l *label) Display() error {
	return nil
}

func (l *label) SetText(str string, c color.RGBA) {
	for i := range l.buf {
		l.buf[i] = 0
	}

	tinyfont.WriteLine(l, &freemono.Regular9pt7b, 3, l.fontHeight, str, c)
}
