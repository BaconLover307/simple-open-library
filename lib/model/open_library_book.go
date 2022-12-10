package model

type OpenLibraryBook struct {
	Key               string   `json:"key"`
	Title             string   `json:"title"`
	EditionCount      int      `json:"edition_count"`
	// CoverID           int      `json:"cover_id"`
	// CoverEditionKey   string   `json:"cover_edition_key"`
	// Subject           []string `json:"subject"`
	// IaCollection      []string `json:"ia_collection"`
	// Lendinglibrary    bool     `json:"lendinglibrary"`
	// Printdisabled     bool     `json:"printdisabled"`
	// LendingEdition    string   `json:"lending_edition"`
	// LendingIdentifier string   `json:"lending_identifier"`
	Authors           []OpenLibraryAuthor `json:"authors"`
	// FirstPublishYear int    `json:"first_publish_year"`
	// Ia               string `json:"ia"`
	// PublicScan       bool   `json:"public_scan"`
	// HasFulltext      bool   `json:"has_fulltext"`
	// Availability     struct {
	// 	Status              string      `json:"status"`
	// 	AvailableToBrowse   bool        `json:"available_to_browse"`
	// 	AvailableToBorrow   bool        `json:"available_to_borrow"`
	// 	AvailableToWaitlist bool        `json:"available_to_waitlist"`
	// 	IsPrintdisabled     bool        `json:"is_printdisabled"`
	// 	IsReadable          bool        `json:"is_readable"`
	// 	IsLendable          bool        `json:"is_lendable"`
	// 	IsPreviewable       bool        `json:"is_previewable"`
	// 	Identifier          string      `json:"identifier"`
	// 	Isbn                interface{} `json:"isbn"`
	// 	Oclc                interface{} `json:"oclc"`
	// 	OpenlibraryWork     string      `json:"openlibrary_work"`
	// 	OpenlibraryEdition  string      `json:"openlibrary_edition"`
	// 	LastLoanDate        interface{} `json:"last_loan_date"`
	// 	NumWaitlist         interface{} `json:"num_waitlist"`
	// 	LastWaitlistDate    interface{} `json:"last_waitlist_date"`
	// 	IsRestricted        bool        `json:"is_restricted"`
	// 	IsBrowseable        bool        `json:"is_browseable"`
	// 	Src                 string      `json:"__src__"`
	// } `json:"availability"`
}