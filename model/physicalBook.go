package model

import "time"

type PhysicalBookInfo struct {
	Id           int
	InLibrary    bool
	Section      SectionInfo
	IssueDate    *time.Time
	CompleteDate *time.Time
	ReaderId     *int
	ReaderName   *string
}
