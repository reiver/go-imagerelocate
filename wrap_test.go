package imagerelocate_test

import (
	"github.com/reiver/go-imagerelocate"

	"github.com/reiver/go-pel"

	"image"
	"math/rand"
	"time"

	"testing"
)

func TestWrap(t *testing.T) {

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
				t.Errorf("For test #%d, this should not happens — the original-y is outside of the original bounds.", testNumber)
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
