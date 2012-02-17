package textproc

/*
#cgo CFLAGS: -I/usr/local/include/pango-1.0 -I/usr/X11/include/cairo
#cgo CFLAGS: -I/usr/local/include/glib-2.0
#cgo CFLAGS: -I/usr/local/Cellar/glib/2.30.2/lib/glib-2.0/include
#cgo LDFLAGS: -L/usr/local/lib -L/usr/X11/lib -L/usr/local/Cellar/cairo/1.10.2/lib
#cgo LDFLAGS: -lcairo -lpango-1.0 -lpangocairo-1.0 -lgobject-2.0 -lglib-2.0
#include <stdlib.h>
#include <cairo.h>
#include <cairo-pdf.h>
#include <pango/pango.h>
#include <pango/pangocairo.h>

PangoFontFamily *indexFamily(PangoFontFamily **array, int i) { return array[i]; }

extern cairo_status_t GoWriteToStream(void *closure, const unsigned char *data, unsigned int length);

cairo_surface_t *	gocairo_pdf_surface_create_for_stream (
														 void *closure,
														 double width_in_points,
														 double height_in_points)
{
	return cairo_pdf_surface_create_for_stream(&GoWriteToStream, closure, width_in_points, height_in_points);
}
*/
import "C"

import(
	"unsafe"
	"io"
)

// Return a list of available fonts
func ListFontFamilies() []string {
	var names []string
	var families **C.PangoFontFamily
	var nfam C.int
	var fontmap *C.PangoFontMap
	fontmap = C.pango_cairo_font_map_get_default()
	C.pango_font_map_list_families(fontmap, &families, &nfam)
	for i := 0; i < int(nfam); i++ {
		family := C.indexFamily(families, C.int(i))
		familyname := C.pango_font_family_get_name(family)
		names = append(names, C.GoString(familyname))
	}
	C.g_free(C.gpointer(families))
	return names
}

//export GoWriteToStream
func GoWriteToStream(closure unsafe.Pointer, data *C.char, length C.uint) C.cairo_status_t {
	stream := *(*io.Writer)(closure)
	bytes := C.GoBytes(unsafe.Pointer(data), C.int(length))
	stream.Write(bytes)
	return C.cairo_status_t(0)
}

type TypesettingProps struct {
	Fontname string
	Fontsize float64
	Baselineskip float64
}

type TextObject interface {
	WriteAt(text string, props TypesettingProps, baseline float64) error
	Close()
}

type PDFStreamTextObject struct {
	surface *C.cairo_surface_t
	context *C.cairo_t
}

func (t *PDFStreamTextObject) WriteAt(text string, props TypesettingProps, x float64, y float64) error {
	var layout *C.PangoLayout
	var font_description *C.PangoFontDescription

	font_description = C.pango_font_description_new()
	cfontname := C.CString(props.Fontname)
	defer C.free(unsafe.Pointer(cfontname))
	C.pango_font_description_set_family(font_description, cfontname)
	C.pango_font_description_set_weight(font_description, C.PANGO_WEIGHT_NORMAL)
	C.pango_font_description_set_absolute_size(font_description, C.double(props.Fontsize)*C.PANGO_SCALE)

	layout = C.pango_cairo_create_layout(t.context)
	C.pango_layout_set_font_description(layout, font_description)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.pango_layout_set_text(layout, ctext, -1)

	C.cairo_set_source_rgb(t.context, 0.0, 0.0, 0.0)
	C.cairo_move_to(t.context, C.double(x), C.double(y))
	skip := props.Baselineskip
	nlines := int(C.pango_layout_get_line_count(layout))
	for i := 0; i < nlines; i++ {
		C.cairo_move_to(t.context, C.double(x), C.double(y+float64(i)*skip))
		C.pango_cairo_show_layout_line(t.context, C.pango_layout_get_line(layout, C.int(i)))
	}

	C.g_object_unref(C.gpointer(layout))
	C.pango_font_description_free(font_description)
	return nil
}

func (t *PDFStreamTextObject) Close() {
	C.cairo_destroy(t.context)
	C.cairo_surface_destroy(t.surface)
	t.context = nil
	t.surface = nil
}

func MakePDFStreamTextObject(writer io.Writer, width, height float64) *PDFStreamTextObject {
	var t PDFStreamTextObject
	t.surface = C.gocairo_pdf_surface_create_for_stream(unsafe.Pointer(&writer), C.double(width), C.double(height))
	t.context = C.cairo_create(t.surface)
	return &t
}
