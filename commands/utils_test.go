package commands

import "testing"

func TestWhich(t *testing.T) {

	_, err := Which("unknowncommand")
	if err == nil {
		t.Errorf("got %s did not want that", err)
	}

}
