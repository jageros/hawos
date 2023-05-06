package docker

import (
	"context"
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	cmds := []string{"nginx", "-s", "reload"}
	rsp, err := Exec(context.Background(), "nginx", cmds)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp)
}
