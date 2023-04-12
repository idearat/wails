//go:build linux

package application

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <webkit2/webkit2.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>

typedef struct Window
{
    uint id;
} Window;

GtkApplication *GTKAPPLICATION(void *pointer)
{
    return GTK_APPLICATION(pointer);
}

GtkWidget *GTKWIDGET(void *pointer)
{
    return GTK_WIDGET(pointer);
}

GtkWindow *GTKWINDOW(void *pointer)
{
    return GTK_WINDOW(pointer);
}

GtkContainer *GTKCONTAINER(void *pointer)
{
    return GTK_CONTAINER(pointer);
}

GtkBox *GTKBOX(void *pointer)
{
    return GTK_BOX(pointer);
}

WebKitWebView *WEBKITWEBVIEW(void *pointer) {
   return WEBKIT_WEB_VIEW(pointer);
}

WebKitUserContentManager* WEBKITUSERCONTENTMANAGER(void *pointer) {
   return (WebKitUserContentManager *)pointer;
}

extern void processURLRequest(uint, void*);

void processRequest(void *request, gpointer data) {
    struct Window *window = data;
    printf("processRequest: %d\n", window->id);
    processURLRequest(window->id, request);
}
*/
import "C"
import (
	"fmt"
	"net/url"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

var showDevTools = func(window unsafe.Pointer) {}

func gtkBool(input bool) C.gboolean {
	if input {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

type linuxWebviewWindow struct {
	id          uint
	application unsafe.Pointer
	window      unsafe.Pointer
	webview     unsafe.Pointer
	parent      *WebviewWindow
	menubar     *C.GtkWidget
	vbox        *C.GtkWidget
	metadata    C.Window
	/*
		menu                          *menu.Menu
		accels                                   *C.GtkAccelGroup
	*/
	minWidth, minHeight, maxWidth, maxHeight int
}

func (w *linuxWebviewWindow) newWebview(gpuPolicy int) unsafe.Pointer {
	manager := C.webkit_user_content_manager_new()
	external := C.CString("external")
	C.webkit_user_content_manager_register_script_message_handler(manager, external)
	C.free(unsafe.Pointer(external))
	webview := C.webkit_web_view_new_with_user_content_manager(manager)
	//gtk_container_add(GTK_CONTAINER(window), webview);
	w.metadata.id = C.uint(w.parent.id)
	wails := C.CString("wails")
	C.webkit_web_context_register_uri_scheme(C.webkit_web_context_get_default(), wails,
		C.WebKitURISchemeRequestCallback(C.processRequest), C.gpointer(&w.metadata), nil)
	C.free(unsafe.Pointer(wails))
	//    g_signal_connect(G_OBJECT(webview), "load-changed", G_CALLBACK(webviewLoadChanged), window);
	// if (hideWindowOnClose)
	// {
	//     g_signal_connect(GTKWIDGET(window), "delete-event", G_CALLBACK(gtk_widget_hide_on_delete), NULL);
	// }
	// else
	// {
	//     g_signal_connect(GTKWIDGET(window), "delete-event", G_CALLBACK(close_button_pressed), window);
	// }

	settings := C.webkit_web_view_get_settings(C.WEBKITWEBVIEW(unsafe.Pointer(webview)))
	wails_io := C.CString("wails.io")
	empty := C.CString("")
	C.webkit_settings_set_user_agent_with_application_details(settings, wails_io, empty)
	C.free(unsafe.Pointer(wails_io))
	C.free(unsafe.Pointer(empty))

	switch gpuPolicy {
	case 0:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS)
		break
	case 1:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND)
		break
	case 2:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_NEVER)
		break
	default:
		C.webkit_settings_set_hardware_acceleration_policy(settings, C.WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND)
	}
	return unsafe.Pointer(webview)
}

