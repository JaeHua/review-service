package data

import (
	"context"
	"review-service/internal/data/model"
	"review-service/internal/data/query"
	"time"

	"review-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type reviewRepo struct {
	data *Data
	log  *log.Helper
}

// NewReviewRepo .
func NewReviewRepo(data *Data, logger log.Logger) biz.ReviewRepo {
	return &reviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// SaveReview 保存评价
func (r *reviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Save(review)
	return review, err
}

// GetReviewByOrderID 根据订单ID获取评价列表
func (r *reviewRepo) GetReviewByOrderID(ctx context.Context, orderID int64) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.OrderID.Eq(orderID)).
		Find()
}

// GetReview 根据评价ID获取评价
func (r *reviewRepo) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(reviewID)).
		First()
}

// SaveReply 保存评价回复
func (r *reviewRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {

	// 3. 更新数据库中的数据(评价表和回复表同时更新/事务操作)
	r.data.query.Transaction(func(tx *query.Query) error {
		// 回复表插入一条数据
		if err := tx.ReviewReplyInfo.
			WithContext(ctx).Save(reply); err != nil {
			r.log.WithContext(ctx).Errorf("Save review reply failed,err:%v", err)
			return err
		}
		// 更新评价表的回复状态
		if _, err := tx.ReviewInfo.WithContext(ctx).
			Where(tx.ReviewInfo.ReviewID.Eq(*reply.ReviewID)).
			Update(tx.ReviewInfo.HasReply, 1); err != nil {
			r.log.WithContext(ctx).Errorf("Update reply failed,err:%v", err)
			return err
		}
		return nil
	})
	// 4.返回
	return reply, nil
}

// GetReviewReply 根据评价ID获取评价回复
func (r *reviewRepo) GetReviewReply(ctx context.Context, reviewID int64) (*model.ReviewReplyInfo, error) {
	return r.data.query.ReviewReplyInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewReplyInfo.ReviewID.Eq(reviewID)).
		First()
}

// AuditReview 审核评价
func (r *reviewRepo) AuditReview(ctx context.Context, param *biz.AuditParam) error {
	updateData := &model.ReviewInfo{
		Status:    param.Status,
		OpUser:    param.OpUser,
		OpReason:  param.OpReason,
		OpRemarks: param.OpRemark,
	}
	_, err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(param.ReviewID)).
		UpdateColumns(updateData)
	return err
}

// AppealReview 申诉评价
func (r *reviewRepo) AppealReview(ctx context.Context, param *biz.AppealParam) error {

	// 更新申诉相关字段
	updateData := map[string]interface{}{
		"op_reason":  param.Reason,
		"op_remarks": param.Content, // 使用 content 作为备注
		"pic_info":   param.PicInfo,
		"video_info": param.VideoInfo,
		"update_by":  "system",   // 设置更新者标识
		"update_at":  time.Now(), // 更新时间
	}

	_, err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(param.ReviewID)).
		Updates(updateData)
	return err
}

// AuditAppeal 审核申诉
func (r *reviewRepo) AuditAppeal(ctx context.Context, param *biz.AuditAppealParam) error {
	updateData := &model.ReviewInfo{
		Status:    param.Status,
		OpUser:    param.OpUser,
		OpReason:  param.OpReason,
		OpRemarks: param.OpRemark,
	}

	_, err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(param.ReviewID)).
		UpdateColumns(updateData)
	return err
}

// ListReviewByUserID 根据用户ID分页获取评价列表
func (r *reviewRepo) ListReviewByUserID(ctx context.Context, userID int64, offset, limit int) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.UserID.Eq(userID)).
		Order(r.data.query.ReviewInfo.CreateAt.Desc()). // 按创建时间倒序
		Offset(offset).
		Limit(limit).
		Find()
}

// CountReviewByUserID 统计用户评价总数
func (r *reviewRepo) CountReviewByUserID(ctx context.Context, userID int64) (int32, error) {
	count, err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.UserID.Eq(userID)).
		Count()
	return int32(count), err
}
