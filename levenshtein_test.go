package levenshtein

import (
	"testing"
	"reflect"
)

func TestSanity(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"", "hello", 5},
		{"hello", "", 5},
		{"hello", "hello", 0},
		{"ab", "aa", 1},
		{"ab", "aaa", 2},
		{"bbb", "a", 3},
		{"kitten", "sitting", 3},
		{"distance", "difference", 5},
		{"levenshtein", "frankenstein", 6},
		{"resume and cafe", "resumes and cafes", 2},
	}
	for i, d := range tests {
		n := ComputeDistance([]rune(d.a), []rune(d.b))
		if n != d.want {
			t.Errorf("Test[%d]: ComputeDistance(%q,%q) returned %v, want %v",
				i, d.a, d.b, n, d.want)
		}
	}
}

func TestUnicode(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		// Testing acutes and umlauts
		{"resumé and café", "resumés and cafés", 2},
		{"resume and cafe", "resumé and café", 2},
		{"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson", 4},
		// Only 2 characters are less in the 2nd string
		{"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ", 2},
	}
	for i, d := range tests {
		n := ComputeDistance([]rune(d.a), []rune(d.b))
		if n != d.want {
			t.Errorf("Test[%d]: ComputeDistance(%q,%q) returned %v, want %v",
				i, d.a, d.b, n, d.want)
		}
	}
}

func eqStats(e1, e2 EditStats) bool {
	if !reflect.DeepEqual(e1.Subs, e2.Subs)	{
		return false
	} else if !reflect.DeepEqual(e1.Ins, e2.Ins)	{
		return false
	} else if !reflect.DeepEqual(e1.Dels, e2.Dels)	{
		return false
	}
	return true
}

func TestReconstruction(t *testing.T) {
	tests := []struct {
		a, b string
		wantScore int
		wantStats EditStats
	} {
		{"", 
		 "hgghg",
		 5,
		 EditStats{
			 Ins : map[string]int {"h":2, "g":3},
			 Dels : map[string]int{},
			 Subs : map[string]int{},
		 },
		},
		{"hgghg",
		 "",
	 	 5,
		 EditStats{
			 Dels : map[string]int {"h":2, "g":3},
			 Ins : map[string]int{},
			 Subs : map[string]int{},
		 }, 
	 	},
		{"hello",
		 "heIIa",
		 3,
		 EditStats{
			 Ins : map[string]int{},
			 Dels : map[string]int{},			 
			 Subs : map[string]int {"lI":2, "oa":1},
		 }, 
	   },
	   {"heahhhllo",
		"hello",
		4,
		EditStats{
			Ins : map[string]int{},
			Dels : map[string]int{"a":1, "h":3},			 
			Subs : map[string]int {},
		}, 
	   },
	   {"hello",
		"nnnnnnnhello",
		7,
		EditStats{
			Ins : map[string]int{"n":7},
			Dels : map[string]int{},			 
			Subs : map[string]int {},
		}, 
	   },
	}
	for i, d := range tests {
		n, e := ComputeDistanceWithConstruction([]rune(d.a), []rune(d.b))
		if n != d.wantScore {
			t.Errorf("Test[%d]: ComputeDistance(%q,%q) returned %v, want %v",
				i, d.a, d.b, n, d.wantScore)
		} else if !eqStats(d.wantStats, e) {
			t.Errorf("Test[%d]: ComputeDistance(%q,%q)",i, d.a, d.b)
			t.Log("Want: ", d.wantStats)
			t.Log("Got: ", e)			
		}
	}
	
}
