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
	CommentSubject = InitCommentSubject(invoker.Logger, invoker.CommentDb)
	CommentIndex = InitCommentIndex(invoker.Logger, invoker.CommentDb)
	CommentContent = InitCommentContent(invoker.Logger, invoker.CommentDb)
	return nil
}
