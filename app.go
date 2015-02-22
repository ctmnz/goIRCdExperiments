package main

import (
	"fmt"
	"net"
	"os"
//	"strconv"
	"strings"
)


// Server
type Server struct {
	name string
	port uint
}

// User

type User struct {
	nick string
	username string
	realname string
	hostname string
	server string
	channels []string
	userconn net.Conn
	isIn bool
}

func (u *User) SetNick (nick string) {
	u.nick = nick
}

func (u *User) SetUsername (username string) {
	u.username = username
}

func (u *User) SetRealname (realname string) {
	u.realname = realname
}

func (u *User) SetHostname (hostname string) {
	u.hostname = hostname
}

func (u *User) SetServer (server string) {
	u.server = server
}

func (u *User) JoinChannel (channel string) {

}

func (u *User) Send (msg []byte) {
	// u.userconn.Write([]byte(joinmsg))
	u.userconn.Write(msg)
}



func main() {
	ln, err := net.Listen("tcp", ":6667")
	if err != nil {
		fmt.Println("ERR! 1: ", err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("ERR! 2")
			os.Exit(1)
		}
		go handleConnection(conn)

	}
}


func handleConnection(conn net.Conn) {

	fmt.Println("Client connection handler!")

	u := new(User)
	u.userconn = conn

	// userid
	idbuf := make([]byte, 1024)
	for {
		//
		if (u.isIn != true) {
			u.userconn.Read(idbuf)
			idcmd := HandleUserCommad(u, string(idbuf))
			if idcmd[0] == "NICK" {
				usernick := strings.TrimSpace(idcmd[1])
				username := strings.TrimSpace(idcmd[3])
				userserver := strings.TrimSpace(idcmd[5])
		//		usernick := idcmd[1]
				//fmt.Println(idcmd)
				u.SetNick(usernick)
				u.SetServer(userserver)
				u.SetUsername(username)
				// Force channel join
				joinmsg := ":" + u.nick + "!ufo@blackhole/from/space JOIN #go \n"
				//fmt.Println(joinmsg)
				//u.userconn.Write([]byte(joinmsg))
				u.Send([]byte(joinmsg))

				u.isIn = true
				break
		        }

		}
	}


	//conn.Write([]byte(":ctmnz!a@a.com JOIN #go \n"))

	buf := make([]byte, 1024)

	for {

		_, err := conn.Read(buf)

		if err != nil {
			fmt.Println("error!: ", err.Error())
		        return

		}


		ucmd := HandleUserCommad(u, string(buf))

		fmt.Println("Message  from the user: ", ucmd)
	}


	conn.Close()

}

func HandleUserCommad(u *User, cmd string) []string {
	cmd = strings.Replace(cmd,"\n"," ",-1)
	usercmd := strings.Split(cmd, " ")
	fmt.Println("User command!: ",cmd)
	//fmt.Println("User command received: ", usercmd)
	if usercmd[0] == "NICK" {
		fmt.Println("command NICK")
	}

	if usercmd[0] == "JOIN" {
		fmt.Println("command JOIN")
	}
	// Experiments 
	if usercmd[0] == "WHO" {
		fmt.Println("command WHO")
		channelmsg := ":" + u.server + " 332 " + u.nick + " #go :This channel topic! \n"
                channelmsg2 := ":" + u.server + " 333 " + u.nick + " #go Daniel 1417218973 \n"
                channelmsg3 := ":" + u.server + " 353 " + u.nick + " = #go :" + u.nick + " @FakeOper1 +FakeOper2 Fakeuser1 Fakeuser2 Fakeuser3 \n"
                channelmsg4 := ":" + u.server + " 366 " + u.nick + " #go :End of /NAMES list.\n"
                channelmsg5 := ":" + u.server + " 324 " + u.nick + " #go +cnt\n"
                channelmsg6 := ":" + u.server + " 329 " + u.nick + " #go 1194785698\n"
                channelmsg7 := ":daniel!daniel@/black/hole PRIVMSG #go :This is fake message\n"
		channelmsg8 := ":FakeOper2!~freeformz@c-76-115-27-201.hsd1.or.comcast.net QUIT :Quit: My MacBook has gone to sleep. ZZZzzz… \n"
		channelmsg9 := ":Daniel@Daniel!id@somewhere.on.the.earth JOIN #go\n"
		channelmsg10 := ":ctmnz!a@a.com MODE #go +o ctmnz \n"
		u.Send([]byte(channelmsg))
                u.Send([]byte(channelmsg2))
                u.Send([]byte(channelmsg3))
                u.Send([]byte(channelmsg4))
                u.Send([]byte(channelmsg5))
                u.Send([]byte(channelmsg6))
                u.Send([]byte(channelmsg7))
		u.Send([]byte(channelmsg8))
		u.Send([]byte(channelmsg9))
		u.Send([]byte(channelmsg10))


	}

	return usercmd
}