func (w *linuxWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {
	// Create the menu
	thisMenu := newMenuImpl(menu)
	thisMenu.update()
	fmt.Println("linux.openContextMenu()")
	//C.windowShowMenu(w.nsWindow, thisMenu.nsMenu, C.int(data.X), C.int(data.Y))
}

func (w *linuxWebviewWindow) getZoom() float64 {
	return float64(C.webkit_web_view_get_zoom_level(C.WEBKITWEBVIEW(w.webview)))
}

func (w *linuxWebviewWindow) setZoom(zoom float64) {
	C.webkit_web_view_set_zoom_level(C.WEBKITWEBVIEW(w.webview), C.double(zoom))
}

func (w *linuxWebviewWindow) setFrameless(frameless bool) {
	if frameless {
		C.gtk_window_set_decorated(C.GTKWINDOW(w.window), C.gboolean(0))
	} else {
		C.gtk_window_set_decorated(C.GTKWINDOW(w.window), C.gboolean(1))
		// TODO: Deal with transparency for the titlebar if possible
		//       Perhaps we just make it undecorated and add a menu bar inside?
	}
}

func (w *linuxWebviewWindow) getScreen() (*Screen, error) {
	return getScreenForWindow(w)
}

func (w *linuxWebviewWindow) show() {
	fmt.Println("linuxWebviewWindow.show()")
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_widget_show_all(C.GTKWIDGET(w.window))
	})
}

func (w *linuxWebviewWindow) hide() {
	fmt.Println("linuxWebviewWindow.hide()")
	C.gtk_widget_hide(C.GTKWIDGET(w.window))
}

func (w *linuxWebviewWindow) setFullscreenButtonEnabled(enabled bool) {
	//	C.setFullscreenButtonEnabled(w.nsWindow, C.bool(enabled))
	fmt.Println("setFullscreenButtonEnabled")
}

func (w *linuxWebviewWindow) disableSizeConstraints() {
	fmt.Println("disableSizeConstraints")
	//C.windowDisableSizeConstraints(w.nsWindow)
}

func (w *linuxWebviewWindow) unfullscreen() {
	fmt.Println("unfullscreen")
}

func (w *linuxWebviewWindow) fullscreen() {
	fmt.Println("fullscreen")
}

func (w *linuxWebviewWindow) unminimise() {
	fmt.Println("unminimise")
}

func (w *linuxWebviewWindow) unmaximise() {
	fmt.Println("unmaximise")
}

func (w *linuxWebviewWindow) maximise() {
	fmt.Println("maximise")
}

func (w *linuxWebviewWindow) minimise() {
	fmt.Println("minimise")
}

func (w *linuxWebviewWindow) on(eventID uint) {
	// Don't think this is correct!
	// GTK Events are strings
	fmt.Println("on()", eventID)
	//C.registerListener(C.uint(eventID))
}

func (w *linuxWebviewWindow) zoom() {
	fmt.Println("zoom")
	//C.windowZoom(w.nsWindow)
}

func (w *linuxWebviewWindow) windowZoom() {
	fmt.Println("windowZoom")
	//	C.windowZoom(w.nsWindow)
}

func (w *linuxWebviewWindow) close() {
	fmt.Println("close")
	C.gtk_window_close(C.GTKWINDOW(w.window))
}

func (w *linuxWebviewWindow) zoomIn() {
	//	lvl := C.webkit_web_view_get_zoom_level()
	fmt.Println("zoomIn")
}

func (w *linuxWebviewWindow) zoomOut() {
	fmt.Println("zoomOut")
}

func (w *linuxWebviewWindow) zoomReset() {
	fmt.Println("zoomReset")
}

func (w *linuxWebviewWindow) toggleDevTools() {
	showDevTools(w.window)
}

func (w *linuxWebviewWindow) reload() {
	// TODO: This should be a constant somewhere I feel
	uri := C.CString("wails://")
	C.webkit_web_view_load_uri(C.WEBKITWEBVIEW(w.window), uri)
	C.free(unsafe.Pointer(uri))
}

func (w *linuxWebviewWindow) forceReload() {
	//TODO: Implement
	println("forceReload called on WebviewWindow", w.parent.id)
}

func (w *linuxWebviewWindow) center() {
	//C.windowCenter(w.nsWindow)
	fmt.Println("center")
}

func (w *linuxWebviewWindow) isMinimised() bool {
	return w.syncMainThreadReturningBool(func() bool {
		// FIXME: add in check here
		return false
		//return bool(C.windowIsMinimised(w.nsWindow))
	})
}

func (w *linuxWebviewWindow) isMaximised() bool {
	// return w.syncMainThreadReturningBool(func() bool {
	// 	return bool(C.windowIsMaximised(w.nsWindow))
	// })
	fmt.Println("isMaximised")
	return false
}

func (w *linuxWebviewWindow) isFullscreen() bool {
	fmt.Println("isFullscreen")
	// return w.syncMainThreadReturningBool(func() bool {
	// 	return bool(C.windowIsFullscreen(w.nsWindow))
	// })
	return false
}

