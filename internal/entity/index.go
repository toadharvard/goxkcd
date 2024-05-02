package entity

type Token = string

type Index interface {
	Search(Token) []int
}
