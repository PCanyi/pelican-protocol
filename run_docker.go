package pelican

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

// this is the current docker image used to run sshd inside of and connect to.
// We don't actually connect to a real sshd anymore, but we did during protocol
// development. It might be useful to test against the real sshd in the future,
// so we keep it around.
var DockerHubTestImage string = "jaten/pelican04"

// StartDockerImage starts a docker container based on image. It runs /sbin/my_init
// as the process, thereby assuming that /sbin/my_init is available inside the
// named image.
func StartDockerImage(image string) {
	cmd := exec.Command("/usr/bin/docker", "run", image, "/sbin/my_init")

	var oe bytes.Buffer
	cmd.Stdout = &oe
	cmd.Stderr = &oe

	err := cmd.Start()
	if err != nil {
		//fmt.Fprintf(os.Stderr, "err '%s' in StartDockerImage(). Out: '%s'", err, string(out))
		panic(err)
	}
	// wait until it is up, 3000 msec
	tries := 30
	waittm := 100 * time.Millisecond
	up := false
	for i := 0; i < tries; i++ {
		if bytes.Contains(oe.Bytes(), []byte("*** Runit started as PID")) {
			fmt.Printf("found desired header in '%s'\n", string(oe.Bytes()))
			up = true
			break
		}
		time.Sleep(waittm)
	}

	if !up {
		panic(fmt.Sprintf("StartDockerImage() could not detect docker running after %v, output: '%s'\n", time.Duration(tries)*waittm, string(oe.Bytes())))
	}

	fmt.Printf("StartDockerImage() done.\n")
}

/*
func (h *KnownHosts) SshAsRootIntoDocker(cmd []string) ([]byte, error) {

	dockerip := GetDockerIP()

	fullcmd := strings.Join(cmd, " ")
	out, err := h.SshConnect("root", "dot.ssh/id_rsa_docker_root", dockerip, 22, fullcmd)
	if err != nil {
		panic(err)
	}

	VPrintf("running '%s' produced: '%s'\n", fullcmd, string(out))

	// examples:
	// make this actually use the "code.google.com/p/go.crypto/ssh"
	// https://godoc.org/golang.org/x/crypto/ssh/agent
	// http://kukuruku.co/hub/golang/ssh-commands-execution-on-hundreds-of-servers-via-go
	// http://gitlab.cslabs.clarkson.edu/meshca/golang-ssh-example/commit/556eb3c3bcb58ad457920d894a696e9266bbad36

	//return exec.Command("make", fmt.Sprintf("ARGS='%s'", strings.Join(cmd, " ")), "sshroot").CombinedOutput()
	return out, err
}
*/

// TrimRightNewline removes the trailing byte of slice if it is a newline '\n' character.
func TrimRightNewline(slice []byte) []byte {
	n := len(slice)
	if n > 0 {
		if slice[n-1] == '\n' {
			slice = slice[:n-1]
		}
	}
	return slice
}

// RunningDockerId runs 'docker ps -q -n=1 -f status=running' and returns the output and any error.
func RunningDockerId() ([]byte, error) {
	out, err := exec.Command("/usr/bin/docker", "ps", "-q", "-n=1", "-f", "status=running").CombinedOutput()
	out = TrimRightNewline(out)
	return out, err
}

// StopAllDockers calls 'docker stop' on all containers determined by
// successive calls to RunningDockerId().
func StopAllDockers() {
	for {
		out, err := RunningDockerId()
		if err != nil {
			panic(err)
		}
		if len(out) == 0 {
			return
		}
		fmt.Printf("StopAllDockers() is stopping '%s'\n", string(out))
		_, err = exec.Command("/usr/bin/docker", "stop", string(out)).CombinedOutput()
		if err != nil {
			panic(err)
		}
	}
}

// GetDockerIP returns the IP address bound by the container returned by RunningDockerId().
func GetDockerIP() string {
	id, err := RunningDockerId()
	if err != nil {
		panic(err)
	}

	ip, err := exec.Command("docker", "inspect", "-f", "{{ .NetworkSettings.IPAddress }}", string(id)).CombinedOutput()
	if err != nil {
		panic(err)
	}
	return string(ip)
}
