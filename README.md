# JSON Data Manager - https://blog.lowlevelforest.com/

## Overview

JSON Data Manager is a Go library designed to efficiently manage and filter JSON data from files. It provides functionalities to:
- Load JSON data from a file.
- Automatically reload the file at regular intervals.
- Filter data based on dynamic conditions.
- Optimize Go runtime settings for performance.

## Features

- **Load JSON Data**: Load data from a JSON file and organize it by a specified key.
- **Auto-Reload**: Automatically reload data from the file every 5 seconds.
- **Dynamic Filtering**: Filter JSON data based on customizable conditions.
- **Runtime Optimization**: Optimize Go runtime settings for improved performance.

## Installation

To use this library, you need to have Go installed on your machine. You can then include this package in your project by cloning this repository or importing it as a module.

```sh
git clone https://github.com/coffeecms/coffee_json_data_manager.git
```

## Usage

### 1. Importing the Package

Import the package into your Go file:

```go
import (
    "log"
    "time"
    "github.com/coffeecms/coffee_json_data_manager"
)
```

### 2. Creating an Instance

Create an instance of `DataManager`:

```go
dataManager := coffee_json_filter.NewDataManager()
```

### 3. Loading Data

Load data from a JSON file and specify the key name used for organizing the data:

```go
err := dataManager.LoadData("users.json", "username")
if err != nil {
    log.Fatal(err)
}
```

### 4. Auto-Reload

Set up auto-reload to refresh the data every 5 seconds:

```go
dataManager.AutoReload("users.json", "username")
```

### 5. Filtering Data

Define filtering conditions and apply them to the data:

```go
conditions := []coffee_json_filter.FilterCondition{
    {Key: "age", ValueType: "int", Operator: ">", Value: 30},
    {Key: "fullname", ValueType: "string", Operator: "contains", Value: "James"},
}

filteredUsers := dataManager.FilterDataByConditions(conditions)
log.Println("Filtered Users:", filteredUsers)
```

### Example

Assume you have a file `users.json` with the following content:

```json
{"username": "user1", "age": 25, "fullname": "James Brown"}
{"username": "user2", "age": 35, "fullname": "Alice Jameson"}
{"username": "user3", "age": 40, "fullname": "James Smith"}
```

With the above configuration, the `FilterDataByConditions` function will return:

```go
[
    {"username": "user2", "age": 35, "fullname": "Alice Jameson"},
    {"username": "user3", "age": 40, "fullname": "James Smith"}
]
```

## Integration into Other Systems

You can integrate this library into other systems by:

1. **Incorporating it as a Module**: Include this library in your Go project using `go get` or by cloning it and importing it.

2. **Custom Data Sources**: Modify the `LoadData` function to work with different data sources or formats if needed.

3. **Advanced Filtering**: Customize the filtering logic based on specific requirements of your application, including adding new operators or value types.

4. **Runtime Configuration**: Utilize the `OptimizeGoRuntime` function to tailor Go runtime settings according to the needs of your system for better performance.

## Contributing

Feel free to contribute to this project by submitting issues, creating pull requests, or suggesting improvements. Please make sure to follow the contribution guidelines provided in the repository.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For any questions or support, please reach out to [lowlevelforest@gmail.com](mailto:lowlevelforest@gmail.com).

```
