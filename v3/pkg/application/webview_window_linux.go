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

GdkWindow *GDKWINDOW(void *pointer)
{
    return GDK_WINDOW(pointer);
}

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
    uint *window = data;
    printf("processRequest: %d\n", *window);
    processURLRequest(*window, request);
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
	/*
		menu                          *menu.Menu
		accels                                   *C.GtkAccelGroup
	*/
	lastWidth  int
	lastHeight int
}

func (w *linuxWebviewWindow) newWebview(gpuPolicy int) unsafe.Pointer {
	manager := C.webkit_user_content_manager_new()
	external := C.CString("external")
	C.webkit_user_content_manager_register_script_message_handler(manager, external)
	C.free(unsafe.Pointer(external))
	webview := C.webkit_web_view_new_with_user_content_manager(manager)
	//gtk_container_add(GTK_CONTAINER(window), webview); // do we need this?

	wails := C.CString("wails")
	id := C.uint(w.parent.id)
	C.webkit_web_context_register_uri_scheme(C.webkit_web_context_get_default(), wails,
		C.WebKitURISchemeRequestCallback(C.processRequest), C.gpointer(&id), nil)
	C.free(unsafe.Pointer(wails))

	//C.g_signal_connect(C.GTKWIDGET(w.window), "delete-event", G_CALLBACK(gtk_widget_hide_on_delete), NULL)
	//    g_signal_connect(G_OBJECT(webview), "load-changed", G_CALLBACK(webviewLoadChanged), window);
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
	mx, my, width, height, scale := w.getCurrentMonitorGeometry()
	return &Screen{
		ID:        fmt.Sprintf("%d", w.id),            // A unique identifier for the display
		Name:      w.parent.Name(),                    // The name of the display
		Scale:     float32(scale),                     // The scale factor of the display
		X:         mx,                                 // The x-coordinate of the top-left corner of the rectangle
		Y:         my,                                 // The y-coordinate of the top-left corner of the rectangle
		Size:      Size{Width: width, Height: height}, // The size of the display
		Bounds:    Rect{},                             // The bounds of the display
		WorkArea:  Rect{},                             // The work area of the display
		IsPrimary: false,                              // Whether this is the primary display
		Rotation:  0.0,                                // The rotation of the display
	}, nil
}

func (w *linuxWebviewWindow) show() {
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_widget_show_all(C.GTKWIDGET(w.window))
	})
}

func (w *linuxWebviewWindow) hide() {
	C.gtk_widget_hide(C.GTKWIDGET(w.window))
}

func (w *linuxWebviewWindow) setFullscreenButtonEnabled(enabled bool) {
	//	C.setFullscreenButtonEnabled(w.nsWindow, C.bool(enabled))
	fmt.Println("setFullscreenButtonEnabled - not implemented")
}

func (w *linuxWebviewWindow) disableSizeConstraints() {
	x, y, width, height, scale := w.getCurrentMonitorGeometry()
	w.setMinMaxSize(x, y, width*scale, height*scale)
}

func (w *linuxWebviewWindow) unfullscreen() {
	fmt.Println("unfullscreen")
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_window_unfullscreen(C.GTKWINDOW(w.window))
	})
	w.unmaximise()
}

func (w *linuxWebviewWindow) fullscreen() {
	w.maximise()
	w.lastWidth, w.lastHeight = w.size()
	globalApplication.dispatchOnMainThread(func() {
		x, y, width, height, scale := w.getCurrentMonitorGeometry()
		if x == -1 && y == -1 && width == -1 && height == -1 {
			return
		}
		fmt.Println("fullscreen", x, y, width, height, scale)
		w.setMinMaxSize(0, 0, width*scale, height*scale)
		w.setSize(width*scale, height*scale)
		C.gtk_window_fullscreen(C.GTKWINDOW(w.window))
	})
	w.setPosition(0, 0)
}

func (w *linuxWebviewWindow) unminimise() {
	C.gtk_window_present(C.GTKWINDOW(w.window))
	// gtk_window_unminimize (C.GTKWINDOW(w.window)) /// gtk4
}

func (w *linuxWebviewWindow) unmaximise() {
	C.gtk_window_unmaximize(C.GTKWINDOW(w.window))
}

func (w *linuxWebviewWindow) maximise() {
	C.gtk_window_maximize(C.GTKWINDOW(w.window))
}

func (w *linuxWebviewWindow) minimise() {
	C.gtk_window_iconify(C.GTKWINDOW(w.window))
}

func (w *linuxWebviewWindow) on(eventID uint) {
	// Don't think this is correct!
	// GTK Events are strings
	fmt.Println("on()", eventID)
	//C.registerListener(C.uint(eventID))
}

func (w *linuxWebviewWindow) zoom() {
	w.zoomIn()
}

