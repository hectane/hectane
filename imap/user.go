package imap

import (
	"strings"

	"github.com/emersion/go-imap/backend"
	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
)

// user maintains user information for an IMAP session.
type user struct {
	imap *IMAP
	user *db.User
}

// Username returns the user's username.
func (u *user) Username() string {
	return u.user.Username
}

// ListMailboxes lists all of the user's folders.
func (u *user) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
	folders := []*db.Folder{}
	if err := db.C.Where("user_id = ?", u.user.ID).Find(&folders).Error; err != nil {
		return nil, err
	}
	m := make([]backend.Mailbox, len(folders))
	for i, f := range folders {
		m[i] = &mailbox{
			imap:   u.imap,
			folder: f,
		}
	}
	return m, nil
}

// GetMailbox retrieves a mailbox by name.
func (u *user) GetMailbox(name string) (backend.Mailbox, error) {
	create := strings.ToLower(name) == strings.ToLower(db.FolderInbox)
	f, err := db.GetFolderTransaction(u.user.ID, name, create)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, backend.ErrNoSuchMailbox
		}
		return nil, err
	}
	return &mailbox{
		imap:   u.imap,
		folder: f,
	}, nil
}

// CreateMailbox creates a new folder.
func (u *user) CreateMailbox(name string) error {
	f := &db.Folder{
		Name:   name,
		UserID: u.user.ID,
	}
	return db.C.Create(f).Error
}

// DeleteMailbox permanently deletes a folder.
func (u *user) DeleteMailbox(name string) error {
	return db.C.
		Where("user_id = ?", u.user.ID).
		Where("name = ?", name).
		Delete(&db.Folder{}).
		Error
}

// RenameMailbox attempts to change the name of a mailbox.
func (u *user) RenameMailbox(existingName, newName string) error {
	return db.C.
		Model(&db.Folder{}).
		Where("user_id = ?", u.user.ID).
		Where("name = ?", existingName).
		Updates(map[string]interface{}{
			"name": newName,
		}).
		Error
}

// Logout doesn't do anything.
func (u *user) Logout() error {
	return nil
}
