package redis

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
)

type handleServiceImpl struct {
	//allHandles map[string]types.ReteHandle
	types.NwServiceImpl
	prefix string
	config map[string]interface{}
}

func NewHandleCollection(nw types.Network, config map[string]interface{}) types.HandleService {
	hc := handleServiceImpl{}
	hc.Nw = nw
	hc.config = config
	//hc.allHandles = make(map[string]types.ReteHandle)
	return &hc
}

func (hc *handleServiceImpl) Init() {
	hc.prefix = hc.Nw.GetPrefix() + ":h:"
	reteCfg := hc.config["rete"].(map[string]interface{})
	jtRef := reteCfg["jt-ref"].(string)
	jts := hc.config["jts"].(map[string]interface{})
	redisCfg := jts[jtRef].(map[string]interface{})
	redisutils.InitService(redisCfg)
}

func (hc *handleServiceImpl) RemoveHandle(tuple model.Tuple) types.ReteHandle {
	rkey := hc.prefix + tuple.GetKey().String()
	redisutils.GetRedisHdl().Del(rkey)
	//TODO: Dummy handle
	h := newReteHandleImpl(hc.GetNw(), tuple)
	return h

}

func (hc *handleServiceImpl) GetHandle(tuple model.Tuple) types.ReteHandle {
	return hc.GetHandleByKey(tuple.GetKey())
}

func (hc *handleServiceImpl) GetHandleByKey(key model.TupleKey) types.ReteHandle {
	rkey := hc.prefix + key.String()

	m := redisutils.GetRedisHdl().HGetAll(rkey)
	if len(m) == 0 {
		return nil
	} else {
		tuple := hc.Nw.GetTupleStore().GetTupleByKey(key)
		if tuple == nil {
			//TODO: error handling
			return nil
		}
		h := newReteHandleImpl(hc.GetNw(), tuple)
		return h
	}
}

func (hc *handleServiceImpl) GetOrCreateHandle(nw types.Network, tuple model.Tuple) types.ReteHandle {

	key := hc.prefix + tuple.GetKey().String()

	m := redisutils.GetRedisHdl().HGetAll(key)
	if len(m) == 0 {
		m := make(map[string]interface{})
		m["k"] = "v"
		redisutils.GetRedisHdl().HSetAll(key, m)
	}

	h := newReteHandleImpl(nw, tuple)
	return h
}