package recursivewalk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalk(t *testing.T) {
	type Sub2 struct {
		ID      int
		Hoge    string
		private string
	}
	type Sub1 struct {
		ID      int
		Hoge    string
		Sub2     Sub2
		private string
	}
	type Sample struct {
		ID        int
		Hoge      string
		Fuga      bool
		Sub       Sub1
		SubP      *Sub1
		SubSlice  []Sub1
		SubSliceP []*Sub1
		SubMap    map[string]Sub1
		SubMapP   map[string]*Sub1
		private   string
	}

	t.Run("ok", func(t *testing.T) {
		testData := Sample{
			ID:   1,
			Hoge: "a",
			Fuga: true,
			Sub: Sub1{
				ID:   3,
				Hoge: "c",
			},
			SubP: &Sub1{
				ID:   4,
				Hoge: "d",
			},
			SubSlice: []Sub1{
				{
					ID:   5,
					Hoge: "e",
				},
			},
			SubSliceP: []*Sub1{
				{
					ID:   6,
					Hoge: "f",
				},
			},
			SubMap: map[string]Sub1{
				"1": {
					ID:   7,
					Hoge: "g",
				},
			},
			SubMapP: map[string]*Sub1{
				"1": {
					ID:   8,
					Hoge: "h",
				},
			},
		}

		expectedFileds := map[string]interface{}{
			"ID":             1,
			"Hoge":           "a",
			"Fuga":           true,
			"Sub.ID":         3,
			"Sub.Hoge":       "c",
			"SubP.ID":        4,
			"SubP.Hoge":      "d",
			"SubSlice.ID":    5,
			"SubSlice.Hoge":  "e",
			"SubSliceP.ID":   6,
			"SubSliceP.Hoge": "f",
			"SubMap.ID":      7,
			"SubMap.Hoge":    "g",
			"SubMapP.ID":     8,
			"SubMapP.Hoge":   "h",
		}
		var count int
		err := Walk(testData, func(meta WalkMeta) {
			fmt.Println(meta)
			//val, ok := expectedFileds[meta.FieldPath]
			//assert.True(t, ok)
			//assert.Equal(t, val, meta.Value)
			count++
		})
		assert.NoError(t, err)
		assert.Equal(t, len(expectedFileds), count)
	})
}
