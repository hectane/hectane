package util

import (
	"golang.org/x/sys/windows"

	"unsafe"
)

const (
	errorSuccess = 0

	seFileObject                     = 1
	ownerSecurityInformation         = 0x1
	daclSecurityInformation          = 0x4
	protectedDaclSecurityInformation = 0x80000000

	genericAll    = 0x10000000
	grantAccess   = 1
	noInheritance = 0

	trusteeIsSID = 0
)

var (
	advapi32              = windows.MustLoadDLL("advapi32.dll")
	getNamedSecurityInfoW = advapi32.MustFindProc("GetNamedSecurityInfoW")
	setEntriesInAclW      = advapi32.MustFindProc("SetEntriesInAclW")
	setNamedSecurityInfoW = advapi32.MustFindProc("SetNamedSecurityInfoW")

	kernel32           = windows.MustLoadDLL("kernel32.dll")
	getModuleFileNameW = kernel32.MustFindProc("GetModuleFileNameW")
)

type trustee struct {
	pMultipleTrustee         unsafe.Pointer
	MultipleTrusteeOperation int32
	TrusteeForm              int32
	TrusteeType              int32
	ptstrName                unsafe.Pointer
}

type explicitAccessW struct {
	grfAccessPermissions uint32
	grfAccessMode        int32
	grfInheritance       uint32
	Trustee              trustee
}

// Ensure that the specified path is only accessible to the current user. This
// is accomplished by setting a single ACE in the path's ACL for the owner.
func SecurePath(path string) error {
	var (
		pSID windows.Handle
		pSD  windows.Handle
	)
	ret, _, err := getNamedSecurityInfoW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(seFileObject),
		uintptr(ownerSecurityInformation),
		uintptr(unsafe.Pointer(&pSID)),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(unsafe.Pointer(&pSD)),
	)
	// TODO: for some reason the error isn't returned
	if ret != errorSuccess {
		return err
	}
	defer windows.LocalFree(pSD)
	var (
		ea = []explicitAccessW{
			{
				grfAccessPermissions: genericAll,
				grfAccessMode:        grantAccess,
				grfInheritance:       noInheritance,
				Trustee: trustee{
					TrusteeForm: trusteeIsSID,
					ptstrName:   unsafe.Pointer(pSID),
				},
			},
		}
		pACL windows.Handle
	)
	ret, _, err = setEntriesInAclW.Call(
		uintptr(len(ea)),
		uintptr(unsafe.Pointer(&ea[0])),
		uintptr(0),
		uintptr(unsafe.Pointer(&pACL)),
	)
	if ret != errorSuccess {
		return err
	}
	defer windows.LocalFree(pACL)
	ret, _, err = setNamedSecurityInfoW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(seFileObject),
		uintptr(daclSecurityInformation|protectedDaclSecurityInformation),
		uintptr(0),
		uintptr(0),
		uintptr(pACL),
		uintptr(0),
	)
	if ret != errorSuccess {
		return err
	}
	return nil
}

// Retrieve the full path to the current executable using the Windows API.
func Executable() (string, error) {
	s := make([]uint16, windows.MAX_PATH)
	ret, _, err := getModuleFileNameW.Call(
		0,
		uintptr(unsafe.Pointer(&s[0])),
		uintptr(len(s)),
	)
	if ret != 0 {
		return windows.UTF16ToString(s[:ret]), nil
	} else {
		return "", err
	}
}
