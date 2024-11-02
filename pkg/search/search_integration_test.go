package search

// import (
// 	"fmt"
// 	"testing"
// )

// // TestExtractBaseCurrency checks the functionality of extractBaseCurrency
// func TestExtractBaseCurrency(t *testing.T) {
// 	tests := []struct {
// 		filename string
// 		expected string
// 	}{
// 		{"LANCER__XNSE_INR.csv", "INR"},
// 		{"ZZZ_USDT.csv", "USD"},
// 		{"LAU__XASX_AUD.csv", "AUD"},
// 		{"ZENI__1_USDT.csv", "USD"},
// 		{"LGBBROSLTD__XBOM_INR.csv", "INR"},
// 	}

// 	for _, test := range tests {
// 		got := extractBaseCurrency(test.filename)
// 		fmt.Printf("Filename: %s, Extracted Currency: %s\n", test.filename, got)
// 		if got != test.expected {
// 			t.Errorf("extractBaseCurrency(%q) = %q; expected %q", test.filename, got, test.expected)
// 		}
// 	}
// }
