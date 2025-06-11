package types

type Precision struct {
    A float64
    B float64
}

type FieldType struct {
    name      string
    dataTaype string
    nullable  bool
    length    int
    precision Precision
}

type Indexes struct {
    name string
    ddl  string
}

type Table struct {
    name   string
    fields []FieldType
}