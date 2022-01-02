package main

import (
	"fmt"
	"github.com/ztrue/shutdown"
	"murphysec-cli-simple/cmd"
	"murphysec-cli-simple/util/output"
	"os"
	"strings"
)

func main() {
	go func() {
		shutdown.Listen()
		output.Error("User request interrupt.")
		os.Exit(-2)
	}()
	output.Debug(fmt.Sprintf("CLI arguments: %s", strings.Join(os.Args, " ")))
	r := cmd.RootCmd()
	_ = r.Execute()
}

//func main(){
//	cmd:=exec.Command("cmd", "/c", "echo aa")
//	stdout, e:=cmd.StdoutPipe()
//	must.Must(e)
//	stderr, e:=cmd.StderrPipe()
//	must.Must(e)
//	cmd.Start()
//	time.Sleep(time.Second)
//	wg:=sync.WaitGroup{}
//	wg.Add(2)
//	go func() {
//		fmt.Println("o", string(must.Byte(ioutil.ReadAll(stdout))))
//		wg.Done()
//	}()
//	go func() {
//		fmt.Println("e", string(must.Byte(ioutil.ReadAll(stderr))))
//		wg.Done()
//	}()
//
//	//must.Must(cmd.Start())
//	wg.Wait()
//	must.Must(cmd.Wait())
//}
