package extract

// RLCS Get a list of substrings whose length is greater than the minLength (without deduplication)
func RLCS(s1, s2 string, minLength int) (listLCS []string) {
	s1Len, s2Len := len(s1), len(s2)
	m := make([][]int, s1Len+1)
	for i := range m {
		m[i] = make([]int, s2Len+1)
	}

	longest := 0
	xLongest, yLongest := -1, -1
	for x := 1; x <= s1Len; x++ {
		for y := 1; y <= s2Len; y++ {
			if s1[x-1] == s2[y-1] {
				m[x][y] = m[x-1][y-1] + 1
				if m[x][y] >= longest {
					longest = m[x][y]
					xLongest = x
					yLongest = y
				}
			} else {
				m[x][y] = 0
			}
		}
	}
	if xLongest == -1 || yLongest == -1 {
		return listLCS
	}

	if xLongest > longest && yLongest > longest {
		listLCS = append(listLCS, RLCS(s1[:xLongest-longest], s2[:yLongest-longest], minLength)...)
	}
	sLCS := s1[xLongest-longest : xLongest]
	sLCSLen := len(sLCS)
	if sLCSLen > minLength {
		listLCS = append(listLCS, sLCS)
	}
	if xLongest < s1Len-1 && yLongest < s2Len-1 {
		listLCS = append(listLCS, RLCS(s1[xLongest-longest:s1Len-1], s2[yLongest-longest:s2Len-1], minLength)...)
	}
	return listLCS
}

// LASER wrong now. don't work
func LASER(s1, s2 []byte, threshold int) (lcs [][]byte) {
	reverse(s1)
	reverse(s2)
	s1Len, s2Len := len(s1), len(s2)
	m := make([][]int, s1Len)
	for i := range m {
		m[i] = make([]int, s2Len)
	}
	for i := range m {
		for j := range m[i] {
			if i == 0 || j == 0 {
				m[i][j] = 0
			} else if s1[i] == s2[j] {
				m[i][j] = 1
			} else if s1[i] != s2[j-1] || s1[i-1] != s2[j] {
				m[i][j] = 2
			} else {
				m[i][j] = 3
			}
		}
		// fmt.Println(m[i])
	}
	i, j := s1Len-1, s2Len-1
	var str []byte
	for m[i][j] != 0 {
		if m[i][j] == 3 {
			j--
		} else if m[i][j] == 2 {
			i--
		} else if m[i][j] == 1 {
			str = append(str, s1[i])
		}
		if m[i-1][j-1] != 1 && len(str) > threshold {
			t := make([]byte, len(str))
			copy(t, str)
			str = str[:0]
			lcs = append(lcs, t)
		}
		i--
		j--
	}
	return lcs
}

func reverse(s []byte) {
	for low, high := 0, len(s)-1; low < high; low, high = low+1, high-1 {
		s[low], s[high] = s[high], s[low]
	}
}
