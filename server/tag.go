package server

import (
	"context"
	"github.com/zhaoziliang2019/tag_service/pkg/bapi"
	"github.com/zhaoziliang2019/tag_service/pkg/errcode"
	pb "github.com/zhaoziliang2019/tag_service/proto"
)

type TagServer struct {
}

func NewTagServer() *TagServer {
	return &TagServer{}
}
func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	api := bapi.NewAPI("127.0.0.1:8000/api/v1/tag?name=" + r.GetName())
	_, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ERROR_GET_TAG_LIST_FAIL)
	}

	return nil, err
}
