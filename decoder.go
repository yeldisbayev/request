package request

type Decoder interface {
	Decode(interface{}) error
}
