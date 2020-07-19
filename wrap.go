package imagerelocate

import (
	"image"
)

// Wrap returns an image.Image that is just like ‘img’,
// except relocated to (‘x’, ‘y’).
func Wrap(x,y int, img image.Image) image.Image{
	return internalImage{
		x:x,
		y:y,
		img:img,
	}
}
