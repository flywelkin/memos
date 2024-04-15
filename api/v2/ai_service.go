package v2

import (
	"context"
	"github.com/lithammer/shortuuid/v4"
	gml "github.com/usememos/memos/plugin/ai"
	apiv2pb "github.com/usememos/memos/proto/gen/api/v2"
	"github.com/usememos/memos/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APIV2Service) AiChat(ctx context.Context, request *apiv2pb.AiRequest) (*apiv2pb.AIResponse, error) {
	println(request.Content)
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user")
	}
	if user == nil {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
	message := gml.ChatRequest{
		Role:    "user",
		Content: request.Content,
	}
	res, err := gml.Chat(message)

	create := &store.Memo{
		ResourceName: shortuuid.New(),
		CreatorID:    user.ID,
		Content:      res,
		Visibility:   store.Private,
	}
	memo, err := s.Store.CreateMemo(ctx, create)
	if err != nil {
		return nil, err
	}

	response := &apiv2pb.AIResponse{
		Content: memo.Content,
	}
	return response, nil
}
