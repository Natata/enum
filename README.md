# Introduce

enum is a tool for generating enum struct 

# Install 

```
go install github.com/Natata/enum
```

# How to use

## Write enum file

File format

```
type:
    [int|string]
name:
    [package name]
list:
    [enum element] = [value]
```

filename doesn't matter

## Generate enum.go

```
enum -fp=exam.enum
```

## Example

### int enum
```
type: 
    int
name: 
    direction
list:
    West = 0
    east = 1
    North = 2
    South = 3
```

result:
```
package direction

// Alias hide the real type of the enum 
// and users can use it to define the var for accepting enum 
type Alias = int

type list struct { 
    West Alias
    East Alias
    North Alias
    South Alias
}

// Enum for public use
var Enum = &list{ 
    West: 0,
    East: 1,
    North: 2,
    South: 3,
}
```

### int enum use iota keyword
```
type: 
    int
name: 
    direction
list:
    West = 9
    east = 8
    North = iota
    South
```

result:
```
package direction

// Alias hide the real type of the enum 
// and users can use it to define the var for accepting enum 
type Alias = int

type list struct { 
    West Alias
    East Alias
    North Alias
    South Alias
}

// Enum for public use
var Enum = &list{ 
    West: 9,
    East: 8,
    North: 0,
    South: 1,
}
```

### string enum
```
type: 
    string
name: 
    direction
list:
    West = "W"
    East = "E"
    North = "N"
    South = "S"
```

result:
```
package direction

// Alias hide the real type of the enum 
// and users can use it to define the var for accepting enum 
type Alias = string

type list struct { 
    West Alias
    East Alias
    North Alias
    South Alias
}

// Enum for public use
var Enum = &list{ 
    West: "W",
    East: "E",
    North: "N",
    South: "S",
}
```
