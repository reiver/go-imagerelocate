package imagerelocate

import (
	"image"
	"image/color"
)

type internalImage struct {
	img image.Image
	x,y int
}

func (receiver internalImage) At(x, y int) color.Color {
	x -= receiver.x
	y -= receiver.y

	return receiver.img.At(x,y)
}

func (receiver internalImage) Bounds() image.Rectangle {
	bounds := receiveri.img.Bounds()

	bounds.Min.X += receiver.x
	bounds.Min.Y += receiver.y

	bounds.Max.X += receiver.x
	bounds.Max.Y += receiver.y

	return bounds
}

func (receiver internalImage) ColorModel() color.Model {
	return receiver.img.ColorModel()
}
