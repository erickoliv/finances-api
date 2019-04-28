package model

import "testing"

func TestBaseModel_IsNew(t *testing.T) {
	new := BaseModel{}

	if !new.IsNew() {
		t.Errorf("A zero-value model should return true.")
	}
}
