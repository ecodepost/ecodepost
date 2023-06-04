package cache

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	"gorm.io/gorm"

	commonv1 "ecodepost/pb/common/v1"

	"github.com/ego-component/eredis"
)

const redisPrefixKey = "file:info:%s"

type fileCache struct{}

var FileCache *fileCache

func init() {
	FileCache = &fileCache{}
}

const expireTime time.Duration = 60 * time.Second

// 缓存文件类型
var fileTypeList = []commonv1.FILE_TYPE{commonv1.FILE_TYPE_DOCUMENT, commonv1.FILE_TYPE_QUESTION, commonv1.FILE_TYPE_COLUMN}

// SetInfo 先给一天的有效期
func (f *fileCache) SetInfo(ctx context.Context, cache *mysql.FileCache) (err error) {
	redisKey := fmt.Sprintf(redisPrefixKey, cache.Guid)
	client := invoker.Redis.Client()
	if err = client.HSet(ctx, redisKey, cache.ToMap()).Err(); err != nil {
		return fmt.Errorf("file cache set info fail, err: %w", err)
	}
	if err = client.Expire(ctx, redisKey, expireTime).Err(); err != nil {
		return fmt.Errorf("file cache set info expire faile, err: %w", err)
	}
	return
}

// SetIsAllowCreateComment 先给一天的有效期
// func (f *fileCache) SetIsAllowCreateComment(ctx context.Context, guid string, isAllowCreateComment int32) (err error) {
//	redisKey := fmt.Sprintf(redisPrefixKey, guid)
//	client := invoker.Redis.Client()
//	err = client.HSet(ctx, redisKey, "isACC", isAllowCreateComment).Err()
//	if err != nil {
//		return fmt.Errorf("file cache set info fail, err: %w", err)
//	}
//	err = client.Expire(ctx, redisKey, time.Hour).Err()
//	if err != nil {
//		return fmt.Errorf("file cache set info expire faile, err: %w", err)
//	}
//	return
// }

// GetInfo 从缓存中根据guid查询fileCache,如果缓存中没有，则从DB查询，并更新到缓存
func (f *fileCache) GetInfo(ctx context.Context, guid string) (cache *mysql.FileCache, err error) {
	redisKey := fmt.Sprintf(redisPrefixKey, guid)
	cache = &mysql.FileCache{}
	err = invoker.Redis.Client().HGetAll(ctx, redisKey).Scan(cache)
	if err != nil && !errors.Is(err, eredis.Nil) {
		return nil, fmt.Errorf("file cache get info fail, err: %w", err)
	}
	// 如果数据为空
	if errors.Is(err, eredis.Nil) || cache.Guid == "" {
		// 从DB中查询
		fileInfo, e := mysql.FileInfoByGuid(invoker.Db.WithContext(ctx), guid)
		if e != nil {
			return nil, fmt.Errorf("file cache get mysql info fail, err: %w", e)
		}
		if fileInfo.Id == 0 {
			return nil, fmt.Errorf("file cache get mysql info not exist")
		}
		fmt.Printf("fileInfo--------------->"+"%+v\n", fileInfo)
		// 如果是文档，则查询summary和isReadMore
		cache = fileInfo.ToCache()
		//if slice.Contains(fileTypeList, fileInfo.FileType) {
		//	output, err := invoker.AliOss.GetObject(fileInfo.ContentKey)
		//	if err != nil {
		//		elog.Error("file cache get oss fail", elog.FieldName(guid), elog.FieldErr(err))
		//	}
		//	// summary, isReadMore, err := GetSummaryAndIsReadMore(output, fileInfo.FileFormat)
		//	// if err != nil {
		//	//	elog.Error("file cache get summary fail", elog.FieldName(guid), elog.FieldErr(err))
		//	// }
		//	cache.WithContent(string(output))
		//}
		// 写到缓存
		if err = f.SetInfo(ctx, cache); err != nil {
			return nil, fmt.Errorf("file cache set info fail, err: %w", err)
		}
	}
	return
}

func (f *fileCache) SetSiteTopField(ctx context.Context, db *gorm.DB, guid string, value int32) (err error) {
	err = f.setField(ctx, db, guid, "isST", value)
	if err != nil {
		return fmt.Errorf("setSiteTopField fail, err: %w", err)
	}
	return nil
}

func (f *fileCache) SetRecommend(ctx context.Context, db *gorm.DB, guid string, value int32) (err error) {
	err = f.setField(ctx, db, guid, "isR", value)
	if err != nil {
		return fmt.Errorf("setSiteTopField fail, err: %w", err)
	}
	return nil
}

func (f *fileCache) SetIsAllowCreateComment(ctx context.Context, db *gorm.DB, guid string, value int32) (err error) {
	err = f.setField(ctx, db, guid, "isACC", value)
	if err != nil {
		return fmt.Errorf("setSiteTopField fail, err: %w", err)
	}
	return nil
}

func (f *fileCache) SetParentGuid(ctx context.Context, db *gorm.DB, guid string, value string) (err error) {
	err = f.setField(ctx, db, guid, "pGuid", value)
	if err != nil {
		return fmt.Errorf("setSiteTopField fail, err: %w", err)
	}
	return nil
}

