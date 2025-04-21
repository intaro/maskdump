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
  "cache_flush_count": 1000,
  "skip_insert_into_table_list": "/path/to/skip_table_list.txt",
  "masking": {
    "email": {
      "target": "username:2-",
      "value": "hash:6"
    },
    "phone": {
      "target": "2,3,5,6,8,10",
      "value": "hash"
    }
  }
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
- The characters to be replaced depend on the settings. By default, the configuration preserves only the first character and the domain, while all other characters in the email username are replaced with the first 6 characters of the MD5 hash of the entire original email string.
- Possible email masking configuration options:
  - `target="2-5"` and `value="hash"` - Numbers separated by a hyphen indicate character positions to replace; "hash" means replacement with characters from the MD5 hash of the original email. The range can be open-ended: "2-" replaces the second and all subsequent characters; "-5" replaces only the first five characters in the email.
  - `target="1~1"` and `value="hash"` - Numbers separated by a tilde indicate how many characters to keep unchanged at the start and end of the string (all others are replaced); "hash" means replacement with characters from the MD5 hash of the original email. The range can be open-ended: "2~" keeps the first and second characters unchanged, replacing all subsequent characters; "~1" replaces all characters except the last one.
  - `target="1,3,5,7"` and `value="hash"` - Comma-separated numbers indicate specific character positions to change; "hash" means replacement with characters from the MD5 hash of the original email.
  - Target can include modifiers: "username:" - modify only the left part of the email (before @) and "domain:" - modify only the right part of the email (after @). For example, `target="username:2-"` means replacing the second and all subsequent characters in the left part of the email, while everything else (first character, @ symbol, and right part of the email) remains unchanged.
  - `value="*"` - means replacement with asterisk characters

### Phone Numbers (`light-mask`)
- Preserves the original phone number format
- Replaces specific digits with digits from the SHA256 hash. Which digits get replaced is determined by settings. By default, digits at these positions are replaced: 2, 3, 5, 6, 8, and 10.

## A quick example of the work

The input is a typical database dump string. The output is the same dump, but with changed email and phone numbers:
```bash
$ echo "INSERT INTO users (id, email, phone) VALUES (123, 't098f6b@example.com', '+7 (904) 111-22-33'), (124, 'admin@site.org', '8-900-000-00-00');" | ./maskdump  --mask-email=light-hash --mask-phone=light-mask --no-cache
```
Result:
```bash
$ INSERT INTO users (id, email, phone) VALUES (123, 'ta6f5ce@example.com', '+7 (354) 101-72-53'), (124, 'a21232f@site.org', '8-700-160-90-20');
```

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
  "cache_flush_count": 1000,
  "skip_insert_into_table_list": "/path/to/skip_table_list.txt",
  "masking": {
    "email": {
      "target": "username:2-",
      "value": "hash:6"
    },
    "phone": {
      "target": "2,3,5,6,8,10",
      "value": "hash"
    }
  }
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
- Заменяемые символы зависят от настроек. По-умолчанию заданы настройки, при которых сохраняется только первый символ и домен, а все прочие символы имени пользователя из email заменяются на первые 6 символов MD5 хэша от всей исходной строки email
- Возможные варианты настроек макировки email:
  - `target="2-5"` и `value="hash"` — числа через дефис — это номера позиций символов, которые заменяем; "hash" — означает, что заменяем символами из MD5 хэша от исходного email. Диапазон заменяемых позиций символов может быть открытым: "2-" — замена второго и всех последующих симолов; "-5" — замена только первых пяти символов в email.
  - `target="1~1"` и `value="hash"` — числа через тильду — это количество символов, которые оставляем неизменными с начала и с конца строки (а все прочие заменяются); "hash" — означает, что заменяем символами из MD5 хэша от исходного email. Диапазон может быть открытым: "2~" — первый и второй символ остаются, в все последующие символы заменяются; "~1" — замена всех символов, кроме последнего.
  - `target="1,3,5,7"` и `value="hash"` — числа через запятую — это номера поиций символов, которые изменям; "hash" — означает, что заменяем символами из MD5 хэша от исходного email.
  - для target могут быть модификаторы: "username:" — изменять только левую часть email и "domain:" — изменять только правую часть email. Например, `target="username:2-"` означает замену второго и всех последующих символов левой части email, а всё остальное (первый символ, знак "@" и правая часть email) остаётся неизменным
  - `value="*"` — означает замену на символ звёздочки

### Телефоны (`light-mask`)
- Сохраняет исходный формат номера
- Заменяет определённые цифры на цифры из SHA256 хэша. Что попадает под замену — определяется настройками. По-умолчанию, заменяются цифры на этих номерах позиций: 2, 3, 5, 6, 8 и 10.

## Быстрый пример работы

На вход подаём строку типичного дампа базы данных. На выходе получаем этот же дамп, но с изменёнными email и телефонами:
```bash
$ echo "INSERT INTO users (id, email, phone) VALUES (123, 't098f6b@example.com', '+7 (904) 111-22-33'), (124, 'admin@site.org', '8-900-000-00-00');" | ./maskdump  --mask-email=light-hash --mask-phone=light-mask --no-cache
```
Результат:
```bash
$ INSERT INTO users (id, email, phone) VALUES (123, 'ta6f5ce@example.com', '+7 (354) 101-72-53'), (124, 'a21232f@site.org', '8-700-160-90-20');
```
