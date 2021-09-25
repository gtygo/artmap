package artmap

type ArtMap struct {

}

func New() *ArtMap {
	return &ArtMap{

	}
}

func (m *ArtMap) Set(key string, value interface{}) {

}

func (m *ArtMap) Get(key string) (interface{}, bool) {
	return nil, false
}

func (m *ArtMap) Count() uint64 {
	return 0
}

func (m *ArtMap) Remove(key string) {

}

func (m *ArtMap) Pop(key string) {

}

func (m *ArtMap) IsEmpty() bool {
	return false
}
