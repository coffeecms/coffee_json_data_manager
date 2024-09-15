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

type DataManager struct {
    data map[string]map[string]interface{} // map theo key name
    mu   sync.RWMutex
}

func NewDataManager() *DataManager {
    return &DataManager{
        data: make(map[string]map[string]interface{}),
    }
}

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

        if key, ok := record[keyName].(string); ok {
            tempData[key] = record
        }
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    dm.mu.Lock()
    dm.data = tempData
    dm.mu.Unlock()

    return nil
}

func (dm *DataManager) GetDataByKey(key string) (map[string]interface{}, bool) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    data, exists := dm.data[key]
    return data, exists
}

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

type FilterCondition struct {
    Key       string      //  (ví dụ: "age", "fullname")
    ValueType string      // (ví dụ: "int", "string")
    Operator  string      // (ví dụ: ">", "==", "contains")
    Value     interface{} //  (ví dụ: 30, "James")
}

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

func OptimizeGoRuntime() {
    runtime.GOMAXPROCS(runtime.NumCPU()) 
    debug.SetGCPercent(100)              
}

func main() {
    OptimizeGoRuntime()

    dataManager := NewDataManager()

    err := dataManager.LoadData("users.json", "username") 
    if err != nil {
        log.Fatal(err)
    }

    dataManager.AutoReload("users.json", "username")

    conditions := []FilterCondition{
        {Key: "age", ValueType: "int", Operator: ">", Value: 30},
        {Key: "fullname", ValueType: "string", Operator: "contains", Value: "James"},
    }

    time.Sleep(6 * time.Second)
    filteredUsers := dataManager.FilterDataByConditions(conditions)

    log.Println("Filtered Users:", filteredUsers)

    select {}
}
