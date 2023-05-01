package mhz19c

import (
	"fmt"
	"machine"
)

const Recive_Size = 9
const Cmd_Size = 9

type MHZ19c struct {
	uart    *machine.UART
	cmdData []byte
	rcvData []byte
}

// New : 初期化関数
func New(uart *machine.UART) (MHZ19c, error) {
	mhz19c := MHZ19c{uart: uart}
	err := mhz19c.uart.Configure(machine.UARTConfig{BaudRate: 9600})
	if err != nil {
		return MHZ19c{}, err
	}
	mhz19c.cmdData = make([]byte, Cmd_Size)
	mhz19c.rcvData = make([]byte, Recive_Size)
	return mhz19c, nil
}

// SetCalibration : AutoCalibrationの設定
func (mhz19c *MHZ19c) AutoCalibration(mode bool) {
	if mode {
		mhz19c.cmdData = []byte{0xff, 0x01, 0x79, 0xa0, 0x00, 0x00, 0x00, 0x00, 0x00}
	} else {
		mhz19c.cmdData = []byte{0xff, 0x01, 0x79, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	}
	mhz19c.uart.Write(mhz19c.cmdData)
}

// GetData : Co2濃度データの取得関数
func (mhz19c *MHZ19c) GetData() int {
	// データ取得コマンド
	mhz19c.cmdData = []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}

	// 取得データを送信
	_, err := mhz19c.uart.Write(mhz19c.cmdData)
	if err != nil {
		// エラーなら濃度を0に設定
		return 9999
	}

	//Co2センサからCo2濃度データを受信
	if mhz19c.uart.Buffered() == Recive_Size {
		for i := 0; i < Recive_Size; i++ {
			mhz19c.rcvData[i], err = mhz19c.uart.ReadByte()
			if err != nil {
				// エラーなら濃度を0に設定
				return 9999
			}
		}
	}
	if mhz19c.rcvData[0] == 0xff && mhz19c.rcvData[1] == 0x86 {
		// 濃度を計算して戻り値を設定
		return int(mhz19c.rcvData[2])*256 + int(mhz19c.rcvData[3])
	} else {
		fmt.Print(mhz19c.rcvData)
		return 9999
	}

}
