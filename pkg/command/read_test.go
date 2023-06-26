package command

import (
	"fmt"
	"testing"
)

func TestReadExe(t *testing.T) {
	var exe = "lrwxrwxrwx 1 lxd docker 0 Jun 15 10:47 /proc/1572104/exe -> /usr/sbin/mysqld\nlrwxrwxrwx 1 lxd docker 0 Jun 15 10:47 /proc/1572104/exe -> /usr/sbin/mysqld"
	a1, b1 := ReadField([]byte(exe), 11)
	fmt.Println(string(a1), "|", string(b1))

}
