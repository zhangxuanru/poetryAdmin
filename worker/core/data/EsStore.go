package data

import "poetryAdmin/worker/core/define"

type EsStore struct {
}

func NewEsStore() *EsStore {
	return new(EsStore)
}

//保存作者信息到 ES中
func (e *EsStore) SaveAuthorData(data *define.PoetryAuthorDetail) {

}
