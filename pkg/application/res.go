package application

import (
	"log/slog"
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
	_, err = file.Write(data)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	pixmap := gui.NewQPixmap3(file.Name(), "PNG", 0)
	icon = gui.NewQIcon2(pixmap)
	if icon == nil {
		slog.Error("error loading png icon: icon is nil: " + filename)
		panic("error loading png icon: icon is nil")
	}
	iconMapMutex.Lock()
	iconMap[filename] = icon
	iconMapMutex.Unlock()
	return icon, nil
}

// func loadSVGIcon(filename string) *gui.QIcon {
// 	data, err := res.ReadFile("res/" + filename)
// 	if err != nil {
// 		slog.Error("error", "err", err)
// 		return nil
// 	}
// 	image := gui.QImage_FromData(data, len(data), "SVG")
// 	image.Rect().SetSize(core.NewQSize2(36, 36))
// 	pixmap := gui.QPixmap_FromImage(image, core.Qt__AutoColor)
// 	return gui.NewQIcon2(pixmap)
// }
