package receiver

import (
	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
)

// getFolder attempts to retrieve the specified folder for the specified user,
// creating it if it does not exist.
func getFolder(c *gorm.DB, userID int64, name string) (*db.Folder, error) {
	var (
		f   = &db.Folder{}
		err = c.FirstOrInit(f, &db.Folder{
			Name:   name,
			UserID: userID,
		}).Error
	)
	if err != nil {
		return nil, err
	}
	return f, nil
}
