package model

// CompanyFactUpsertedEvent is published to the "company-fact-upserted" Dapr pub/sub topic
// after every successful upsert of a company_fact row.
// The full CompanyFact is embedded so downstream consumers (e.g. eiker-company-es)
// can synchronise their indexes without an additional database query.
type CompanyFactUpsertedEvent struct {
	// Action is "create" when a new row was inserted, or "update" when an existing row was modified.
	Action string `json:"action"`
	// Fact is the complete company_fact entity as it exists in the database after the upsert.
	Fact *CompanyFact `json:"fact"`
}
