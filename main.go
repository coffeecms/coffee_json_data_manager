package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// DataManager manages JSON data either in memory or in split mode
type DataManager struct {
	data         map[string]map[string]interface{}
	mu           sync.RWMutex
	maxRAMUsage  int64 // Max memory usage in bytes (default: 2GB)
	currentUsage int64
	mode         string // "InMemory" or "Split"
	index        map[string]map[string]int // Index for optimized search
	wg           sync.WaitGroup
}

// FilterCondition describes a filtering condition
type FilterCondition struct {
	Key       string      // Field name (e.g., "age", "fullname")
	ValueType string      // Data type (e.g., "int", "string", "datetime", "date", "bool")
	Operator  string      // Comparison operator (e.g., ">", "==", "contains", "<")
	Value     interface{} // Value to compare (e.g., 30, "James", "2024-01-01", true)
}

// NewDataManager creates a new DataManager instance
func NewDataManager(maxRAMUsage int64, mode string) *DataManager {
	return &DataManager{
		data:        make(map[string]map[string]interface{}),
		maxRAMUsage: maxRAMUsage,
		mode:        mode,
		index:       make(map[string]map[string]int),
	}
}

// LoadDataInMemory loads the entire JSON file into memory and creates index
func (dm *DataManager) LoadDataInMemory(filePath string, keyName string) error {
	if dm.mode != "InMemory" {
		return errors.New("Invalid mode for this operation")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	tempData := make(map[string]map[string]interface{})
	tempIndex := make(map[string]map[string]int)

	for scanner.Scan() {
		var record map[string]interface{}
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return err
		}

		if key, ok := record[keyName].(string); ok {
			tempData[key] = record

			// Create index for optimized search on keyName
			if _, exists := tempIndex[keyName]; !exists {
				tempIndex[keyName] = make(map[string]int)
			}
			// Assuming line numbers are keys for this example (you could use offsets)
			tempIndex[keyName][key] = len(tempData)
		}

		// Simulate RAM usage tracking
		dm.currentUsage += int64(len(line))
		if dm.currentUsage > dm.maxRAMUsage {
			return errors.New("Memory usage exceeds the maximum allowed limit")
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	dm.mu.Lock()
	dm.data = tempData
	dm.index = tempIndex
	dm.mu.Unlock()

	return nil
}

// LoadDataInSplitMode reads the JSON file in parts and filters data based on conditions
func (dm *DataManager) LoadDataInSplitMode(filePath string, conditions []FilterCondition) ([]map[string]interface{}, error) {
	if dm.mode != "Split" {
		return nil, errors.New("Invalid mode for this operation")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var filteredData []map[string]interface{}

	for scanner.Scan() {
		var record map[string]interface{}
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, err
		}

		// Apply filter conditions on each record
		if dm.matchConditions(record, conditions) {
			filteredData = append(filteredData, record)
		}

		// Track memory usage to ensure it doesn't exceed the limit
		dm.currentUsage += int64(len(line))
		if dm.currentUsage > dm.maxRAMUsage {
			return nil, errors.New("Memory usage exceeds the maximum allowed limit")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return filteredData, nil
}

// applyIntCondition applies integer-based filter conditions
func applyIntCondition(fieldValue interface{}, operator string, value interface{}) bool {
	fieldVal, ok := fieldValue.(float64) // JSON numbers are float64 by default
	if !ok {
		return false
	}
	compareVal, ok := value.(int)
	if !ok {
		return false
	}

	switch operator {
	case ">":
		return fieldVal > float64(compareVal)
	case ">=":
		return fieldVal >= float64(compareVal)
	case "<":
		return fieldVal < float64(compareVal)
	case "<=":
		return fieldVal <= float64(compareVal)
	case "==":
		return fieldVal == float64(compareVal)
	default:
		return false
	}
}

// applyStringCondition applies string-based filter conditions
func applyStringCondition(fieldValue interface{}, operator string, value interface{}) bool {
	fieldVal, ok := fieldValue.(string)
	if !ok {
		return false
	}
	compareVal, ok := value.(string)
	if !ok {
		return false
	}

	switch operator {
	case "contains":
		return strings.Contains(fieldVal, compareVal)
	case "==":
		return fieldVal == compareVal
	default:
		return false
	}
}

// applyDateTimeCondition applies datetime-based filter conditions
func applyDateTimeCondition(fieldValue interface{}, operator string, value interface{}) bool {
	fieldValStr, ok := fieldValue.(string)
	if !ok {
		return false
	}
	fieldVal, err := time.Parse("2006-01-02 15:04:05", fieldValStr)
	if err != nil {
		return false
	}

	compareValStr, ok := value.(string)
	if !ok {
		return false
	}
	compareVal, err := time.Parse("2006-01-02 15:04:05", compareValStr)
	if err != nil {
		return false
	}

	switch operator {
	case ">":
		return fieldVal.After(compareVal)
	case ">=":
		return fieldVal.After(compareVal) || fieldVal.Equal(compareVal)
	case "<":
		return fieldVal.Before(compareVal)
	case "<=":
		return fieldVal.Before(compareVal) || fieldVal.Equal(compareVal)
	case "==":
		return fieldVal.Equal(compareVal)
	default:
		return false
	}
}

// applyDateCondition applies date-based filter conditions
func applyDateCondition(fieldValue interface{}, operator string, value interface{}) bool {
	fieldValStr, ok := fieldValue.(string)
	if !ok {
		return false
	}
	fieldVal, err := time.Parse("2006-01-02", fieldValStr)
	if err != nil {
		return false
	}

	compareValStr, ok := value.(string)
	if !ok {
		return false
	}
	compareVal, err := time.Parse("2006-01-02", compareValStr)
	if err != nil {
		return false
	}

	switch operator {
	case ">":
		return fieldVal.After(compareVal)
	case ">=":
		return fieldVal.After(compareVal) || fieldVal.Equal(compareVal)
	case "<":
		return fieldVal.Before(compareVal)
	case "<=":
		return fieldVal.Before(compareVal) || fieldVal.Equal(compareVal)
	case "==":
		return fieldVal.Equal(compareVal)
	default:
		return false
	}
}

// applyBoolCondition applies boolean-based filter conditions
func applyBoolCondition(fieldValue interface{}, operator string, value interface{}) bool {
	fieldVal, ok := fieldValue.(bool)
	if !ok {
		return false
	}
	compareVal, ok := value.(bool)
	if !ok {
		return false
	}

	switch operator {
	case "==":
		return fieldVal == compareVal
	default:
		return false
	}
}

// matchConditions checks if a record matches the given filter conditions
func (dm *DataManager) matchConditions(record map[string]interface{}, conditions []FilterCondition) bool {
	for _, condition := range conditions {
		fieldValue, exists := record[condition.Key]
		if !exists {
			return false
		}

		switch condition.ValueType {
		case "int":
			if !applyIntCondition(fieldValue, condition.Operator, condition.Value) {
				return false
			}
		case "string":
			if !applyStringCondition(fieldValue, condition.Operator, condition.Value) {
				return false
			}
		case "datetime":
			if !applyDateTimeCondition(fieldValue, condition.Operator, condition.Value) {
				return false
			}
		case "date":
			if !applyDateCondition(fieldValue, condition.Operator, condition.Value) {
				return false
			}
		case "bool":
			if !applyBoolCondition(fieldValue, condition.Operator, condition.Value) {
				return false
			}
		default:
			return false
		}
	}

	return true
}

// Main function
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize DataManager
	dataManager := NewDataManager(2*1024*1024*1024, "Split") // Max 2GB RAM usage

	// Define filter conditions
	conditions := []FilterCondition{
		{Key: "age", ValueType: "int", Operator: ">", Value: 30},
		{Key: "fullname", ValueType: "string", Operator: "contains", Value: "James"},
		{Key: "status", ValueType: "bool", Operator: "==", Value: false},
		{Key: "ent_dt", ValueType: "datetime", Operator: "==", Value: "2024-09-03 09:00:00"},
	}

	// Load and filter data depending on the mode
	var filteredData []map[string]interface{}
	var err error
	if dataManager.mode == "InMemory" {
		// Load data into memory and apply filtering
		err = dataManager.LoadDataInMemory("users.json", "username")
		if err == nil {
			for _, record := range dataManager.data {
				if dataManager.matchConditions(record, conditions) {
					filteredData = append(filteredData, record)
				}
			}
		}
	} else if dataManager.mode == "Split" {
		// Process data in Split mode
		filteredData, err = dataManager.LoadDataInSplitMode("users.json", conditions)
	}

	if err != nil {
		log.Println("Error loading data:", err)
	} else {
		fmt.Println("Filtered Data:", filteredData)
	}

	// Wait for all goroutines to finish
	dataManager.wg.Wait()
}
