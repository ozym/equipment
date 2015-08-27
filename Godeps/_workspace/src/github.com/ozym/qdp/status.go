package qdp

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func QTime(t uint32) time.Time {
	return time.Unix((int64)(t)+946684800, 0)
}

type Status struct {
	Header struct {
		DriftTol   uint16
		UserMsgCnt uint16
		LastReboot uint32
		Spare      [2]uint32
		BitMap     uint32
	}

	Global struct {
		Aqctr          uint16
		ClockQual      uint16
		ClockLoss      uint16
		CurrentVoltage uint16
		SecOffset      uint32
		UsecOffset     uint32
		TotalTime      uint32
		TotalPower     uint32
		LastResync     uint32
		Resyncs        uint32
		GpsStat        uint16
		CalStat        uint16
		SensorMap      uint16
		CurVCO         uint16
		DataSeq        uint16
		PLLFlag        uint16
		StatInp        uint16
		MiscInp        uint16
		CurSequence    uint32
	}

	Boom struct {
		Booms       [6]int16
		AmbPos      uint16
		AmbNeg      uint16
		Supply      uint16
		SysTemp     int16
		MainCur     int16
		AntCur      int16
		Seis1Temp   int16
		Seis2Temp   int16
		CalTimeouts uint32
	}

	Ether struct {
		Check    uint32
		IOErrors uint32
		PhyNum   uint16
		Spare    uint16
		Unreach  uint32
		Quench   uint32
		Echo     uint32
		Redirect uint32
		Runt     uint32
		CRCErr   uint32
		BCast    uint32
		UCast    uint32
		Good     uint32
		Jabber   uint32
		OutWin   uint32
		TXOK     uint32
		Miss     uint32
		Collide  uint32
		LinkStat uint16
		Spare2   uint16
		Spare3   uint32
	}

	GPS struct {
		GPSTime  uint16
		GPSOn    uint16
		SatUsed  uint16
		SatView  uint16
		Time     [10]byte
		Date     [12]byte
		Fix      [6]byte
		Height   [12]byte
		Lat      [14]byte
		Lon      [14]byte
		LastGood uint32
		CheckErr uint32
	}

	LPort [4]struct {
		Sent     uint32
		Resends  uint32
		Fill     uint32
		Seq      uint32
		PackUsed uint32
		LastAck  uint32
		PhyNum   uint16
		LogNum   uint16
		Retran   uint16
		Spare    uint16
	}

	/*
		Pwr struct {
			Phase      uint16
			BatTemp    int16
			Capacity   uint16
			Depth      uint16
			BatVolt    uint16
			InpVolt    uint16
			BatCur     int16
			Absorption uint16
			Floating   uint16
			Spare      uint16
		}
	*/
}

