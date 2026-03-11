//go:build darwin

package notify

/*
#cgo darwin LDFLAGS: -framework UserNotifications -framework Foundation -framework AppKit
#include <stdlib.h>

extern void SetDarwinAppIcon(const void *data, int len);
extern void SendDarwinNotification(const char *title, const char *message);
*/
import "C" //nolint:dupImport // cgo requires separate import block

import (
	"sync"
	"unsafe"
)

type darwinNotifier struct{}

var appIconOnce sync.Once

func init() {
	platform = &darwinNotifier{}
}

func (d *darwinNotifier) Send(title, message string, icon []byte) error {
	// Set the process's application icon so macOS uses it for notifications.
	appIconOnce.Do(func() {
		C.SetDarwinAppIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
	})

	cTitle := C.CString(title)
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cMessage))

	C.SendDarwinNotification(cTitle, cMessage)
	return nil
}
