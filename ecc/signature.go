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

func (s *Signature) Der() []byte {
	encodedSection := make([]byte, 0)

	toEncode := [][]byte{s.r.num.Bytes(), s.s.num.Bytes()}
	for _, field := range toEncode {
		encodedSection = append(encodedSection, 0x02)
		if field[0] >= 0x80 {
			encodedSection = append(encodedSection, byte(len(field)+1))
			encodedSection = append(encodedSection, 0x00)
		} else {
			encodedSection = append(encodedSection, byte(len(field)))
		}
		encodedSection = append(encodedSection, field...)
	}

	encodedSection = append([]byte{0x30, byte(len(encodedSection))}, encodedSection...)
	return encodedSection
}
