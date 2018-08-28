package meta

type IAmount interface{
	getInt() int
	getFloat() float32
	getString() string
}
