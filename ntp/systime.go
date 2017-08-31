package ntp

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"
)

var ostype = runtime.GOOS

func linux(tm time.Time) error {

	settime := fmt.Sprintf("date -s \"%d/%d/%d %02d:%02d:%02d.%03d\"",
		tm.Year(), tm.Month(), tm.Day(),
		tm.Hour(), tm.Minute(), tm.Second(),
		time.Duration(tm.Nanosecond())/time.Millisecond)

	log.Println(settime)

	log.Println("set time to [", tm, "]")

	cmd := exec.Command("/bin/bash", "-c", settime)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func win32(tm time.Time) error {
	settime := fmt.Sprintf("date %d/%d/%d && time %02d:%02d:%02d.%02d",
		tm.Year(), tm.Month(), tm.Day(),
		tm.Hour(), tm.Minute(), tm.Second(),
		time.Duration(tm.Nanosecond())/(10*time.Millisecond))

	log.Println(settime)

	log.Println("set time to [", tm, "]")

	cmd := exec.Command("cmd.exe", "-c", settime)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func SetTimeToOs(tm time.Time) error {

	switch ostype {
	case "windows":
		{
			return win32(tm)
		}
	case "linux":
		{
			return linux(tm)
		}
	default:
		{
			return errors.New("not support local host os " + ostype)
		}
	}
}
