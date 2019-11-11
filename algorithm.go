package git

import "math"

const (
	scalingFactorStandard = 0.1
	commonPrefixLimiter   = 4
)

// Get m variable
func mAndTVariable(r1, r2 []rune, r1Length, r2Length int) (float64, float64) {

	// match range = max(len(a), len(b)) / 2 - 1
	d := int(math.Floor(math.Max(float64(r1Length), float64(r2Length))/2.0)) - 1
	d = int(math.Max(0, float64(d-1)))
	var m, t float64
	tFlag := make([]bool, r2Length)

	for i := 0; i < r1Length; i++ {
		start := int(math.Max(0, float64(i-d)))
		end := int(math.Min(float64(r2Length-1), float64(i+d)))

		for j := start; j <= end; j++ {
			if tFlag[j] {
				continue
			}

			if r1[i] == r2[j] {
				if i != j {
					t = t + 1
				}
				m = m + 1
				tFlag[j] = true
				break
			}
		}
	}

	return m, t / 2
}

// Get l variable
func lVariable(r1, r2 []rune, r1Length int) float64 {
	var l int
	for i := 0; i < r1Length; i++ {
		if i > commonPrefixLimiter {
			return float64(l)
		}
		if r1[i] == r2[i] {
			if i == int(l) && i <= commonPrefixLimiter {
				l = l + 1
			}
		}
	}
	return float64(l)
}

// JaroWinkler function returns value between 0 and 1 range, estimating the similarity of 2 input strings (s1,s2)
func JaroWinkler(s1, s2 string, threshold float64) float64 {
	// Count input data (rune) length in any languages
	r1, r2 := []rune(s1), []rune(s2)
	if len(r1) > len(r2) {
		r1, r2 = r2, r1
	}

	r1Length, r2Length := len(r1), len(r2)
	if r1Length == 0 || r2Length == 0 {
		return 0
	}
	l := lVariable(r1, r2, r1Length)

	dwMax := 2.0/3 + float64(r1Length/(3*r2Length)) + float64(l*scalingFactorStandard/3)*(1-float64(r1Length/r2Length))
	if dwMax < threshold {
		return 0
	}
	m, t := mAndTVariable(r1, r2, r1Length, r2Length)

	// Format equations
	s1Eq, s2Eq, tEq := float64(r1Length), float64(r2Length), math.Floor(t)

	// Jaro distance formula
	dj := (m/s1Eq + m/s2Eq + (m-tEq)/m) / 3

	// Jaro and Winkler distance formula
	dw := dj + (l * scalingFactorStandard * (1 - dj))
	if math.IsNaN(dw) {
		return 0
	}
	return dw
}
