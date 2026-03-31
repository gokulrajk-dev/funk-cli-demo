<p align='center'>

```
                     $$$$$$\                      $$\       
                    $$  __$$\                     $$ |      
                    $$ /  \__|$$\   $$\ $$$$$$$\  $$ |  $$\ 
                    $$$$\     $$ |  $$ |$$  __$$\ $$ | $$  |
                    $$  _|    $$ |  $$ |$$ |  $$ |$$$$$$  / 
                    $$ |      $$ |  $$ |$$ |  $$ |$$  _$$<  
                    $$ |      \$$$$$$  |$$ |  $$ |$$ | \$$\ 
                    \__|       \______/ \__|  \__|\__|  \__|
```

</p>

---

> Lightweight suite of useful CLI tools to solve specific annoying problems faced by developers and Linux users alike, written in Go.

---

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Installation](#installation)
- [Commands](#commands)
  - [1. Unit Converter (`conv`)](#1-unit-converter-conv)
  - [2. Countdown Timer (`timer`)](#2-countdown-timer-timer)
  - [3. Todo Manager (`todo`)](#3-todo-manager-todo)
  - [4. Internet Speed Test (`internet_speed`)](#4-internet-speed-test-internet_speed)
  - [5. File Scanner (`sift`)](#5-file-scanner-sift)
- [Data Storage](#data-storage)
- [Project Structure](#project-structure)

---

## Overview

**funk** is a terminal multi-tool. a collection of small, focused utilities that solve everyday annoyances for developers, sysadmins, and terminal enthusiasts. Instead of Googling unit conversions, writing one-off scripts, or hunting for bloated files, you just run a funky `funk` command.

Funk was built to be simple, lightweight and easy to use so that no one has to open a slow chrome tab to perform simple tasks, get quick information. More tools will be added and with our unique installation script, users can pick and choose which tools they would like to install for their workflow. 


---

## Tech Stack

| Component        | Technology                                      |
| ---------------- | ----------------------------------------------- |
| Language         | Go 1.25.7                                       |
| CLI Framework    | [urfave/cli v3](https://github.com/urfave/cli)  |
| Database         | SQLite via `mattn/go-sqlite3`                    |
| Terminal UI      | `nsf/termbox-go`, `golang.org/x/term`           |
| Output Styling   | `fatih/color`, `alexeyco/simpletable`, `olekukonko/tablewriter` |

---

## Installation

```bash
git clone https://github.com/HUMANS-ORG/funk-cli.git
cd funk-cli
go build -o funk .
```

---

## Commands

### 1. Unit Converter (`conv`)

Converts between various units across multiple categories.

```
funk conv [INPUT FLAGS] <value> [CONVERSION FLAGS]
```

#### Distance

| Input Flag      | Alias | Description              |
| --------------- | ----- | ------------------------ |
| `--miles`       | `-M`  | Enter value in miles     |
| `--km`          | `-k`  | Enter value in kilometers|
| `--m`           | `-m`  | Enter value in meters    |

| Conversion Flag | Alias | Description              |
| --------------- | ----- | ------------------------ |
| `--to-km`       | `-tk` | Convert to kilometers    |
| `--to-miles`    | `-tM` | Convert to miles         |
| `--to-m`        | `-tm` | Convert to meters        |
| `--to-cm`       | `-tc` | Convert to centimeters   |

**Examples:**
```bash
funk conv -k 7 -tm      # 7 km →  7000 meters
funk conv -M 3 -tk      # 3 miles → 4.83 kilometers
funk conv -m 1500 -tk   # 1500 meters → 1.50 kilometers
```

#### Weight

| Input Flag | Alias | Description              |
| ---------- | ----- | ------------------------ |
| `--lbs`    | `-p`  | Enter value in pounds    |
| `--kg`     | `-w`  | Enter value in kilograms |

| Conversion Flag | Alias | Description         |
| --------------- | ----- | ------------------- |
| `--to-kg`       | `-tw` | Convert to kilograms|
| `--to-lbs`      | `-tp` | Convert to pounds   |
| `--to-gm`       | `-tg` | Convert to grams    |

**Examples:**
```bash
funk conv -p 150 -tw    # 150 lbs → 68.04 kg
funk conv -w 70 -tp     # 70 kg →  154.32 lbs
```

#### Temperature

| Input Flag      | Alias | Description                |
| --------------- | ----- | -------------------------- |
| `--celsius`     | `-c`  | Enter value in Celsius     |
| `--fahrenheit`  | `-f`  | Enter value in Fahrenheit  |

| Conversion Flag | Alias | Description                    |
| --------------- | ----- | ------------------------------ |
| `--to-f`        | `-tf` | Convert Celsius → Fahrenheit   |
| `--to-c`        | `-tC` | Convert Fahrenheit → Celsius   |

**Examples:**
```bash
funk conv -c 100 -tf    # 100°C → 212 °F
funk conv -f 72 -tC     # 72°F → 22.22 °C
```

#### Number System (Binary / Hexadecimal)

| Input Flag  | Alias | Description                  |
| ----------- | ----- | ---------------------------- |
| `--binary`  | `-b`  | Enter binary values (quoted) |
| `--hex`     | `-H`  | Enter hex values             |

| Conversion Flag | Alias | Description           |
| --------------- | ----- | --------------------- |
| `--to-hex`      | `-tH` | Binary → Hexadecimal  |
| `--to-binary`   | `-tb` | Hexadecimal → Binary  |

**Examples:**
```bash
funk conv -b "11010110" -tH    # Binary → Hex
funk conv -H "0xFF" -tb       # Hex → Binary
```

---

### 2. Countdown Timer (`timer`)

Set a countdown timer with a visual terminal display. Timer history is persisted to a local SQLite database.

```
funk timer [flags] [task name]
```

| Flag           | Alias      | Description                          |
| -------------- | ---------- | ------------------------------------ |
| `--sec`        | `-s`       | Timer duration in seconds            |
| `--min`        | `-m`       | Timer duration in minutes            |
| `--hr`         |            | Timer duration in hours              |
| `--his`        |            | Show timer history                   |
| `--del`        | `-d`, `-rm`| Delete a timer record by ID          |
| `--delete_all` |            | Delete all timer records             |

#### Features

- **Visual display**: Large block-digit characters rendered in the center of the terminal.
- **Color feedback**: Green digits turn red in the last 4 seconds.
- **Pause/Resume**: Press `Space` to pause or resume the timer.
- **Stop early**: Press `q` or `Ctrl+C` to stop the timer and save elapsed time.
- **History**: All timer sessions are saved to SQLite with date and optional task name.
- **Windows support**: Uses PowerShell + BurntToast notifications on Windows.

**Examples:**
```bash
funk timer --sec 30                 # 30 second timer
funk timer --min 5 "Study session"  # 5 min timer with a task name
funk timer --sec 90 --min 1         # 2 min 30 sec timer (combinable)
funk timer --his                    # Show timer history
funk timer --del 3                  # Delete record with ID 3
funk timer --delete_all             # Delete all records
```

---

### 3. Todo Manager (`todo`)

A simple task manager that persists tasks to a local SQLite database.

```
funk todo [flags]
```

| Flag    | Alias      | Description                       |
| ------- | ---------- | --------------------------------- |
| `--add` | `-a`, `-i` | Add a new task                    |
| `--del` | `-d`, `-rm`| Delete a task by ID               |
| `--task`| `-t`       | Show all tasks                    |
| `--done`| `-c`       | Mark a task as complete by ID     |

#### Features

- Tasks are stored with an auto-incremented ID, creation date, and status (`Incomplete` / `complete`).
- Listing tasks sorts incomplete tasks first.
- Marking a task as done prompts for confirmation (`y/n`).

**Examples:**
```bash
funk todo --add "Fix login bug"   # Add a task
funk todo --task                  # List all tasks
funk todo --done 2                # Mark task #2 as complete
funk todo --del 1                 # Delete task #1
```

---

### 4. Internet Speed Test (`internet_speed`)

Measures download and upload speeds using Cloudflare's speed test endpoints.

```
funk internet_speed --start
```
Alias: `funk is --start`

| Flag      | Description            |
| --------- | ---------------------- |
| `--start` | Start the speed test   |

#### Features

- **Multi-threaded**: Uses 8 parallel workers for both download and upload.
- **10-second test window**: Each direction (download/upload) runs for 10 seconds.
- **Live progress**: Real-time Mbps readout updated every second.
- **Download test**: Fetches 10 MB chunks from `speed.cloudflare.com`.
- **Upload test**: Sends 1 MB chunks to `speed.cloudflare.com`.

**Example:**
```bash
funk is --start
# Output:
# 📥 Downloading... 45.32 Mbps
# ✅ Final Download: 47.10 Mbps
# 📤 Uploading... 12.50 Mbps
# ✅ Final Upload: 13.20 Mbps
```

---

### 5. File Scanner (`sift`)

Scans a directory and reports on files matching various filters. Outputs are formatted in styled tables.

```
funk sift [flags]
```

| Flag    | Alias | Description                              |
| ------- | ----- | ---------------------------------------- |
| `--emt` | `-e`  | Find empty (0-byte) files                |
| `--rec` | `-r`  | Find files modified in last N days       |
| `--lrg` | `-l`  | Find files larger than N GB              |
| `--top` | `-t`  | Show top N largest files                 |
| `--ext` | `-x`  | Show file extension summary              |
| `--dup` | `-d`  | Find duplicate files (by name)           |
| `--log` | `-L`  | Find all `.log` files                    |
| `-p`    | `-P`  | Directory path to scan (default: `.`)    |

#### Features

- All flags can be combined in a single command.
- Results are displayed in formatted ASCII tables.
- Recursively walks the directory tree.
- Styled section headers and dividers for clean output.

**Examples:**
```bash
funk sift --emt                        # Find empty files in current dir
funk sift --rec 7 -p /var/log          # Files modified in last 7 days
funk sift --lrg 1.5                    # Files larger than 1.5 GB
funk sift --top 10                     # Top 10 largest files
funk sift --ext                        # Extension breakdown
funk sift --dup -p ~/Downloads         # Find duplicate filenames
funk sift --log -p /var                # Find all .log files
funk sift --emt --ext --top 5          # Combine multiple filters
```

---

## Data Storage

Funk uses a local **SQLite database** stored at `~/.funk.db`. Two tables are used:

| Table    | Used By | Columns                                          |
| -------- | ------- | ------------------------------------------------ |
| `Timer`  | `timer` | `id`, `create_at` (date), `timer` (HH:MM:SS), `name` |
| `tasks`  | `todo`  | `id`, `created_at` (datetime), `task`, `status`  |

Tables are auto-created on first use. No external database setup is required.

---

## Project Structure

```
funk-cli/
├── main.go                  # Entry point — registers all subcommands
├── go.mod                   # Go module definition & dependencies
├── go.sum                   # Dependency checksums
├── LICENSE                  # MIT License
├── README.md                # Project README
├── commands/
│   ├── convert.go           # Unit converter (conv)
│   ├── timer.go             # Countdown timer (timer)
│   ├── todo.go              # Task manager (todo)
│   ├── internet_speed.go    # Speed test (internet_speed)
│   └── fdetect.go           # File scanner (sift)
└── sqldb/
    ├── db.go                # SQLite connection + Timer table CRUD
    └── tododb.go            # Tasks table CRUD
```

---


## Who Needs funk?

* _**Developers**_ who spend 80% of their time in terminal (hello, fellow cool kids!).
* _**System Admins**_ who need quick solutions without writing 50-line scripts.
* _**Terminal Enthusiasts**_ who think GUIs are overrated.
* _**Power Users**_ who want more control over their files without leaving the terminal.
* _**Ricers of r/unix**_ who are looking for the next cool CLI tool for their thinkpad with Arch(btw). 

