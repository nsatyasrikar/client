package client

import "testing"

type embedsMarker struct{ TypeMarker }

func TestTypeMarkerSatisfiesType(t *testing.T) {
	var _ Type = embedsMarker{}
}
