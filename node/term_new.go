package main

import (
	"log"
	"os/exec"

	"github.com/creack/pty"
)

func (c *sshClient) bridgeWSAndPTY() {
	defer c.conn.Close()

	wdSize, err := c.getWindowSize()
	if err != nil {
		log.Println("bridgeWSAndPTY: getWindowSize:", err)
		return
	}

	// Запускаем локальную оболочку с использованием PTY
	cmd := exec.Command("/bin/bash")

	// Создаем PTY и запускаем команду
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("bridgeWSAndPTY: pty.Start:", err)
		return
	}
	defer func() {
		_ = ptmx.Close()
	}()

	// Устанавливаем размер окна терминала
	if err := pty.Setsize(ptmx, &pty.Winsize{
		Rows: uint16(wdSize.High),
		Cols: uint16(wdSize.Width),
	}); err != nil {
		log.Println("bridgeWSAndPTY: pty.Setsize:", err)
		return
	}

	// Настраиваем ввод и вывод
	c.sessIn = ptmx
	c.sessOut = ptmx
	// c.se

	log.Println("Запущена локальная оболочка")
	defer log.Println("Локальная оболочка закрыта")

	go func() {
		if err := c.wsRead(); err != nil {
			log.Println("bridgeWSAndPTY: wsRead:", err)
		}
	}()

	go func() {
		if err := c.wsWrite(); err != nil {
			log.Println("bridgeWSAndPTY: wsWrite:", err)
		}
	}()

	<-c.closeSig
}