func (w *linuxWebviewWindow) syncMainThreadReturningBool(fn func() bool) bool {
	var wg sync.WaitGroup
	wg.Add(1)
	var result bool
	globalApplication.dispatchOnMainThread(func() {
		result = fn()
		wg.Done()
	})
	wg.Wait()
	return result
}

func (w *linuxWebviewWindow) restore() {
	// restore window to normal size
	fmt.Println("restore")
}

func (w *linuxWebviewWindow) restoreWindow() {
	fmt.Println("restoreWindow") // what's the difference with "restore"?
}

func (w *linuxWebviewWindow) execJS(js string) {
	value := C.CString(js)
	C.webkit_web_view_evaluate_javascript(C.WEBKITWEBVIEW(w.webview),
		value,
		C.long(len(js)),
		nil,
		C.CString(""),
		nil,
		nil,
		nil)
	C.free(unsafe.Pointer(value))
}

func (w *linuxWebviewWindow) setURL(uri string) {
	if uri != "" {
		url, err := url.Parse(uri)
		if err == nil && url.Scheme == "" && url.Host == "" {
			// TODO handle this in a central location, the scheme and host might be platform dependant.
			url.Scheme = "wails"
			url.Host = "wails"
			uri = url.String()
		}
	}
	target := C.CString(uri)
	C.webkit_web_view_load_uri(C.WEBKITWEBVIEW(w.webview), target)
	C.free(unsafe.Pointer(target))
}

func (w *linuxWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	C.gtk_window_set_keep_above(C.GTKWINDOW(w.window), gtkBool(alwaysOnTop))
}

func newWindowImpl(parent *WebviewWindow) *linuxWebviewWindow {
	return &linuxWebviewWindow{
		application: (globalApplication.impl).(*linuxApp).application,
		parent:      parent,
	}
}

func (w *linuxWebviewWindow) setTitle(title string) {
	if !w.parent.options.Frameless {
		cTitle := C.CString(title)
		C.gtk_window_set_title(C.GTKWINDOW(w.window), cTitle)
		C.free(unsafe.Pointer(cTitle))
	}
}

func (w *linuxWebviewWindow) setSize(width, height int) {
	C.gtk_window_set_default_size(C.GTKWINDOW(w.window), C.int(width), C.int(height))
}

func (w *linuxWebviewWindow) setMinSize(width, height int) {
	size := C.GdkGeometry{min_width: C.int(width), min_height: C.int(height)}
	C.gtk_window_set_geometry_hints(C.GTKWINDOW(w.window), nil, &size, C.GDK_HINT_MIN_SIZE)
}
func (w *linuxWebviewWindow) setMaxSize(width, height int) {
	size := C.GdkGeometry{}
	size.max_width = 0
	size.max_height = 0

	// C.gtk_window_set_geometry_hints(C.GTKWINDOW(w.window), NULL, &size, flags);
}

func (w *linuxWebviewWindow) setResizable(resizable bool) {
	C.gtk_window_set_resizable(C.GTKWINDOW(w.window), 1)
}

func (w *linuxWebviewWindow) enableDevTools() {
	settings := C.webkit_web_view_get_settings(C.WEBKITWEBVIEW(w.webview))
	C.webkit_settings_set_enable_developer_extras(settings, C.int(1))
}

func (w *linuxWebviewWindow) size() (int, int) {
	var width, height C.int
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_window_get_size(C.GTKWINDOW(w.window), &width, &height)
		wg.Done()
	})
	wg.Wait()
	return int(width), int(height)
}

func (w *linuxWebviewWindow) setPosition(x, y int) {
	fmt.Println("setPosition")
	//	C.windowSetPosition(w.nsWindow, C.int(x), C.int(y))
}

func (w *linuxWebviewWindow) width() int {
	width, _ := w.size()
	return width
}

func (w *linuxWebviewWindow) height() int {
	_, height := w.size()
	return height
}

