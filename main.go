package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"
	"strings"
)

// Errors
var (
	ErrInvalidArgumentCount = errors.New("invalid number of arguments")
)

// GetImage creates an image based on values in info
func GetImage(info ImageInfo) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, info.Width, info.Height))
	for x := 0; x < info.Width; x++ {
		for y := 0; y < info.Height; y++ {
			img.Set(x, y, info.BackgroundColor)
		}
	}
	return img
}

// ImageInfo groups basic image details
type ImageInfo struct {
	Height          int
	Width           int
	ForegroundColor color.Color
	BackgroundColor color.Color
	Text            string
}

// byteToHex converts byte in hex range to hex
func byteToHex(b byte) uint8 {
	if '0' <= b && b <= '9' {
		return b - '0'
	}
	if 'a' <= b && b <= 'f' {
		return 10 + b - 'a'
	}
	if 'A' <= b && b <= 'F' {
		return 10 + b - 'A'
	}
	return 0
}

// parseColor parses the hex color code to color type
func parseColor(code string, def color.Color) color.Color {
	c := color.RGBA{A: 0xff}
	switch len(code) {
	case 3:
		c.R = byteToHex(code[0])<<4 + byteToHex(code[0])
		c.G = byteToHex(code[1])<<4 + byteToHex(code[1])
		c.B = byteToHex(code[2])<<4 + byteToHex(code[2])
	case 6:
		c.R = byteToHex(code[0])<<4 + byteToHex(code[1])
		c.G = byteToHex(code[2])<<4 + byteToHex(code[3])
		c.B = byteToHex(code[4])<<4 + byteToHex(code[5])
	default:
		return def
	}
	return c
}

// parseDimensions parses WxH format to width and height
func parseDimensions(dim string) (int, int) {
	dParts := strings.Split(dim, "x")
	var width, height int
	switch len(dParts) {
	case 1:
		width, _ = strconv.Atoi(dParts[0])
		height, _ = strconv.Atoi(dParts[0])
	case 2:
		width, _ = strconv.Atoi(dParts[0])
		height, _ = strconv.Atoi(dParts[1])
	}
	return width, height
}

// Decode decodes the url path to image info
func (i *ImageInfo) Decode(path string) error {
	if path == "" {
		return ErrInvalidArgumentCount
	}
	pathParts := strings.Split(path, "/")
	switch len(pathParts) {
	case 1:
		i.BackgroundColor = color.White
		i.ForegroundColor = color.Black
		i.Width, i.Height = parseDimensions(pathParts[0])
	case 2:
		i.ForegroundColor = color.White
		i.BackgroundColor = parseColor(pathParts[1], color.Black)
		i.Width, i.Height = parseDimensions(pathParts[0])
	case 3:
		i.ForegroundColor = parseColor(pathParts[2], color.Black)
		i.BackgroundColor = parseColor(pathParts[1], color.White)
		i.Width, i.Height = parseDimensions(pathParts[0])
	default:
		return ErrInvalidArgumentCount
	}
	return nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var info ImageInfo
		if err := (&info).Decode(r.URL.Path[1:]); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		info.Text = r.FormValue("text")
		img := GetImage(info)
		png.Encode(w, img)
	})
	http.ListenAndServe(":8080", nil)
}
