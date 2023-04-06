package handlers

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func WriteHeaders(w http.ResponseWriter, extras []string) {
	// We want to write headers for each request, the content type and
	// the CORS settings
	extra := strings.Join(extras, ",")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, "+extra)
}

func getHexKey() string {
	// Returns a 'hex' key
	pos := "0123456789abcdef"
	key := ""
	for i := 0; i < 16; i++ {
		key += string(pos[rand.Intn(len(pos)-1)])
	}
	return key
}

func roundFloat(val float64, precision uint) float64 {
	// Round floats with precision
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func lenOfSyncMap(mp sync.Map) int {
	// We want to get the length of a sync map
	var i int
	mp.Range(func(key any, value any) bool {
		i++
		return true
	})
	return i
}

func isInArr(arr []string, needle string) bool {
	// Returns if a string is in an array
	for _, s := range arr {
		if s == needle {
			return true
		}
	}
	return false
}

func getWindowsStats() Stats {
	// Runs stats for windows
	ps, _ := exec.LookPath("powershell.exe")
	cpu := exec.Command(ps, "Get-CimInstance win32_processor | Measure-Object -Property LoadPercentage -Average")
	ramTotal := exec.Command(ps, "wmic ComputerSystem get TotalPhysicalMemory")
	ramFree := exec.Command(ps, "wmic OS get FreePhysicalMemory")

	//CPU
	var cpuout bytes.Buffer
	cpu.Stdout = &cpuout
	cpuerr := cpu.Run()
	if cpuerr != nil {
		fmt.Println(cpuerr)
	}

	cpuavreg, _ := regexp.Compile("Average  : [0-9]*")
	cpuav := cpuavreg.FindString(cpuout.String())
	cpureg, _ := regexp.Compile("[0-9]+")
	cpuStat, _ := strconv.ParseFloat(cpureg.FindString(cpuav), 64)

	//RAM
	var ramTotalVal bytes.Buffer
	var ramFreeVal bytes.Buffer

	//regex
	ramreg, _ := regexp.Compile("[0-9]+")

	ramTotal.Stdout = &ramTotalVal
	ramFree.Stdout = &ramFreeVal

	ramterr := ramTotal.Run()

	if ramterr != nil {
		fmt.Println(ramterr)
	}
	ramferr := ramFree.Run()
	if ramferr != nil {
		fmt.Println(ramferr)
	}

	ramFreeF, _ := strconv.ParseFloat(ramreg.FindString(ramFreeVal.String()), 32)
	ramTotalF, _ := strconv.ParseFloat(ramreg.FindString(ramTotalVal.String()), 32)

	ramPerc := roundFloat((((ramTotalF/1000)-ramFreeF)/(ramTotalF/1000))*100, 1)

	return Stats{cpuStat, ramPerc}
}

func getLinuxStats() Stats {
	// Runs stats for linux
	cpu := exec.Command("top", "-b", "-n", "1")
	//CPU
	var cpuout bytes.Buffer
	cpu.Stdout = &cpuout
	cpuerr := cpu.Run()
	if cpuerr != nil {
		fmt.Println(cpuerr)
	}

	cpuavreg, _ := regexp.Compile(",[ ]*[0-9.]+ id")
	cpuavregnum, _ := regexp.Compile("[0-9.]+")
	cpuav := cpuavreg.FindString(cpuout.String())
	cpuidle, _ := strconv.ParseFloat(cpuavregnum.FindString(cpuav), 32)

	cpuused := 100 - cpuidle

	//RAM
	ramavreg, _ := regexp.Compile("MiB Mem.*?free")
	ramavfreereg, _ := regexp.Compile("[0-9.]+ free")
	ramavtotalreg, _ := regexp.Compile("[0-9.]+ total")
	ramavnumreg, _ := regexp.Compile("[0-9.]+")

	ramfreeav := ramavreg.FindString(cpuout.String())
	ramfreestr := ramavfreereg.FindString(ramfreeav)
	ramfree, _ := strconv.ParseFloat(ramavnumreg.FindString(ramfreestr), 32)

	ramtotalav := ramavreg.FindString(ramfreeav)
	ramtotalstr := ramavtotalreg.FindString(ramtotalav)
	ramtotal, _ := strconv.ParseFloat(ramavnumreg.FindString(ramtotalstr), 32)

	ramperc := roundFloat((((ramtotal)-ramfree)/(ramtotal))*100, 1)
	return Stats{cpuused, ramperc}
}

func getMacStats() Stats {
	// Placeholder stats for Mac
	return Stats{
		1.1,
		2.3,
	}
}
