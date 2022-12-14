package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

type CreateRequest struct {
	RootDrivePath string `json:"root_image_path"`
	KernelPath    string `json:"kernel_path"`
}

type CreateResponse struct {
	IpAddress string `json:"ip_address"`
	ID        string `json:"id"`
}

type DeleteRequest struct {
	ID string `json:"id"`
}

var (
	runningVMs      = make(map[string]RunningFirecracker)
	ipByte     byte = 3
)

func main() {
	http.HandleFunc("/create", createRequestHandler)
	http.HandleFunc("/delete", deleteRequestHandler)
	defer cleanup()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func cleanup() {
	for _, running := range runningVMs {
		shutDown(running)
	}
}

func shutDown(running RunningFirecracker) {
	running.machine.StopVMM()
	os.Remove(running.image)
}

func makeIso(cloudInitPath string) (string, error) {
	image := "/tmp/cloud-init.iso"
	metaDataPath := "/tmp/my-meta-data.yml"
	err := os.WriteFile(metaDataPath, []byte("instance-id: i-litchi12345"), 0o644)
	if err != nil {
		return "", fmt.Errorf("Failed to create metadata file: %s", err)
	}
	if err := exec.Command("cloud-localds", image, cloudInitPath, metaDataPath).Run(); err != nil {
		return "", fmt.Errorf("cloud-localds failed: %s", err)
	}
	return image, nil
}

func deleteRequestHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("failed to read body, %s", err)
	}
	var req DeleteRequest
	json.Unmarshal([]byte(body), &req)
	if err != nil {
		log.Fatalf(err.Error())
	}

	running := runningVMs[req.ID]
	shutDown(running)
	delete(runningVMs, req.ID)
}

func createRequestHandler(w http.ResponseWriter, r *http.Request) {
	ipByte += 1
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("failed to read body, %s", err)
	}
	var req CreateRequest
	json.Unmarshal([]byte(body), &req)
	opts := getOptions(ipByte, req)
	running, err := opts.createVMM(context.Background())
	if err != nil {
		log.Fatalf(err.Error())
	}

	id := pseudo_uuid()
	resp := CreateResponse{
		IpAddress: opts.FcIP,
		ID:        id,
	}
	response, err := json.Marshal(&resp)
	if err != nil {
		log.Fatalf("failed to marshal json, %s", err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)

	runningVMs[id] = *running

	go func() {
		defer running.cancelCtx()
		// there's an error here but we ignore it for now because we terminate
		// the VM on /delete and it returns an error when it's terminated
		running.machine.Wait(running.ctx)
	}()
}

func pseudo_uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("failed to generate uuid, %s", err)
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func getOptions(id byte, req CreateRequest) options {
	fc_ip := net.IPv4(172, 102, 0, id).String()
	gateway_ip := "172.102.0.1"
	docker_mask_long := "255.255.255.0"
	bootArgs := "ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules random.trust_cpu=on "
	bootArgs = bootArgs + fmt.Sprintf("ip=%s::%s:%s::eth0:off", fc_ip, gateway_ip, docker_mask_long)
	return options{
		FcBinary:        "/usr/bin/firecracker",
		Request:         req,
		FcKernelCmdLine: bootArgs,
		FcSocketPath:    fmt.Sprintf("/tmp/firecracker-%d.sock", id),
		TapMacAddr:      fmt.Sprintf("02:FC:00:00:00:%02x", id),
		TapDev:          fmt.Sprintf("fc-tap-%d", id),
		FcIP:            fc_ip,
		FcCPUCount:      1,
		FcMemSz:         512,
	}
}

type RunningFirecracker struct {
	ctx       context.Context
	cancelCtx context.CancelFunc
	image     string
	machine   *firecracker.Machine
}

