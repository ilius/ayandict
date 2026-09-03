package application

import (
	"log/slog"
	"os"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/logging"
	qt "github.com/mappu/miqt/qt6"
)

// setFFmpegHwDisabled disables Qt Multimedia's FFmpeg hardware-acceleration
// backends unless the user explicitly configured them. On startup the FFmpeg
// backend probes VA-API/VDPAU, which opens its own X11 connection; on systems
// without a working GPU driver (e.g. VMs) that probing can trigger a fatal
// X11 error (BadWindow) and crash the app. ayandict only plays short audio
// clips, where hardware acceleration is never beneficial anyway.
func setFFmpegHwDisabled() {
	for _, key := range []string{
		"QT_FFMPEG_DECODING_HW_DEVICE_TYPES",
		"QT_FFMPEG_ENCODING_HW_DEVICE_TYPES",
	} {
		if os.Getenv(key) == "" {
			// empty list (",") disables all hardware backends
			_ = os.Setenv(key, ",")
		}
	}
}

func Run(query string) {
	setFFmpegHwDisabled()
	qt.QCoreApplication_SetApplicationName(appinfo.APP_DESC)
	qt.QGuiApplication_SetDesktopFileName(appinfo.APP_NAME)

	if config.PrivateMode {
		qt.QCoreApplication_SetApplicationName(appinfo.APP_DESC + " (private mode)")
	}

	// following line changes window title to "AyanDict — ayandict" on Linux
	// (not tested on Mac or Windows)
	// qt.QGuiApplication_SetApplicationDisplayName(appinfo.APP_NAME)

	app := &Application{
		QApplication: qt.NewQApplication(os.Args),
		window:       qt.NewQMainWindow(nil),
	}
	logging.ShowErrorDialog = showErrorMessage
	app.style = qt.QApplication_Style()
	app.bottomBoxStyleOpt = qt.NewQStyleOptionButton()

	cacheDir := config.Paths.CacheDir()
	if cacheDir == "" {
		slog.Error("cacheDir is empty")
	}
	{
		err := os.MkdirAll(cacheDir, 0o755)
		if err != nil {
			slog.Error("error in MkdirAll: " + err.Error())
		}
	}
	{
		err := os.MkdirAll(stateDir, 0o755)
		if err != nil {
			slog.Error("error in MkdirAll: " + err.Error())
		}
	}
	app.Run(query)
}
