package server

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/watchfire-io/watchfire/internal/daemon/revision"
	pb "github.com/watchfire-io/watchfire/proto"
)

type revisionService struct {
	pb.UnimplementedRevisionServiceServer
	manager *revision.Manager
}

func (s *revisionService) ListRevisions(_ context.Context, req *pb.ListRevisionsRequest) (*pb.RevisionList, error) {
	projectPath, err := getProjectPath(req.ProjectId)
	if err != nil {
		return nil, err
	}

	revisions, err := s.manager.ListRevisions(projectPath)
	if err != nil {
		return nil, err
	}

	list := &pb.RevisionList{Revisions: make([]*pb.Revision, 0, len(revisions))}
	for _, r := range revisions {
		list.Revisions = append(list.Revisions, modelToProtoRevision(r))
	}
	return list, nil
}

func (s *revisionService) GetRevision(_ context.Context, req *pb.RevisionId) (*pb.Revision, error) {
	projectPath, err := getProjectPath(req.ProjectId)
	if err != nil {
		return nil, err
	}

	r, err := s.manager.GetRevision(projectPath, int(req.RevisionNumber))
	if err != nil {
		return nil, err
	}
	return modelToProtoRevision(r), nil
}

func (s *revisionService) CreateRevision(_ context.Context, req *pb.CreateRevisionRequest) (*pb.Revision, error) {
	projectPath, err := getProjectPath(req.ProjectId)
	if err != nil {
		return nil, err
	}

	r, err := s.manager.CreateRevision(projectPath, revision.CreateOptions{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return nil, err
	}
	return modelToProtoRevision(r), nil
}

func (s *revisionService) UpdateRevision(_ context.Context, req *pb.UpdateRevisionRequest) (*pb.Revision, error) {
	projectPath, err := getProjectPath(req.ProjectId)
	if err != nil {
		return nil, err
	}

	opts := revision.UpdateOptions{RevisionNumber: int(req.RevisionNumber)}
	if req.Title != nil {
		opts.Title = req.Title
	}
	if req.Content != nil {
		opts.Content = req.Content
	}
	if req.Complete != nil {
		opts.Complete = req.Complete
	}

	r, err := s.manager.UpdateRevision(projectPath, opts)
	if err != nil {
		return nil, err
	}
	return modelToProtoRevision(r), nil
}

func (s *revisionService) DeleteRevision(_ context.Context, req *pb.RevisionId) (*emptypb.Empty, error) {
	projectPath, err := getProjectPath(req.ProjectId)
	if err != nil {
		return nil, err
	}

	if err := s.manager.DeleteRevision(projectPath, int(req.RevisionNumber)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
