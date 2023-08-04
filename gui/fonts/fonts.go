package fonts

import (
	_ "embed"
	"unicode"

	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed RedHatMono/RedHatMono-Medium.ttf
var redhatMonoMedium []byte

var (
	RedhatMonoMedium12 *text.Atlas
	RedhatMonoMedium14 *text.Atlas
	RedhatMonoMedium18 *text.Atlas
)

var fontsInit = map[**text.Atlas]struct {
	binary []byte
	size   int
}{
	&RedhatMonoMedium12: {redhatMonoMedium, 12},
	&RedhatMonoMedium14: {redhatMonoMedium, 14},
	&RedhatMonoMedium18: {redhatMonoMedium, 18},
}

func init() {
	var (
		ttFont   *truetype.Font
		ttOpts   truetype.Options
		fontFace font.Face
		err      error
	)

	for atlasPtr, fontDesc := range fontsInit {
		ttFont, err = truetype.Parse(fontDesc.binary)
		if err != nil {
			panic(err)
		}

		ttOpts.Size = float64(fontDesc.size)

		fontFace = truetype.NewFace(ttFont, &ttOpts)
		*atlasPtr = text.NewAtlas(fontFace, text.ASCII, text.RangeTable(unicode.Latin))
	}
}
