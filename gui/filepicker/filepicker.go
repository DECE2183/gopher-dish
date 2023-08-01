package filepicker

/*
#cgo LDFLAGS: -lstdc++ -lwindowsapp -luuid
#cgo CXXFLAGS: -std=c++14
#include "filepicker.h"
#include <stdlib.h>
#include <string.h>
*/
import "C"
import "unsafe"

func OpenDir(title string) string {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	cpath := C.malloc(1024)
	defer C.free(cpath)

	if C.open_directory(ctitle, (*C.char)(cpath)) {
		return C.GoString((*C.char)(cpath))
	} else {
		return ""
	}
}

func OpenFile(title string) string {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	cpath := C.malloc(1024)
	defer C.free(cpath)

	if C.open_file(ctitle, (*C.char)(cpath)) {
		return C.GoString((*C.char)(cpath))
	} else {
		return ""
	}
}

func SaveFile(title string, initialName string) string {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	cpath := C.malloc(1024)
	defer C.free(cpath)

	n := C.CString(initialName)
	C.memcpy(cpath, unsafe.Pointer(n), C.size_t(len(initialName)+1))
	C.free(unsafe.Pointer(n))

	if C.save_file(ctitle, (*C.char)(cpath)) {
		return C.GoString((*C.char)(cpath))
	} else {
		return ""
	}
}
