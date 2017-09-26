// +build linux

package capsicum

/*
#cgo LDFLAGS: -lcaprights -lffi
#include <assert.h>
#include <errno.h>
#include <ffi.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <sys/capsicum.h>

bool doCapRightsAny(void (*fn)(void), ffi_type *rtype, void *rvalue, cap_rights_t *rights, uint64_t values[], int c) {
    ffi_cif cif;
    ffi_status status;
    int i;
    uint64_t zero = 0;

    ffi_type **arg_types = calloc(c + 2, sizeof(ffi_type *));
    void **arg_values = calloc(c + 2, sizeof(void *));

    arg_types[0] = &ffi_type_pointer;
    for (i = 1; i < c + 2; i++) {
        arg_types[i] = &ffi_type_uint64;
    }

    status = ffi_prep_cif_var(&cif, FFI_DEFAULT_ABI, 1, c+2, rtype, arg_types);

    if (status != FFI_OK) {
        return false;
    }

    arg_values[0] = &rights;
    for (i = 0; i < c; i++) {
        arg_values[i + 1] = &values[i];
    }
    arg_values[c + 1] = &zero;

    ffi_call(&cif, fn, rvalue, arg_values);

    free(arg_types);
    free(arg_values);

    return true;
}

cap_rights_t *doCapRightsPtr(void (*fn)(void), cap_rights_t *rights, uint64_t values[], int c) {
    cap_rights_t *result;

    if (!doCapRightsAny(fn, &ffi_type_pointer, &result, rights, values, c))
        return NULL;

    return result;
}

cap_rights_t *doCapRightsInit(cap_rights_t *rights, uint64_t values[], int c) {
    ffi_cif cif;
    ffi_status status;
    int i;
    uint64_t version = CAP_RIGHTS_VERSION;
    uint64_t zero = 0;
    cap_rights_t *result;

    ffi_type **arg_types = calloc(c + 3, sizeof(ffi_type *));
    void **arg_values = calloc(c + 3, sizeof(void *));

    arg_types[0] = &ffi_type_uint64;
    arg_types[1] = &ffi_type_pointer;
    for (i = 2; i < c + 3; i++) {
        arg_types[i] = &ffi_type_uint64;
    }

    status = ffi_prep_cif_var(&cif, FFI_DEFAULT_ABI, 2, c+3, &ffi_type_pointer, arg_types);

    if (status != FFI_OK) {
        return false;
    }

    arg_values[0] = &version;
    arg_values[1] = &rights;
    for (i = 0; i < c; i++) {
        arg_values[i + 2] = &values[i];
    }
    arg_values[c + 2] = &zero;

    ffi_call(&cif, FFI_FN(_cap_rights_init), &result, arg_values);

    free(arg_types);
    free(arg_values);

    return rights;
}

cap_rights_t *doCapRightsSet(cap_rights_t *rights, uint64_t values[], int c) {
    return doCapRightsPtr(FFI_FN(_cap_rights_set), rights, values, c);
}

cap_rights_t *doCapRightsClear(cap_rights_t *rights, uint64_t values[], int c) {
    return doCapRightsPtr(FFI_FN(_cap_rights_clear), rights, values, c);
}

bool doCapRightsIsSet(cap_rights_t *rights, uint64_t values[], int c) {
    unsigned ret;

    bool r = doCapRightsAny(FFI_FN(_cap_rights_is_set), &ffi_type_uint, &ret, rights, values, c);
    assert(r);

    return ret;
}

*/
import "C"
import "syscall"

const (
	ECAPMODE    = syscall.Errno(C.ECAPMODE)
	ENOTCAPABLE = syscall.Errno(C.ENOTCAPABLE)
)

const (
	CAP_CREATE = C.CAP_CREATE
	CAP_EVENT  = C.CAP_EVENT
	CAP_LISTEN = C.CAP_LISTEN
	CAP_LOOKUP = C.CAP_LOOKUP
	CAP_PDWAIT = C.CAP_PDWAIT
	CAP_READ   = C.CAP_READ
	CAP_WRITE  = C.CAP_WRITE
)

type CapRights C.struct_cap_rights

func CapEnter() error {
	ok, err := C.cap_enter()
	if ok == 0 {
		return nil
	}
	return err
}

func CapRightsInit(rights ...uint64) (*CapRights, error) {
	var r C.struct_cap_rights
	ret, err := C.doCapRightsInit(&r, (*C.uint64_t)(&rights[0]), C.int(len(rights)))
	return (*CapRights)(ret), err
}

// FIXME: should take a File, not an fd?
func CapRightsLimitFd(fd uintptr, r *CapRights) error {
	ret, err := C.cap_rights_limit(C.int(fd), (*C.struct_cap_rights)(r))
	if ret == 0 {
		return nil
	}
	return err
}

func CapRightsGetFd(fd uintptr) (*CapRights, error) {
	var r C.struct_cap_rights
	ret, err := C.cap_rights_get(C.int(fd), &r)
	if ret == 0 {
		return (*CapRights)(&r), nil
	}
	return nil, err
}

func CapRightsSet(r *CapRights, rights ...uint64) error {
	_, err := C.doCapRightsSet((*C.struct_cap_rights)(r), (*C.uint64_t)(&rights[0]), C.int(len(rights)))
	return err
}

func CapRightsClear(r *CapRights, rights ...uint64) error {
	_, err := C.doCapRightsClear((*C.struct_cap_rights)(r), (*C.uint64_t)(&rights[0]), C.int(len(rights)))
	return err
}

func CapRightsIsSet(r *CapRights, rights ...uint64) (bool, error) {
	ret, err := C.doCapRightsIsSet((*C.struct_cap_rights)(r), (*C.uint64_t)(&rights[0]), C.int(len(rights)))
	if err != nil {
		return false, err
	}
	return bool(ret), nil
}
