#include <pango/pango.h>
#include <cairo.h>
#include <cairo-pdf.h>

extern cairo_status_t GoWriteToStream(void *closure, unsigned char *data, unsigned int length);

PangoFontFamily *indexFamily(PangoFontFamily **array, int i) 
{ 
    return array[i]; 
}

cairo_surface_t *gocairo_pdf_surface_create_for_stream (
        void *closure, 
        double width_in_points, 
        double height_in_points) 
{
	typedef cairo_status_t (*WriteFn)(void *, const unsigned char *, unsigned int);
	return cairo_pdf_surface_create_for_stream((WriteFn)&GoWriteToStream, closure, width_in_points, height_in_points);
}
