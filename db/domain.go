package db

// Domain represents a FQDN used for routing incoming email and validating
// outgoing email.
type Domain struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func migrateDomainsTable(t *Token) error {
	_, err := t.exec(
		`
CREATE TABLE IF NOT EXISTS Domains (
    ID   SERIAL PRIMARY KEY,
    Name VARCHAR(80) NOT NULL UNIQUE
)
        `,
	)
	return err
}

// Domains retrieves a list of all domains in the database.
func Domains(t *Token) ([]*Domain, error) {
	r, err := t.query(
		`
SELECT ID, Name
FROM Domains ORDER BY Name
        `,
	)
	if err != nil {
		return nil, err
	}
	domains := make([]*Domain, 0, 1)
	for r.Next() {
		d := &Domain{}
		if err := r.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, nil
}

// Save persists changes to the domain. If ID is set to zero, a new domain is
// created and its ID updated.
func (d *Domain) Save(t *Token) error {
	if d.ID == 0 {
		err := t.queryRow(
			`
INSERT INTO Domains (Name)
VALUES ($1)
            `,
			d.Name,
		).Scan(&d.ID)
		if err != nil {
			return err
		}
		return nil
	} else {
		_, err := t.exec(
			`
UPDATE Domains SET Name=$1
WHERE ID = $2
            `,
			d.Name,
			d.ID,
		)
		return err
	}
}

// Delete the domain from the database.
func (d *Domain) Delete(t *Token) error {
	_, err := t.exec(
		`
DELETE FROM Domains WHERE ID = $1
        `,
		d.ID,
	)
	return err
}
