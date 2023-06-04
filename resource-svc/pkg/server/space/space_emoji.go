package space

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	spacev1 "ecodepost/pb/space/v1"

	"github.com/samber/lo"
)

func (GrpcServer) EmojiList(ctx context.Context, req *spacev1.EmojiListReq) (*spacev1.EmojiListRes, error) {
	spaceList, err := mysql.SpaceListByUser(invoker.Db.WithContext(ctx), req.GetUid())
	if err != nil {
		return nil, fmt.Errorf("space tree get space list fail, err: %w", err)
	}
	output := lo.Map(spaceList, func(t *mysql.Space, i int) *spacev1.SpaceEmojiList {
		return &spacev1.SpaceEmojiList{
			SpaceGuid: t.Guid,
			EmojiList: mysql.EmojiList(),
		}
	})
	return &spacev1.EmojiListRes{
		SpaceList: output,
	}, nil
}
