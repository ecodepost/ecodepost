package service

type BucketName string

var (
	Resource *resource
	Subject  *subjectService
	Index    *indexService
	Space    *space
	File     *file
	Audit    *audit
	Comment  *comment
)

func Init() error {
	Comment = &comment{}
	Subject = &subjectService{}
	Space = InitSpace()
	Index = &indexService{}
	Resource = initResource()
	File = InitFile()
	Audit = &audit{}
	return nil
}
