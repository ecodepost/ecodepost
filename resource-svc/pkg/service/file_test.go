package service

import (
	"log"
	"testing"

	"ecodepost/resource-svc/pkg/model/mysql"

	"github.com/samber/lo"
)

func loMapToSlice() []int64 {
	var cacheMap = map[string]*mysql.FileCache{
		"x1": {Guid: "x1", CreatedUid: 101},
		"x2": {Guid: "x2", CreatedUid: 102},
		"x3": {Guid: "x3", CreatedUid: 103},
	}
	// uids := make([]int64, 0)
	// for _, value := range cacheMap {
	// 	uids = append(uids, value.CreatedUid)
	// }
	uids := lo.MapToSlice(cacheMap, func(_ string, v *mysql.FileCache) int64 { return v.CreatedUid })

	return uids
}

func loMap() []string {
	var mysqlList = []mysql.FileGuid{
		{Guid: "x1", CntCollect: 1},
		{Guid: "x2", CntCollect: 2},
		{Guid: "x3", CntCollect: 3},
	}
	guids := lo.Map(mysqlList, func(v mysql.FileGuid, _ int) string { return v.Guid })
	return guids
}

func TestTT(t *testing.T) {
	res1 := loMapToSlice()
	log.Printf("res1--------------->"+"%+v\n", res1)

	res2 := loMap()
	log.Printf("res2--------------->"+"%+v\n", res2)
}
