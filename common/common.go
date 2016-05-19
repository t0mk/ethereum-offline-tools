package common

import (
	"bufio"
	"fmt"
	"log"

	"os"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
)

// http://stackoverflow.com/questions/2137357/getpasswd-functionality-in-go
func GetPassword(prompt string) string {
	fmt.Print(prompt)

	// Common settings and variables for both stty calls.
	attrs := syscall.ProcAttr{
		Dir:   "",
		Env:   []string{},
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys:   nil}
	var ws syscall.WaitStatus

	// Disable echoing.
	pid, err := syscall.ForkExec(
		"/bin/stty",
		[]string{"stty", "-echo"},
		&attrs)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the stty process to complete.
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Echo is disabled, now grab the data.
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// Re-enable echo.
	pid, err = syscall.ForkExec(
		"/bin/stty",
		[]string{"stty", "echo"},
		&attrs)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the stty process to complete.
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(text)
}

func GetCheckedPassword(prompt string) string {
	p := GetPassword(prompt + " first time: ")
	fmt.Println()
	p2 := GetPassword(prompt + " second time: ")
	fmt.Println()
	if p != p2 {
		log.Fatal("password verfication check failed")
	}
	return p
}

func GetAccountForAddr(am *accounts.Manager, addr common.Address) accounts.Account {
	for a := range am.Accounts() {
		acc, err := am.AccountByIndex(a)
		if err != nil {
			log.Fatal(err)
		}
		if acc.Address == addr {
			return acc
		}
	}
	log.Fatal("couldnt find account for address", addr)
	return accounts.Account{}
}
