package application

import (
	"os"
	"sync"

	"github.com/ilius/qt/gui"
)

var (
	iconMap      = map[string]*gui.QIcon{}
	iconMapMutex sync.RWMutex
)

func loadPNGIcon(filename string) (*gui.QIcon, error) {
	iconMapMutex.RLock()
	icon, ok := iconMap[filename]
	iconMapMutex.RUnlock()
	if ok {
		return icon, nil
	}
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
	icon = gui.NewQIcon2(pixmap)
	iconMapMutex.Lock()
	iconMap[filename] = icon
	iconMapMutex.Unlock()
	return icon, nil
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
