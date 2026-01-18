package main

import (
	"bytes"
	"image"
	// "image/color"
	"image/gif"

	// "image/draw"
	"golang.org/x/image/math/fixed"

	"golang.org/x/image/font"
    "golang.org/x/image/font/opentype"
    "golang.org/x/image/font/gofont/goregular"
)

func addCaption(gifBytes []byte, text string) []byte {
    reader := bytes.NewReader(gifBytes)
    g, err := gif.DecodeAll(reader)
    if err != nil {
        panic(err)
    }
    
	fullWidth := g.Config.Width
	fullHeight := g.Config.Height

    // Loop through each frame and add text
    for _, frame := range g.Image {
        addTextToFrame(frame, text, fullWidth, fullHeight)
    }
    
    // Re-encode to gif
    var buf bytes.Buffer
    err = gif.EncodeAll(&buf, g)
    if err != nil {
        panic(err)
    }
    
    return buf.Bytes()
}

func addTextToFrame(img *image.Paletted, text string, fullWidth, fullHeight int) {
    ttf, _ := opentype.Parse(goregular.TTF)
    face, _ := opentype.NewFace(ttf, &opentype.FaceOptions{
        Size: 24,
        DPI:  72,
    })
    
    // Actually measure the text
    drawer := font.Drawer{Face: face}
    textWidth := drawer.MeasureString(text).Ceil()
    
    x := (fullWidth - textWidth) / 2
    y := fullHeight - 20  // Closer to bottom
    
    // Draw outline
    offsets := []struct{dx, dy int}{
        {-1, -1}, {0, -1}, {1, -1},
        {-1,  0},          {1,  0},
        {-1,  1}, {0,  1}, {1,  1},
    }
    
    for _, offset := range offsets {
        d := &font.Drawer{
            Dst:  img,
            Src:  image.Black,
            Face: face,
            Dot:  fixed.P(x+offset.dx, y+offset.dy),
        }
        d.DrawString(text)
    }
    
    // Draw white text
    d := &font.Drawer{
        Dst:  img,
        Src:  image.White,
        Face: face,
        Dot:  fixed.P(x, y),
    }
    d.DrawString(text)
}
