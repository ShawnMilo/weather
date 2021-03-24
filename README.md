# Weather

Simple Go library using the [OpenWeather](https://openweathermap.org/) API, with caching.

Example usage:

```go
package main

import (
	"fmt"

	"github.com/shawnmilo/weather"
)

var zips = []string{"12345", "90210", "19406", "08861", "08772-2109", "potato", "00000"}

func main() {
	for _, zip := range zips {
		w, err := weather.Get(zip)
		if err != nil {
			fmt.Printf("%s: %s\n", zip, err)
			continue
		}
		fmt.Printf("%s: %.2f°, wind %.2f MPH at %d°\n", zip, w.Temperature, w.WindSpeed, w.WindDirection)
	}
}
```

Output:

```
12345: 52.86°, wind 3.44 MPH at 210°  
90210: 65.59°, wind 11.50 MPH at 340° 
19406: 55.00°, wind 5.75 MPH at 110°  
08861: 47.26°, wind 5.37 MPH at 110°  
08772-2109: no data for zip code      
potato: invalid zip code              
00000: invalid zip code               
```
