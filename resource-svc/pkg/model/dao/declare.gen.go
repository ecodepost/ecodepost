package dao

import (
	"ecodepost/resource-svc/pkg/invoker"
)

var (
	CommentSubject *commentSubject
	CommentIndex   *commentIndex
	CommentContent *commentContent
)

func InitGen() error {
	CommentSubject = InitCommentSubject(invoker.Logger, invoker.Db)
	CommentIndex = InitCommentIndex(invoker.Logger, invoker.Db)
	CommentContent = InitCommentContent(invoker.Logger, invoker.Db)
	return nil
}
