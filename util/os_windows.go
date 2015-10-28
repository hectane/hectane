package util

import (
	"golang.org/x/sys/windows"

	"unsafe"
)

const (
	errorSuccess = 0

	seFileObject                     = 1
	daclSecurityInformation          = 0x4
	protectedDaclSecurityInformation = 0x80000000

	genericAll                     = 0x10000000
	grantAccess                    = 1
	subContainersAndObjectsInherit = 0x3

	trusteeIsSID = 0
)

var (
	advapi32              = windows.MustLoadDLL("advapi32.dll")
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

// Ensure that the specified path is only accessible to the system user. This
// is a rather tedious operation that begins with obtaining a PSID for the
// LocalSystem user. An EXPLICIT_ACCESS entry is created for the user granting
// them GENERIC_ALL access. This then becomes the single entry in the path's
// ACL, preventing all other users from accessing it.
func SecurePath(path string) error {
	var pSID *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		1,
		windows.SECURITY_LOCAL_SYSTEM_RID,
		0, 0, 0, 0, 0, 0, 0,
		&pSID,
	)
	if err != nil {
		return err
	}
	defer windows.FreeSid(pSID)
	var (
		ea = []explicitAccessW{
			{
				grfAccessPermissions: genericAll,
				grfAccessMode:        grantAccess,
				grfInheritance:       subContainersAndObjectsInherit,
				Trustee: trustee{
					TrusteeForm: trusteeIsSID,
					ptstrName:   unsafe.Pointer(pSID),
				},
			},
		}
		pACL windows.Handle
	)
	ret, _, err := setEntriesInAclW.Call(
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

// Retrieve the full path to the current executable.
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
