package main


import "fmt"

type ListOpts []string
// The string flag list, implemented flag.Value interface
type StrFlags []string

func (opts *ListOpts) String() string {
	return fmt.Sprint(*opts)
}

func (opts *ListOpts) Set(value string) error {
	*opts = append(*opts, value)
	return nil
}




func (s *StrFlags) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *StrFlags) Set(value string) error {
	*s = append(*s, value)
	return nil
}