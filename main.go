package main

import (
    "bufio"
    "encoding/json"
    "log"
    "os"
    "strings"
    "sync"
    "time"
    "runtime"
    "runtime/debug"
)

// DataManager quản lý dữ liệu theo key name
type DataManager struct {
    data map[string]map[string]interface{} // map theo key name
    mu   sync.RWMutex
}

// Hàm khởi tạo DataManager
func NewDataManager() *DataManager {
    return &DataManager{
        data: make(map[string]map[string]interface{}),
    }
}

// LoadData đọc file JSON và cập nhật dữ liệu vào map theo key name
func (dm *DataManager) LoadData(filePath string, keyName string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    tempData := make(map[string]map[string]interface{})

    for scanner.Scan() {
        var record map[string]interface{}
        line := scanner.Text()
        if err := json.Unmarshal([]byte(line), &record); err != nil {
            return err
        }

        // Lấy giá trị key name và lưu vào map
        if key, ok := record[keyName].(string); ok {
            tempData[key] = record
        }
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    // Cập nhật dữ liệu với mutex để tránh race condition
    dm.mu.Lock()
    dm.data = tempData
    dm.mu.Unlock()

    return nil
}

// GetDataByKey lấy dữ liệu theo key
func (dm *DataManager) GetDataByKey(key string) (map[string]interface{}, bool) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    data, exists := dm.data[key]
    return data, exists
}

// AutoReload tự động tải lại file sau mỗi 5 giây
func (dm *DataManager) AutoReload(filePath string, keyName string) {
    ticker := time.NewTicker(5 * time.Second)

    go func() {
        for range ticker.C {
            err := dm.LoadData(filePath, keyName)
            if err != nil {
                log.Println("Error loading data:", err)
            } else {
                log.Println("Data reloaded successfully")
            }
        }
    }()
}

// FilterCondition mô tả điều kiện lọc
type FilterCondition struct {
    Key       string      // Tên trường cần lọc (ví dụ: "age", "fullname")
    ValueType string      // Kiểu dữ liệu (ví dụ: "int", "string")
    Operator  string      // Điều kiện lọc (ví dụ: ">", "==", "contains")
    Value     interface{} // Giá trị so sánh (ví dụ: 30, "James")
}

// FilterDataByConditions lọc dữ liệu dựa trên điều kiện tùy chọn
func (dm *DataManager) FilterDataByConditions(conditions []FilterCondition) []map[string]interface{} {
    dm.mu.RLock()
    defer dm.mu.RUnlock()

    var results []map[string]interface{}
    for _, record := range dm.data {
        match := true
        for _, condition := range conditions {
            fieldValue, exists := record[condition.Key]
            if !exists {
                match = false
                break
            }

            // Áp dụng điều kiện lọc theo kiểu dữ liệu
            switch condition.ValueType {
            case "int":
                if !applyIntCondition(fieldValue, condition.Operator, condition.Value) {
                    match = false
                }
            case "string":
                if !applyStringCondition(fieldValue, condition.Operator, condition.Value) {
                    match = false
                }
            default:
                match = false
            }

            if !match {
                break
            }
        }

        if match {
            results = append(results, record)
        }
    }

    return results
}

// applyIntCondition áp dụng điều kiện lọc cho kiểu int
func applyIntCondition(fieldValue interface{}, operator string, compareValue interface{}) bool {
    fieldInt, ok := fieldValue.(float64) // JSON trả về float64 cho số
    compareInt, ok2 := compareValue.(int)
    if !ok || !ok2 {
        return false
    }

    switch operator {
    case ">":
        return fieldInt > float64(compareInt)
    case "<":
        return fieldInt < float64(compareInt)
    case "==":
        return fieldInt == float64(compareInt)
    default:
        return false
    }
}

// applyStringCondition áp dụng điều kiện lọc cho kiểu string
func applyStringCondition(fieldValue interface{}, operator string, compareValue interface{}) bool {
    fieldStr, ok := fieldValue.(string)
    compareStr, ok2 := compareValue.(string)
    if !ok || !ok2 {
        return false
    }

    switch operator {
    case "contains":
        return strings.Contains(fieldStr, compareStr)
    case "==":
        return fieldStr == compareStr
    default:
        return false
    }
}

// Tối ưu hóa Go Runtime
func OptimizeGoRuntime() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // Sử dụng tối đa số CPU có sẵn
    debug.SetGCPercent(100)              // Điều chỉnh quá trình thu gom rác
}

func main() {
    // Tối ưu hóa Go Runtime
    OptimizeGoRuntime()

    // Khởi tạo DataManager
    dataManager := NewDataManager()

    // Load dữ liệu ban đầu từ file users.json
    err := dataManager.LoadData("users.json", "username") // Key name là "username"
    if err != nil {
        log.Fatal(err)
    }

    // Tự động đọc lại file sau mỗi 5 giây
    dataManager.AutoReload("users.json", "username")

    // Định nghĩa các điều kiện lọc (age > 30 và fullname chứa "James")
    conditions := []FilterCondition{
        {Key: "age", ValueType: "int", Operator: ">", Value: 30},
        {Key: "fullname", ValueType: "string", Operator: "contains", Value: "James"},
    }

    // Truy xuất và lọc dữ liệu theo các điều kiện tùy chọn
    time.Sleep(6 * time.Second) // Đợi dữ liệu được tải lại
    filteredUsers := dataManager.FilterDataByConditions(conditions)

    // In ra kết quả lọc
    log.Println("Filtered Users:", filteredUsers)

    // Chương trình tiếp tục chạy
    select {} // Giữ chương trình chạy để ticker hoạt động
}
