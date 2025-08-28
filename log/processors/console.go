package processors

import (
	"fmt"
)

type Console struct{}

func (p *Console) Process(record []byte) {
	fmt.Println(string(record))
}
