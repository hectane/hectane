package exec

import (
)


// Run the application either using signals or as a Windows service, depending
// on whether the -service flag was provided.
func Exec() {
	execSignal()
}
