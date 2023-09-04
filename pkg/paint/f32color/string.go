package f32color

import (
	"image/color"
	"strconv"
	"strings"
)

func HexString(c color.Color) string {
	cc := color.NRGBAModel.Convert(c).(color.NRGBA)
	return "#" + toHex(cc.R) + toHex(cc.G) + toHex(cc.B) + toHex(cc.A)
}

func toHex(v uint8) string {
	hex := strconv.FormatUint(uint64(v), 16)
	if len(hex) == 1 {
		return "0" + strings.ToUpper(hex)
	}
	return strings.ToUpper(hex)
}
