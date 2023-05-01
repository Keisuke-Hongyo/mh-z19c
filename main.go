package main

import (
	"fmt"
	"machine"
	"mh-z19c/mh_z19c"
	"time"
)

const BufferSize = 9

// Go routine
/* UARTからの受信処理 */
func DataRecv(ch chan<- int, sensor mhz19c.MHZ19c) {
	for {
		ch <- sensor.GetData()
		time.Sleep(2 * time.Second)
	}
}

func main() {
	var ppm int
	sensor, err := mhz19c.New(machine.UART1)
	if err != nil {
		panic(err)
	}
	sensorCh := make(chan int, 1)
	sensor.AutoCalibration(true)
	go DataRecv(sensorCh, sensor)
	for {
		select {
		case ppm = <-sensorCh:
			fmt.Printf("Co2 %d[ppm]\n", ppm)
			break

		}
	}
}
