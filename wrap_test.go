package imagerelocate_test

import (
	"github.com/reiver/go-imagerelocate"

	"github.com/reiver/go-imgstr"
	"github.com/reiver/go-palette2048"
	"github.com/reiver/go-pel"
	"github.com/reiver/go-sprite8x8"

	"fmt"
	"image"
	"math/rand"
	"time"

	"testing"
)

func TestWrap_pel(t *testing.T) {

	randomness := rand.New(rand.NewSource( time.Now().UTC().UnixNano() ))

	for testNumber:=0; testNumber<10; testNumber++ {

		// Randomly pick an (x,y)-coordinate.
		//
		// This will be the (x,y)-coordinate of the pel / pixel.
		var originalX, originalY int
		{
			originalX = randomness.Int()
			if 0 == randomness.Int()%2 {
				originalX *= -1
			}

			originalY = randomness.Int()
			if 0 == randomness.Int()%2 {
				originalY *= -1
			}

		}

		// Randomly pick an r,g,b,a for an rgba color.
		//
		// This will be the color of the pel / pixel.
		var r,g,b,a uint8
		{
			r = uint8(randomness.Intn(256))
			g = uint8(randomness.Intn(256))
			b = uint8(randomness.Intn(256))
			a = 255
		}

		// Create the pel / pixel.
		var pixel pel.RGBA = pel.RGBA{
			X:originalX,
			Y:originalY,

			R: r,
			G: g,
			B: b,
			A: a,
		}

		originalBounds := pixel.Bounds()

		// Even though this is a test for this package, we just want to douple check that
		// pel.RGBA is working correctly.
		//
		// (Just so we know if this test fails, it isn't because pel.RGBA isn't working correctly.)
		//
		// Here we check to make sure that the bounds on the pel.RGBA is what we expect.
		{
			if originalX < originalBounds.Min.X || originalBounds.Max.X <= originalX {
				t.Errorf("For test #%d, this should not happens — the original-x is outside of the original bounds.", testNumber)
				t.Logf("ORIGINAL-X: %d", originalX)
				t.Logf("ORIGINAL-X BOUNDS: [%d,%d)", originalBounds.Min.X, originalBounds.Max.X)
				continue
			}
			if originalY < originalBounds.Min.Y || originalBounds.Max.Y <= originalY {
				t.Errorf("For test #%d, this should not happens — the original-y is outside of the original bounds.", testNumber)
				t.Logf("ORIGINAL-Y: %d", originalY)
				t.Logf("ORIGINAL-Y BOUNDS: [%d,%d)", originalBounds.Min.Y, originalBounds.Max.Y)
				continue
			}

			expected := image.Rectangle{
				Min:image.Point{
					X: originalX,
					Y: originalY,
				},
				Max:image.Point{
					X: originalX+1,
					Y: originalY+1,
				},
			}

			if actual := originalBounds; expected != actual {
				t.Errorf("For test #%d, the actual bounds are not what was expected.", testNumber)
				t.Logf("EXPECTED: %#v", expected)
				t.Logf("ACTUAL:  %#v", actual)
				continue
			}
		}

		// Get the pixel colors from:
		//
		//
		//	(x-1, y-1) (x, , y-1) (x+1, y-1)
		//	(x-1, y  ) (x  , y  ) (x+1, y  )
		//	(x-1, y+1) (x  , y+1) (x+1, y+1)
		var originalSampleXYs[9][2]int
		{
			originalSampleXYs[0] = [2]int{originalX-1, originalY-1}
			originalSampleXYs[1] = [2]int{originalX  , originalY-1}
			originalSampleXYs[2] = [2]int{originalX+1, originalY-1}

			originalSampleXYs[3] = [2]int{originalX-1, originalY  }
			originalSampleXYs[4] = [2]int{originalX  , originalY  }
			originalSampleXYs[5] = [2]int{originalX+1, originalY  }

			originalSampleXYs[6] = [2]int{originalX-1, originalY+1}
			originalSampleXYs[7] = [2]int{originalX  , originalY+1}
			originalSampleXYs[8] = [2]int{originalX+1, originalY+1}
		}

		var originalSampleRGBAs [9][4]uint32
		for i, originalSampleXY := range originalSampleXYs {
				originalSampleX := originalSampleXY[0]
				originalSampleY := originalSampleXY[1]

				originalSampleR, originalSampleG, originalSampleB, originalSampleA := pixel.At(originalSampleX, originalSampleY).RGBA()

				originalSampleRGBAs[i] = [4]uint32{originalSampleR, originalSampleG, originalSampleB, originalSampleA}
		}

		// This is where we start to test our package.
		//
		// We use ‘imagerelocate.Wrap()’ to move ‘pel.RGBA’ by (x,y)=(‘dx’,‘dy’).
		//
		// So, for example, if ‘pel.RGBA’ was originally a single pel / pixel at
		// (x,y)=(10,20), then after wrapping it it should now be a single pel / pixel
		// at (x,y)=(11,21).
		//
		// But we need to confirm that both .Bounds(), and .At() reflect that.
		var dx,dy int = 1,1
		var img image.Image = imagerelocate.Wrap(dx,dy, pixel)

		newX := originalX + dx
		newY := originalY + dy

		// Confirm that .Bounds() reflects that ‘pel.RGBA’ has been relocated.
		{
			expected := image.Rectangle{
				Min:image.Point{
					X:newX,
					Y:newY,
				},
				Max:image.Point{
					X:newX+1,
					Y:newY+1,
				},
			}

			actual := img.Bounds()

			if expected != actual {
				t.Errorf("For test #%d, the actual relocated bounds is not what was expected.", testNumber)
				t.Logf("EXPECTED: %#v", expected)
				t.Logf("ACTUAL:   %#v", actual)
				continue
			}
		}

		// Get the relocated pixel colors from:
		//
		//
		//	(x-1, y-1) (x, , y-1) (x+1, y-1)
		//	(x-1, y  ) (x  , y  ) (x+1, y  )
		//	(x-1, y+1) (x  , y+1) (x+1, y+1)
		//
		// Where ‘x’ & ‘y’ here are the relocated (x,y),
		// not the original (x,y).
		var newSampleXYs[9][2]int
		{
			newSampleXYs[0] = [2]int{newX-1, newY-1}
			newSampleXYs[1] = [2]int{newX  , newY-1}
			newSampleXYs[2] = [2]int{newX+1, newY-1}

			newSampleXYs[3] = [2]int{newX-1, newY  }
			newSampleXYs[4] = [2]int{newX  , newY  }
			newSampleXYs[5] = [2]int{newX+1, newY  }

			newSampleXYs[6] = [2]int{newX-1, newY+1}
			newSampleXYs[7] = [2]int{newX  , newY+1}
			newSampleXYs[8] = [2]int{newX+1, newY+1}
		}

		var newSampleRGBAs [9][4]uint32
		for i, newSampleXY := range newSampleXYs {
				newSampleX := newSampleXY[0]
				newSampleY := newSampleXY[1]

				newSampleR, newSampleG, newSampleB, newSampleA := img.At(newSampleX, newSampleY).RGBA()

				newSampleRGBAs[i] = [4]uint32{newSampleR, newSampleG, newSampleB, newSampleA}
		}

		// Confirm the colors!
		if expected, actual := originalSampleRGBAs, newSampleRGBAs; expected != actual {
			t.Errorf("For test #%d, the actual colors at & near the (x,y)s was not what was expected.", testNumber)
			t.Logf("ORIGINAL (x,y) = (%d,%d)", originalX, originalY)
			t.Logf("       (dx,dy) = (%d,%d)", dx,dy)
			t.Logf("NEW      (x,y) = (%d,%d)", newX, newY)
			t.Logf("rgba(%d,%d,%d,%d)", r,g,b,a)
			t.Log("EXPECTED:")
			for i, originalSampleXY := range originalSampleXYs {
				originalSampleRGBA := actual[i]
				t.Logf("(%d,%d)->(%d,%d,%d,%d)", originalSampleXY[0], originalSampleXY[1], originalSampleRGBA[0], originalSampleRGBA[1], originalSampleRGBA[2], originalSampleRGBA[3])
			}
			t.Logf("ACTUAL:")
			for i, newSampleXY := range newSampleXYs {
				newSampleRGBA := actual[i]
				t.Logf("(%d,%d)->(%d,%d,%d,%d)", newSampleXY[0], newSampleXY[1], newSampleRGBA[0], newSampleRGBA[1], newSampleRGBA[2], newSampleRGBA[3])
			}
			continue
		}
	}

}

