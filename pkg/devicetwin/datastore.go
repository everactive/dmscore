package devicetwin

type DataStore interface {
}

type DataStoreImpl struct{}

func NewDataStore() *DataStoreImpl {
	return &DataStoreImpl{}
}
