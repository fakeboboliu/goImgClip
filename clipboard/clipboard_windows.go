package clipboard

import (
	"errors"
	"fmt"
	"image"
	"reflect"
	"syscall"
	"time"
	"unsafe"

	"github.com/lxn/win"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32")
	lstrcpy  = kernel32.NewProc("lstrcpyW")
)

// Thanks to github.com/vova616/screenshot/blob/master/screenshot_windows.go
// and github.com/atotto/clipboard/

func GetImageFromClipBoard() (*image.RGBA, error) {
	err := waitOpenClipboard()
	if err != nil {
		return nil, err
	}
	defer win.CloseClipboard()

	src_hBmp := win.GetClipboardData(win.CF_BITMAP)
	if src_hBmp == 0 {
		return nil, nil
	}

	// Get DC and Select bmp into DC
	hDC := win.GetDC(0)
	if hDC == 0 {
		return nil, lastError()
	}
	defer win.ReleaseDC(0, hDC)
	m_hDC := win.CreateCompatibleDC(hDC)
	if m_hDC == 0 {
		return nil, lastError()
	}
	defer win.DeleteDC(m_hDC)
	win.SelectObject(m_hDC, win.HGDIOBJ(src_hBmp))

	// Get info of src_hBmp
	bmp := win.BITMAP{}
	if win.GetObject(win.HGDIOBJ(src_hBmp), unsafe.Sizeof(bmp), unsafe.Pointer(&bmp)) == 0 {
		return nil, err
	}
	x, y := int(bmp.BmWidth), int(bmp.BmHeight)
	// bmi is used to GetDIBits
	bmi := win.BITMAPINFO{}
	bmi.BmiHeader.BiSize = uint32(reflect.TypeOf(bmi.BmiHeader).Size())
	bmi.BmiHeader.BiWidth = int32(x)
	bmi.BmiHeader.BiHeight = int32(-y)
	bmi.BmiHeader.BiPlanes = 1
	bmi.BmiHeader.BiBitCount = 32
	bmi.BmiHeader.BiCompression = win.BI_RGB

	// Load data from memory
	slice := make([]byte, x*y*4)
	r := win.GetDIBits(m_hDC, win.HBITMAP(src_hBmp), 0, uint32(bmi.BmiHeader.BiHeight), &slice[0], &bmi, win.DIB_RGB_COLORS)
	if r == 0 {
		return nil, lastError()
	}

	// Change pixel order from BGRA
	imageBytes := make([]byte, len(slice))
	for i := 0; i < len(imageBytes); i += 4 {
		imageBytes[i], imageBytes[i+2], imageBytes[i+1], imageBytes[i+3] = slice[i+2], slice[i], slice[i+1], slice[i+3]
	}

	img := &image.RGBA{Pix: imageBytes, Stride: 4 * x, Rect: image.Rect(0, 0, x, y)}
	return img, nil
}

func SetTextToClipboard(text string) error {
	err := waitOpenClipboard()
	if err != nil {
		return err
	}
	defer win.CloseClipboard()

	ok := win.EmptyClipboard()
	if !ok {
		return lastError()
	}

	data, err := syscall.UTF16FromString(text)
	if err != nil {
		return err
	}

	// "If the hMem parameter identifies a memory object, the object must have
	// been allocated using the function with the GMEM_MOVEABLE flag."
	h := win.GlobalAlloc(win.GMEM_MOVEABLE, uintptr(len(data)*int(unsafe.Sizeof(data[0]))))
	if h == 0 {
		return lastError()
	}
	defer func() {
		if h != 0 {
			win.GlobalFree(h)
		}
	}()

	l := win.GlobalLock(h)
	if l == nil {
		return lastError()
	}

	var r uintptr
	r, _, err = lstrcpy.Call(uintptr(l), uintptr(unsafe.Pointer(&data[0])))
	if r == 0 {
		return err
	}

	ok = win.GlobalUnlock(h)
	if !ok {
		if err.(syscall.Errno) != 0 {
			return lastError()
		}
	}

	h2 := win.SetClipboardData(win.CF_UNICODETEXT, win.HANDLE(h))
	if uintptr(h2) == 0 {
		return lastError()
	}
	h = 0 // suppress deferred cleanup
	return nil
}

func waitOpenClipboard() error {
	started := time.Now()
	limit := started.Add(300 * time.Millisecond)
	for time.Now().Before(limit) {
		done := win.OpenClipboard(0)
		if done {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
	return errors.New("cannot open clipboard")
}

func lastError() error {
	return errors.New(fmt.Sprint("error:", win.GetLastError()))
}
