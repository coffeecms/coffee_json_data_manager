# JSON Data Manager - https://blog.lowlevelforest.com/

## Overview

The JSON Data Manager is a flexible and efficient tool designed to manage, filter, and process JSON data. It supports two modes of operation: "InMemory" and "Split". The system can handle various data types including `int`, `string`, `datetime`, `date`, and `bool`. It also provides robust memory management and concurrency features to handle large datasets.

## Features

- **Modes of Operation**:
  - **InMemory**: Loads the entire JSON file into memory, creates an index for optimized searches, and applies filter conditions.
  - **Split**: Reads the JSON file in chunks, applies filter conditions on each chunk, and processes data efficiently without loading the entire file into memory.

- **Data Types Supported**:
  - **Integer**: Supports comparison operators such as `>`, `<`, `>=`, `<=`, `==`.
  - **String**: Supports equality check (`==`) and substring search (`contains`).
  - **Datetime**: Uses the format `yyyy-MM-dd HH:mm:ss`. Supports comparison operators such as `>`, `<`, `>=`, `<=`, `==`.
  - **Date**: Uses the format `yyyy-MM-dd`. Supports comparison operators such as `>`, `<`, `>=`, `<=`, `==`.
  - **Boolean**: Supports equality check (`==`).

- **Memory Management**:
  - Limits memory usage based on user-defined settings (default: 2GB).
  - Efficient memory tracking and error handling to avoid exceeding memory limits.

- **Concurrency**:
  - Utilizes all CPU cores with Goroutines for efficient data processing and filtering.
  - Handles concurrent operations using `sync.WaitGroup`.

## Installation

1. Ensure you have Go installed. You can download it from [golang.org](https://golang.org/dl/).

2. Clone this repository:
   ```bash
   git clone https://github.com/coffeecms/coffee_json_data_manager.git
   ```

3. Navigate to the project directory:
   ```bash
   cd coffee_json_filter
   ```

4. Build the project:
   ```bash
   go build -o coffee_json_filter
   ```

## Usage

### Configuration

Before running the program, you need to configure the `DataManager` instance and set the mode of operation. Modify the `main()` function in `main.go` to suit your needs.

### Running the Program

1. Place your JSON data file (e.g., `users.json`) in the project directory.

2. Update the `main()` function in `main.go` with the appropriate file path and filter conditions.

3. Run the program:
   ```bash
   ./coffee_json_filter
   ```

### Example

#### Configuration for `InMemory` Mode

```go
dataManager := NewDataManager(2*1024*1024*1024, "InMemory") // 2GB RAM limit
```

#### Configuration for `Split` Mode

```go
dataManager := NewDataManager(2*1024*1024*1024, "Split") // 2GB RAM limit
```

#### Filter Conditions

```go
conditions := []FilterCondition{
    {Key: "age", ValueType: "int", Operator: ">", Value: 30},
    {Key: "fullname", ValueType: "string", Operator: "contains", Value: "James"},
    {Key: "created_at", ValueType: "datetime", Operator: ">", Value: "2024-01-01 00:00:00"},
    {Key: "birthday", ValueType: "date", Operator: "<", Value: "2000-01-01"},
    {Key: "active", ValueType: "bool", Operator: "==", Value: true},
}
```

### Notes

- Ensure the JSON file is properly formatted and contains the expected fields.
- The date and datetime formats must match the specified formats for accurate parsing and comparison.
- Adjust memory limits (`maxRAMUsage`) and file paths according to your requirements.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes. For bug reports and feature requests, open an issue on GitHub.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or further information, please contact [lowlevelforest@gmail.com](mailto:lowlevelforest@gmail.com).

