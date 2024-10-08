package req

type Decoder interface {
	Decode(interface{}) error
}
