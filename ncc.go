package lookup

import (
	"fmt"
	"math"
)

// GPoint represents a match of a template inside an image.
type GPoint struct {
	X, Y int
	G    float64
}

func lookupAll(imgBin *imageBinary, x1, y1, x2, y2 int, templateBin *imageBinary, m float64, all bool) ([]GPoint, error) {
	var list []GPoint

	templateWidth := templateBin.width
	templateHeight := templateBin.height
	for x := x1; x <= x2-templateWidth+1; x++ {
		for y := y1; y <= y2-templateHeight+1; y++ {
			g, err := lookup(imgBin, templateBin, x, y, m)
			if err != nil {
				return nil, err
			}
			if g != nil {
				list = append(list, *g)
				if !all {
					return list, nil
				}
			}
		}
	}
	return list, nil
}

//	Normalized Cross Correlation algorithm
//	1) mean && stddev
//	2) image1(x,y) - mean1 && image2(x,y) - mean2
//	3) [3] = (image1(x,y) - mean)(x,y) * (image2(x,y) - mean)(x,y)
//	4) [4] = mean([3])
//	5) [4] / (stddev1 * stddev2)
//
// See http://www.fmwconcepts.com/imagemagick/similar/index.php
func lookup(img *imageBinary, template *imageBinary, x int, y int, m float64) (*GPoint, error) {
	ci := img.channels
	ct := template.channels

	ii := min(len(ci), len(ct))
	g := math.MaxFloat64

	for i := 0; i < ii; i++ {
		cct := ct[i]
		cci := ci[i]
		if cct.channelType != cci.channelType {
			return nil, fmt.Errorf("incompatible channels %d <> %d", cct.channelType, cci.channelType)
		}
		gg := gamma(cci, cct, x, y)
		if gg < m {
			return nil, nil
		}
		g = math.Min(g, gg)
	}
	return &GPoint{X: x, Y: y, G: g}, nil
}

func gamma(img *imageBinaryChannel, template *imageBinaryChannel, xx int, yy int) float64 {
	d := denominator(img, template, xx, yy)
	if d == 0 {
		return -1
	}

	n := numerator(img, template, xx, yy)
	return n / d
}

func denominator(img *imageBinaryChannel, template *imageBinaryChannel, xx int, yy int) float64 {
	di := img.dev2nRect(xx, yy, xx+template.width-1, yy+template.height-1)
	dt := template.dev2n()
	return math.Sqrt(di * dt)
}

func numerator(img *imageBinaryChannel, template *imageBinaryChannel, offsetX int, offsetY int) float64 {
	imgWidth := img.width
	imgArray := img.zeroMeanImage
	templateWidth := template.width
	templateHeight := template.height
	templateArray := template.zeroMeanImage
	var sum float64
	for x := 0; x < templateWidth; x++ {
		for y := 0; y < templateHeight; y++ {
			value := imgArray[(offsetY+y)*imgWidth+offsetX+x] * templateArray[y*templateWidth+x]
			sum += value
		}
	}
	return sum
}
