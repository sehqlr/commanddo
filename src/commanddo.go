package main

import (
        "bytes"
        "fmt"
        "os"
        "os/exec"
        "log"
        "io/ioutil"

        "gopkg.in/yaml.v2"
)


type ShellScript struct {
        Cmd string
        Opts map[string]string
        Args []string

}

func main() {
        filepath := os.Args[1]

        data, io_err := ioutil.ReadFile(filepath)
        if io_err != nil {
                log.Fatal(io_err)
        } else {
                fmt.Println("stdin:")
                fmt.Println(string(data))
        }


        sh_slice := []ShellScript{}

        yaml_err := yaml.Unmarshal(data, &sh_slice)
        if yaml_err != nil {
                log.Fatal(yaml_err)
        }

        var opts []string
        var args []string
        var out bytes.Buffer

        for _,sh := range sh_slice {
                for k,v := range sh.Opts {
                        if len(k) == 1 {
                                opts = append(opts, "-" + k + v)
                        } else {
                                opts = append(opts, "--" + k + "=" + v)
                        }
                }

                for _,v := range sh.Args {
                        args = append(args, os.ExpandEnv(v))
                }
                
                all_args := append(opts, args...)
                cmd := exec.Command(string(sh.Cmd), all_args...)
                cmd.Stdout = &out

                cmd_err := cmd.Run()
                if cmd_err != nil {
                        log.Fatal(cmd_err)
                }
        }


        fmt.Printf("stdout:\n%s\n", out.String())
}
