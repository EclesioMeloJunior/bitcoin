package ecc

import "fmt"

type Signature struct {
	r *FieldElement
	s *FieldElement
}

func NewSignature(r, s *FieldElement) *Signature {
	return &Signature{r, s}
}

func (s *Signature) String() string {
	return fmt.Sprintf("Sig(r: {%s}, s: {%s})", s.r, s.s)
}
