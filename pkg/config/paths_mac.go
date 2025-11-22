//go:build darwin
// +build darwin

package config

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#include <Foundation/Foundation.h>
#include <stdlib.h>

char* getDir(NSSearchPathDirectory dir) {
    NSArray *paths = NSSearchPathForDirectoriesInDomains(
        dir,
        NSUserDomainMask,
        YES
    );
    if ([paths count] == 0) return NULL;
    NSString *path = [paths objectAtIndex:0];
    return strdup([path UTF8String]);
}
*/
import "C"

import (
	"os"
	"path/filepath"
	"unsafe"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
)

func foundationDir(dir C.NSSearchPathDirectory) string {
	cpath := C.getDir(dir)
	if cpath == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cpath))
	return C.GoString(cpath)
}

func platformConfigDir() string {
	lib := foundationDir(C.NSLibraryDirectory)
	if lib == "" {
		lib = filepath.Join(os.Getenv(S_HOME), "Library")
	}
	return filepath.Join(lib, "Preferences", appinfo.APP_DESC)
}

func GetCacheDir() string {
	parent := foundationDir(C.NSCachesDirectory)
	if parent == "" {
		parent = filepath.Join(os.Getenv(S_HOME), "Library", "Caches")
	}
	return filepath.Join(parent, appinfo.APP_DESC)
}

func GetStateDir() string {
	parent := foundationDir(C.NSApplicationSupportDirectory)
	if parent == "" {
		parent = filepath.Join(
			os.Getenv(S_HOME),
			"Library",
			"Application Support",
		)
	}
	return filepath.Join(parent, appinfo.APP_DESC)
}
