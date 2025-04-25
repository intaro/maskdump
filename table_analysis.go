package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// TableInfo хранит информацию о структуре таблицы
type TableInfo struct {
	Name   string
	Fields []FieldInfo
}

// FieldInfo хранит информацию о поле таблицы
type FieldInfo struct {
	Name     string
	Type     string
	Position int
}

// Глобальные переменные для хранения состояния
var (
	tableInfos      = make(map[string]*TableInfo)
	currentTable    *TableInfo
	processingTable bool
	mutex           sync.Mutex
)

// Регулярные выражения для парсинга SQL
var (
	createTableRegex = regexp.MustCompile(`CREATE TABLE ` + "`" + `(.+?)` + "`")
	fieldRegex       = regexp.MustCompile("`" + `(.+?)` + "`" + `\s+([^\s,]+)`)
	endTableRegex    = regexp.MustCompile(`\)[^)]*;`)
)

// ParseTableStructure анализирует строку дампа и собирает информацию о таблицах
func ParseTableStructure(line string) {
	mutex.Lock()
	defer mutex.Unlock()

	line = strings.TrimSpace(line)

	// Проверяем начало новой таблицы
	if matches := createTableRegex.FindStringSubmatch(line); matches != nil {
		tableName := matches[1]
		currentTable = &TableInfo{
			Name:   tableName,
			Fields: make([]FieldInfo, 0),
		}
		processingTable = true
		return
	}

	// Если мы в процессе обработки таблицы
	if processingTable && currentTable != nil {
		// Проверяем строки с определением полей
		if matches := fieldRegex.FindStringSubmatch(line); matches != nil {
			fieldName := matches[1]
			fieldType := matches[2]
			fieldPos := len(currentTable.Fields) + 1

			currentTable.Fields = append(currentTable.Fields, FieldInfo{
				Name:     fieldName,
				Type:     fieldType,
				Position: fieldPos,
			})
			return
		}

		// Проверяем конец определения таблицы
		if endTableRegex.MatchString(line) {
			tableInfos[currentTable.Name] = currentTable
			currentTable = nil
			processingTable = false
			return
		}
	}
}

// GetTableInfo возвращает информацию о таблице по имени
func GetTableInfo(tableName string) (*TableInfo, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	info, exists := tableInfos[tableName]
	return info, exists
}

// GetAllTables возвращает информацию о всех найденных таблицах
func GetAllTables() map[string]*TableInfo {
	mutex.Lock()
	defer mutex.Unlock()

	// Создаем копию для безопасности
	copy := make(map[string]*TableInfo)
	for k, v := range tableInfos {
		copy[k] = v
	}
	return copy
}

// ProcessDumpLine обрабатывает строку дампа и возвращает модифицированную строку
func ProcessDumpLine(line string, config MaskConfig, cache *Cache) string {
	// Проверяем, является ли строка INSERT запросом
	matches := insertRegex.FindStringSubmatch(line)
	if len(matches) != 3 {
		return line // Не INSERT запрос, возвращаем как есть
	}

	tableName := matches[1]
	valuesPart := matches[2]

	// Проверяем, нужно ли обрабатывать эту таблицу
	tableConfig, ok := ProcessingTables.Tables[tableName]

	if !ok {
		return line // Таблица не в конфиге, пропускаем
	}

	// Получаем информацию о полях таблицы
	tableInfo, ok := tableInfos[tableName]
	if !ok {
		return line // Нет информации о таблице, пропускаем
	}

	// Создаем карты позиций полей
	emailFields := make(map[int]bool)
	phoneFields := make(map[int]bool)

	if config.emailAlgorithm == "light-hash" {
		// Заполняем позиции полей для email
		for _, fieldName := range tableConfig.Email {
			for _, field := range tableInfo.Fields {
				if field.Name == fieldName {
					emailFields[field.Position-1] = true // переводим в 0-based индекс
					break
				}
			}
		}
	}

	if config.phoneAlgorithm == "light-mask" {
		// Заполняем позиции полей для phone
		for _, fieldName := range tableConfig.Phone {
			for _, field := range tableInfo.Fields {
				if field.Name == fieldName {
					phoneFields[field.Position-1] = true // переводим в 0-based индекс
					break
				}
			}
		}
	}

	if tableName == "b_user" {
		Log("DEBUG! b_user")
		LogStruct("TableConfig", tableConfig)
		Log(fmt.Sprintf("TableInfo for %s: %+v", tableName, tableInfo))
		Log(fmt.Sprintf("emailFields: %v", emailFields))
		Log(fmt.Sprintf("phoneFields: %v", phoneFields))
	}

	// Обрабатываем все кортежи в строке
	modified := false
	modifiedValues := tupleRegex.ReplaceAllStringFunc(valuesPart, func(tuple string) string {
		// Извлекаем значения из кортежа
		values := parseTuple(tuple)
		if len(values) == 0 {
			return tuple
		}

		if config.emailAlgorithm == "light-hash" {
			// Обрабатываем email поля
			for pos := range emailFields {
				if pos < len(values) && values[pos] != "" && values[pos] != "NULL" {
					//masked := maskEmailWithRules(values[pos], cache)
					masked := EmailRegex.ReplaceAllStringFunc(values[pos], func(email string) string {
						return maskEmailWithRules(email, cache)
					})
					if masked != values[pos] {
						values[pos] = masked
						modified = true
					}
				}
			}
		}

		if config.phoneAlgorithm == "light-mask" {
			// Обрабатываем phone поля
			for pos := range phoneFields {
				if pos < len(values) && values[pos] != "" && values[pos] != "NULL" {
					//masked := maskPhoneWithRules(values[pos], cache)
					masked := PhoneRegex.ReplaceAllStringFunc(values[pos], func(phone string) string {
						return maskPhoneWithRules(phone, cache)
					})
					if masked != values[pos] {
						values[pos] = masked
						modified = true
					}
				}
			}
		}

		// Собираем кортеж обратно
		return "(" + strings.Join(values, ",") + ")"
	})

	if !modified {
		return line // Ничего не изменилось, возвращаем оригинал
	}

	// Собираем модифицированную строку
	return "INSERT INTO `" + tableName + "` VALUES " + modifiedValues
}

// parseTuple разбирает кортеж значений на отдельные элементы
func parseTuple(tuple string) []string {
	tuple = strings.TrimPrefix(tuple, "(")
	tuple = strings.TrimSuffix(tuple, ")")

	var values []string
	var current strings.Builder
	inQuotes := false
	escape := false

	for _, c := range tuple {
		switch {
		case escape:
			current.WriteRune(c)
			escape = false
		case c == '\\':
			escape = true
		case c == '\'':
			inQuotes = !inQuotes
			current.WriteRune(c)
		case c == ',' && !inQuotes:
			values = append(values, current.String())
			current.Reset()
		default:
			current.WriteRune(c)
		}
	}

	if current.Len() > 0 {
		values = append(values, current.String())
	}

	return values
}
