package main

import "goRedis/tcp"

var banner = `
   ______          ___
  / ____/___  ____/ (_)____
 / / __/ __ \/ __  / / ___/
/ /_/ / /_/ / /_/ / (__  )
\____/\____/\__,_/_/____/
`

func main() {
	print(banner)
	tcp.ListenAndServeDemo(":8000")
}
