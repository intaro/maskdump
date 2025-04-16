# MaskDump - Database Anonymization Tool

## Description

MaskDump is a powerful tool for database anonymization and data masking designed to protect Personally Identifiable Information (PII) in database dumps. It provides secure PII obfuscation while maintaining data structure and format integrity.

Key features:
- Email and phone number masking with configurable algorithms
- White-list support for preserving specific values
- Caching system for consistent masking across multiple runs
- Regular expression customization for different data formats
- Pipeline-friendly design for integration with existing workflows

Use cases:
- Creating safe development/test environments from production data
- GDPR/CCPA compliance for data sharing
- Database sanitization before analytics processing
- Data masking for non-production environments

## Installation

### Build from source

1. Ensure you have Go installed (version 1.16+ recommended)
2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/maskdump.git
   cd maskdump
   ```
3. Build the binary:
   ```bash
   go build -o maskdump .
   ```

## Usage

### Basic pipeline usage

```bash
mysqldump dbname | ./maskdump --mask-email=light-hash --mask-phone=light-mask > anonymized_dump.sql
```

### Command-line options

| Option           | Description                                      | Default       |
|------------------|--------------------------------------------------|--------------|
| `--mask-email`   | Email masking algorithm (`light-hash`)           | (disabled)   |
| `--mask-phone`   | Phone masking algorithm (`light-mask`)           | (disabled)   |
| `--no-cache`     | Disable caching of masked values                 | false        |
| `--config`       | Path to configuration file                      | (autodetect) |

### Configuration file

Create `maskdump.conf` in the same directory as the binary or specify path with `--config`:

```json
{
  "cache_path": "/path/to/cache.json",
  "email_regex": "\\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}\\b",
  "phone_regex": "(?:\\+7|7|8)?(?:[\\s\\-\\(\\)]*\\d){10}",
  "email_white_list": "/path/to/white_list_email.txt",
  "phone_white_list": "/path/to/white_list_phone.txt",
  "memory_limit_mb": 1024,
  "cache_flush_count": 1000
}
```

### White lists

Create text files with one value per line to exclude from masking:

`white_list_email.txt` example:
```
admin@example.com
support@company.org
```

`white_list_phone.txt` example:
```
+79001234567
88005553535
```

## Masking Algorithms

### Email (`light-hash`)
- Preserves first character before @ and domain
- Hashes remaining local part with MD5 (first 6 chars of hash)

### Phone (`light-mask`)
- Preserves original phone number format
- Replaces specific digits (2,3,5,6,8,10) with SHA256 hash digits

---

# MaskDump - Инструмент анонимизации баз данных

## Описание

MaskDump - мощный инструмент для анонимизации баз данных и маскировки информации, предназначенный для защиты персональных данных (PII) в дампах БД. Обеспечивает безопасное преобразование данных с сохранением структуры и формата.

Основные возможности:
- Маскировка email и номеров телефонов с настраиваемыми алгоритмами
- Поддержка белых списков для исключения определённых значений
- Система кэширования для согласованного преобразования
- Настройка регулярных выражений для разных форматов данных
- Интеграция в существующие процессы обработки данных

Применение:
- Создание безопасных сред разработки/тестирования
- Обеспечение соответствия GDPR/CCPA
- Очистка данных перед аналитикой
- Маскировка данных для непродуктивных сред

## Установка

### Сборка из исходников

1. Убедитесь, что установлен Go (версия 1.16+)
2. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/yourusername/maskdump.git
   cd maskdump
   ```
3. Соберите бинарник:
   ```bash
   go build -o maskdump .
   ```

## Использование

### Использование в пайплайне

```bash
mysqldump dbname | ./maskdump --mask-email=light-hash --mask-phone=light-mask > anonymized_dump.sql
```

### Параметры запуска

| Параметр        | Описание                                       | По умолчанию |
|-----------------|-----------------------------------------------|--------------|
| `--mask-email`  | Алгоритм маскировки email (`light-hash`)      | (отключено)  |
| `--mask-phone`  | Алгоритм маскировки телефонов (`light-mask`) | (отключено)  |
| `--no-cache`    | Отключить кэширование                        | false        |
| `--config`      | Путь к конфигурационному файлу               | (автопоиск) |

### Конфигурационный файл

Создайте `maskdump.conf` в той же директории или укажите путь через `--config`:

```json
{
  "cache_path": "/path/to/cache.json",
  "email_regex": "\\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}\\b",
  "phone_regex": "(?:\\+7|7|8)?(?:[\\s\\-\\(\\)]*\\d){10}",
  "email_white_list": "/path/to/white_list_email.txt",
  "phone_white_list": "/path/to/white_list_phone.txt",
  "memory_limit_mb": 1024,
  "cache_flush_count": 1000
}
```

### Белые списки

Создайте текстовые файлы со значениями, которые не нужно маскировать:

Пример `white_list_email.txt`:
```
admin@example.com
support@company.org
```

Пример `white_list_phone.txt`:
```
+79001234567
88005553535
```

## Алгоритмы маскировки

### Email (`light-hash`)
- Сохраняет первый символ и домен
- Хеширует остальную часть с помощью MD5 (первые 6 символов от хэша)

### Телефоны (`light-mask`)
- Сохраняет исходный формат номера
- Заменяет определённые цифры (2,3,5,6,8,10) на цифры из SHA256 хэша