// reformatted subset of stats
type SOH struct {
	Header struct {
		DriftTol   uint16    `json:"drift_tol"`
		UserMsgCnt uint16    `json:"user_msg_cnt"`
		LastReboot time.Time `json:"last_reboot"`
		BitMap     string    `json:"bitmap"`
	} `json:"header"`

	Global struct {
		ClockQual      uint16    `json:"clock_qual"`
		ClockLoss      uint16    `json:"clock_loss"`
		CurrentVoltage uint16    `json:"current_voltage"`
		SecOffset      uint32    `json:"sec_offset"`
		UsecOffset     uint32    `json:"usec_offset"`
		TotalTime      uint32    `json:"total_time"`
		TotalPower     uint32    `json:"total_power"`
		LastResync     time.Time `json:"last_resync"`
		Resyncs        uint32    `json:"resyncs"`
		GpsStat        uint16    `json:"gps_stat"`
		CalStat        uint16    `json:"cal_stat"`
		SensorMap      uint16    `json:"sensor_map"`
		CurVCO         uint16    `json:"cur_vco"`
		DataSeq        uint16    `json:"data_seq"`
		PLLFlag        uint16    `json:"pll_flag"`
		StatInp        uint16    `json:"stat_inp"`
		MiscInp        uint16    `json:"misc_inp"`
		CurSequence    uint32    `json:"cur_sequence"`
	} `json:"global"`

	Boom struct {
		Booms       [6]int16 `json:"booms"`
		AmbPos      uint16   `json:"amb_pos"`
		AmbNeg      uint16   `json:"amb_neg"`
		Supply      uint16   `json:"supply"`
		SysTemp     int16    `json:"sys_temp"`
		MainCur     int16    `json:"main_cur"`
		AntCur      int16    `json:"ant_cur"`
		Seis1Temp   int16    `json:"seis1_temp"`
		Seis2Temp   int16    `json:"seis2_temp"`
		CalTimeouts uint32   `json:"cal_timeouts"`
	} `json:"boom"`

	Ether struct {
		Check    uint32 `json:"check"`
		IOErrors uint32 `json:"io_erros"`
		PhyNum   uint16 `json:"phy_num"`
		Unreach  uint32 `json:"unreach"`
		Quench   uint32 `json:"quench"`
		Echo     uint32 `json:"echo"`
		Redirect uint32 `json:"redirect"`
		Runt     uint32 `json:"runt"`
		CRCErr   uint32 `json:"crc_err"`
		BCast    uint32 `json:"bcast"`
		UCast    uint32 `json:"ucast"`
		Good     uint32 `json:"good"`
		Jabber   uint32 `json:"jabber"`
		OutWin   uint32 `json:"outwin"`
		TXOK     uint32 `json:"tx_ok"`
		Miss     uint32 `json:"miss"`
		Collide  uint32 `json:"collide"`
		LinkStat uint16 `json:"link_stat"`
	} `json:"ether"`

	GPS struct {
		GPSTime  uint16    `json:"gps_time"`
		GPSOn    uint16    `json:"gps_on"`
		SatUsed  uint16    `json:"sat_used"`
		SatView  uint16    `json:"sat_view"`
		Time     string    `json:"time"`
		Date     string    `json:"date"`
		Fix      string    `json:"fix"`
		Height   string    `json:"height"`
		Lat      string    `json:"lat"`
		Lon      string    `json:"lon"`
		LastGood time.Time `json:"last_good"`
		CheckErr uint32    `json:"check_err"`
	} `json:"gps"`

	LPort [4]struct {
		Port     uint8  `json:"port"`
		Sent     uint32 `json:"sent"`
		Resends  uint32 `json:"resends"`
		Fill     uint32 `json:"fill"`
		Seq      uint32 `json:"seq"`
		PackUsed uint32 `json:"pack_used"`
		LastAck  uint32 `json:"last_ack"`
		PhyNum   uint16 `json:"phy_num"`
		LogNum   uint16 `json:"log_num"`
		Retran   uint16 `json:"retran"`
	} `json:"lport"`

	Timestamp time.Time `json:"timestamp"`
}

func (s *SOH) String() (string, error) {

	r, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	return (string)(r), nil
}

