// file: gotop.go
package main

import (
  //"log"
  //"time"
  //"sort"
  "fmt"
  "github.com/shirou/gopsutil/process"
)


type ProcInfo struct{
  Name  string
  Usage float64
}

type ByUsage []ProcInfo

func (a ByUsage) Len() int      { return len(a) }
func (a ByUsage) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUsage) Less(i, j int) bool {
  return a[i].Usage > a[j].Usage
}


func main() {

  var c []float64


  for n := 0 ;n<2; n++{
    processes, _ := process.Processes()
    for _, p := range processes{
      a, _ := p.CPUPercent()
      n, _ := p.Name()
  		if n == "usage"{
  			c = append(c, a)
        break
  		}
    }

  }

  m := 0
  for i, e := range c {
    if i==0 || e > m {
        m = e
    }
  }

  fmt.Println(m)

}
