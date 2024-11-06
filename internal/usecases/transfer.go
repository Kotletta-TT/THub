package usecases

import (
	"io"
	"github.com/Kotletta-TT/THub/logger"
)

type NodeRepoGet interface {
	Get(id string) (io.ReadWriter, error)
}

type TransferUseCase struct {
	nr NodeRepoGet
	l  logger.Logger
}

func NewTransferUseCase(nr NodeRepoGet, l logger.Logger) *TransferUseCase {
	return &TransferUseCase{
		nr: nr,
		l:  l,
	}
}

func (tuc *TransferUseCase) userToNode(user io.Reader, node io.Writer) error {
	_, err := io.Copy(node, user)
	return err
}

func (tuc *TransferUseCase) nodeToUser(node io.Reader, user io.Writer) error {
	_, err := io.Copy(user, node)
	return err
}

func (tuc *TransferUseCase) Transfer(user io.ReadWriter, nodeId string) error {
	node, err := tuc.nr.Get(nodeId)
	if err != nil {
		return err
	}
	go tuc.userToNode(user, node)
	//TODO ошибка userToNode не обработана
	return tuc.nodeToUser(node, user)
}
