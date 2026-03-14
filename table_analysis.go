package main

import (
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

// Регулярные выражения для парсинга SQL
var (
	createTableRegex = regexp.MustCompile(`CREATE TABLE ` + "`" + `(.+?)` + "`")
	fieldRegex       = regexp.MustCompile("`" + `(.+?)` + "`" + `\s+([^\s,]+)`)
	endTableRegex    = regexp.MustCompile(`\)[^)]*;`)
)

// TableParser keeps SQL dump parsing state isolated from package globals.
type TableParser struct {
	runtime         *Runtime
	tableInfos      map[string]*TableInfo
	currentTable    *TableInfo
	processingTable bool
	mutex           sync.Mutex
}

// NewTableParser creates an isolated parser state for selective dump processing.
func NewTableParser(runtime *Runtime) *TableParser {
	return &TableParser{
		runtime:    runtime,
		tableInfos: make(map[string]*TableInfo),
	}
}

// ParseTableStructure анализирует строку дампа и собирает информацию о таблицах
func (p *TableParser) ParseTableStructure(line string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	line = strings.TrimSpace(line)

	// Проверяем начало новой таблицы
	if matches := createTableRegex.FindStringSubmatch(line); matches != nil {
		tableName := matches[1]
		p.currentTable = &TableInfo{
			Name:   tableName,
			Fields: make([]FieldInfo, 0),
		}
		p.processingTable = true
		return
	}

	// Если мы в процессе обработки таблицы
	if p.processingTable && p.currentTable != nil {
		// Проверяем строки с определением полей
		if matches := fieldRegex.FindStringSubmatch(line); matches != nil {
			fieldName := matches[1]
			fieldType := matches[2]
			fieldPos := len(p.currentTable.Fields) + 1

			p.currentTable.Fields = append(p.currentTable.Fields, FieldInfo{
				Name:     fieldName,
				Type:     fieldType,
				Position: fieldPos,
			})
			return
		}

		// Проверяем конец определения таблицы
		if endTableRegex.MatchString(line) {
			p.tableInfos[p.currentTable.Name] = p.currentTable
			p.currentTable = nil
			p.processingTable = false
			return
		}
	}
}

// GetTableInfo возвращает информацию о таблице по имени
func (p *TableParser) GetTableInfo(tableName string) (*TableInfo, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	info, exists := p.tableInfos[tableName]
	return info, exists
}

// GetAllTables возвращает информацию о всех найденных таблицах
func (p *TableParser) GetAllTables() map[string]*TableInfo {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Создаем копию для безопасности
	copy := make(map[string]*TableInfo)
	for k, v := range p.tableInfos {
		copy[k] = v
	}
	return copy
}

// ProcessDumpLine обрабатывает строку дампа и возвращает модифицированную строку
func (p *TableParser) ProcessDumpLine(line string, config MaskConfig, cache *Cache) string {
	// Проверяем, является ли строка INSERT запросом
	matches := insertRegex.FindStringSubmatch(line)
	if len(matches) != 3 {
		return line // Не INSERT запрос, возвращаем как есть
	}

	tableName := matches[1]
	valuesPart := matches[2]

	// Проверяем, нужно ли обрабатывать эту таблицу
	tableConfig, ok := p.runtime.ProcessingTables[tableName]
	if !ok {
		return line
	}

	if !ok {
		return line // Таблица не в конфиге, пропускаем
	}

	// Получаем информацию о полях таблицы
	tableInfo, ok := p.tableInfos[tableName]
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
					masked := p.runtime.EmailRegex.ReplaceAllStringFunc(values[pos], func(email string) string {
						return p.runtime.MaskEmailWithRules(email, cache)
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
					masked := p.runtime.PhoneRegex.ReplaceAllStringFunc(values[pos], func(phone string) string {
						return p.runtime.MaskPhoneWithRules(phone, cache)
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

// ParseTableStructure keeps the legacy package-level API available.
func ParseTableStructure(line string) {
	defaultTableParser.ParseTableStructure(line)
}

// GetTableInfo keeps the legacy package-level API available.
func GetTableInfo(tableName string) (*TableInfo, bool) {
	return defaultTableParser.GetTableInfo(tableName)
}

// GetAllTables keeps the legacy package-level API available.
func GetAllTables() map[string]*TableInfo {
	return defaultTableParser.GetAllTables()
}

// ProcessDumpLine keeps the legacy package-level API available.
func ProcessDumpLine(line string, config MaskConfig, cache *Cache) string {
	return defaultTableParser.ProcessDumpLine(line, config, cache)
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
