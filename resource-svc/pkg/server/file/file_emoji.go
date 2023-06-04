package file

import (
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	filev1 "ecodepost/pb/file/v1"

	"github.com/samber/lo"
)

// MyEmojiListByFileGuids 根据文章guids和用户信息，返回他的状态信息，emoji，收藏等信息
func (GrpcServer) MyEmojiListByFileGuids(ctx context.Context, req *filev1.MyEmojiListByFileGuidsReq) (*filev1.MyEmojiListByFileGuidsRes, error) {
	// 找到我的这些文章的emoji的所有数据
	// 一个文章对应多个emoji列表
	list, err := mysql.MyEmojiList(invoker.Db.WithContext(ctx), req.GetGuids(), req.GetUid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("my emoji list, err: " + err.Error())
	}

	// 分堆
	groupMap := lo.GroupBy(list, func(t mysql.FileEmojiStatics) string { return t.Guid })
	// 根据请求顺序，遍历guids
	myEmojiList := lo.Map(req.GetGuids(), func(guid string, i int) *filev1.MyEmojiInfo {
		list, flag := groupMap[guid]
		if !flag {
			return &filev1.MyEmojiInfo{
				Guid: guid,
				List: nil,
			}
		}
		// 拿到列表，转成PB
		emojiList := lo.Map(list, func(t mysql.FileEmojiStatics, i int) *commonv1.EmojiInfo { return t.ToPb() })
		return &filev1.MyEmojiInfo{
			Guid: guid,
			List: emojiList,
		}
	})

	return &filev1.MyEmojiListByFileGuidsRes{
		List: myEmojiList,
	}, nil
}
