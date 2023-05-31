package collectx

// JDiff j 交集， d1 a独有的， d2 b独有的
func JDiff(a []int64, b []int64) (j []int64, d1 []int64, d2 []int64) {
	j = make([]int64, 0)
	d1 = make([]int64, 0)
	d2 = make([]int64, 0)

	m := make(map[int64]bool)
	for _, i := range a {
		m[i] = false
	}

	for _, k := range b {
		_, ok := m[k]
		if ok {
			j = append(j, k)
			m[k] = true
		} else {
			d2 = append(d2, k)
		}
	}

	for k, v := range m {
		if !v {
			d1 = append(d1, k)
		}
	}

	return j, d1, d2
}

func JDiffString(a []string, b []string) (j []string, d1 []string, d2 []string) {
	j = make([]string, 0)
	d1 = make([]string, 0)
	d2 = make([]string, 0)

	m := make(map[string]bool)
	for _, i := range a {
		m[i] = false
	}

	for _, k := range b {
		_, ok := m[k]
		if ok {
			j = append(j, k)
			m[k] = true
		} else {
			d2 = append(d2, k)
		}
	}

	for k, v := range m {
		if !v {
			d1 = append(d1, k)
		}
	}

	return j, d1, d2
}
