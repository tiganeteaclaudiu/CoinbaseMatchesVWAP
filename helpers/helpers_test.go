package helpers

import "testing"

func TestGetSubscribeToMatchesMessage(t *testing.T) {
	expectedResult := `
	{
	   "type":"subscribe",
	   "channels":[
		  {
			 "name":"matches",
			 "product_ids":["BTC-USD","ETH-USD","ETH-BTC"]
		  }
	   ]
	}
`
	pairs := []string{"BTC-USD", "ETH-USD", "ETH-BTC"}
	result := GetSubscribeToMatchesMessage(pairs)
	if result != expectedResult {
		t.Errorf("got %s, expected %s", result, expectedResult)
	}
}

func TestGetMaxFloat(t *testing.T) {
	floats := []float64{1.23, 32, 1, 4, 2.32}
	result := GetMaxFloat(floats)
	if result != 32 {
		t.Errorf("got %f, expected 32", result)
	}
}

func TestGetMinFloat(t *testing.T) {
	floats := []float64{1.23, 32, 1, 4, 2.32}
	result := GetMinFloat(floats)
	if result != 1 {
		t.Errorf("got %f, expected 1", result)
	}
}
