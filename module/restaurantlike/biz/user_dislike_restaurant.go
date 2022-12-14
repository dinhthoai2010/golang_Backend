package restaurantlikebiz

import (
	"context"
	"food-delivery-service/common"
	restaurantlikemodel "food-delivery-service/module/restaurantlike/model"
	"food-delivery-service/pubsub"
)

type UserDislikeRestaurantStore interface {
	Delete(ctx context.Context, userId, restaurantId int) error
	FindUserLike(ctx context.Context, userId, restaurantId int) (*restaurantlikemodel.Like, error)
}

//type DecLikedCountResStore interface {
//	DecreaseLikeCount(ctx context.Context, id int) error
//}

type userDislikeRestaurantBiz struct {
	store UserDislikeRestaurantStore
	//decStore DecLikedCountResStore
	pb pubsub.PubSub
}

func NewUserDislikeRestaurantBiz(
	store UserDislikeRestaurantStore,
	//decStore DecLikedCountResStore,
	pb pubsub.PubSub,
) *userDislikeRestaurantBiz {
	return &userDislikeRestaurantBiz{
		store: store,
		pb:    pb,
		//decStore: decStore,
	}
}

func (biz *userDislikeRestaurantBiz) DislikeRestaurant(
	ctx context.Context,
	userId,
	restaurantId int,
) error {
	oldData, err := biz.store.FindUserLike(ctx, userId, restaurantId)

	if oldData == nil {
		return restaurantlikemodel.ErrCannotDidNotlikeRestaurant(err)
	}

	err = biz.store.Delete(ctx, userId, restaurantId)

	if err != nil {
		return restaurantlikemodel.ErrCannotDislikeRestaurant(err)
	}

	// Side effect
	//go func() {
	//	defer common.Recover()
	//	job := asyncjob.NewJob(func(ctx context.Context) error {
	//		if err := biz.decStore.DecreaseLikeCount(ctx, restaurantId); err != nil {
	//			logger.GetCurrent().GetLogger("user.dislike.restaurant").Errorln(err)
	//			return err
	//		}
	//
	//		return nil
	//	}, asyncjob.WithName("DecreaseLikeCount"))
	//
	//	if err := asyncjob.NewGroup(false, job).Run(ctx); err != nil {
	//		logger.GetCurrent().GetLogger("user.dislike.restaurant").Errorln(err)
	//	}
	//}()

	newMessage := pubsub.NewMessage(map[string]interface{}{
		"user_id":       userId,
		"restaurant_id": restaurantId,
	})

	_ = biz.pb.Publish(ctx, common.TopicUserDislikeRestaurant, newMessage)

	return nil
}
