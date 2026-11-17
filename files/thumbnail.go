package files

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

// ThumbnailResult holds the base64-encoded thumbnail and the original
// image dimensions. The server expects `thumb_inline` (base64 string),
// and the REAL width/height of the source image (not arbitrary values
// like 200x200). Without `thumb_inline` the server rejects Image-type
// sendMessage requests with INVALID_INPUT.
type ThumbnailResult struct {
	ThumbInline string // base64-encoded image bytes (small, resized)
	Width       int    // real image width (pixels)
	Height      int    // real image height (pixels)
}

// maxThumbDim caps the longest side of the generated thumbnail. The
// server rejects requests whose `thumb_inline` field is too large
// (HTTP 413 Request Entity Too Large from nginx). ParsRubika caps the
// longest side at 100px, which gives a small payload with good visual
// quality. For a 900x900 photo this yields a ~100x100 JPEG of about
// 3-5 KB — well within nginx's default 1 MB limit.
const maxThumbDim = 100

// GenerateThumbnail reads `src` (path or bytes), decodes the image,
// resizes it to a small thumbnail (max 100px on the longest side,
// preserving aspect ratio), re-encodes it as JPEG (quality 85 — much
// smaller than PNG for photos), and returns a ThumbnailResult.
//
// Returns nil if the image cannot be decoded — callers should fall
// back to type="File" in that case.
func GenerateThumbnail(src Source) (*ThumbnailResult, error) {
	var r io.Reader
	if src.Path != "" {
		f, err := os.Open(src.Path)
		if err != nil {
			return nil, fmt.Errorf("files: open %s: %w", src.Path, err)
		}
		defer f.Close()
		r = f
	} else if src.Bytes != nil {
		r = bytes.NewReader(src.Bytes)
	} else {
		return nil, fmt.Errorf("files: source must be a path or bytes")
	}

	// Decode the image. The blank imports above register the
	// gif/jpeg/png decoders so `image.Decode` can pick the right one.
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("files: decode image: %w", err)
	}

	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Resize to a small thumbnail. We use a simple area-averaging
	// downscale (nearest-neighbor upscale for tiny images) which is
	// good enough for a preview and avoids pulling in
	// golang.org/x/image. The thumbnail preserves the aspect ratio.
	thumb := resizeImage(img, origWidth, origHeight, maxThumbDim)

	// Encode as JPEG (quality 85). JPEG is dramatically smaller than
	// PNG for photographic content — a 100x100 JPEG is ~3-5 KB vs
	// 50-100 KB for PNG. This keeps the total sendMessage request
	// body well under nginx's default client_max_body_size.
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 85}); err != nil {
		return nil, fmt.Errorf("files: encode thumbnail: %w", err)
	}

	thumbB64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return &ThumbnailResult{
		ThumbInline: thumbB64,
		Width:       origWidth,
		Height:      origHeight,
	}, nil
}

// resizeImage scales `img` (origWidth x origHeight) so that its
// longest side is at most `maxDim` pixels, preserving the aspect
// ratio. It uses area-averaging when downscaling (good quality for
// thumbnails) and nearest-neighbor when upscaling (rare for
// thumbnails). This is a self-contained implementation that relies
// only on the standard library `image` package — no external deps.
func resizeImage(img image.Image, origW, origH, maxDim int) image.Image {
	if origW <= 0 || origH <= 0 {
		return img
	}
	// Compute the thumbnail dimensions preserving aspect ratio.
	newW, newH := origW, origH
	if newW > newH {
		if newW > maxDim {
			newH = newH * maxDim / newW
			newW = maxDim
		}
	} else {
		if newH > maxDim {
			newW = newW * maxDim / newH
			newH = maxDim
		}
	}
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}
	// If the image is already smaller than maxDim, return as-is.
	if newW == origW && newH == origH {
		return img
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	// Area-averaging downscale: each destination pixel is the
	// average of the source pixels it covers. This gives much
	// better quality than nearest-neighbor for large downscale
	// factors (e.g. 900 -> 100).
	xRatio := float64(srcW) / float64(newW)
	yRatio := float64(srcH) / float64(newH)

	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			// Source region covered by this destination pixel.
			sx0 := float64(x) * xRatio
			sy0 := float64(y) * yRatio
			sx1 := sx0 + xRatio
			sy1 := sy0 + yRatio
			r, g, b, a := averageRegion(img, srcBounds.Min.X, srcBounds.Min.Y, srcW, srcH, sx0, sy0, sx1, sy1)
			dst.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}
	return dst
}

// averageRegion computes the average RGBA color of the source image
// over the floating-point region [sx0,sx1) x [sy0,sy1). Each source
// pixel is weighted by the fraction of its area that falls inside
// the region — this is the standard area-averaging algorithm.
func averageRegion(img image.Image, offX, offY, srcW, srcH int, sx0, sy0, sx1, sy1 float64) (uint8, uint8, uint8, uint8) {
	x0 := int(sx0)
	y0 := int(sy0)
	x1 := int(sx1)
	if x1 <= x0 {
		x1 = x0 + 1
	}
	y1 := int(sy1)
	if y1 <= y0 {
		y1 = y0 + 1
	}
	if x1 > srcW {
		x1 = srcW
	}
	if y1 > srcH {
		y1 = srcH
	}
	if x0 < 0 {
		x0 = 0
	}
	if y0 < 0 {
		y0 = 0
	}

	var rSum, gSum, bSum, aSum, wSum float64
	for sy := y0; sy < y1; sy++ {
		for sx := x0; sx < x1; sx++ {
			sr, sg, sb, sa := img.At(offX+sx, offY+sy).RGBA()
			// Weight: fraction of this source pixel that is inside
			// the destination pixel's region. For the common case
			// where the downscale factor is an integer, this is 1.0.
			wx := 1.0
			wy := 1.0
			if float64(sx) < sx0 {
				wx = float64(sx+1) - sx0
			}
			if float64(sx+1) > sx1 {
				wx = sx1 - float64(sx)
			}
			if float64(sy) < sy0 {
				wy = float64(sy+1) - sy0
			}
			if float64(sy+1) > sy1 {
				wy = sy1 - float64(sy)
			}
			w := wx * wy
			// `RGBA()` returns pre-multiplied alpha values in the
			// range [0, 0xFFFF]. Convert to [0, 255] for averaging.
			if sa > 0 {
				rSum += float64(sr) / float64(sa) * float64(sa>>8) * w
				gSum += float64(sg) / float64(sa) * float64(sa>>8) * w
				bSum += float64(sb) / float64(sa) * float64(sa>>8) * w
			}
			aSum += float64(sa>>8) * w
			wSum += w
		}
	}
	if wSum == 0 {
		return 0, 0, 0, 0
	}
	return uint8(rSum / wSum), uint8(gSum / wSum), uint8(bSum / wSum), uint8(aSum / wSum)
}
