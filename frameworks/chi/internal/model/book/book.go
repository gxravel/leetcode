package book

import "database/sql"

const limit = 100

type Language uint

const (
	LangEN Language = iota
	LangRU
	LangFR
)

type Contributor struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

type Edition struct {
	ID    int    `json:"ID"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Book struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	WrittenAt    int           `json:"written_at,omitempty"`
	Language     Language      `json:"lang,omitempty"`
	Contributors []Contributor `json:"contributors"`
	Edition      Edition       `json:"edition"`
}

type Wrapper struct {
	db *sql.DB
}

func New(db *sql.DB) Wrapper {
	return Wrapper{db: db}
}

func (w Wrapper) All() ([]Book, error) {
	bRows, err := w.db.Query(`select b.id, b.title, b.written_at, b.lang, e.id, e.name, e.image
	from book as b
	inner join edition as e on e.id=b.default_edition
	limit by ?`, limit)
	if err != nil {
		return nil, err
	}
	defer bRows.Close()
	var result []Book
	for bRows.Next() {
		var b = Book{}
		var e = Edition{}
		err := bRows.Scan(&b.ID, &b.Title, &b.WrittenAt, &b.Language, &e.ID, &e.Name, &e.Image)
		if err != nil {
			return nil, err
		}
		b.Edition = e
		if err != nil {
			return nil, err
		}
		b.Contributors, err = w.contributorsByID(b.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	return result, nil
}

func (w Wrapper) ByID(id int) (*Book, error) {
	row := w.db.QueryRow(`select b.id, b.title, b.written_at, b.lang, e.id, e.name, e.image
	from book as b
	inner join edition as e on e.id=b.default_edition
	where b.id=?`, id)
	var b = &Book{}
	var e = Edition{}
	err := row.Scan(&b.ID, &b.Title, &b.WrittenAt, &b.Language, &e.ID, &e.Name, &e.Image)
	if err != nil {
		return nil, err
	}
	b.Edition = e
	b.Contributors, err = w.contributorsByID(b.ID)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (w Wrapper) contributorsByID(id int) ([]Contributor, error) {
	cRows, err := w.db.Query(`select c.id, c.name, c.type
	from contributor as c
	inner join contributor_book as cb on c.id=cb.contributor_id
	inner join book as b on cb.book_id=b.id
	where b.id=?`, id)
	if err != nil {
		return nil, err
	}
	var result []Contributor
	defer cRows.Close()
	for cRows.Next() {
		var c = Contributor{}
		err = cRows.Scan(&c.ID, &c.Name, &c.Type)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}
