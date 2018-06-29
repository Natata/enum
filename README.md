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

Example 1 (int enum)
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

Example 2 (int enum use iota keyword)
```
type: 
    int
name: 
    direction
list:
    West = iota
    East
    North
    South
```

Example 3 (string enum)
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

## Generate enum.go

```
enum -fp=exam.enum
```
