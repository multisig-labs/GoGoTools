package version

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
)

// GitCommit returns the git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Version returns the main version number that is being run at the moment.
const Version = "0.0.1"

// BuildDate returns the date the binary was built
var BuildDate = ""

// GoVersion returns the version of the go runtime used to compile the binary
var GoVersion = runtime.Version()

// OsArch returns the os and arch used to build the binary
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)

func PanicHandler() {
	if panicPayload := recover(); panicPayload != nil {
		stack := string(debug.Stack())
		fmt.Fprintln(os.Stderr, "================================================================================")
		fmt.Fprintln(os.Stderr, "            Encountered a fatal error. This is a bug!")
		fmt.Fprintln(os.Stderr, "================================================================================")
		fmt.Fprintf(os.Stderr, "Version:           %s\n", Version)
		fmt.Fprintf(os.Stderr, "Build Date:        %s\n", BuildDate)
		fmt.Fprintf(os.Stderr, "Git Commit:        %s\n", GitCommit)
		fmt.Fprintf(os.Stderr, "Go Version:        %s\n", GoVersion)
		fmt.Fprintf(os.Stderr, "OS / Arch:         %s\n", OsArch)
		fmt.Fprintf(os.Stderr, "Panic:             %s\n\n", panicPayload)
		fmt.Fprintln(os.Stderr, stack)
		os.Exit(1)
	}
}
