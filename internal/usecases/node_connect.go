package usecases

import (
	"context"
	"io"
	"github.com/Kotletta-TT/THub/internal/entity"
	"github.com/Kotletta-TT/THub/internal/usecases/repo"
	"github.com/Kotletta-TT/THub/logger"
)

type NodeRepoStore interface {
	Set(repo.NodeInterface)
}

type ConnectUseCase struct {
	nrs NodeRepoStore
	l   logger.Logger
}

func NewConnectUseCase(nrs NodeRepoStore, l logger.Logger) *ConnectUseCase {
	return &ConnectUseCase{
		nrs: nrs,
		l:   l,
	}
}

func (cuc *ConnectUseCase) NodeConnect(id string, conn io.ReadWriter) (context.Context, error) {
	cuc.nrs.Set(entity.NewNode(id, conn))
	ctx := context.Background()
	return ctx, nil
}

//TODO move Disconnect here
