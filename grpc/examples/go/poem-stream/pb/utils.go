package pb

import (
	"fmt"
)

func Serialize(p *Poem) string {
	s := fmt.Sprintf("%s\n%s\n", p.GetTitle(), p.GetAuthor())
	for _, content := range p.GetContents() {
		s += fmt.Sprintf("%s\n", content)
	}
	return s
}
