//go:build darwin
// +build darwin

package ospaths

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
	"path/filepath"
	"unsafe"
)

func foundationDir(dir C.NSSearchPathDirectory) string {
	cpath := C.getDir(dir)
	if cpath == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cpath))
	return C.GoString(cpath)
}

func (p *Paths) ConfigDir() string {
	lib := foundationDir(C.NSLibraryDirectory)
	if lib == "" {
		lib = filepath.Join(p.Home, "Library")
	}
	return filepath.Join(lib, "Preferences", p.AppNameCap)
}

func (p *Paths) CacheDir() string {
	parent := foundationDir(C.NSCachesDirectory)
	if parent == "" {
		parent = filepath.Join(p.Home, "Library", "Caches")
	}
	return filepath.Join(parent, p.AppNameCap)
}

func (p *Paths) StateDir() string {
	parent := foundationDir(C.NSApplicationSupportDirectory)
	if parent == "" {
		parent = filepath.Join(
			p.Home,
			"Library",
			"Application Support",
		)
	}
	return filepath.Join(parent, p.AppNameCap)
}
