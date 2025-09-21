package service

import (
	"context"
	"fmt"
	pb "review-service/api/review/v1"
	"review-service/internal/biz"
	"review-service/internal/data/model"
)

type ReviewService struct {
	pb.UnimplementedReviewServer
	uc *biz.ReviewUsecase
}

func NewReviewService(uc *biz.ReviewUsecase) *ReviewService {
	return &ReviewService{
		uc: uc,
	}
}

// CreateReview 创建评价
func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	fmt.Printf("[service] CreateReview,req:%v", req)
	// 参数转化
	// 调用biz层
	var anonymous int32
	if req.Anonymous {
		anonymous = 1
	}
	review, err := s.uc.CreateReview(ctx, &model.ReviewInfo{
		UserID:       &req.UserID,
		OrderID:      req.OrderID,
		Score:        req.Score,
		ServiceScore: req.ServiceScore,
		ExpressScore: req.ExpressScore,
		Content:      req.Content,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		Anonymous:    anonymous,
		Status:       0,
	})
	// 拼装返回结果
	if err != nil {
		return nil, err
	}
	return &pb.CreateReviewReply{ReviewID: *review.ReviewID}, nil
}

// GetReview 获取评价
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {

	fmt.Printf("[service] GetReview,req:%v", req)
	review, err := s.uc.GetReview(ctx, req.ReviewID)
	if err != nil {
		return nil, err
	}
	// 将 biz.Review 转换为 pb.ReviewInfo
	var anonymous bool
	if review.Anonymous == 1 {
		anonymous = true
	}
	reviewInfo := &pb.ReviewInfo{
		ReviewID:     *review.ReviewID,
		UserID:       *review.UserID,
		OrderID:      review.OrderID,
		Score:        review.Score,
		ServiceScore: review.ServiceScore,
		ExpressScore: review.ExpressScore,
		Content:      review.Content,
		PicInfo:      review.PicInfo,
		VideoInfo:    review.VideoInfo,
		Anonymous:    anonymous,
		Status:       review.Status,
	}
	return &pb.GetReviewReply{Data: reviewInfo}, nil
}

// AuditReview 管理员审核评价
func (s *ReviewService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	return &pb.AuditReviewReply{}, nil
}

// ReplyReview 回复评价
func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	fmt.Printf("[service] ReplyReview,req:%v\n", req)
	// 调用biz层
	reply, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID:  req.GetReviewID(),
		StoreID:   req.GetStoreID(),
		Content:   req.GetContent(),
		PicInfo:   req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewReply{ReplyID: *reply.ReplyID}, nil
}

// AppealReview 申诉评价
func (s *ReviewService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
	return &pb.AppealReviewReply{}, nil
}

// AuditAppeal 审核申诉
func (s *ReviewService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error) {
	return &pb.AuditAppealReply{}, nil
}

// ListReviewByUserID 用户查看自己的评价列表
func (s *ReviewService) ListReviewByUserID(ctx context.Context, req *pb.ListReviewByUserIDRequest) (*pb.ListReviewByUserIDReply, error) {
	fmt.Printf("[service] ListReviewByUserID,req:%v\n", req)

	// 调用业务逻辑层
	reviews, total, err := s.uc.ListReviewByUserID(ctx, req.UserID, req.PageNumber, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 转换数据格式
	var reviewInfos []*pb.ReviewInfo
	for _, review := range reviews {
		var anonymous bool
		if review.Anonymous == 1 {
			anonymous = true
		}

		reviewInfo := &pb.ReviewInfo{
			ReviewID:     *review.ReviewID,
			UserID:       *review.UserID,
			OrderID:      review.OrderID,
			Score:        review.Score,
			ServiceScore: review.ServiceScore,
			ExpressScore: review.ExpressScore,
			Content:      review.Content,
			PicInfo:      review.PicInfo,
			VideoInfo:    review.VideoInfo,
			Anonymous:    anonymous,
			Status:       review.Status,
		}
		reviewInfos = append(reviewInfos, reviewInfo)
	}

	return &pb.ListReviewByUserIDReply{
		Data:  reviewInfos,
		Total: total,
	}, nil
}
