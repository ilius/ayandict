package application

import (
	"os"

	"github.com/therecipe/qt/gui"
)

func loadPNGIcon(filename string) (*gui.QIcon, error) {
	data, err := res.ReadFile("res/" + filename)
	if err != nil {
		return nil, err
	}
	file, err := os.CreateTemp("", filename)
	if err != nil {
		return nil, err
	}
	file.Write(data)
	file.Close()
	pixmap := gui.NewQPixmap3(file.Name(), "PNG", 0)
	return gui.NewQIcon2(pixmap), nil
}

// func loadSVGIcon(filename string) *gui.QIcon {
// 	data, err := res.ReadFile("res/" + filename)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	image := gui.QImage_FromData(data, len(data), "SVG")
// 	image.Rect().SetSize(core.NewQSize2(36, 36))
// 	pixmap := gui.QPixmap_FromImage(image, core.Qt__AutoColor)
// 	return gui.NewQIcon2(pixmap)
// }