func TestWrap_sprite8x8(t *testing.T) {


	randomness := rand.New(rand.NewSource( time.Now().UTC().UnixNano() ))

	Loop: for testNumber:=0; testNumber<10; testNumber++ {

		const black = 0 // Index into the color palette.
		const blackR =   1
		const blackG =   2
		const blackB =   3
		const blackA = 255

		const green = 2 // Index into the color palette.
		const greenR =  57
		const greenG = 181
		const greenB =  74
		const greenA = 255

		const yellow = 3 // Index into the color palette.
		const yellowR = 255
		const yellowG = 199
		const yellowB =   6
		const yellowA = 255

		const blue = 4 // Index into the color palette.
		const blueR =   0
		const blueG = 111
		const blueB = 184
		const blueA = 255

		const white = 15 // Index into the color palette.
		const whiteR = 250
		const whiteG = 251
		const whiteB = 252
		const whiteA = 255


		// Set the colors in the palette.
		var palettedBuffer [palette2048.ByteSize]uint8

		var palette palette2048.Slice = palette2048.Slice(palettedBuffer[:])

		palette.SetColorRGBA(black,  blackR,  blackG,  blackB,  blackA)
		palette.SetColorRGBA(green,  greenR,  greenG,  greenB,  greenA)
		palette.SetColorRGBA(yellow, yellowR, yellowG, yellowB, yellowA)
		palette.SetColorRGBA(blue,   blueR,   blueG,   blueB,   blueA)
		palette.SetColorRGBA(white,  whiteR,  whiteG,  whiteB,  whiteA)

		// This is the (x,y)-coordinate of the top-left corner of the sprite.
		//
		// These sprites start off having their top-left corner at (x,y)=(0,0)
		var originalX, originalY int = 0,0

		// This is the sprite.
		var pix [8*8]uint8 = [8*8]uint8{
			white,  black,  black,  black,  black,  black,  black, yellow,
			black, yellow,  black,  black,  black,  black,  blue,   black,
			black,  black,   blue,  black,  black,  green,  black,  black,
			black,  black,  black,  green,  white,  black,  black,  black,
			black,  black,  black,   blue, yellow,  black,  black,  black,
			black,  black, yellow,  black,  black,  white,  black,  black,
			black,  white,  black,  black,  black,  black,  green,  black,
			green,  black,  black,  black,  black,  black,  black,  blue,
		}

		category := fmt.Sprintf("should not matter - category %d", randomness.Int())
		id       := uint8(randomness.Intn(256))

		var sprite sprite8x8.Paletted = sprite8x8.Paletted{
			Pix: pix[:],
			Palette: palette,
			Category: category,
			ID:       id,
		}

		{
			expected := "IMAGE:iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIEAIAAAAb/fWfAAAAcklEQVR4nGL59ev37z9/2MGAAS/49+/oUWZmBkZGJiZm5v//jx9nY4OwMUkGhvz8HTugbFwSENLScutWLy8UNZjmQRRBnIopy4Lp1uPHd+3Ky2NlFRLi4sLiFWTdmD7BtIcBlwQunzBgegt/WAECAAD//wHnOPHocYAXAAAAAElFTkSuQmCC"
			actual := imgstr.ImageString(sprite)

			if expected != actual {
				t.Errorf("For test #%d, the actual image for the sprite is not what was expected.", testNumber)
				t.Logf("EXPECTED: %q", expected)
				t.Logf("ACTUAL:   %q", actual)
				continue
			}
		}

		originalBounds := sprite.Bounds()

		// Even though this is a test for this package, we just want to douple check that
		// sprite8x8.Paletted is working correctly.
		//
		// (Just so we know if this test fails, it isn't because sprite8x8.Paletted isn't working correctly.)
		//
		// Here we check to make sure that the bounds on the pel.RGBA is what we expect.
		{
			if originalX < originalBounds.Min.X || originalBounds.Max.X <= originalX {
				t.Errorf("For test #%d, this should not happens — the original-x is outside of the original bounds.", testNumber)
				t.Logf("ORIGINAL-X: %d", originalX)
				t.Logf("ORIGINAL-X BOUNDS: [%d,%d)", originalBounds.Min.X, originalBounds.Max.X)
				continue
			}
			if originalY < originalBounds.Min.Y || originalBounds.Max.Y <= originalY {
				t.Errorf("For test #%d, this should not happens — the original-y is outside of the original bounds.", testNumber)
				t.Logf("ORIGINAL-Y: %d", originalY)
				t.Logf("ORIGINAL-Y BOUNDS: [%d,%d)", originalBounds.Min.Y, originalBounds.Max.Y)
				continue
			}

			expected := image.Rectangle{
				Min:image.Point{
					X: originalX,
					Y: originalY,
				},
				Max:image.Point{
					X: originalX+8,
					Y: originalY+8,
				},
			}

			if actual := originalBounds; expected != actual {
				t.Errorf("For test #%d, the actual bounds are not what was expected.", testNumber)
				t.Logf("EXPECTED: %#v", expected)
				t.Logf("ACTUAL:  %#v", actual)
				continue
			}
		}

		// Relocate the sprite, and then relocate it back. The result should that .At() and .Bounds() act just like the original.
		for subTestNumber:=0; subTestNumber<10; subTestNumber++ {
			var dx, dy int = randomness.Intn(200), randomness.Intn(200)
			if 0 == randomness.Int()%2 {
				dx *= -1
			}
			if 0 == randomness.Int()%2 {
				dy *= -1
			}


			var img image.Image = imagerelocate.Wrap(
						-dx, -dy,
						imagerelocate.Wrap(
							dx, dy,
							sprite,
						),
			)

			if expected, actual := sprite.Bounds(), img.Bounds(); expected != actual {
				t.Errorf("For test #%d:%d, the actual bounds is not what was expected.", testNumber, subTestNumber)
				t.Logf("EXPECTED: %#v", expected)
				t.Logf("ACTUAL:   %#v", actual)
				continue
			}

			{
				bounds := sprite.Bounds()

				minY := bounds.Min.Y
				maxY := bounds.Max.Y

				minX := bounds.Min.X
				maxX := bounds.Max.X

				for y:=minY; y<maxY; y++ {
					for x:=minX; x<maxX; x++ {
						eR, eG, eB, eA := sprite.At(x,y).RGBA()
						aR, aG, aB, aA := img.At(x,y).RGBA()

						if eR != aR || eG != aG || eB != aB || eA != aA {
							t.Errorf("For test #%d:%d, the actual color at (x,y)=(%d,%d) is not what was expected.", testNumber, subTestNumber, x,y)
							t.Logf("EXPECTED: (r,g,b,a)=(%d,%d,%d,%d)", eR, eG, eB, eA)
							t.Logf("ACTUAL:   (t,g,b,a)=(%d,%d,%d,%d)", aR, aG, aB, aA)
							continue Loop
						}
					}
				}
			}
		}


		// Get the pixel colors from:
		//
		//	(x-1, y-1) (x, , y-1) (x+1, y-1) (x+2, y-1) (x+3, y-1) (x+4, y-1) (x+5, y-1) (x+6, y-1) (x+7, y-1) (x+8, y-1)
		//	(x-1, y  ) (x  , y  ) (x+1, y  ) (x+2, y  ) (x+3, y  ) (x+4, y  ) (x+5, y  ) (x+6, y  ) (x+7, y  ) (x+8, y  )
		//	(x-1, y+1) (x  , y+1) (x+1, y+1) (x+2, y+1) (x+3, y+1) (x+4, y+1) (x+5, y+1) (x+6, y+1) (x+7, y+1) (x+8, y+1)
		//	(x-1, y+2) (x  , y+2) (x+1, y+2) (x+2, y+2) (x+3, y+2) (x+4, y+2) (x+5, y+2) (x+6, y+2) (x+7, y+2) (x+8, y+2)
		//	(x-1, y+3) (x  , y+3) (x+1, y+3) (x+2, y+3) (x+3, y+3) (x+4, y+3) (x+5, y+3) (x+6, y+3) (x+7, y+3) (x+8, y+3)
		//	(x-1, y+4) (x  , y+4) (x+1, y+4) (x+2, y+4) (x+3, y+4) (x+4, y+4) (x+5, y+4) (x+6, y+4) (x+7, y+4) (x+8, y+4)
		//	(x-1, y+5) (x  , y+5) (x+1, y+5) (x+2, y+5) (x+3, y+5) (x+4, y+5) (x+5, y+5) (x+6, y+5) (x+7, y+5) (x+8, y+5)
		//	(x-1, y+6) (x  , y+6) (x+1, y+6) (x+2, y+6) (x+3, y+6) (x+4, y+6) (x+5, y+6) (x+6, y+6) (x+7, y+6) (x+8, y+6)
		//	(x-1, y+7) (x  , y+7) (x+1, y+7) (x+2, y+7) (x+3, y+7) (x+4, y+7) (x+5, y+7) (x+6, y+7) (x+7, y+7) (x+8, y+7)
		//	(x-1, y+8) (x  , y+8) (x+1, y+8) (x+2, y+8) (x+3, y+8) (x+4, y+8) (x+5, y+8) (x+6, y+8) (x+7, y+8) (x+8, y+8)
		//
		// I.e., we get the pixels in the sprite, plus a 1-pixel-thick border.
		var originalSampleXYs[10*10][2]int
		{
			originalSampleXYs[ 0] = [2]int{originalX-1, originalY-1}
			originalSampleXYs[ 1] = [2]int{originalX  , originalY-1}
			originalSampleXYs[ 2] = [2]int{originalX+1, originalY-1}
			originalSampleXYs[ 3] = [2]int{originalX+2, originalY-1}
			originalSampleXYs[ 4] = [2]int{originalX+3, originalY-1}
			originalSampleXYs[ 5] = [2]int{originalX+4, originalY-1}
			originalSampleXYs[ 6] = [2]int{originalX+5, originalY-1}
			originalSampleXYs[ 7] = [2]int{originalX+6, originalY-1}
			originalSampleXYs[ 8] = [2]int{originalX+7, originalY-1}
			originalSampleXYs[ 9] = [2]int{originalX+8, originalY-1}

			originalSampleXYs[10] = [2]int{originalX-1, originalY  }
			originalSampleXYs[11] = [2]int{originalX  , originalY  }
			originalSampleXYs[12] = [2]int{originalX+1, originalY  }
			originalSampleXYs[13] = [2]int{originalX+2, originalY  }
			originalSampleXYs[14] = [2]int{originalX+3, originalY  }
			originalSampleXYs[15] = [2]int{originalX+4, originalY  }
			originalSampleXYs[16] = [2]int{originalX+5, originalY  }
			originalSampleXYs[17] = [2]int{originalX+6, originalY  }
			originalSampleXYs[18] = [2]int{originalX+7, originalY  }
			originalSampleXYs[19] = [2]int{originalX+8, originalY  }

			originalSampleXYs[20] = [2]int{originalX-1, originalY+1}
			originalSampleXYs[21] = [2]int{originalX  , originalY+1}
			originalSampleXYs[22] = [2]int{originalX+1, originalY+1}
			originalSampleXYs[23] = [2]int{originalX+2, originalY+1}
			originalSampleXYs[24] = [2]int{originalX+3, originalY+1}
			originalSampleXYs[25] = [2]int{originalX+4, originalY+1}
			originalSampleXYs[26] = [2]int{originalX+5, originalY+1}
			originalSampleXYs[27] = [2]int{originalX+6, originalY+1}
			originalSampleXYs[28] = [2]int{originalX+7, originalY+1}
			originalSampleXYs[29] = [2]int{originalX+8, originalY+1}

			originalSampleXYs[30] = [2]int{originalX-1, originalY+2}
			originalSampleXYs[31] = [2]int{originalX  , originalY+2}
			originalSampleXYs[32] = [2]int{originalX+1, originalY+2}
			originalSampleXYs[33] = [2]int{originalX+2, originalY+2}
			originalSampleXYs[34] = [2]int{originalX+3, originalY+2}
			originalSampleXYs[35] = [2]int{originalX+4, originalY+2}
			originalSampleXYs[36] = [2]int{originalX+5, originalY+2}
			originalSampleXYs[37] = [2]int{originalX+6, originalY+2}
			originalSampleXYs[38] = [2]int{originalX+7, originalY+2}
			originalSampleXYs[39] = [2]int{originalX+8, originalY+2}

			originalSampleXYs[40] = [2]int{originalX-1, originalY+3}
			originalSampleXYs[41] = [2]int{originalX  , originalY+3}
			originalSampleXYs[42] = [2]int{originalX+1, originalY+3}
			originalSampleXYs[43] = [2]int{originalX+2, originalY+3}
			originalSampleXYs[44] = [2]int{originalX+3, originalY+3}
			originalSampleXYs[45] = [2]int{originalX+4, originalY+3}
			originalSampleXYs[46] = [2]int{originalX+5, originalY+3}
			originalSampleXYs[47] = [2]int{originalX+6, originalY+3}
			originalSampleXYs[48] = [2]int{originalX+7, originalY+3}
			originalSampleXYs[49] = [2]int{originalX+8, originalY+3}

			originalSampleXYs[50] = [2]int{originalX-1, originalY+4}
			originalSampleXYs[51] = [2]int{originalX  , originalY+4}
			originalSampleXYs[52] = [2]int{originalX+1, originalY+4}
			originalSampleXYs[53] = [2]int{originalX+2, originalY+4}
			originalSampleXYs[54] = [2]int{originalX+3, originalY+4}
			originalSampleXYs[55] = [2]int{originalX+4, originalY+4}
			originalSampleXYs[56] = [2]int{originalX+5, originalY+4}
			originalSampleXYs[57] = [2]int{originalX+6, originalY+4}
			originalSampleXYs[58] = [2]int{originalX+7, originalY+4}
			originalSampleXYs[59] = [2]int{originalX+8, originalY+4}

			originalSampleXYs[60] = [2]int{originalX-1, originalY+5}
			originalSampleXYs[61] = [2]int{originalX  , originalY+5}
			originalSampleXYs[62] = [2]int{originalX+1, originalY+5}
			originalSampleXYs[63] = [2]int{originalX+2, originalY+5}
			originalSampleXYs[64] = [2]int{originalX+3, originalY+5}
			originalSampleXYs[65] = [2]int{originalX+4, originalY+5}
			originalSampleXYs[66] = [2]int{originalX+5, originalY+5}
			originalSampleXYs[67] = [2]int{originalX+6, originalY+5}
			originalSampleXYs[68] = [2]int{originalX+7, originalY+5}
			originalSampleXYs[69] = [2]int{originalX+8, originalY+5}

			originalSampleXYs[70] = [2]int{originalX-1, originalY+6}
			originalSampleXYs[71] = [2]int{originalX  , originalY+6}
			originalSampleXYs[72] = [2]int{originalX+1, originalY+6}
			originalSampleXYs[73] = [2]int{originalX+2, originalY+6}
			originalSampleXYs[74] = [2]int{originalX+3, originalY+6}
			originalSampleXYs[75] = [2]int{originalX+4, originalY+6}
			originalSampleXYs[76] = [2]int{originalX+5, originalY+6}
			originalSampleXYs[77] = [2]int{originalX+6, originalY+6}
			originalSampleXYs[78] = [2]int{originalX+7, originalY+6}
			originalSampleXYs[79] = [2]int{originalX+8, originalY+6}

			originalSampleXYs[80] = [2]int{originalX-1, originalY+7}
			originalSampleXYs[81] = [2]int{originalX  , originalY+7}
			originalSampleXYs[82] = [2]int{originalX+1, originalY+7}
			originalSampleXYs[83] = [2]int{originalX+2, originalY+7}
			originalSampleXYs[84] = [2]int{originalX+3, originalY+7}
			originalSampleXYs[85] = [2]int{originalX+4, originalY+7}
			originalSampleXYs[86] = [2]int{originalX+5, originalY+7}
			originalSampleXYs[87] = [2]int{originalX+6, originalY+7}
			originalSampleXYs[88] = [2]int{originalX+7, originalY+7}
			originalSampleXYs[89] = [2]int{originalX+8, originalY+7}

			originalSampleXYs[90] = [2]int{originalX-1, originalY+8}
			originalSampleXYs[91] = [2]int{originalX  , originalY+8}
			originalSampleXYs[92] = [2]int{originalX+1, originalY+8}
			originalSampleXYs[93] = [2]int{originalX+2, originalY+8}
			originalSampleXYs[94] = [2]int{originalX+3, originalY+8}
			originalSampleXYs[95] = [2]int{originalX+4, originalY+8}
			originalSampleXYs[96] = [2]int{originalX+5, originalY+8}
			originalSampleXYs[97] = [2]int{originalX+6, originalY+8}
			originalSampleXYs[98] = [2]int{originalX+7, originalY+8}
			originalSampleXYs[99] = [2]int{originalX+8, originalY+8}
		}

		var originalSampleRGBAs [10*10][4]uint32
		for i, originalSampleXY := range originalSampleXYs {
				originalSampleX := originalSampleXY[0]
				originalSampleY := originalSampleXY[1]

				originalSampleR, originalSampleG, originalSampleB, originalSampleA := sprite.At(originalSampleX, originalSampleY).RGBA()

				originalSampleRGBAs[i] = [4]uint32{originalSampleR, originalSampleG, originalSampleB, originalSampleA}
		}

		// This is where we start to test our package on its own.
		//
		// We use ‘imagerelocate.Wrap()’ to move ‘sprite8x8.Paletted’ by (x,y)=(‘dx’,‘dy’).
		//
		// We need to confirm that both .Bounds(), and .At() reflect that.
		var dx,dy int = 7,7
		var img image.Image = imagerelocate.Wrap(dx,dy, sprite)

		{
			expected := "IMAGE:iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIEAIAAAAb/fWfAAAAcklEQVR4nGL59ev37z9/2MGAAS/49+/oUWZmBkZGJiZm5v//jx9nY4OwMUkGhvz8HTugbFwSENLScutWLy8UNZjmQRRBnIopy4Lp1uPHd+3Ky2NlFRLi4sLiFWTdmD7BtIcBlwQunzBgegt/WAECAAD//wHnOPHocYAXAAAAAElFTkSuQmCC"
			actual := imgstr.ImageString(img)

			if expected != actual {
				t.Errorf("For test #%d, the actual image for the sprite is not what was expected.", testNumber)
				t.Logf("EXPECTED: %q", expected)
				t.Logf("ACTUAL:   %q", actual)
				continue
			}
		}

		newX := originalX + dx
		newY := originalY + dy

		// Confirm that .Bounds() reflects that ‘sprite8x8.Paletted’ has been relocated.
		{
			expected := image.Rectangle{
				Min:image.Point{
					X:newX,
					Y:newY,
				},
				Max:image.Point{
					X:newX+8,
					Y:newY+8,
				},
			}

			actual := img.Bounds()

			if expected != actual {
				t.Errorf("For test #%d, the actual relocated bounds is not what was expected.", testNumber)
				t.Logf("EXPECTED: %#v", expected)
				t.Logf("ACTUAL:   %#v", actual)
				continue
			}
		}

		// Get the relocated pixel colors from:
		//
		//	(x-1, y-1) (x, , y-1) (x+1, y-1) (x+2, y-1) (x+3, y-1) (x+4, y-1) (x+5, y-1) (x+6, y-1) (x+7, y-1) (x+8, y-1)
		//	(x-1, y  ) (x  , y  ) (x+1, y  ) (x+2, y  ) (x+3, y  ) (x+4, y  ) (x+5, y  ) (x+6, y  ) (x+7, y  ) (x+8, y  )
		//	(x-1, y+1) (x  , y+1) (x+1, y+1) (x+2, y+1) (x+3, y+1) (x+4, y+1) (x+5, y+1) (x+6, y+1) (x+7, y+1) (x+8, y+1)
		//	(x-1, y+2) (x  , y+2) (x+1, y+2) (x+2, y+2) (x+3, y+2) (x+4, y+2) (x+5, y+2) (x+6, y+2) (x+7, y+2) (x+8, y+2)
		//	(x-1, y+3) (x  , y+3) (x+1, y+3) (x+2, y+3) (x+3, y+3) (x+4, y+3) (x+5, y+3) (x+6, y+3) (x+7, y+3) (x+8, y+3)
		//	(x-1, y+4) (x  , y+4) (x+1, y+4) (x+2, y+4) (x+3, y+4) (x+4, y+4) (x+5, y+4) (x+6, y+4) (x+7, y+4) (x+8, y+4)
		//	(x-1, y+5) (x  , y+5) (x+1, y+5) (x+2, y+5) (x+3, y+5) (x+4, y+5) (x+5, y+5) (x+6, y+5) (x+7, y+5) (x+8, y+5)
		//	(x-1, y+6) (x  , y+6) (x+1, y+6) (x+2, y+6) (x+3, y+6) (x+4, y+6) (x+5, y+6) (x+6, y+6) (x+7, y+6) (x+8, y+6)
		//	(x-1, y+7) (x  , y+7) (x+1, y+7) (x+2, y+7) (x+3, y+7) (x+4, y+7) (x+5, y+7) (x+6, y+7) (x+7, y+7) (x+8, y+7)
		//	(x-1, y+8) (x  , y+8) (x+1, y+8) (x+2, y+8) (x+3, y+8) (x+4, y+8) (x+5, y+8) (x+6, y+8) (x+7, y+8) (x+8, y+8)
		//
		// Where ‘x’ & ‘y’ here are the relocated (x,y),
		// not the original (x,y).
		var newSampleXYs[10*10][2]int
		{
			newSampleXYs[ 0] = [2]int{newX-1, newY-1}
			newSampleXYs[ 1] = [2]int{newX  , newY-1}
			newSampleXYs[ 2] = [2]int{newX+1, newY-1}
			newSampleXYs[ 3] = [2]int{newX+2, newY-1}
			newSampleXYs[ 4] = [2]int{newX+3, newY-1}
			newSampleXYs[ 5] = [2]int{newX+4, newY-1}
			newSampleXYs[ 6] = [2]int{newX+5, newY-1}
			newSampleXYs[ 7] = [2]int{newX+6, newY-1}
			newSampleXYs[ 8] = [2]int{newX+7, newY-1}
			newSampleXYs[ 9] = [2]int{newX+8, newY-1}

			newSampleXYs[10] = [2]int{newX-1, newY  }
			newSampleXYs[11] = [2]int{newX  , newY  }
			newSampleXYs[12] = [2]int{newX+1, newY  }
			newSampleXYs[13] = [2]int{newX+2, newY  }
			newSampleXYs[14] = [2]int{newX+3, newY  }
			newSampleXYs[15] = [2]int{newX+4, newY  }
			newSampleXYs[16] = [2]int{newX+5, newY  }
			newSampleXYs[17] = [2]int{newX+6, newY  }
			newSampleXYs[18] = [2]int{newX+7, newY  }
			newSampleXYs[19] = [2]int{newX+8, newY  }

			newSampleXYs[20] = [2]int{newX-1, newY+1}
			newSampleXYs[21] = [2]int{newX  , newY+1}
			newSampleXYs[22] = [2]int{newX+1, newY+1}
			newSampleXYs[23] = [2]int{newX+2, newY+1}
			newSampleXYs[24] = [2]int{newX+3, newY+1}
			newSampleXYs[25] = [2]int{newX+4, newY+1}
			newSampleXYs[26] = [2]int{newX+5, newY+1}
			newSampleXYs[27] = [2]int{newX+6, newY+1}
			newSampleXYs[28] = [2]int{newX+7, newY+1}
			newSampleXYs[29] = [2]int{newX+8, newY+1}

			newSampleXYs[30] = [2]int{newX-1, newY+2}
			newSampleXYs[31] = [2]int{newX  , newY+2}
			newSampleXYs[32] = [2]int{newX+1, newY+2}
			newSampleXYs[33] = [2]int{newX+2, newY+2}
			newSampleXYs[34] = [2]int{newX+3, newY+2}
			newSampleXYs[35] = [2]int{newX+4, newY+2}
			newSampleXYs[36] = [2]int{newX+5, newY+2}
			newSampleXYs[37] = [2]int{newX+6, newY+2}
			newSampleXYs[38] = [2]int{newX+7, newY+2}
			newSampleXYs[39] = [2]int{newX+8, newY+2}

			newSampleXYs[40] = [2]int{newX-1, newY+3}
			newSampleXYs[41] = [2]int{newX  , newY+3}
			newSampleXYs[42] = [2]int{newX+1, newY+3}
			newSampleXYs[43] = [2]int{newX+2, newY+3}
			newSampleXYs[44] = [2]int{newX+3, newY+3}
			newSampleXYs[45] = [2]int{newX+4, newY+3}
			newSampleXYs[46] = [2]int{newX+5, newY+3}
			newSampleXYs[47] = [2]int{newX+6, newY+3}
			newSampleXYs[48] = [2]int{newX+7, newY+3}
			newSampleXYs[49] = [2]int{newX+8, newY+3}

			newSampleXYs[50] = [2]int{newX-1, newY+4}
			newSampleXYs[51] = [2]int{newX  , newY+4}
			newSampleXYs[52] = [2]int{newX+1, newY+4}
			newSampleXYs[53] = [2]int{newX+2, newY+4}
			newSampleXYs[54] = [2]int{newX+3, newY+4}
			newSampleXYs[55] = [2]int{newX+4, newY+4}
			newSampleXYs[56] = [2]int{newX+5, newY+4}
			newSampleXYs[57] = [2]int{newX+6, newY+4}
			newSampleXYs[58] = [2]int{newX+7, newY+4}
			newSampleXYs[59] = [2]int{newX+8, newY+4}

			newSampleXYs[60] = [2]int{newX-1, newY+5}
			newSampleXYs[61] = [2]int{newX  , newY+5}
			newSampleXYs[62] = [2]int{newX+1, newY+5}
			newSampleXYs[63] = [2]int{newX+2, newY+5}
			newSampleXYs[64] = [2]int{newX+3, newY+5}
			newSampleXYs[65] = [2]int{newX+4, newY+5}
			newSampleXYs[66] = [2]int{newX+5, newY+5}
			newSampleXYs[67] = [2]int{newX+6, newY+5}
			newSampleXYs[68] = [2]int{newX+7, newY+5}
			newSampleXYs[69] = [2]int{newX+8, newY+5}

			newSampleXYs[70] = [2]int{newX-1, newY+6}
			newSampleXYs[71] = [2]int{newX  , newY+6}
			newSampleXYs[72] = [2]int{newX+1, newY+6}
			newSampleXYs[73] = [2]int{newX+2, newY+6}
			newSampleXYs[74] = [2]int{newX+3, newY+6}
			newSampleXYs[75] = [2]int{newX+4, newY+6}
			newSampleXYs[76] = [2]int{newX+5, newY+6}
			newSampleXYs[77] = [2]int{newX+6, newY+6}
			newSampleXYs[78] = [2]int{newX+7, newY+6}
			newSampleXYs[79] = [2]int{newX+8, newY+6}

			newSampleXYs[80] = [2]int{newX-1, newY+7}
			newSampleXYs[81] = [2]int{newX  , newY+7}
			newSampleXYs[82] = [2]int{newX+1, newY+7}
			newSampleXYs[83] = [2]int{newX+2, newY+7}
			newSampleXYs[84] = [2]int{newX+3, newY+7}
			newSampleXYs[85] = [2]int{newX+4, newY+7}
			newSampleXYs[86] = [2]int{newX+5, newY+7}
			newSampleXYs[87] = [2]int{newX+6, newY+7}
			newSampleXYs[88] = [2]int{newX+7, newY+7}
			newSampleXYs[89] = [2]int{newX+8, newY+7}

			newSampleXYs[90] = [2]int{newX-1, newY+8}
			newSampleXYs[91] = [2]int{newX  , newY+8}
			newSampleXYs[92] = [2]int{newX+1, newY+8}
			newSampleXYs[93] = [2]int{newX+2, newY+8}
			newSampleXYs[94] = [2]int{newX+3, newY+8}
			newSampleXYs[95] = [2]int{newX+4, newY+8}
			newSampleXYs[96] = [2]int{newX+5, newY+8}
			newSampleXYs[97] = [2]int{newX+6, newY+8}
			newSampleXYs[98] = [2]int{newX+7, newY+8}
			newSampleXYs[99] = [2]int{newX+8, newY+8}
		}

		// Make sure the newSampleXYs are what we expect.
		for i, originalSampleXY := range originalSampleXYs {
			newSampleXY := newSampleXYs[i]

			expectedX := originalSampleXY[0] + dx
			expectedY := originalSampleXY[1] + dy

			actualX := newSampleXY[0]
			actualY := newSampleXY[1]

			if expectedX != actualX || expectedY != actualY {
				t.Errorf("For test #%d, the actual (x,y) sample coordinate is not what was expected.", testNumber)
				t.Logf("EXPECTED (x,y)=(%d,%d)", expectedX, expectedY)
				t.Logf("ACTUAL   (x,y)=(%d,%d)", actualX,   actualY)
				continue
			}
		}


		var newSampleRGBAs [10*10][4]uint32
		for i, newSampleXY := range newSampleXYs {
				newSampleX := newSampleXY[0]
				newSampleY := newSampleXY[1]

				newSampleR, newSampleG, newSampleB, newSampleA := img.At(newSampleX, newSampleY).RGBA()

				newSampleRGBAs[i] = [4]uint32{newSampleR, newSampleG, newSampleB, newSampleA}
		}

		// Confirm the colors!
		if expected, actual := originalSampleRGBAs, newSampleRGBAs; expected != actual {
			t.Errorf("For test #%d, the actual colors at & near the (x,y)s was not what was expected.", testNumber)
			t.Logf("ORIGINAL (x,y) = (%d,%d)", originalX, originalY)
			t.Logf("       (dx,dy) = (%d,%d)", dx,dy)
			t.Logf("NEW      (x,y) = (%d,%d)", newX, newY)
			t.Logf(" black rgba(%d,%d,%d,%d) -> (%d,%d,%d,%d)",  blackR,  blackG,  blackB,  blackA,   uint32( blackR)*(0xffff/0xff), uint32( blackG)*(0xffff/0xff), uint32( blackB)*(0xffff/0xff), uint32( blackA)*(0xffff/0xff))
			t.Logf(" green rgba(%d,%d,%d,%d) -> (%d,%d,%d,%d)",  greenR,  greenG,  greenB,  greenA,   uint32( greenR)*(0xffff/0xff), uint32( greenG)*(0xffff/0xff), uint32( greenB)*(0xffff/0xff), uint32( greenA)*(0xffff/0xff))
			t.Logf("yellow rgba(%d,%d,%d,%d) -> (%d,%d,%d,%d)", yellowR, yellowG, yellowB, yellowA,   uint32(yellowR)*(0xffff/0xff), uint32(yellowG)*(0xffff/0xff), uint32(yellowB)*(0xffff/0xff), uint32(yellowA)*(0xffff/0xff))
			t.Logf("  blue rgba(%d,%d,%d,%d) -> (%d,%d,%d,%d)",   blueR,   blueG,   blueB,   blueA,   uint32(  blueR)*(0xffff/0xff), uint32(  blueG)*(0xffff/0xff), uint32(  blueB)*(0xffff/0xff), uint32(  blueA)*(0xffff/0xff))
			t.Logf(" white rgba(%d,%d,%d,%d) -> (%d,%d,%d,%d)",  whiteR,  whiteG,  whiteB,  whiteA,   uint32( whiteR)*(0xffff/0xff), uint32( whiteG)*(0xffff/0xff), uint32( whiteB)*(0xffff/0xff), uint32( whiteA)*(0xffff/0xff))
			t.Log("EXPECTED & ACTUAL")
			for i, originalSampleXY := range originalSampleXYs {
				newSampleXY := newSampleXYs[i]
				originalSampleRGBA := originalSampleRGBAs[i]
				newSampleRGBA := newSampleRGBAs[i]
				t.Logf("orig(%d,%d) & new(%d,%d)-> orig(%d,%d,%d,%d) & new(%d,%d,%d,%d)",
					originalSampleXY[0], originalSampleXY[1],
					newSampleXY[0],      newSampleXY[1],
					originalSampleRGBA[0], originalSampleRGBA[1], originalSampleRGBA[2], originalSampleRGBA[3],
					     newSampleRGBA[0],      newSampleRGBA[1],      newSampleRGBA[2],      newSampleRGBA[3],
				)
				if originalSampleRGBA != newSampleRGBA {
					t.Logf("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
				}
			}
			continue
		}
	}
}
