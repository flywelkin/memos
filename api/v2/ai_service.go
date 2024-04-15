package v2

import (
	"context"
	gml "github.com/usememos/memos/plugin/ai"
	apiv2pb "github.com/usememos/memos/proto/gen/api/v2"
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
	response := &apiv2pb.AIResponse{
		Content: res,
	}
	return response, nil
}
