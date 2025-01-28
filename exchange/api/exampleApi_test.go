package api

import "testing"

func TestCallPokeapi(t *testing.T) {
	response := CallPokeapi()
	if response.Name != "kanto" {
		t.Errorf("Expected name to be kanto, but got %s", response.Name)
	}
}
