package node

import (
    //"bufio"
    //"fmt"
    "io/ioutil"
    //"os"
)

func Write(b []byte, p string)  error{
  err := ioutil.WriteFile(p, b, 0644)
  return err
}

func Read(p string) (error, string){
  dat, err := ioutil.ReadFile(p)
  return err, string(dat)
}
