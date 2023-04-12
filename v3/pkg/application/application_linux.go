package application

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <webkit2/webkit2.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>

extern void processApplicationEvent(uint);

static void activate (GtkApplication* app, gpointer data) {
   // FIXME: should likely emit a WAILS specific code
   // events.Mac.EventApplicationDidFinishLaunching == 1032
   processApplicationEvent(1032);
}

static GtkApplication* init(char* name) {
   return gtk_application_new(name, G_APPLICATION_DEFAULT_FLAGS);
}

static int run(void *app) {
  g_signal_connect (app, "activate", G_CALLBACK (activate), NULL);
  g_application_hold(app);  // allows it to run without a window
  int status = g_application_run (G_APPLICATION (app), 0, NULL);
  g_application_release(app);
  g_object_unref (app);
  return status;
}

extern GtkApplication *GTKAPPLICATION(void *pointer);
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func init() {
	// Set GDK_BACKEND=x11 if currently unset and XDG_SESSION_TYPE is unset, unspecified or x11 to prevent warnings
	if os.Getenv("GDK_BACKEND") == "" && (os.Getenv("XDG_SESSION_TYPE") == "" || os.Getenv("XDG_SESSION_TYPE") == "unspecified" || os.Getenv("XDG_SESSION_TYPE") == "x11") {
		_ = os.Setenv("GDK_BACKEND", "x11")
	}
}

type linuxApp struct {
	application     unsafe.Pointer
	applicationMenu unsafe.Pointer
	parent          *App
}

func (m *linuxApp) hide() {

	//	C.hide()
}

func (m *linuxApp) show() {
	//	C.show()
}

func (m *linuxApp) on(eventID uint) {
	log.Println("linuxApp.on()", eventID)
	// TODO: Setup signal handling as appropriate
	// Note: GTK signals seem to be strings!
}

func (m *linuxApp) setIcon(icon []byte) {
	//	C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}

func (m *linuxApp) name() string {
	return "not implemented"
	// appName := C.getAppName()
	// defer C.free(unsafe.Pointer(appName))
	// return C.GoString(appName)
}

func (m *linuxApp) getCurrentWindowID() uint {
	window := C.gtk_application_get_active_window(C.GTKAPPLICATION(m.application))
	if window != nil {
		//		return uint(window.id)
		fmt.Println("getCurrentWindowID", window)
	}
	return uint(0)
}

func (m *linuxApp) setApplicationMenu(menu *Menu) {
	if menu == nil {
		// Create a default menu
		menu = defaultApplicationMenu()
	}
	menu.Update()
	// Convert impl to linuxMenu object
	//	m.applicationMenu = (menu.impl).(*linuxMenu).menuModel
	//C.gtk_application_set_app_menu(m.application, m.applicationMenu);
}

func (m *linuxApp) run() error {

	// Add a hook to the ApplicationDidFinishLaunching event
	// FIXME: add Wails specific events - i.e. Shouldn't platform specific ones be translated to Wails events?
	m.parent.On(events.Mac.ApplicationDidFinishLaunching, func() {
		// Do we need to do anything now?
	})

	C.run(m.application)
	return nil
}

func (m *linuxApp) destroy() {
	C.gtk_main_quit()
}

func newPlatformApp(parent *App) *linuxApp {
	name := strings.ToLower(strings.Replace(parent.options.Name, " ", "", -1))
	if name == "" {
		name = "undefined"
	}
	nameC := C.CString(fmt.Sprintf("org.wails.%s", name))
	app := &linuxApp{
		parent:      parent,
		application: unsafe.Pointer(C.init(nameC)),
	}
	C.free(unsafe.Pointer(nameC))
	return app
}

//export processApplicationEvent
func processApplicationEvent(eventID C.uint) {
	// TODO: add translation to Wails events
	//       currently reusing Mac specific values
	applicationEvents <- uint(eventID)
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	windowEvents <- &WindowEvent{
		WindowID: uint(windowID),
		EventID:  uint(eventID),
	}
}

//export processMessage
func processMessage(windowID C.uint, message *C.char) {
	windowMessageBuffer <- &windowMessage{
		windowId: uint(windowID),
		message:  C.GoString(message),
	}
}

//export processURLRequest
func processURLRequest(windowID C.uint, wkUrlSchemeTask unsafe.Pointer) {
	webviewRequests <- &webViewAssetRequest{
		Request:    webview.NewRequest(wkUrlSchemeTask),
		windowId:   uint(windowID),
		windowName: globalApplication.getWindowForID(uint(windowID)).Name(),
	}
}

//export processDragItems
func processDragItems(windowID C.uint, arr **C.char, length C.int) {
	var filenames []string
	// Convert the C array to a Go slice
	goSlice := (*[1 << 30]*C.char)(unsafe.Pointer(arr))[:length:length]
	for _, str := range goSlice {
		filenames = append(filenames, C.GoString(str))
	}
	windowDragAndDropBuffer <- &dragAndDropMessage{
		windowId:  uint(windowID),
		filenames: filenames,
	}
}

//export processMenuItemClick
func processMenuItemClick(menuID C.uint) {
	menuItemClicked <- uint(menuID)
}

func setIcon(icon []byte) {
	if icon == nil {
		return
	}
	//C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}
