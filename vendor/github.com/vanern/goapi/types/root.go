package types

type ConfigHook func() error
type Hook string

const (
	GlobalAfterConfigHook Hook = "GlobalAfterConfigHook"
	ServerAfterConfigHook Hook = "ServerAfterConfigHook"
	TokenAfterConfigHook  Hook = "TokenAfterConfigHook"
)

type ListResult[T any] struct {
	Data  []T
	Page  int32
	Total int32
}
