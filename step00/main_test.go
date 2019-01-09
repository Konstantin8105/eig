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
		_ = e
	})
	t.Run("2", func(t *testing.T) {
		e, err := pm([][]float64{
			{4, 5},
			{6, 5},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})

	t.Run("E3", func(t *testing.T) {
		e, err := pm([][]float64{
			{1, 0},
			{0, -1},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})
	t.Run("E4", func(t *testing.T) {
		e, err := pm([][]float64{
			{2, 0, 0},
			{0, 2, 0},
			{0, 0, 1},
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = e
	})

}
