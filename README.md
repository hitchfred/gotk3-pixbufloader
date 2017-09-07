Example code that shows that a finalizer can be run prematurely when
the call to runtime.KeepAlive() is commented out.

The program can be run with the included sample data:
```shell
  $ gotk3-pixbufloader pix/*
```

There is an -iterations switch that'll run the code in a loop to stress
things a bit and occasionally have GTK spit out error messages. 
```shell
  $ gotk3-pixbufloader -iterations 500 pix/*
```

Rarely, the finalizer will be called and unref() the loader object
before gtk_loader_get_pixbuf() has a chance to do anything, resulting
in a nil Pixbuf being returned.

```
(gotk3-pixbufloader.exe:7732): GdkPixbuf-CRITICAL **: gdk_pixbuf_loader_get_pixbuf: assertion 'GDK_IS_PIXBUF_LOADER (loader)' failed
```
