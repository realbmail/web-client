package main

import (
	"fmt"
	"github.com/blockchainstamp/fdlimit"
	"github.com/realbmail/web-client/srv"
	"github.com/realbmail/web-client/utils"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

type SysParam struct {
	version bool
	baseDir string
}

const (
	PidFileName = "pid"
)

var (
	param = &SysParam{}
)

func init() {
	flags := rootCmd.Flags()

	flags.BoolVarP(&param.version, "version",
		"v", false, "web_mta -v")
}

var rootCmd = &cobra.Command{
	Use: "web_mta",

	Short: "web_mta is a mail translate agent for web client",

	Long: `usage description::TODO::`,

	Run: mainRun,
}

func initSystem() error {

	if err := os.Setenv("GODEBUG", "netdns=go"); err != nil {
		return err
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	limit, err := fdlimit.Maximum()
	if err != nil {
		return fmt.Errorf("failed to retrieve file descriptor allowance:%s", err)
	}
	_, err = fdlimit.Raise(uint64(limit))
	if err != nil {
		return fmt.Errorf("failed to raise file descriptor allowance:%s", err)
	}
	return nil
}

func waitShutdownSignal() {

	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("\n>>>>>>>>>>ninja node start at pid(%s)<<<<<<<<<<\n", pid)
	path := filepath.Join(utils.WorkHome(param.baseDir), string(filepath.Separator), PidFileName)
	if err := os.WriteFile(path, []byte(pid), 0644); err != nil {
		fmt.Print("failed to write running pid", err)
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2)

	sig := <-sigCh
	srv.Inst().ShutDown()
	fmt.Printf("\n>>>>>>>>>>process finished(%s)<<<<<<<<<<\n", sig)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
func mainRun(_ *cobra.Command, _ []string) {
	if param.version {
		fmt.Println("\n==================================================")
		fmt.Printf("Version:\t%s\n", utils.Version)
		fmt.Printf("Build:\t\t%s\n", utils.BuildTime)
		fmt.Printf("Commit:\t\t%s\n", utils.Commit)
		fmt.Println("==================================================")
		return
	}

	if err := initSystem(); err != nil {
		panic(err)
	}
	if err := srv.Inst().Start(); err != nil {
		panic(err)
	}

	waitShutdownSignal()
}
