package domain

import "testing"

func TestValidateDeviceID(t *testing.T) {
	valid := []string{"oled-001", "OLED_2", "device.a"}
	invalid := []string{"", "../oled", "a/b", " device", "device?x=1"}

	for _, id := range valid {
		if err := ValidateDeviceID(id); err != nil {
			t.Fatalf("expected %q to be valid: %v", id, err)
		}
	}

	for _, id := range invalid {
		if err := ValidateDeviceID(id); err == nil {
			t.Fatalf("expected %q to be invalid", id)
		}
	}
}
