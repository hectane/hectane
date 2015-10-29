package util

import (
	"golang.org/x/sys/windows"

	"unsafe"
)

const (
	errorSuccess = 0

	securityMaxSidSize = 68

	winServiceSid               = 12
	winBuiltinAdministratorsSid = 26

	genericAll                     = 0x10000000
	grantAccess                    = 1
	subContainersAndObjectsInherit = 0x3

	seFileObject                     = 1
	daclSecurityInformation          = 0x4
	protectedDaclSecurityInformation = 0x80000000

	trusteeIsSID = 0
)

var (
	advapi32              = windows.MustLoadDLL("advapi32.dll")
	createWellKnownSid    = advapi32.MustFindProc("CreateWellKnownSid")
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

// Create a SID for the specified user or group.
func createSid(sidType int32) ([]byte, error) {
	var (
		sid    = make([]byte, securityMaxSidSize)
		sidLen = uint32(securityMaxSidSize)
	)
	ret, _, err := createWellKnownSid.Call(
		uintptr(sidType),
		uintptr(0),
		uintptr(unsafe.Pointer(&sid[0])),
		uintptr(unsafe.Pointer(&sidLen)),
	)
	if ret == 0 {
		return nil, err
	} else {
		return sid, nil
	}
}

// Ensure that the specified path is only accessible to the administrator
// builtin group. This is a rather tedious operation that begins with obtaining
// a PSID for the LocalSystem user. An EXPLICIT_ACCESS entry is created for the
// user granting them GENERIC_ALL access. This then becomes the single entry in
// the path's ACL, preventing all other users from accessing it.
func SecurePath(path string) error {
	var (
		sidService []byte
		sidAdmins  []byte
	)
	sidService, err := createSid(winServiceSid)
	if err != nil {
		return err
	}
	sidAdmins, err = createSid(winBuiltinAdministratorsSid)
	if err != nil {
		return err
	}
	var (
		ea = []explicitAccessW{
			{
				grfAccessPermissions: genericAll,
				grfAccessMode:        grantAccess,
				grfInheritance:       subContainersAndObjectsInherit,
				Trustee: trustee{
					TrusteeForm: trusteeIsSID,
					ptstrName:   unsafe.Pointer(&sidService[0]),
				},
			},
			{
				grfAccessPermissions: genericAll,
				grfAccessMode:        grantAccess,
				grfInheritance:       subContainersAndObjectsInherit,
				Trustee: trustee{
					TrusteeForm: trusteeIsSID,
					ptstrName:   unsafe.Pointer(&sidAdmins[0]),
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
