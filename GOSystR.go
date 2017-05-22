package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

type HdwrGroup struct {
	Os       string
	MemTotal string
	MemFreed string
	MemUsed  string

	TotDisk     string
	UsedDisk    string
	FreeDisk    string
	UsedDiskPct string

	CPUIndexNo string
	VendorID   string
	CPUFamily  string
	NoCores    string
	CPUModel   string
	CPUSpeed   string

	Hostname string
	Uptime   string
	NumProcs string
	OsType   string
	Platform string
	HostID   string

	HdwrAddr    string
	CPUPrct     []string
	NetStat     []string

}

func ErrFunc(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

func GetHardwareData(w http.ResponseWriter, r *http.Request) {
	runtimeOS := runtime.GOOS
	// memory
	vmStat, err := mem.VirtualMemory()
	ErrFunc(err)

	// disk - start from "/" mount point for Linux
	// might have to change for Windows!!
	// don't have a Window to test this out, if detect OS == windows
	// then use "\" instead of "/"

	diskStat, err := disk.Usage("/")
	ErrFunc(err)

	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	ErrFunc(err)
	percentage, err := cpu.Percent(0, true)
	ErrFunc(err)

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	ErrFunc(err)

	// get interfaces MAC/hardware address
	interfStat, err := net.Interfaces()
	ErrFunc(err)

	cpuPrcnt, err := json.MarshalIndent(percentage, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	log.Println(string(cpuPrcnt))

	intPld, err := json.MarshalIndent(interfStat, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	log.Println(string(intPld))

	hardware := HdwrGroup{
		Os:       runtimeOS,
		MemTotal: strconv.FormatUint(vmStat.Total, 10),
		MemFreed: strconv.FormatUint(vmStat.Free, 10),
		MemUsed:  strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64),

		TotDisk:     strconv.FormatUint(diskStat.Total, 10),
		UsedDisk:    strconv.FormatUint(diskStat.Used, 10),
		FreeDisk:    strconv.FormatUint(diskStat.Free, 10),
		UsedDiskPct: strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64),

		CPUIndexNo: strconv.FormatInt(int64(cpuStat[0].CPU), 10),
		VendorID:   cpuStat[0].VendorID,
		CPUFamily:  cpuStat[0].Family,
		NoCores:    strconv.FormatInt(int64(cpuStat[0].Cores), 10),
		CPUModel:   cpuStat[0].ModelName,
		CPUSpeed:   strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64),

		Hostname: hostStat.Hostname,
		Uptime:   strconv.FormatUint(hostStat.Uptime, 10),
		NumProcs: strconv.FormatUint(hostStat.Procs, 10),
		OsType:   hostStat.OS,
		Platform: hostStat.Platform,
		HostID:   hostStat.HostID,
	}
	//set a test value to the struct hardware
	//hardware.HdwrAddr = "HOLLER"

	HdwrPayload, err := json.MarshalIndent(hardware, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}

	for idx, cpupercent := range percentage {

		log.Println(int(idx))
		log.Println(percentage)
		log.Println(cpupercent)

		idxPld, err := json.MarshalIndent(percentage, "", "    ")
		if err != nil {
			fmt.Println("error:", err)
		}
		log.Println(string(idxPld))

	}

	pimple := []byte(HdwrPayload)

	log.Println(string(HdwrPayload))
	w.Write(pimple)
} ///end of the api handler

func NetHandler(w http.ResponseWriter, r *http.Request) {
	// get interfaces MAC/hardware address
	interfStat, err := net.Interfaces()
	ErrFunc(err)
	intPld, err := json.MarshalIndent(interfStat, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	log.Println(string(intPld))
	w.Write([]byte(intPld))
}

func cpuInfo(w http.ResponseWriter, r *http.Request) {
	// cpu - get CPU number of cores and speed

	percentage, err := cpu.Percent(0, true)
	ErrFunc(err)
	idxPld, err := json.MarshalIndent(percentage, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	log.Println(string(idxPld))
	w.Write([]byte(idxPld))
}




func WriteOut(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func main() {
	mux := http.NewServeMux()
	// mux.HandleFunc("/", SayName)
	mux.HandleFunc("/system", GetHardwareData)
	mux.HandleFunc("/net", NetHandler)
	mux.HandleFunc("/cpu", cpuInfo)
  // mux.HandleFunc("/disk", cpuInfo)
  mux.Handle("/", http.FileServer(http.Dir("./static/")))

  http.Handle("/", mux)
    err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
	// http.ListenAndServe(":8080", mux)
}