func (w *linuxWebviewWindow) windowZoom() {
	w.zoom() // FIXME> This should be removed
}

func (w *linuxWebviewWindow) close() {
	C.gtk_window_close(C.GTKWINDOW(w.window))
}

func (w *linuxWebviewWindow) zoomIn() {
	lvl := C.webkit_web_view_get_zoom_level(C.WEBKITWEBVIEW(w.webview))
	C.webkit_web_view_set_zoom_level(C.WEBKITWEBVIEW(w.webview), lvl+0.5)
}

func (w *linuxWebviewWindow) zoomOut() {
	lvl := C.webkit_web_view_get_zoom_level(C.WEBKITWEBVIEW(w.webview))
	C.webkit_web_view_set_zoom_level(C.WEBKITWEBVIEW(w.webview), lvl-0.5)
}

func (w *linuxWebviewWindow) zoomReset() {
	C.webkit_web_view_set_zoom_level(C.WEBKITWEBVIEW(w.webview), 0.0)
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
	w.reload()
}

func (w linuxWebviewWindow) getCurrentMonitor() *C.GdkMonitor {
	// Get the monitor that the window is currently on
	display := C.gtk_widget_get_display(C.GTKWIDGET(w.window))
	gdk_window := C.gtk_widget_get_window(C.GTKWIDGET(w.window))
	if gdk_window == nil {
		return nil
	}
	return C.gdk_display_get_monitor_at_window(display, gdk_window)
}

func (w linuxWebviewWindow) getCurrentMonitorGeometry() (x int, y int, width int, height int, scale int) {
	monitor := w.getCurrentMonitor()
	if monitor == nil {
		return -1, -1, -1, -1, 1
	}
	var result C.GdkRectangle
	C.gdk_monitor_get_geometry(monitor, &result)
	scale = int(C.gdk_monitor_get_scale_factor(monitor))
	return int(result.x), int(result.y), int(result.width), int(result.height), scale
}

func (w *linuxWebviewWindow) center() {
	x, y, width, height, _ := w.getCurrentMonitorGeometry()
	if x == -1 && y == -1 && width == -1 && height == -1 {
		return
	}

	var windowWidth C.int
	var windowHeight C.int
	C.gtk_window_get_size(C.GTKWINDOW(w.window), &windowWidth, &windowHeight)

	newX := C.int(((width - int(windowWidth)) / 2) + x)
	newY := C.int(((height - int(windowHeight)) / 2) + y)

	// Place the window at the center of the monitor
	C.gtk_window_move(C.GTKWINDOW(w.window), newX, newY)
}

func (w *linuxWebviewWindow) isMinimised() bool {
	gdkwindow := C.gtk_widget_get_window(C.GTKWIDGET(w.window))
	state := C.gdk_window_get_state(gdkwindow)
	return state&C.GDK_WINDOW_STATE_ICONIFIED > 0
}

func (w *linuxWebviewWindow) isMaximised() bool {
	return w.syncMainThreadReturningBool(func() bool {
		gdkwindow := C.gtk_widget_get_window(C.GTKWIDGET(w.window))
		state := C.gdk_window_get_state(gdkwindow)
		return state&C.GDK_WINDOW_STATE_MAXIMIZED > 0 && state&C.GDK_WINDOW_STATE_FULLSCREEN == 0
	})
}

func (w *linuxWebviewWindow) isFullscreen() bool {
	return w.syncMainThreadReturningBool(func() bool {
		gdkwindow := C.gtk_widget_get_window(C.GTKWIDGET(w.window))
		state := C.gdk_window_get_state(gdkwindow)
		return state&C.GDK_WINDOW_STATE_FULLSCREEN > 0
	})
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
	// FIXME: never called!  - remove from webviewImpl interface
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
	C.gtk_window_resize(C.GTKWINDOW(w.window), C.gint(width), C.gint(height))
}

func (w *linuxWebviewWindow) setMinMaxSize(minWidth, minHeight, maxWidth, maxHeight int) {
	// if minWidth == 0 {
	// 	minWidth = -1
	// }
	// if minHeight == 0 {
	// 	minHeight = -1
	// }
	// if maxWidth == 0 {
	// 	maxWidth = -1
	// }
	// if maxHeight == 0 {
	// 	maxHeight = -1
	// }
	size := C.GdkGeometry{
		min_width:  C.int(minWidth),
		min_height: C.int(minHeight),
		max_width:  C.int(maxWidth),
		max_height: C.int(maxHeight),
	}
	C.gtk_window_set_geometry_hints(C.GTKWINDOW(w.window), nil, &size, C.GDK_HINT_MAX_SIZE|C.GDK_HINT_MIN_SIZE)
}

