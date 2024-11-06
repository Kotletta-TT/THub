package usecases

import "github.com/Kotletta-TT/THub/logger"

type NodeRepoRemove interface {
	Remove(id string) error
}

type DisconnectUseCase struct {
	l   logger.Logger
	nrr NodeRepoRemove
}

func NewDisconnectUseCase(nrr NodeRepoRemove, l logger.Logger) *DisconnectUseCase {
	return &DisconnectUseCase{
		nrr: nrr,
		l:   l,
	}
}

func (nd *DisconnectUseCase) Disconnect(id string) error {
	return nd.nrr.Remove(id)
}
