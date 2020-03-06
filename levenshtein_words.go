package levenshtein

import (
	"hash/fnv"
	"strings"
)


func hash(s string) uint64 {
        h := fnv.New64()
        h.Write([]byte(s))
        return h.Sum64()
}

func ComputeDistance64(s1, s2 []uint64) int {
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


func ComputeDistWithCon64(s1, s2 []uint64, r1, r2 map[uint64]string) (int, EditStats) {

	if len(s1) == 0 {
		e := NewEditStats()
		for _, c := range s2 {
			e.Ins[r1[c]] += 1
		}
		return len(s2), e
	}

	if len(s2) == 0 {
		e := NewEditStats()	
		for _, c := range s1 {
			e.Dels[r2[c]] += 1
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
	
	return int(d[lenS1][lenS2]), reconstruct64(d, s1, s2, r1, r2)

}

func reconstruct64(d [][]uint16, s1, s2 []uint64, r1, r2 map[uint64]string) EditStats {
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
				e.Subs[r1[s1[i-1]] + r2[s2[j-1]]] += 1
			}
			j -= 1
			i -= 1
		} else if d[i-1][j] + 1 <= d[i][j-1] + 1 {
			e.Dels[r1[s1[i-1]]] += 1			
			i -=  1
		} else {
			e.Ins[r2[s2[j-1]]] += 1			
			j -= 1
		}
	}
	
	if i > 0 {
		for k := i; k > 0; k-- {
			e.Dels[r1[s1[k-1]]] += 1						
		} 
	}
	
	if j > 0 {
		for k := j; j > 0; j-- {
			e.Ins[r2[s2[k-1]]] += 1						
		} 		
	}
	
	return e
}


func ComputeWordDistance(s1, s2 []rune) (int) {
	words1 := strings.Fields(string(s1))
	words2 := strings.Fields(string(s2))
	
	hashes1 := make([]uint64, len(words1))
	hashes2 := make([]uint64, len(words2))
	
	for i, w := range words1 {
		hashes1[i] = hash(w)
	}
	
	for i, w := range words2 {
		hashes2[i] = hash(w)
	}
	
	return ComputeDistance64(hashes1, hashes2)
}


func ComputeWordDistCon(s1, s2 []rune) (int, EditStats) {
    words1 := strings.Fields(string(s1))
	words2 := strings.Fields(string(s2))
	
	hashes1 := make([]uint64, len(words1))
	hashes2 := make([]uint64, len(words2))
	
    revDict1 := make(map[uint64]string)
    revDict2 := make(map[uint64]string)
    
	for i, w := range words1 {
		hashes1[i] = hash(w)
        revDict1[hashes1[i]] = w
	}
	
	for i, w := range words2 {
		hashes2[i] = hash(w)
        revDict2[hashes2[i]] = w
	}
    
    return ComputeDistWithCon64(hashes1, hashes2, revDict1, revDict2)
}





