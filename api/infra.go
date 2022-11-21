package api

func checkNotNull(a any) {
	if a == nil {
		panic("value is null")
	}
}

func checkNotZeroInt(a int) {
	if a == 0 {
		panic("value is 0")
	}
}
