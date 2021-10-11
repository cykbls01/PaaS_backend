package util

import (
	"time"
)


const timeTemplate = "2006/01/02 15:04:05"


func StringToTime(timeString string)time.Time{
	t,_ := time.Parse(timeTemplate,timeString)
	return t
}

func TimeToString(t time.Time)string{
	timeString := t.Format(timeTemplate)
	return timeString
}

func GetNowTime()string{
	var cstZone = time.FixedZone("CST", 8*3600)
	timeString := time.Now().In(cstZone).Format(timeTemplate)
	return timeString
}

func CompareTime(a,b time.Time)bool{
	return a.Before(b)
}
