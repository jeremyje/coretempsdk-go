# CoreTemp Go Bindings

Go bindings for ALCPU's Core Temp [shared memory interface](https://www.alcpu.com/CoreTemp/developers.html).

_For this library to work Core Temp must be running in the background and the [SDK's GetCoreTempInfo.dll](https://www.alcpu.com/CoreTemp/main_data/CoreTempSDK.zip) must be located in the same directory as the Go binary._

## Example

```go
package main

import (
 "log"

 "github.com/jeremyje/coretempsdk-go"
)

func main() {
 info, err := coretempsdk.GetCoreTempInfo()
 if err != nil {
  log.Printf("ERROR: %s", err)
  return
 }
 log.Printf("CPU: %s", info.CPUName)
 log.Printf("Temperatures: %v", info.Temperature)
}
```

## SDK DLL

[Download the Core Temp SDK](https://www.alcpu.com/CoreTemp/main_data/CoreTempSDK.zip) which contains the `GetCoreTempInfo.dll` to interact with Core Temp.
The preferred DLL is `Shared Memory Framework for Native\Shared Memory Libraries\x64\GetCoreTempInfo.dll`.

You can visit the [developer page](https://www.alcpu.com/CoreTemp/developers.html) for more details.
