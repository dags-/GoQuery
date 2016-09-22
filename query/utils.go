package goquery

import (
	"bytes"
	"encoding/json"
	"image"
)

type Data map[string]interface{}

type Set map[string]bool

func (data Data) Put(key string, val interface{}) {
	if key != "" && val != "" && val != nil {
		data[key] = val
	}
}

func (data Data) Retain(keys ...string) Data {
	result := Data{}
	for _, k := range keys {
		result.Put(k, data[k])
	}
	return result
}

func (set Set) Add(val string) {
	set[val] = true
}

func (set Set) Contains(val string) bool {
	return len(set) == 0 || set[val]
}

func ToJson(input interface{}, pretty bool) string {
	var prefix, indent = "", ""
	if pretty {
		indent = "    "
	}
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent(prefix, indent)
	err := encoder.Encode(input)
	if err != nil {
		return "{}"
	}
	return string(buffer.Bytes())
}

func scaleImage(source image.Image, mask image.Rectangle, scale int) *image.RGBA {
	maskWidth := mask.Max.X - mask.Min.X
	maskHeight := mask.Max.Y - mask.Min.Y
	target := image.NewRGBA(image.Rect(0, 0, maskWidth * scale, maskHeight * scale))
	for x := 0; x < maskWidth; x++ {
		for y := 0; y < maskHeight; y++ {
			color := source.At(mask.Min.X + x, mask.Min.Y + y)
			offsetX, offsetY := x * scale, y * scale
			for toX := 0; toX < scale; toX++ {
				for toY := 0; toY < scale; toY++ {
					target.Set(offsetX + toX, offsetY + toY, color)
				}
			}
		}
	}
	return target
}