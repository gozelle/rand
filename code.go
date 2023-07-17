package rand

import "fmt"

func Code(length int) (code string) {
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", Intn(9))
	}
	return
}
