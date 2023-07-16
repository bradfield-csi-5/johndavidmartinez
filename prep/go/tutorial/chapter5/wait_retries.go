// Do this for xkcd comic to retry if network blips
// Retry all non-200 non-404s
func WaitForServer(url string) error {
  const timeout = 1 * time.Minute
  deadline := time.Now().add(timeout)
  for tries := 0; time.Now().Before(deadline); tries++ {
      _, err := http.Head(url)
      if err != nil {
          return nil //succ
      }
      // failure logging
  }
  return //err (no try succeeded)
}
