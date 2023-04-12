//go:build linux

package application

/*
#include <stdlib.h>
#include <stdbool.h>

typedef struct Screen {
	const char* id;
	const char* name;
	int p_width;
	int p_height;
	int width;
	int height;
	int x;
	int y;
	int w_width;
	int w_height;
	int w_x;
	int w_y;
	float scale;
	double rotation;
	bool isPrimary;
} Screen;


int GetNumScreens(){
    return 0;
}

*/
import "C"
import (
	"fmt"
)

func cScreenToScreen(screen C.Screen) *Screen {

	return &Screen{
		Size: Size{
			Width:  int(screen.p_width),
			Height: int(screen.p_height),
		},
		Bounds: Rect{
			X:      int(screen.x),
			Y:      int(screen.y),
			Height: int(screen.height),
			Width:  int(screen.width),
		},
		WorkArea: Rect{
			X:      int(screen.w_x),
			Y:      int(screen.w_y),
			Height: int(screen.w_height),
			Width:  int(screen.w_width),
		},
		Scale:     float32(screen.scale),
		ID:        C.GoString(screen.id),
		Name:      C.GoString(screen.name),
		IsPrimary: bool(screen.isPrimary),
		Rotation:  float32(screen.rotation),
	}
}

func getPrimaryScreen() (*Screen, error) {
	//	cScreen := C.GetPrimaryScreen()

	return nil, fmt.Errorf("not implemented")
}

func getScreens() ([]*Screen, error) {
	// cScreens := C.getAllScreens()
	// defer C.free(unsafe.Pointer(cScreens))
	// numScreens := int(C.GetNumScreens())
	// displays := make([]*Screen, numScreens)
	// cScreenHeaders := (*[1 << 30]C.Screen)(unsafe.Pointer(cScreens))[:numScreens:numScreens]
	// for i := 0; i < numScreens; i++ {
	// 	displays[i] = cScreenToScreen(cScreenHeaders[i])
	// }
	return nil, fmt.Errorf("not implemented")
}

func getScreenForWindow(window *linuxWebviewWindow) (*Screen, error) {
	//	cScreen := C.getScreenForWindow(window.nsWindow)
	return nil, fmt.Errorf("not implemented")
}
