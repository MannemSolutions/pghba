package hba

type Line interface {
    String()   string
    Comments() Comments
    Bare()     string
    Less()     bool
    Compare()  int
}

type Lines []Line
