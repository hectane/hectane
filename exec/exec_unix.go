// +build !windows

package exec

// Run the application until terminated by a signal.
func Exec() error {
	execSignal()
	return nil
}
