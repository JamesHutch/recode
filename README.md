# recode #

recode is a library to map structs with `interface{}` Fields to structs with proper types.

## Usage ##

```go
import "github.com/JamesHutch/recode"
```

recode is used in scenarios where you have input data of unknown type which is mapped
into a struct and you wish to remap this into a struct with more type safety

Exports a single `Recode()` function. Example:

```go
type Input struct {
	Data interface{}
}

type Output struct {
	Data string
}

i := Input{"Hello World"}
recode.Recode(i, &o)
fmt.Print(o.String) // will output: "Hello World"
i = Input{1}
recode.Recode(i, &o)
fmt.Print(o.String) // will output: "1"
```

(Notable missing features are bool and complex)

## License ##

Apache v2
