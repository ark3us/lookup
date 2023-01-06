package lookup

import "image"

// Lookup implements a image search algorithm based on Normalized Cross Correlation.
// For an overview of the algorithm, see http://www.fmwconcepts.com/imagemagick/similar/index.php
type Lookup struct {
	imgBin *imageBinary
}

// NewLookup creates a Lookup object. Lookup objects created by these function will do a gray
// scale search of the templates. Note that if the image or any template is not GrayScale,
// it will be converted automatically. This could cause some performance hits
func NewLookup(image image.Image) *Lookup {
	imgGray := ensureGrayScale(image)
	return &Lookup{
		imgBin: newImageBinary(imgGray),
	}
}

// NewLookupColor creates a Lookup object that works with all color channels of the image (RGB).
// Keep in mind that Lookup objects created with this function can only accept color
// templates as parameters to the FindAll method
func NewLookupColor(image image.Image) *Lookup {
	return &Lookup{
		imgBin: newImageBinary(image),
	}
}

// FindInRect searches for all occurrences of template only inside a part of the image.
// This can be used to speed up the search if you know the region of the
// image that the template should appear in.
func (l *Lookup) FindInRect(template image.Image, rect image.Rectangle, threshold float64, all bool) ([]GPoint, error) {
	if len(l.imgBin.channels) == 1 {
		template = ensureGrayScale(template)
	}
	tb := newImageBinary(template)
	return lookupAll(l.imgBin, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y, tb, threshold, all)
}

// FindAll searches for all occurrences of template inside the whole image.
func (l *Lookup) FindAll(template image.Image, threshold float64) ([]GPoint, error) {
	return l.FindInRect(template, image.Rect(0, 0, l.imgBin.width-1, l.imgBin.height-1), threshold, true)
}

// FindOne searches for the first occurrence of template inside the whole image.
func (l *Lookup) FindOne(template image.Image, threshold float64) ([]GPoint, error) {
	return l.FindInRect(template, image.Rect(0, 0, l.imgBin.width-1, l.imgBin.height-1), threshold, false)
}
