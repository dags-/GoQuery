package goquery

import (
	"bytes"
	"encoding/json"
	"image"
)

type Data map[string] interface{}

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

func drawScaledImage(target *image.RGBA, source image.Image, mask image.Rectangle) {
	maskWidth := mask.Max.X - mask.Min.X
	maskHeight := mask.Max.Y - mask.Min.Y
	scaleX := target.Bounds().Max.X / maskWidth
	scaleY := target.Bounds().Max.X / maskHeight
	for x := 0; x < maskWidth; x++ {
		for y := 0; y < maskHeight; y++ {
			color := source.At(mask.Min.X + x, mask.Min.Y + y)
			xPos, yPos := x * scaleX, y * scaleY
			for targetX := 0; targetX < scaleX; targetX++ {
				for targetY := 0; targetY < scaleY; targetY++ {
					target.Set(xPos + targetX, yPos + targetY, color)
				}
			}
		}
	}
}