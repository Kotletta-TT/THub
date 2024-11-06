package usecases

import (
	"github.com/Kotletta-TT/THub/logger"
)

type NodeRepoList interface {
	List() []string
}

type ListNodesUseCase struct {
	nr NodeRepoList
	l  logger.Logger
}

func NewListNodesUseCase(nr NodeRepoList, l logger.Logger) *ListNodesUseCase {
	return &ListNodesUseCase{nr: nr, l: l}
}

func (l *ListNodesUseCase) ListNodes() []string {
	return l.nr.List()
}
