package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Server struct {
	name     string
	port     string
	channels []*Channel
	users	[]*User
}

func (s *Server) CreateChannel(channel *Channel) {
	s.channels = append(s.channels, channel)
}




func (s *Server) MsgToChannel (channel *Channel, u *User, m string) {

// ":daniel!daniel@/black/hole PRIVMSG #go :This is fake message\n"

	// get the message
	chmsg := m // TODO: parse illegal chars

	// get sender user host address
	uaddr := u.GetUserAddress()

	// get channel name
	chname := channel.GetChannelName()

	// get channel members
	chmembers := channel.GetChannelMembers()

	fmt.Println("channel members: ", chmembers)

	// send them the message
	sm := ":"+uaddr+" PRIVMSG "+chname+" "+chmsg+"\n"


	for _, member := range chmembers {
		member.Send([]byte(sm))
	}

}


type Channel struct {
	name    string
	members []*User
	topic   string
}

func (c *Channel) AddMember(u *User) {
	c.members = append(c.members, u)
}

func (c *Channel) GetChannelName() string {
	return c.name
}

func (c *Channel) GetChannelMembers() []*User {
	return c.members
}


type User struct {
	nick     string
	username string
	realname string
	hostname string
	channels []*Channel
	server	*Server
	userconn net.Conn
	isIn     bool
}

func (u *User) SetNick(nick string) {
	u.nick = nick
}

func (u *User) SetUsername(username string) {
	u.username = username
}

func (u *User) SetRealname(realname string) {
	u.realname = realname
}

func (u *User) SetHostname(hostname string) {
	u.hostname = hostname

}

func (u *User) GetUserAddress() string {
	// ":daniel!daniel@/black/hole PRIVMSG #go :This is fake message\n"
	uaddress := u.nick + u.hostname
	return uaddress
}

func (u *User) SetServer(server *Server) {
	u.server = server
}

func (u *User) JoinChannel(channel string) {
	joinmsg := ":" + u.nick + u.hostname + " JOIN " + channel + "\n"
	u.Send([]byte(joinmsg))
//	TODO: 
//	u.channels = append(u.channels, channel)
}

func (u *User) Send(msg []byte) {
	// u.userconn.Write([]byte(joinmsg))
	u.userconn.Write(msg)
}



func (u *User) SendNotice(notice string) {
	// :ctmnz!~ctmnz@bg.ibg.bg NOTICE ctmnz :test
	nmsg :=":" + u.nick + u.hostname + " NOTICE " + u.nick + " :" + notice + "\n"
	u.Send([]byte(nmsg))
}


//////////////////

func main() {

	ircServer := Server{name: "pmp.ibg.bg", port: ":6667"}

	initchannel := Channel{name: "#bulgaria", topic: "Welcome! :-)"}

	ln, err := net.Listen("tcp", ircServer.port)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

// --->	initchannel := Channel{name:"#Bulgaria",topic:"Welcome to the channel #Bulgaria"}

	// Server Connection loop

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		go handleUserConnection(conn, &ircServer, &initchannel)
	}

}

func handleUserConnection(conn net.Conn, s *Server, initchannel *Channel) {
	// initialize the user
	ircUser := createUserFromConn(conn, s)
	buf := make([]byte, 1024)
	// User Connection Loop
	uchan := make(chan string)
	ucommand := []string{}

	initchannel.AddMember(ircUser)

	// s.MsgToChannel(initchannel, ircUser, "New user logged in the network!")


	for {

		go func() {
			ircUser.userconn.Read(buf)
			uchan <- string(buf)
		}()

		//fmt.Println(<-uchan)
		ucommand = parseUserCommand(<-uchan)
		fmt.Println(ucommand[0])
		switch {
			case ucommand[0] == "QUIT":
				// remove the user from all channels
				conn.Close()
				return
			case ucommand[0] == "PRIVMSG":
				if string(ucommand[1][0]) == "#" {
					fmt.Println("PRIVMSG channel command: ",ucommand[1:])
					fmt.Println("User channels: ", ircUser.channels)
					umsg := strings.Join(ucommand[2:], " ")
                                        fmt.Println("umsg: ", umsg)
					s.MsgToChannel(initchannel, ircUser, umsg)
				}
			default:
				ircUser.SendNotice("Command '"+  ucommand[0] + "' is not implemented")
		}





	}

	conn.Close()
}

func createUserFromConn(c net.Conn,s *Server) *User {
	u := new(User)
	u.userconn = c
	idbuf := make([]byte, 1024)

	// Identify user loop

	for {
		//
		if u.isIn != true {
			u.userconn.Read(idbuf)
			idcmd := parseUserCommand(string(idbuf))
			if idcmd[0] == "NICK" {
				usernick := strings.TrimSpace(idcmd[1])
				username := strings.TrimSpace(idcmd[3])
			//	userserver := strings.TrimSpace(idcmd[5])
				u.SetNick(usernick)
				u.SetServer(s)
				u.SetUsername(username)
				u.SetHostname("!ufo@blackhole/from/space")
				u.JoinChannel("#bulgaria")
				u.isIn = true
				break
			}

		}
	}

	// We have uer now

	return u

}

func parseUserCommand(c string) []string {
	c = strings.Replace(c, "\n", " ", -1)
	usercmd := strings.Split(c, " ")
	return usercmd
}
