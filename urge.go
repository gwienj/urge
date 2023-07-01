package urge

import (
	"fmt"
	"time"
	"strings"
	"net"
	"crypto/md5"
	//"encoding/json"
)

const (
	cServerPort = 5291
)

type cb func([]byte) ([]byte, bool)

func Sharekey(name string, key string) string {
	data := []byte(name + key)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func Buffer2String(buf []byte) string {
	for i, c := range buf {
		if c == 0 {
			return string(buf[:i])
		}
	}
	return string(buf)
}

func Client() {
	//todo
	var conn net.Conn
	var err error
	var login bool = false
	var prompt, cmd, target, param1, param2, param3, param4, param5, param6 string

	fmt.Println("welcome to urge console (gotoo ver1.0 komi 2023)")
	for {
		if login {
			fmt.Print(prompt)
		} else {
			fmt.Print("gotoo# ")
		}		
		n, _ := fmt.Scanln(&cmd, &target, &param1, &param2, &param3, &param4, &param5, &param6)
		switch cmd {
		case "help":
			fmt.Print("command format:\n  login [serverip] [name] [sharekey]\n  logout\n  exit\n")
		case "login":
			if n == 4 {
				if login {
					conn.Close()
				}

				conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target, cServerPort), 3*time.Second)
				if err != nil {
					fmt.Println("login failed")
				} else {
					login = true
					prompt = fmt.Sprintf("[%s]%s> ", target, param1)

					conn.Write([]byte(fmt.Sprintf("login %s %s", param1, Sharekey(param1, param2))))

					err = conn.SetReadDeadline(time.Now().Add(3*time.Second))
					if err != nil {
						fmt.Println(err)
					}

					buffer := make([]byte, 1024)
					_, e := conn.Read(buffer)
					if e != nil {
						conn.Close()
						login = false
						fmt.Println("no response")
					} else {
						res := Buffer2String(buffer)
						if res == "terminateconnection" {
							conn.Close()
							login = false														
						} else {
							fmt.Println(res)
						}
					}					
				}				
			} else {
				fmt.Println("command format: login [serverip] [username] [sharekey]")
			}
		case "logout":
			if login {
				conn.Close()
				login = false
				fmt.Println("logout")
			}
		case "exit":
			if login {
				conn.Close()
				login = false
			}
			return
		default:
			if login {
				msg := strings.Trim(fmt.Sprintf("%s %s %s %s %s %s %s %s", cmd, target, param1, param2, param3, param4, param5, param6), " ")
				conn.Write([]byte(msg))

				err = conn.SetReadDeadline(time.Now().Add(3*time.Second))
				if err != nil {
					fmt.Println(err)
				}

				buffer := make([]byte, 1024)
				_, e := conn.Read(buffer)
				if e != nil {
					conn.Close()
					login = false
					fmt.Println("no response")
				} else {
					fmt.Printf("%s\n",buffer)
				}
			}
		}
	}
}

func Server(f cb, m int) {
	if m < 3 {
		m = 3
	}
	vty := make(chan bool, m)
	for i := 0; i < m; i++ {
		vty <- true
	}

	srv, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cServerPort))
	if err != nil {
		return
	}
	defer srv.Close()

	for {
		conn, err := srv.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn, f, vty)
	}
}

func handleConnection(conn net.Conn, f cb, vty chan bool) {
	defer conn.Close()

	fmt.Println("client connected")

	select {
	case <- vty:
		for {
			buffer := make([]byte, 1024)
			_, re := conn.Read(buffer)
			if re != nil {
				break
			}
			resp, ok := f(buffer)
			_, we := conn.Write(resp)
			if we != nil {
				break
			}
			time.Sleep(1*time.Second)
			if !ok {
				conn.Write([]byte("terminateconnection"))
				time.Sleep(1*time.Second)
				break
			}
			
			
		}
		conn.Close()
		vty <- true
	default:
		conn.Write([]byte("limit maximum connections"))
		time.Sleep(1*time.Second)
		conn.Close()
	}	
	
	fmt.Println("client closed")	
}