func (w *linuxWebviewWindow) setMinSize(width, height int) {
	size := C.GdkGeometry{min_width: C.int(width), min_height: C.int(height)}
	C.gtk_window_set_geometry_hints(C.GTKWINDOW(w.window), nil, &size, C.GDK_HINT_MIN_SIZE)
}

func (w *linuxWebviewWindow) setMaxSize(width, height int) {
	size := C.GdkGeometry{max_width: C.int(width), max_height: C.int(height)}
	C.gtk_window_set_geometry_hints(C.GTKWINDOW(w.window), C.GTKWIDGET(C.NULL), &size, C.GDK_HINT_MAX_SIZE)
}

func (w *linuxWebviewWindow) setResizable(resizable bool) {
	if resizable {
		C.gtk_window_set_resizable(C.GTKWINDOW(w.window), 1)
	} else {
		C.gtk_window_set_resizable(C.GTKWINDOW(w.window), 0)
	}
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
	mx, my, _, _, _ := w.getCurrentMonitorGeometry()
	globalApplication.dispatchOnMainThread(func() {
		C.gtk_window_move(C.GTKWINDOW(w.window), C.int(x+mx), C.int(y+my))
	})
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
		w.window = unsafe.Pointer(C.gtk_application_window_new(C.GTKAPPLICATION(w.application)))
		C.g_object_ref_sink(C.gpointer(w.window))
		w.webview = w.newWebview(1)
		w.vbox = C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0)
		C.gtk_container_add(C.GTKCONTAINER(w.window), w.vbox)
		if w.menubar != nil {
			C.gtk_box_pack_start(C.GTKBOX(unsafe.Pointer(w.vbox)), w.menubar, 0, 0, 0)
		}
		C.gtk_box_pack_start(C.GTKBOX(unsafe.Pointer(w.vbox)), C.GTKWIDGET(w.webview), 1, 1, 0)
		w.setTitle(w.parent.options.Title)
		w.setAlwaysOnTop(w.parent.options.AlwaysOnTop)
		w.setResizable(!w.parent.options.DisableResize)
		// only set min/max size if actually set
		if w.parent.options.MinWidth != 0 &&
			w.parent.options.MinHeight != 0 &&
			w.parent.options.MaxWidth != 0 &&
			w.parent.options.MaxHeight != 0 {
			w.setMinMaxSize(
				w.parent.options.MinWidth,
				w.parent.options.MinHeight,
				w.parent.options.MaxWidth,
				w.parent.options.MaxHeight,
			)
		}
		w.setSize(w.parent.options.Width, w.parent.options.Height)
		w.setZoom(w.parent.options.Zoom)
		w.enableDevTools()
		w.setBackgroundColour(w.parent.options.BackgroundColour)
		w.setFrameless(w.parent.options.Frameless)

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

func (w *linuxWebviewWindow) setTransparent() {
	screen := C.gtk_widget_get_screen(C.GTKWIDGET(w.window))
	visual := C.gdk_screen_get_rgba_visual(screen)

	if visual != nil && C.gdk_screen_is_composited(screen) == C.int(1) {
		C.gtk_widget_set_app_paintable(C.GTKWIDGET(w.window), C.gboolean(1))
		C.gtk_widget_set_visual(C.GTKWIDGET(w.window), visual)
	}
}

func (w *linuxWebviewWindow) setBackgroundColour(colour *RGBA) {
	if colour == nil {
		return
	}
	if colour.Alpha != 0 {
		w.setTransparent()
	}
	rgba := C.GdkRGBA{C.double(colour.Red) / 255.0, C.double(colour.Green) / 255.0, C.double(colour.Blue) / 255.0, C.double(colour.Alpha) / 255.0}
	C.webkit_web_view_set_background_color(C.WEBKITWEBVIEW(w.webview), &rgba)
}

func (w *linuxWebviewWindow) position() (int, int) {
	var x, y C.int
	var wg sync.WaitGroup
	wg.Add(1)
	go globalApplication.dispatchOnMainThread(func() {
		C.gtk_window_get_position(C.GTKWINDOW(w.window), &x, &y)
		wg.Done()
	})
	wg.Wait()
	return int(x), int(y)
}

func (w *linuxWebviewWindow) destroy() {
	C.gtk_widget_destroy(C.GTKWIDGET(w.window))
	C.gtk_widget_destroy(C.GTKWIDGET(w.webview))
}

func (w *linuxWebviewWindow) setHTML(html string) {
	fmt.Println("setHTML")
	// Convert HTML to C string
	//	cHTML := C.CString(html)
	// Render HTML
	// FIXME: What are we replacing?
	/*
		C.webkit_web_view_load_alternate_html (C.WEBKITWEBVIEW(w.webview),
						const gchar *content,
						const gchar *content_uri,
						const gchar *base_uri);
	*/
}
