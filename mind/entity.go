package mind

type Entity struct {
	ID           string
	BytesInput   int
	BytesOutput  int
	IsConnected  bool
	Base64Input  string
	Base64Output string
}

func NewEntity(entityID string, bytesInput, bytesOutput int) *Entity {
	return &Entity{
		ID:          entityID,
		BytesInput:  bytesInput,
		BytesOutput: bytesOutput,
	}
}
