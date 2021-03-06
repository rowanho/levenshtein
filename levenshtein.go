// Package levenshtein is a Go implementation to calculate Levenshtein Distance.
//
// Implementation taken from
// https://gist.github.com/andrei-m/982927#gistcomment-1931258
package levenshtein


// ComputeDistance computes the levenshtein distance between the two
// strings passed as an argument. The return value is the levenshtein distance
//
// Works on runes (Unicode code points) but does not normalize
// the input strings. See https://blog.golang.org/normalization
// and the golang.org/x/text/unicode/norm pacage.
func ComputeDistance(s1, s2 []rune) int {
	if len(s1) == 0 {
		return len(s2)
	}

	if len(s2) == 0 {
		return len(s1)
	}


	lenS1 := len(s1)
	lenS2 := len(s2)

	// init the row
	x := make([]uint16, lenS1+1)
	// we start from 1 because index 0 is already 0.
	for i := 1; i < len(x); i++ {
		x[i] = uint16(i)
	}

	// make a dummy bounds check to prevent the 2 bounds check down below.
	// The one inside the loop is particularly costly.
	_ = x[lenS1]
	// fill in the rest
	for i := 1; i <= lenS2; i++ {
		prev := uint16(i)
		var current uint16
		for j := 1; j <= lenS1; j++ {
			if s2[i-1] == s1[j-1] {
				current = x[j-1] // match
			} else {
				current = min(min(x[j-1]+1, prev+1), x[j]+1)
			}
			x[j-1] = prev
			prev = current
		}
		x[lenS1] = prev
	}
	return int(x[lenS1])
}

// EditStats stores information about the number of substitutions, insertions and deletions 
// used in the optimal path between the two string
type EditStats = struct {
	Subs map[string]int `json:"subs"`
	Ins map[string]int  `json:"ins"`
	Dels map[string]int `json:"dels"`
}

// NewEditStats returns a new EditStats object
func NewEditStats() EditStats {
	return EditStats{
		Subs: map[string]int{},
		Ins: map[string]int{},
		Dels: map[string]int{},
	}
}

// ComputeDistanceWithConstruction extends ComputeDistance to return 
// information about the number of substitutions, insertions and deletions
// alongside the distance score
func ComputeDistanceWithConstruction(s1, s2 []rune) (int, EditStats) {
	
	if len(s1) == 0 {
		e := NewEditStats()
		for _, c := range s2 {
			e.Ins[string(c)] += 1
		}
		return len(s2), e
	}

	if len(s2) == 0 {
		e := NewEditStats()	
		for _, c := range s1 {
			e.Dels[string(c)] += 1
		}
		return len(s1), e
	}

	lenS1 := len(s1)
	lenS2 := len(s2)
	
	d := make([][]uint16, lenS1 + 1)
	for i := 0 ; i < lenS1 + 1; i++ {
		d[i] = make([]uint16, lenS2 +1)
		d[i][0] = uint16(i)
	}
	
	for j := 0 ; j < lenS2 + 1; j++ {
		d[0][j] = uint16(j)
	}
	
	var s uint16
	for j := 1; j < lenS2 + 1; j++ {
		for i := 1; i < lenS1 + 1; i++ {
			if s1[i-1] == s2[j-1] {
				s = 0
			} else {
				s = 1
			}
			d[i][j] = min(min(d[i-1][j] + 1, d[i][j-1] + 1), d[i-1][j-1] + s)
		}
	}
	
	return int(d[lenS1][lenS2]), reconstruct(d, s1, s2)

}

func reconstruct(d [][]uint16, s1, s2 []rune) EditStats {
	e := NewEditStats()
	i := len(s1)
	j := len(s2)
	var s uint16
	
	for i > 0 && j > 0 {
		if s1[i-1] == s2[j-1] {
			s = 0
		} else {
			s = 1
		}
		if d[i-1][j-1] + s <= min(d[i-1][j] + 1, d[i][j-1] + 1){
			if s == 1 {
				// Mismatch substitution
				e.Subs[string(s1[i-1]) + string(s2[j-1])] += 1
			}
			j -= 1
			i -= 1
		} else if d[i-1][j] + 1 <= d[i][j-1] + 1 {
			e.Dels[string(s1[i-1])] += 1			
			i -=  1
		} else {
			e.Ins[string(s2[j-1])] += 1			
			j -= 1
		}
	}
	
	if i > 0 {
		for k := i; k > 0; k-- {
			e.Dels[string(s1[k-1])] += 1						
		} 
	}
	
	if j > 0 {
		for k := j; j > 0; j-- {
			e.Ins[string(s2[k-1])] += 1						
		} 		
	}
	
	return e
}

func min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}
