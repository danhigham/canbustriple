package main

func padRight(str, pad string, length int) string {
    for {
        str += pad
        if len(str) > length {
            return str[0:length]
        }
    }
}

func padLeft(str, pad string, length int) string {
	for {
		str = pad + str
		if len(str) > length {
			return str[0:length+1]
		}
	}
}
