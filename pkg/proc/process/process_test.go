package process

import (
	"fmt"
	"testing"
)

func TestGetProcessExe(t *testing.T) {
	var mockProcess = []struct {
		Data   *MockProcess
		Result string
	}{
		{
			Data: &MockProcess{
				cwd:     "lrwxrwxrwx 1 root root 0 Jun 20 16:25 /proc/10/cwd -> /home/ubuntu/",
				comm:    "bash",
				cmdline: "bash -c ../sleep.sh",
				pid:     10,
				exe:     "lrwxrwxrwx 1 root root 0 Jun 20 17:29 /proc/3747251/exe -> /usr/bin/bash*",
			},
			Result: "/home/sleep.sh",
		},
		{
			Data: &MockProcess{
				cwd:     "lrwxrwxrwx 1 root root 0 Jun 20 17:28 /proc/11/cwd -> //",
				comm:    "lighttpd",
				cmdline: "lighttpd -D -f /etc/lighttpd/lighttpd.conf",
				pid:     11,
				exe:     "lrwxrwxrwx 1 root root 0 Jun 19 10:52 /proc/2548908/exe -> /usr/sbin/lighttpd",
			},
			Result: "/usr/sbin/lighttpd",
		},
	}
	for idx, process := range mockProcess {
		r, err := GetProcessExe(process.Data)
		if err != nil {
			t.Fatal(err)
		}
		if r == process.Result {
			t.Logf("Test %d: success", idx)
		} else {
			t.Logf("Test fail %d: %s", idx, r)
		}
	}
	//process := MockProcess{cwd: "/home/test", comm: "bash", cmdline: "bash\t-c\t../sleep.sh", pid: 10}
	//fmt.Println(IsPath("abc"))
	//fmt.Println(GetProcessExe(3747251))

	//fmt.Println(bytes.Contains([]byte("/opt/jboss/wildfly/bin/standalone.sh"), []byte("standalone.sh")))
	///opt/jboss/wildfly/bin/standalone.sh
	//standalone.sh
}

func TestNsPids(t *testing.T) {
	fmt.Println(NewProcess(2874135, nil).NsPids())
}
