# Light Web Server

Легковесный одностраничный веб-сервер на Python для просмотра Markdown-документов, исходного кода и директорий в браузере. Поддерживает подсветку синтаксиса (Pygments) и рендеринг LaTeX-формул (KaTeX).

## Возможности

- **Markdown** — полный рендеринг через библиотеку `markdown` с расширениями: заголовки, параграфы, списки, таблицы, фрагменты кода, ссылки, изображения, сноски, admonitions и оглавление (TOC)
- **Подсветка синтаксиса** — 500+ языков через Pygments для код-блоков в Markdown и отдельных исходных файлов
- **LaTeX-формулы** — `$...$` (inline) и `$$...$$` (display) рендерятся на стороне клиента через KaTeX из CDN, без Python-зависимостей
- **Файловый браузер** — навигация по директориям с breadcrumbs и указанием размера файлов
- **Исходный код** — `.py`, `.c`, `.cpp`, `.js`, `.ts`, `.rs`, `.go` и многие другие рендерятся с подсветкой
- **Безопасность** — защита от path traversal (выход за пределы корневой директории запрещён), весь пользовательский контент экранируется через `html.escape()`
- **Единый файл** — весь сервер живёт в одном `ultra_simple.py` (~400 строк), без шаблонизаторов и отдельных HTML-файлов

## Требования

- Python 3.8+
- `markdown >= 3.5`
- `Pygments >= 2.16`

## Быстрый старт

```bash
# Установка зависимостей (рекомендуется venv из-за PEP 668)
python3 -m venv .venv && .venv/bin/pip install -r requirements.txt

# Запуск через скрипт-обёртку (покажет все доступные IP со ссылками)
~/bin/light_web_run.sh
~/bin/light_web_run.sh --ip 127.0.0.1 --port 9090 --directory /path/to/docs

# Или напрямую через venv
.venv/bin/python3 ultra_simple.py --host 0.0.0.0 --port 8080 --directory .
```

Откройте http://localhost:8080 в браузере (или по одному из IP-адресов, которые покажет скрипт).

### Параметры

| Параметр | Краткий | По умолчанию | Описание |
|----------|---------|--------------|----------|
| `--host` | `-h` | `0.0.0.0` | Адрес привязки (все интерфейсы) |
| `--port` | `-p` | `8080` | Порт |
| `--directory` | `-d` | `.` | Корневая директория |

> Примечание: в `light_web_run.sh` для `-h`/`-p`/`-d` используются `--help`, `--port` и `--directory` — там краткие флаги `-i`, `-p`, `-d`.

## Архитектура

```
light_web/
├── ultra_simple.py     # Весь сервер (~400 строк)
├── requirements.txt    # markdown + Pygments
└── light_web_run.sh    # Bash-обёртка для запуска
```

- **Основной класс**: `UltraSimpleHandler(BaseHTTPRequestHandler)` обрабатывает все запросы
- **HTML генерируется inline** — без шаблонизатора
- **LaTeX рендерится на клиенте** через KaTeX из CDN

### Структура кода

- `render_markdown()` — рендер Markdown с подсветкой синтаксиса и LaTeX
- `render_source_code()` — подсветка исходных файлов (`.py`, `.c`, `.cpp`, `.js`, и т.д.)
- `_math_on_html()` — замена `$...$` / `$$...$$` на разметку, совместимую с KaTeX
- `generate_directory_html()` — страница списка директорий
- `generate_markdown_html()` — обёртка Markdown-страницы с KaTeX-скриптами из CDN
- `generate_source_html()` — обёртка страницы с исходным кодом

## Поддерживаемые исходные файлы

`SOURCE_EXTS` в коде определяет распознаваемые расширения:

Python, C/C++, JS/TS, Java/Kotlin, Scala, Clojure, Rust, Go, Ruby, PHP, Swift, Perl, R, Objective-C, Shell/Bash/Zsh/Fish, Lua, Vim, SQL, CSS/SCSS/LESS/Sass, XML, YAML, TOML, JSON, Proto, Gradle, CMake, TeX, Haskell, Erlang, Elixir, Elm, Assembly, D, OCaml, Zig, Nim, Crystal, Racket и другие.

Специальные имена файлов: `Makefile`, `Dockerfile`, `CMakeLists.txt`, `.gitignore`, `.env.example`.

## Поведение по умолчанию

- Сервер по умолчанию привязывается к `0.0.0.0` (доступен из сети)
- Скрытые файлы (с префиксом `.`) исключаются из списка директорий
- Бинарные файлы отдаются как есть (с определённым MIME-типом)
- Markdown-файлы с ошибками кодировки fallback'ятся на бинарную отдачу

## Лицензия

Распространяется под лицензией [Apache-2.0](LICENSE).
