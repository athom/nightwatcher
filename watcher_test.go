package nightwatcher

import (
	"strings"
	"testing"
	"time"
)

func TestGlance(t *testing.T) {
	tests := []struct {
		in string
		//shouldRaiseError bool
		wantStatusCode int
	}{
		{"http://github.com", 301},
		{"https://github.com", 200},
		{"https://github.com/athom", 200},
		{"https://github.com/athom/notfoundproject", 404},
		{"github.com", -1},
		{"https://super-strange-address-must-not-exist.com", -1},
		{"https://console.cloud.tencent.com", 302},
		{"http://119.28.1.61/", 500},
	}

	for _, test := range tests {
		aim, err := glance(test.in)
		if test.wantStatusCode == -1 {
			if err == nil {
				t.Fatalf("visiting %v should raise error but no", test.in)
			}
			return
		}

		if err != nil && test.wantStatusCode != -1 {
			t.Fatalf("visting %v raise error: %v", test.in, err)
			return
		}

		if aim.StatusCode != test.wantStatusCode {
			t.Errorf("visting %v, status code want %v, but got %v", test.in, test.wantStatusCode, aim.StatusCode)
			return
		}

		if strings.Contains(test.in, `https`) {
			if aim.Result.TLSHandshake <= 0 {
				t.Errorf("visting %v should have TLS hanksake, but no", test.in)
				return
			}
		}

		checkResult(t, aim)
	}

}

func checkResult(t *testing.T, aim *Aim) {
	url := aim.TargetURL
	for k, d := range aim.Durations() {
		if k == "TLSHandshake" {
			if d < 0*time.Millisecond {
				t.Fatalf("visiting %v, expect %s to be non-zero, but got %v", url, k, d)
			}
			continue
		}

		if d <= 0*time.Millisecond {
			t.Fatalf("visiting %v, expect %s to be non-zero, but got %v", url, k, d)
		}
	}
}