func (opts *options) createVMM(ctx context.Context) (*RunningFirecracker, error) {
	vmmCtx, vmmCancel := context.WithCancel(ctx)
	rootImagePath, err := copyImage(opts.Request.RootDrivePath)
	opts.Request.RootDrivePath = rootImagePath
	if err != nil {
		return nil, fmt.Errorf("Failed copying root path: %s", err)
	}
	fcCfg, err := opts.getConfig()
	if err != nil {
		return nil, err
	}

	cmd := firecracker.VMCommandBuilder{}.
		WithBin(opts.FcBinary).
		WithSocketPath(fcCfg.SocketPath).
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Build(ctx)

	machineOpts := []firecracker.Opt{
		firecracker.WithProcessRunner(cmd),
	}
	exec.Command("ip", "link", "del", opts.TapDev).Run()
	if err := exec.Command("ip", "tuntap", "add", "dev", opts.TapDev, "mode", "tap").Run(); err != nil {
		return nil, fmt.Errorf("Failed creating ip link: %s", err)
	}
	if err := exec.Command("rm", "-f", opts.FcSocketPath).Run(); err != nil {
		return nil, fmt.Errorf("Failed to delete old socket path: %s", err)
	}
	if err := exec.Command("ip", "link", "set", opts.TapDev, "master", "firecracker0").Run(); err != nil {
		return nil, fmt.Errorf("Failed adding tap device to bridge: %s", err)
	}
	if err := exec.Command("ip", "link", "set", opts.TapDev, "up").Run(); err != nil {
		return nil, fmt.Errorf("Failed creating ip link: %s", err)
	}
	if err := exec.Command("sysctl", "-w", fmt.Sprintf("net.ipv4.conf.%s.proxy_arp=1", opts.TapDev)).Run(); err != nil {
		return nil, fmt.Errorf("Failed doing first sysctl: %s", err)
	}
	if err := exec.Command("sysctl", "-w", fmt.Sprintf("net.ipv6.conf.%s.disable_ipv6=1", opts.TapDev)).Run(); err != nil {
		return nil, fmt.Errorf("Failed doing second sysctl: %s", err)
	}
	m, err := firecracker.NewMachine(vmmCtx, *fcCfg, machineOpts...)
	if err != nil {
		return nil, fmt.Errorf("Failed creating machine: %s", err)
	}
	if err := m.Start(vmmCtx); err != nil {
		return nil, fmt.Errorf("Failed to start machine: %v", err)
	}
	installSignalHandlers(vmmCtx, m)
	return &RunningFirecracker{
		ctx:       vmmCtx,
		image:     rootImagePath,
		cancelCtx: vmmCancel,
		machine:   m,
	}, nil
}

type options struct {
	Id string `long:"id" description:"Jailer VMM id"`
	// maybe make this an int instead
	IpId            byte   `byte:"id" description:"an ip we use to generate an ip address"`
	FcBinary        string `long:"firecracker-binary" description:"Path to firecracker binary"`
	FcKernelCmdLine string `long:"kernel-opts" description:"Kernel commandline"`
	Request         CreateRequest
	FcSocketPath    string `long:"socket-path" short:"s" description:"path to use for firecracker socket"`
	TapMacAddr      string `long:"tap-mac-addr" description:"tap macaddress"`
	TapDev          string `long:"tap-dev" description:"tap device"`
	FcCPUCount      int64  `long:"ncpus" short:"c" description:"Number of CPUs"`
	FcMemSz         int64  `long:"memory" short:"m" description:"VM memory, in MiB"`
	FcIP            string `long:"fc-ip" description:"IP address of the VM"`
}

func (opts *options) getConfig() (*firecracker.Config, error) {
	drives := []models.Drive{
		{
			DriveID:      firecracker.String("1"),
			PathOnHost:   &opts.Request.RootDrivePath,
			IsRootDevice: firecracker.Bool(true),
			IsReadOnly:   firecracker.Bool(false),
		},
	}
	//if opts.Request.CloudInitPath != "" {
	//	isoPath, err := makeIso(opts.Request.CloudInitPath)
	//	if err != nil {
	//		return nil, fmt.Errorf("Failed to create iso: %s", err)
	//	}
	//	drives = append(drives, models.Drive{
	//		DriveID:      firecracker.String("2"),
	//		PathOnHost:   &isoPath,
	//		IsRootDevice: firecracker.Bool(false),
	//		IsReadOnly:   firecracker.Bool(true),
	//	})
	//}

	return &firecracker.Config{
		VMID:            opts.Id,
		SocketPath:      opts.FcSocketPath,
		KernelImagePath: opts.Request.KernelPath,
		KernelArgs:      opts.FcKernelCmdLine,
		Drives:          drives,
		NetworkInterfaces: []firecracker.NetworkInterface{
			{
				StaticConfiguration: &firecracker.StaticNetworkConfiguration{
					MacAddress:  opts.TapMacAddr,
					HostDevName: opts.TapDev,
				},
				// AllowMMDS: allowMMDS,
			},
		},
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(opts.FcCPUCount),
			MemSizeMib: firecracker.Int64(opts.FcMemSz),
			// CPUTemplate: models.CPUTemplate(opts.FcCPUTemplate),
			// HtEnabled: firecracker.Bool(false),
		},
		// JailerCfg: jail,
		// VsockDevices:      vsocks,
		// LogFifo:           opts.FcLogFifo,
		// LogLevel:          opts.FcLogLevel,
		// MetricsFifo:       opts.FcMetricsFifo,
		// FifoLogWriter:     fifo,
	}, nil
}

func copyImage(src string) (string, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return "", err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer source.Close()

	destination, err := ioutil.TempFile("/images", "image")
	if err != nil {
		return "", err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return destination.Name(), err
}

func installSignalHandlers(ctx context.Context, m *firecracker.Machine) {
	// not sure if this is actually really helping with anything
	go func() {
		// Clear some default handlers installed by the firecracker SDK:
		signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		for {
			switch s := <-c; {
			case s == syscall.SIGTERM || s == os.Interrupt:
				log.Printf("Caught SIGINT, requesting clean shutdown")
				m.Shutdown(ctx)
			case s == syscall.SIGQUIT:
				log.Printf("Caught SIGTERM, forcing shutdown")
				m.StopVMM()
			}
		}
	}()
}
