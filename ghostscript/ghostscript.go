// Package ghostscript provides simple, and idiomatic Go bindings for the
// Ghostscript Interpreter C API.
// For more information: http://www.ghostscript.com/doc/current/API.htm
package ghostscript

/*
#include <stdlib.h>
#include <ghostscript/gdevdsp.h>
#include <ghostscript/iapi.h>
#include <ghostscript/ierrors.h>
#cgo LDFLAGS: -lgs
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

const (
	MAX_SUPPORTED_REVISION = 915
	MIN_SUPPORTED_REVISION = 910
	GS_NO_ERRORS           = 0
	GS_FATAL_ERROR         = C.e_Fatal
)

var (
	instantiated bool
	mu           sync.Mutex
)

// Revision contains information about the Ghostscript interpreter.
type Revision struct {
	Product      string
	Copyright    string
	Revision     int
	RevisionDate int
}

// Ghostscript contains a pointer to the global static instance of Ghostscript.
// It should not be initialised manually. i.e. Call NewInstance instead.
// Only one instance of Ghostscript may exist at any time.
type Ghostscript struct {
	instance unsafe.Pointer
}

// CStrings converts a Go string array to a C array of char pointers,
// and returns a pointer to that array.
// It will allocate a new string array of the appropriate length, and should
// be garbage collected using FreeCStrings.
func CStrings(goStrings []string) **C.char {
	// Kids, don't try this at home. Absolutely nasty.
	var char C.char
	length := len(goStrings)
	charArray := C.calloc(C.size_t(unsafe.Sizeof(&char)), C.size_t(length))
	tmp := (*[1 << 30]*C.char)(unsafe.Pointer(charArray))[:length:length]
	for i, str := range goStrings {
		tmp[i] = C.CString(str)
	}
	return (**C.char)(charArray)
}

// FreeCStrings frees memory used to allocate a C array of char pointers.
func FreeCStrings(cStrings **C.char, length int) {
	// Yea, this is real nasty.
	tmp := (*[1 << 30]*C.char)(unsafe.Pointer(cStrings))[:length:length]
	for i, _ := range tmp {
		C.free(unsafe.Pointer(tmp[i]))
	}
	C.free(unsafe.Pointer(cStrings))
}

func instantiate() {
	mu.Lock()
	defer mu.Unlock()
	instantiated = true
}

func uninstantiate() {
	mu.Lock()
	defer mu.Unlock()
	instantiated = false
}

// IsInstantiated returns true if a global static instance of Ghostscript
// already exists.
// This should be used to ensure that only one instance of Ghostscript may
// exist at any time. Ghostscript does not support multiple instances.
func IsInstantiated() bool {
	return instantiated
}

// GetRevision is an implementation of gsapi_revision.
// It returns the version numbers, and strings of the Ghostscript interpreter.
// It is safe to call at any time, and it does not rely on an instance of
// Ghostscript.
// It should be called before any other interpreter library functions to ensure
// that the correct version of Ghostscript interpreter has been loaded.
func GetRevision() (Revision, error) {
	revision := Revision{}

	var gsapiRevision C.gsapi_revision_t
	if err := C.gsapi_revision(&gsapiRevision, C.int(unsafe.Sizeof(gsapiRevision))); err != 0 {
		return revision, fmt.Errorf("revision structure size is incorrect, expected: %+v", err)
	}

	revision.Product = C.GoString(gsapiRevision.product)
	revision.Copyright = C.GoString(gsapiRevision.copyright)
	revision.Revision = int(gsapiRevision.revision)
	revision.RevisionDate = int(gsapiRevision.revisiondate)
	return revision, nil
}

// NewInstance is an implementation of gsapi_new_instance.
// It returns a global static instance of Ghostscript (encapsulated in a
// struct).
// i.e. Do not call NewInstance more than once, otherwise an error will be
// returned.
func NewInstance() (*Ghostscript, error) {
	if IsInstantiated() {
		return nil, fmt.Errorf("unable to create a new instance of Ghostscript, an instance already exists")
	}

	rev, err := GetRevision()
	if err != nil {
		return nil, err
	}
	if rev.Revision < MIN_SUPPORTED_REVISION || rev.Revision > MAX_SUPPORTED_REVISION {
		return nil, fmt.Errorf("Ghostscript interpreter version not supported: %d (must be >%d, and <%d)", rev.Revision, MIN_SUPPORTED_REVISION, MAX_SUPPORTED_REVISION)
	}

	var instance unsafe.Pointer
	if err := C.gsapi_new_instance(&instance, nil); err < GS_NO_ERRORS {
		return nil, fmt.Errorf("unable to create a new instance of Ghostscript: %+v", err)
	}

	defer instantiate()
	return &Ghostscript{instance}, nil
}

// Destroy is an implementation of gsapi_delete_instance.
// It destroys a global static instance of Ghostscript.
// It should be called only after Exit has been called if Init has been called.
func (gs *Ghostscript) Destroy() {
	defer uninstantiate()
	C.gsapi_delete_instance(gs.instance)
}

// Init is an implementation of gsapi_init_with_args.
// It initialises the Ghostscript interpreter given a set of arguments.
// The first argument is ignored, and the arguments that will be used starts
// from index 1.
func (gs *Ghostscript) Init(args []string) error {
	cArgs := CStrings(args)
	defer FreeCStrings(cArgs, len(args))

	err := C.gsapi_init_with_args(gs.instance, C.int(len(args)), cArgs)
	if err <= GS_FATAL_ERROR {
		_ = gs.Exit()
		return fmt.Errorf("unable to initialise Ghostscript interpreter due to a fatal error, exiting: %+v", err)
	}
	if err < GS_NO_ERRORS {
		return fmt.Errorf("unable to initialise Ghostscript interpreter: %+v", err)
	}
	return nil
}

// RunOnString is an implementation of gsapi_run_string_with_length.
// It runs the Ghostscript interpreter against a document in the form of a
// fixed length string.
func (gs *Ghostscript) RunOnString(strDoc string) error {
	gsStrDoc := C.CString(strDoc)
	defer C.free(unsafe.Pointer(gsStrDoc))

	var exitCode C.int
	err := C.gsapi_run_string_with_length(gs.instance, gsStrDoc, C.uint(len(strDoc)), C.int(0), &exitCode)
	if err <= GS_FATAL_ERROR {
		_ = gs.Exit()
		return fmt.Errorf("unable to run Ghostscript interpreter due to a fatal error, exiting: %+v", err)
	}
	if err < GS_NO_ERRORS {
		return fmt.Errorf("unable to run Ghostscript interpreter on string: %+v", err)
	}
	return nil
}

// RunOnFile is an implementation of gsapi_run_file.
// It runs the Ghostscript interpreter against an existing file, given its
// name / path.
func (gs *Ghostscript) RunOnFile(fnDoc string) error {
	gsFnDoc := C.CString(fnDoc)
	defer C.free(unsafe.Pointer(gsFnDoc))

	var exitCode C.int
	err := C.gsapi_run_file(gs.instance, gsFnDoc, C.int(0), &exitCode)
	if err <= GS_FATAL_ERROR {
		_ = gs.Exit()
		return fmt.Errorf("unable to run Ghostscript interpreter due to a fatal error, exiting: %+v", err)
	}
	if err < GS_NO_ERRORS {
		return fmt.Errorf("unable to run Ghostscript interpreter on file name: %+v", err)
	}
	return nil
}

// Exit is an implementation of gsapi_exit.
// It exits the Ghostscript interpreter.
// It must be called if Init has been called, and just before Destroy.
func (gs *Ghostscript) Exit() error {
	if err := C.gsapi_exit(gs.instance); err < GS_NO_ERRORS {
		return fmt.Errorf("unable to exit Ghostscript interpreter: %+v", err)
	}
	return nil
}
