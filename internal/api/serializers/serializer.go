package serializers

type DataSerializer interface {
	Serialize(obj interface{}, many bool) (interface{}, error)
}