func (p *Ping) Status() *SOH {

	if p.PingType != 3 {
		return nil
	}

	s := Status{}

	// decode wire format ...
	binary.Read(bytes.NewReader(p.Data[:]), binary.BigEndian, &s.Header)

	// offsets to encoded data
	var offset int = 20
	var offsets = [...]int{52, 84, 0, 32, 0, 0, 0, 0, 32, 32, 32, 32, 0, 0, 0, 76}

	// encoded messages ...
	for i := (uint)(0); i < 16; i++ {
		if (s.Header.BitMap & (0x01 << i)) != 0x00 {
			switch i {
			case 0:
				binary.Read(bytes.NewReader(p.Data[offset:]), binary.BigEndian, &s.Global)
			case 1:
				binary.Read(bytes.NewReader(p.Data[offset:]), binary.BigEndian, &s.GPS)
			case 3:
				binary.Read(bytes.NewReader(p.Data[offset:]), binary.BigEndian, &s.Boom)
			case 8:
				binary.Read(bytes.NewReader(p.Data[offset:]), binary.BigEndian, &s.LPort[0])
			case 9:
				binary.Read(bytes.NewReader(p.Data[offset:]), binary.BigEndian, &s.LPort[1])
			case 10:
				binary.Read(bytes.NewReader(p.Data[offset:]), binary.BigEndian, &s.LPort[2])
			case 11:
				binary.Read(bytes.NewReader(p.Data[offset:]), binary.BigEndian, &s.LPort[3])
			case 15:
				binary.Read(bytes.NewReader(p.Data[offset:]), binary.BigEndian, &s.Ether)
			}
			offset += offsets[i]
		}
	}

	o := SOH{}

	// reformat a subset of header stats
	o.Header.DriftTol = s.Header.DriftTol
	o.Header.UserMsgCnt = s.Header.UserMsgCnt
	o.Header.LastReboot = QTime(s.Header.LastReboot)
	o.Header.BitMap = fmt.Sprintf("0x%x", s.Header.BitMap)

	// reformat a subset of global stats
	o.Global.ClockQual = s.Global.ClockQual
	o.Global.ClockLoss = s.Global.ClockLoss
	o.Global.CurrentVoltage = s.Global.CurrentVoltage
	o.Global.SecOffset = s.Global.SecOffset
	o.Global.UsecOffset = s.Global.UsecOffset
	o.Global.TotalTime = s.Global.TotalTime
	o.Global.TotalPower = s.Global.TotalPower
	o.Global.LastResync = QTime(s.Global.LastResync)
	o.Global.Resyncs = s.Global.Resyncs
	o.Global.GpsStat = s.Global.GpsStat
	o.Global.CalStat = s.Global.CalStat
	o.Global.SensorMap = s.Global.SensorMap
	o.Global.CurVCO = s.Global.CurVCO
	o.Global.DataSeq = s.Global.DataSeq
	o.Global.PLLFlag = s.Global.PLLFlag
	o.Global.StatInp = s.Global.StatInp
	o.Global.MiscInp = s.Global.MiscInp
	o.Global.CurSequence = s.Global.CurSequence

	// reformat a subset of boom stats
	o.Boom.Booms = s.Boom.Booms
	o.Boom.AmbPos = s.Boom.AmbPos
	o.Boom.AmbNeg = s.Boom.AmbNeg
	o.Boom.Supply = s.Boom.Supply
	o.Boom.SysTemp = s.Boom.SysTemp
	o.Boom.MainCur = s.Boom.MainCur
	o.Boom.AntCur = s.Boom.AntCur
	o.Boom.Seis1Temp = s.Boom.Seis1Temp
	o.Boom.Seis2Temp = s.Boom.Seis2Temp
	o.Boom.CalTimeouts = s.Boom.CalTimeouts

	// reformat a subset of ethernet stats
	o.Ether.Check = s.Ether.Check
	o.Ether.IOErrors = s.Ether.IOErrors
	o.Ether.PhyNum = s.Ether.PhyNum
	o.Ether.Unreach = s.Ether.Unreach
	o.Ether.Quench = s.Ether.Quench
	o.Ether.Echo = s.Ether.Echo
	o.Ether.Redirect = s.Ether.Redirect
	o.Ether.Runt = s.Ether.Runt
	o.Ether.CRCErr = s.Ether.CRCErr
	o.Ether.BCast = s.Ether.BCast
	o.Ether.UCast = s.Ether.UCast
	o.Ether.Good = s.Ether.Good
	o.Ether.Jabber = s.Ether.Jabber
	o.Ether.OutWin = s.Ether.OutWin
	o.Ether.TXOK = s.Ether.TXOK
	o.Ether.Miss = s.Ether.Miss
	o.Ether.Collide = s.Ether.Collide
	o.Ether.LinkStat = s.Ether.LinkStat

	// reformat a subset of gps stats
	o.GPS.GPSOn = s.GPS.GPSOn
	o.GPS.GPSTime = s.GPS.GPSTime
	o.GPS.SatUsed = s.GPS.SatUsed
	o.GPS.SatView = s.GPS.SatView
	o.GPS.Fix = strings.TrimRight((string)(s.GPS.Fix[1:]), "\u0000")
	o.GPS.Time = strings.TrimRight((string)(s.GPS.Time[1:]), "\u0000")
	o.GPS.Date = strings.TrimRight((string)(s.GPS.Date[1:]), "\u0000")
	o.GPS.Height = strings.TrimRight((string)(s.GPS.Height[1:]), "\u0000")
	o.GPS.Lat = strings.TrimRight((string)(s.GPS.Lat[1:]), "\u0000")
	o.GPS.Lon = strings.TrimRight((string)(s.GPS.Lon[1:]), "\u0000")
	o.GPS.LastGood = QTime(s.GPS.LastGood)
	o.GPS.CheckErr = s.GPS.CheckErr

	// reformat a subset of logical port stats
	for lp, _ := range s.LPort {
		o.LPort[lp].Port = (uint8)(lp)
		o.LPort[lp].Sent = s.LPort[lp].Sent
		o.LPort[lp].Resends = s.LPort[lp].Resends
		o.LPort[lp].Fill = s.LPort[lp].Fill
		o.LPort[lp].Seq = s.LPort[lp].Seq
		o.LPort[lp].PackUsed = s.LPort[lp].PackUsed
		o.LPort[lp].LastAck = s.LPort[lp].LastAck
		o.LPort[lp].PhyNum = s.LPort[lp].PhyNum
		o.LPort[lp].LogNum = s.LPort[lp].LogNum
		o.LPort[lp].Retran = s.LPort[lp].Retran
	}

	// tag via timestamp
	o.Timestamp = time.Now().UTC()

	// done
	return &o
}
