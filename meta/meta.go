package meta

type IAmount interface{
	GetInt() int
	GetFloat() float32
	GetString() string
}
