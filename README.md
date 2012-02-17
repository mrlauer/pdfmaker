PDFMaker
========

A tiny webapp that makes PDFs. Under construction, obviously.

Requirements
------------

You need go, updated to a reasonably recent weekly (will be go 1 when it is available).

You need cairo and pango and all that they entail. The build relies on pkg-config to find them.

I'm using OSX 10.6, and used homebrew for pango, and the built-in cairo. That works fine, although it's probably not guaranteed to do so forever.

Building
--------

Use the go tool. You need to set GOPATH manually.

```bash
cd pdfmaker
export GOPATH=`pwd`
go install pdfapp
./bin/pdfapp
```
