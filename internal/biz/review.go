package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	snowflake "review-service/pkg"
)

type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
	GetReview(context.Context, int64) (*model.ReviewInfo, error)
	SaveReply(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
	GetReviewReply(context.Context, int64) (*model.ReviewReplyInfo, error)
	AuditReview(context.Context, *AuditParam) error
	AppealReview(context.Context, *AppealParam) error
	AuditAppeal(context.Context, *AuditAppealParam) error
	ListReviewByUserID(ctx context.Context, userID int64, offset, limit int) ([]*model.ReviewInfo, error)
	CountReviewByUserID(ctx context.Context, userID int64) (int32, error)
}

type ReviewUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

func NewReviewUsecase(repo ReviewRepo, logger log.Logger) *ReviewUsecase {
	return &ReviewUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateReview 创建评价
// service调用该方法
func (uc *ReviewUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debug("[biz] CreateReview,req:%v", review)
	// 1. 数据校验
	// 1.1 参数基础校验 : 正常来说不应该放在这一层，在上一层能拦截住（validate)
	// 1.2 参数业务校验:比如已经评价过了订单，不能重复评价
	reviews, err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, v1.ErrorDbFailed("查询数据库失败")
	}
	if len(reviews) > 0 {
		return nil, v1.ErrorOrderReviewed("订单%d已经评价过了", review.OrderID)
	}
	// 2. 生成review ID (雪花算法)
	id := snowflake.GenID()
	review.ReviewID = &id
	// 3. 查询订单和商品快照信息
	// 实际业务场景下需要查询订单和商家服务（rpc调用服务）
	// 4. 拼装数据入库
	return uc.repo.SaveReview(ctx, review)
}

func (uc *ReviewUsecase) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debug("[biz] GetReview,reviewID:%d", reviewID)
	return uc.repo.GetReview(ctx, reviewID)
}

// CreateReply 创建评价回复
func (uc *ReviewUsecase) CreateReply(ctx context.Context, param *ReplyParam) (*model.ReviewReplyInfo, error) {
	// 1.数据校验
	// 1.1 已经回复的评价允许商家重复回复
	uc.log.WithContext(ctx).Infof("[biz] CreateReply,param:%v", param)
	review, err := uc.repo.GetReview(ctx, param.ReviewID)
	if err != nil {
		return nil, v1.ErrorDbFailed("查询数据库失败")
	}
	if review.HasReply == 1 {
		// 已经回复过了，允许重复回复
		return nil, v1.ErrorReviewAlreadyReplied("评价%d已经回复过了", param.ReviewID)
	}
	// 2.水平越权(A商家只能回复自己的不能回复B商家的)
	if *review.StoreID != param.StoreID {
		return nil, v1.ErrorPermissionDenied("水平越权，不能回复其他商家的评价")
	}
	// 调用data层创建一个评价的回复
	uc.log.WithContext(ctx).Debug("[biz] CreateReply,req:%v", param)
	id := snowflake.GenID()
	reply := &model.ReviewReplyInfo{
		ReplyID:   &id,
		ReviewID:  &param.ReviewID,
		StoreID:   &param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	return uc.repo.SaveReply(ctx, reply)

}

// AuditReview 审核评价
func (uc *ReviewUsecase) AuditReview(ctx context.Context, param *AuditParam) error {
	uc.log.WithContext(ctx).Debug("[biz] AuditReview,param:%v", param)
	return nil
}

// AppealReview 申诉评价
func (uc *ReviewUsecase) AppealReview(ctx context.Context, param *AppealParam) error {
	uc.log.WithContext(ctx).Debug("[biz] AppealReview,param:%v", param)
	return nil
}

// ListReviewByUserID 获取用户评价列表
func (uc *ReviewUsecase) ListReviewByUserID(ctx context.Context, userID int64, pageNumber, pageSize int32) ([]*model.ReviewInfo, int32, error) {
	uc.log.WithContext(ctx).Debug("[biz] ListReviewByUserID,userID:%d,pageNumber:%d,pageSize:%d", userID, pageNumber, pageSize)

	// 2. 计算分页参数
	offset := (pageNumber - 1) * pageSize

	// 3. 查询数据
	reviews, err := uc.repo.ListReviewByUserID(ctx, userID, int(offset), int(pageSize))
	if err != nil {
		return nil, 0, v1.ErrorDbFailed("查询用户评价列表失败")
	}

	// 4. 查询总数
	total, err := uc.repo.CountReviewByUserID(ctx, userID)
	if err != nil {
		return nil, 0, v1.ErrorDbFailed("查询用户评价总数失败")
	}

	return reviews, total, nil
}