// setField 设置某个key的数据
func (f *fileCache) setField(ctx context.Context, db *gorm.DB, guid string, field string, value any) (err error) {
	//  todo 先判断key，是否存在，这里要做个lua判断，保证原子性
	// 如果存在继续续期，不存在，先从DB中查询，然后写入缓存
	redisKey := fmt.Sprintf(redisPrefixKey, guid)
	flag, err := invoker.Redis.Exists(ctx, redisKey)
	if err != nil {
		return fmt.Errorf("redis exist fail, err: %w", err)
	}
	// todo 极低概率刚好还有几秒就过期
	if flag {
		invoker.Redis.Expire(ctx, redisKey, expireTime)
		err = invoker.Redis.HSet(ctx, redisKey, field, value)
		if err != nil {
			return fmt.Errorf("redis hset fail, err: %w", err)
		}
		return nil
	}

	// 从DB中查询
	fileInfo, e := mysql.FileInfoByGuid(db, guid)
	if e != nil {
		return fmt.Errorf("file cache get mysql info fail, err: %w", e)
	}
	if fileInfo.Id == 0 {
		return fmt.Errorf("file cache get mysql info not exist")
	}

	cache := fileInfo.ToCache()
	// 如果是文档，则查询summary和isReadMore
	//if slice.Contains(fileTypeList, fileInfo.FileType) {
	//	output, err := invoker.AliOss.GetObject(fileInfo.ContentKey)
	//	if err != nil {
	//		elog.Error("file cache get oss fail", elog.FieldName(guid), elog.FieldErr(err))
	//		return fmt.Errorf("GetObject output fail, %w", err)
	//	}
	//	// summary, isReadMore, err := GetSummaryAndIsReadMore(output, fileInfo.FileFormat)
	//	// if err != nil {
	//	//	elog.Error("file cache get summary fail", elog.FieldName(guid), elog.FieldErr(err))
	//	//	return fmt.Errorf("GetSummaryAndIsReadMore fail, %w", err)
	//	// }
	//	cache.WithContent(string(output))
	//}
	// 写到缓存
	if err = f.SetInfo(ctx, cache); err != nil {
		return fmt.Errorf("file cache set info fail, err: %w", err)
	}
	return
}

// BatchGetInfo 并发获取，所以没有顺序，需要用map得到数据
func (f *fileCache) BatchGetInfo(ctx context.Context, guids []string) (caches map[string]*mysql.FileCache, err error) {
	caches = make(map[string]*mysql.FileCache, 0)
	mtx := sync.RWMutex{}
	guidLength := len(guids)
	if guidLength <= 5 {
		caches, err = f.batchGetInfoByLengthLt5(ctx, guids)
	} else {
		// 计算并发数，向上取整必须使用float64，改变guid length的类型，否则容易出问题
		concurrentInt := int(math.Ceil(float64(guidLength) / 5))
		wg := sync.WaitGroup{}
		wg.Add(concurrentInt)
		for i := 0; i < concurrentInt; i++ {
			go func(i int, guids []string) {
				defer wg.Done()
				var tmpCaches map[string]*mysql.FileCache
				if i == concurrentInt-1 {
					tmpCaches, err = f.batchGetInfoByLengthLt5(ctx, guids[i*5:])
					if err != nil {
						err = fmt.Errorf("batch get info fail, err: %w", err)
						return
					}
				} else {
					tmpCaches, err = f.batchGetInfoByLengthLt5(ctx, guids[i*5:i*5+5])
					if err != nil {
						err = fmt.Errorf("batch get info fail, err: %w", err)
						return
					}
				}
				mtx.Lock()
				for _, value := range tmpCaches {
					caches[value.Guid] = value
				}
				mtx.Unlock()
			}(i, guids)
		}
		wg.Wait()
	}
	return
}

func (f *fileCache) batchGetInfoByLengthLt5(ctx context.Context, guids []string) (map[string]*mysql.FileCache, error) {
	caches := make(map[string]*mysql.FileCache, 0)
	for _, guid := range guids {
		cacheInfo, err := f.GetInfo(ctx, guid)
		if err != nil {
			return nil, fmt.Errorf("batch get "+guid+"info fail, err: %w", err)
		}
		caches[cacheInfo.Guid] = cacheInfo
	}
	return caches, nil
}

// func GetSummaryAndIsReadMore(content []byte, format commonv1.FILE_FORMAT) (summary string, isReadMore int32, err error) {
//	switch format {
//	// 只有富文本处理 summary
//	case commonv1.FILE_FORMAT_DOCUMENT_RICH:
//		// 从富文本中获取summary
//		summary = trimHtml(string(content))
//		// case commonv1.FILE_FORMAT_DOCUMENT_JSON:
//		//	// 从json中获取summary
//		//	quillInfo, err := quill.Render(content)
//		//	if err != nil {
//		//		return "", 0, fmt.Errorf("get quil summary fail, err: %w", err)
//		//	}
//		//	summary = trimHtml(string(quillInfo))
//		newRune := []rune(summary)
//		if len(newRune) > 1000 {
//			isReadMore = 1
//			summary = string(newRune[:800]) + "..."
//			return
//		}
//	}
//	return
// }
//
// func trimHtml(src string) string {
//	// 将HTML标签全转换成小写
//	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
//	src = re.ReplaceAllStringFunc(src, strings.ToLower)
//	// 去除STYLE
//	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
//	src = re.ReplaceAllString(src, "")
//	// 去除SCRIPT
//	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
//	src = re.ReplaceAllString(src, "")
//	// 去除所有尖括号内的HTML代码，并换成换行符
//	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
//	src = re.ReplaceAllString(src, "\n")
//	// 去除连续的换行符
//	re, _ = regexp.Compile("\\s{2,}")
//	src = re.ReplaceAllString(src, "\n")
//	return strings.TrimSpace(src)
// }
