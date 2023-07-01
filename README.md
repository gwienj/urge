# urge
unsafe remote golang envelope, remote console for golang backgroud services.

# Server example
package main

import (
	"fmt"
	//"strings"
	"github.com/gwienj/urge"
)

func main() {
	fmt.Println("gurge server started")
	urge.Server(command, 3)
}

func command(msg []byte) ([]byte, bool) {
	var cmd, target, param1, param2, param3, param4, param5, param6 string	
	str := urge.Buffer2String(msg)
	fmt.Sscan(str, &cmd, &target, &param1, &param2, &param3, &param4, &param5, &param6)
	switch cmd {
	case "list":
		return []byte("list command response"), true
	case "get":
		return []byte("get command response"), true
	case "set":
		return []byte("set command response"), true
	case "reboot":
		return []byte("reboot command response"), true
	case "reset":
		return []byte("reset command response"), true
	case "ping":
		return []byte("ping command response"), true
	case "login":
		if urge.Sharekey("surge", "abcd") == param1 {
			return []byte("login success"), true
		} else {
			return []byte("authentication failed"), false
		}				
	default:
		return []byte("invalid command"), true
	}
}

# Client example
package main

import (
	"github.com/gwienj/urge"
)

func main() {
	urge.Client()
}
