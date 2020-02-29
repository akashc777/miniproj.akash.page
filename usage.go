package main

import (
"fmt"
"github.com/shirou/gopsutil/cpu"
"time"
)

func main() {
   //Percent calculates the percentage of cpu used either per CPU or combined.
   percent, _ := cpu.Percent(2*time.Second, false)
   fmt.Printf("\npercent val: %v\nlength : %v",percent, len(percent))
   var sum float64
   for _,val := range percent{
     sum += val
   }
   fmt.Printf("\n%v\n", sum/(float64(len(percent))))
//   fmt.Printf("  User: %.2f\n",percent[cpu.CPUser])
//   fmt.Printf("  Nice: %.2f\n",percent[cpu.CPNice])
//   fmt.Printf("   Sys: %.2f\n",percent[cpu.CPSys])
//   fmt.Printf("  Intr: %.2f\n",percent[cpu.CPIntr])
//   fmt.Printf("  Idle: %.2f\n",percent[cpu.CPIdle])
//   fmt.Printf("States: %.2f\n",percent[cpu.CPUStates])
}

