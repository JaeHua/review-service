package biz

// ReplyParam 商家回复评价的参数
type ReplyParam struct {
	ReviewID  int64 // 评价ID
	StoreID   int64
	Content   string
	PicInfo   string
	VideoInfo string
}

// AuditParam 运营审核评价的参数
type AuditParam struct {
	ReviewID int64 // 评价ID
	OpUser   string
	OpReason string
	OpRemark string
	Status   int32 // 1通过 2拒绝
}

// AppealParam 商家申诉评价的参数
type AppealParam struct {
	ReviewID  int64 // 评价ID
	StoreID   int64
	Reason    string
	Content   string
	PicInfo   string
	VideoInfo string
}

// AuditAppealParam 运营审核申诉的参数
type AuditAppealParam struct {
	ReviewID int64 // 评价ID
	OpUser   string
	OpReason string
	OpRemark string
	Status   int32 // 1通过 2拒绝
}