func (w *linuxWebviewWindow) run() {
	for eventId := range w.parent.eventListeners {
		w.on(eventId)
	}
	globalApplication.dispatchOnMainThread(func() {
		// FIXME: Should the application be passed to the window instead?
		app := (globalApplication.impl).(*linuxApp).application
		w.window = unsafe.Pointer(C.gtk_application_window_new(C.GTKAPPLICATION(app)))
		C.g_object_ref_sink(C.gpointer(w.window))
		w.webview = w.newWebview(1)
		w.vbox = C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0)
		C.gtk_container_add(C.GTKCONTAINER(w.window), w.vbox)
		if w.menubar != nil {
			C.gtk_box_pack_start(C.GTKBOX(unsafe.Pointer(w.vbox)), w.menubar, 0, 0, 0)
		}
		C.gtk_box_pack_start(C.GTKBOX(unsafe.Pointer(w.vbox)), C.GTKWIDGET(w.webview), 1, 1, 0)

		w.setSize(w.parent.options.Width, w.parent.options.Height)
		w.setTitle(w.parent.options.Title)
		w.setAlwaysOnTop(w.parent.options.AlwaysOnTop)
		w.setResizable(!w.parent.options.DisableResize)
		if w.parent.options.MinWidth != 0 || w.parent.options.MinHeight != 0 {
			w.setMinSize(w.parent.options.MinWidth, w.parent.options.MinHeight)
		}
		if w.parent.options.MaxWidth != 0 || w.parent.options.MaxHeight != 0 {
			w.setMaxSize(w.parent.options.MaxWidth, w.parent.options.MaxHeight)
		}
		w.setZoom(w.parent.options.Zoom)
		w.enableDevTools()
		w.setBackgroundColour(w.parent.options.BackgroundColour)

		// FIXME: This should either use the Linux options or a common set
		macOptions := w.parent.options.Mac
		switch macOptions.Backdrop {
		case MacBackdropTransparent:
			fmt.Println("MacBackdropTransparent - not implemented")
			//C.windowSetTransparent(w.nsWindow)
			//C.webviewSetTransparent(w.nsWindow)
		case MacBackdropTranslucent:
			fmt.Println("MacBackdropTranslucent - not implemented")
			//C.windowSetTranslucent(w.nsWindow)
			//C.webviewSetTransparent(w.nsWindow)
		}

		//		titleBarOptions := macOptions.TitleBar
		w.setFrameless(w.parent.options.Frameless)
		if macOptions.Appearance != "" {
			//C.windowSetAppearanceTypeByName(w.nsWindow, C.CString(string(macOptions.Appearance)))
			fmt.Println("Appearance: ", macOptions.Appearance)
		}

		if macOptions.InvisibleTitleBarHeight != 0 {

			fmt.Println("InvisibleTitleBarHeight - not implemented")
			//C.windowSetInvisibleTitleBar(w.nsWindow, C.uint(macOptions.InvisibleTitleBarHeight))
		}

		switch w.parent.options.StartState {
		case WindowStateMaximised:
			w.maximise()
		case WindowStateMinimised:
			w.minimise()
		case WindowStateFullscreen:
			w.fullscreen()

		}
		w.center()

		if w.parent.options.URL != "" {
			w.setURL(w.parent.options.URL)
		}
		// We need to wait for the HTML to load before we can execute the javascript
		w.parent.On(events.Mac.WebViewDidFinishNavigation, func(_ *WindowEventContext) {
			if w.parent.options.JS != "" {
				w.execJS(w.parent.options.JS)
			}
			if w.parent.options.CSS != "" {
				//FIXME: How do we do this?
				//C.windowInjectCSS(w.nsWindow, C.CString(w.parent.options.CSS))
			}
		})
		if w.parent.options.HTML != "" {
			w.setHTML(w.parent.options.HTML)
		}
		if w.parent.options.Hidden == false {
			w.show()
		}
	})
}

func (w *linuxWebviewWindow) setBackgroundColour(colour *RGBA) {
	if colour == nil {
		return
	}
	//C.windowSetBackgroundColour(w.nsWindow, C.int(colour.Red), C.int(colour.Green), C.int(colour.Blue), C.int(colour.Alpha))
	fmt.Println("SetBackgroundColour", colour)
}

func (w *linuxWebviewWindow) position() (int, int) {
	var x, y C.int
	var wg sync.WaitGroup
	wg.Add(1)
	go globalApplication.dispatchOnMainThread(func() {
		//C.windowGetPosition(w.nsWindow, &x, &y)
		wg.Done()
	})
	wg.Wait()
	return int(x), int(y)
}

func (w *linuxWebviewWindow) destroy() {
	fmt.Println("destroy")
}

func (w *linuxWebviewWindow) setHTML(html string) {
	fmt.Println("setHTML")
	// Convert HTML to C string
	//	cHTML := C.CString(html)
	// Render HTML
	// FIXME: What are we replacing?
	/*	C.webkit_web_view_load_alternate_html (C.WEBKITWEBVIEW(w.webview),
		const gchar *content,
		const gchar *content_uri,
		const gchar *base_uri);
	*/
}
