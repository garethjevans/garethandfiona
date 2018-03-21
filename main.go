package main

import (
  "os"
)

func main() {
  a := App{}
  a.Initialize(
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
	os.Getenv("NEWRELIC_APP_NAME"),
	os.Getenv("NEWRELIC_LICENSE_KEY"),
  )

  a.Run(":8080")
}
