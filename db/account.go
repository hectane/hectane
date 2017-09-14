package db

import (
	"strconv"

	"github.com/manyminds/api2go/jsonapi"
)

const (
	accountUserRelationship   = "user"
	accountDomainRelationship = "domain"
)

// Account represents an individual email account owned by a user.
type Account struct {
	ID       int64   `json:"-"`
	Name     string  `json:"name" gorm:"type:varchar(40);not null"`
	User     *User   `json:"-" gorm:"ForeignKey:UserID"`
	UserID   int64   `json:"-" sql:"type:int REFERENCES users(id)"`
	Domain   *Domain `json:"-" gorm:"ForeignKey:DomainID"`
	DomainID int64   `json:"-" sql:"type:int REFERENCES domains(id)"`
}

func (a *Account) GetName() string {
	return "accounts"
}

func (a *Account) GetID() string {
	return strconv.FormatInt(a.ID, 10)
}

func (a *Account) SetID(id string) error {
	a.ID, _ = strconv.ParseInt(id, 10, 64)
	return nil
}

func (a *Account) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type:         a.Domain.GetName(),
			Name:         accountDomainRelationship,
			Relationship: jsonapi.ToOneRelationship,
		},
	}
}

func (a *Account) GetReferencedIDs() []jsonapi.ReferenceID {
	return []jsonapi.ReferenceID{
		{
			ID:           a.Domain.GetID(),
			Type:         a.Domain.GetName(),
			Name:         accountDomainRelationship,
			Relationship: jsonapi.ToOneRelationship,
		},
	}
}

func (a *Account) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	return []jsonapi.MarshalIdentifier{
		a.Domain,
	}
}

func (a *Account) SetToOneReferenceID(name, id string) error {
	switch name {
	case accountUserRelationship:
		a.UserID, _ = strconv.ParseInt(id, 10, 64)
	case accountDomainRelationship:
		a.DomainID, _ = strconv.ParseInt(id, 10, 64)
	default:
		return ErrInvalidRelationship
	}
	return nil
}
