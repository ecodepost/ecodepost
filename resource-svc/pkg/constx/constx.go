package constx

import (
	"time"
)

const (
	AuditInit    int8 = 0
	AuditSuccess int8 = 1
	AuditIng     int8 = 2
	AuditDeleted int8 = 3
)

const CommentExpire time.Duration = time.Second * 60 // 86400
const CommentContentExpire time.Duration = time.Second * 3 * 86400
const CommentSubjectCache string = "comment:subject_cache_"
const CommentContentCache string = "comment:content_cache_"
const CommentStarCache string = "comment_start_list_cache_"

// const CommentIndexCache string = "comment:index_cache_"
