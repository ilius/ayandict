package application

import (
	"log/slog"
	"os"
	"sync"

	qt "github.com/mappu/miqt/qt6"
)

var (
	iconMap      = map[string]*qt.QIcon{}
	iconMapMutex sync.RWMutex
)

func loadPNGIcon(filename string) (*qt.QIcon, error) {
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
	pixmap := qt.NewQPixmap6(file.Name(), "PNG")
	icon = qt.NewQIcon2(pixmap)
	if icon == nil {
		slog.Error("error loading png icon: icon is nil: " + filename)
		panic("error loading png icon: icon is nil")
	}
	iconMapMutex.Lock()
	iconMap[filename] = icon
	iconMapMutex.Unlock()
	return icon, nil
}

// func loadSVGIcon(filename string) *qt.QIcon {
// 	data, err := res.ReadFile("res/" + filename)
// 	if err != nil {
// 		slog.Error("error", "err", err)
// 		return nil
// 	}
// 	image := qt.QImage_FromData(data, len(data), "SVG")
// 	image.Rect().SetSize(qt.NewQSize2(36, 36))
// 	pixmap := qt.QPixmap_FromImage(image, qt.AutoColor)
// 	return qt.NewQIcon2(pixmap)
// }
