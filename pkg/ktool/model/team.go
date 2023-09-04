package model

// Team .
type Team struct {
	Model

	Players []struct {
		ID string `json:"_id"`
	}
}
