package main

import "testing"

func Test(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		e, err := pm([][]float64{
			{2, -12},
			{1, -5},
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(e)
	})
}
