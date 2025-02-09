# GoScope & rgs

GoScope (`goscope`) and rgs are two command-line tools designed to help Go developers quickly locate and navigate to functions, variables, constants, and types in Go projects. GoScope analyzes all `.go` files starting from the current directory (recursively), while rgs provides an interactive interface (via [fzf](https://github.com/junegunn/fzf)) to jump straight into your editor at the relevant line.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [GoScope](#goscope)
  - [rgs](#rgs)
- [Examples](#examples)
- [Notes and Tips](#notes-and-tips)
- [References](#references)

## Features

- **Automatic Scanning**: GoScope scans the current directory and all subdirectories for `.go` files.
- **Comprehensive Listing**: It outputs:
  - Declared functions (with file and line number).
  - Declared variables and constants (with file and line number).
  - Declared types (with file and line number).
  - Function call occurrences in the format `Caller.Callee file:line`.
- **Interactive Search**: The `rgs` script leverages `fzf` to allow fuzzy-searching any function, variable, constant, or type, and then opens the file at the exact line in your preferred text editor.

## Requirements

1. **GoScope**
   - [Go](https://go.dev/) 1.23 or later (to compile the tool).
   - A Unix-like shell environment (macOS, Linux, etc.).
2. **rgs**
   - [fzf](https://github.com/junegunn/fzf) for fuzzy searching.
   - [bat](https://github.com/sharkdp/bat) (optional but highly recommended for colorful previews).
   - An editor accessible via the `$EDITOR` environment variable (e.g., `vim`, `nvim`, etc.).
   - A Unix-like shell (e.g., Bash, Zsh).

## Installation

1. **Clone or Download** this repository.
2. **Build GoScope**:
   ```bash
   cd path/to/this/repo
   go build -o goscope
   ```
   This produces an executable named `goscope`.
3. **Install GoScope** (optional but recommended):
   ```bash
   mv goscope /usr/local/bin/
   ```
   Make sure `/usr/local/bin` is in your `$PATH`.
4. **Install rgs**:
   ```bash
   chmod +x rgs
   cp rgs /usr/local/bin/
   ```
   Again, confirm `/usr/local/bin` is in your `$PATH`.

## Usage

### GoScope

- Simply run `goscope` in the terminal.
- It scans the current directory recursively and prints all findings:
  - Functions, variables, constants, and types in the format:
    ```
    FunctionName file.go:line
    VariableName file.go:line
    ConstantName file.go:line
    TypeName file.go:line
    ```
  - Followed by function call mappings in the format:
    ```
    CallerName.CalleeName file.go:line
    ```
- There are **no command-line parameters** for GoScope: it always starts scanning from the current directory.

### rgs

- In the terminal, run:
  ```bash
  rgs
  ```
- This calls `goscope`, then pipes its output into `fzf`.
- You can type partial names to filter results.
- Use the arrow keys or your usual `fzf` navigation to select an item and press <kbd>Enter</kbd>.
- `rgs` then opens the corresponding file at the line where the item is declared (or called), using your editor set by `$EDITOR`.

## Examples

Below is a simple Go code snippet to show how GoScope output might appear:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
    PrintNumber(42)
}

func PrintNumber(num int) {
    fmt.Printf("Number: %d\n", num)
}
```

**Sample GoScope output**:
```
main            main.go:5
PrintNumber     main.go:9
main.Println    main.go:6
main.PrintNumber main.go:7
```

**Using rgs**:
1. In your project directory, type:
   ```bash
   rgs
   ```
2. Type `PrintNum` in the search prompt, and select the result.
3. Press <kbd>Enter</kbd> to open the file at the corresponding line.

## Notes and Tips

- Ensure `goscope` is in your `$PATH` so that `rgs` can invoke it properly.
- If `bat` is not installed, you may remove or modify the preview command in `rgs`.
- If `$EDITOR` is not set, you can either set it before running `rgs`, or update the script to call your favorite editor directly.

## References

- [Official Go Documentation](https://go.dev/doc/)
- [fzf on GitHub](https://github.com/junegunn/fzf)
- [bat on GitHub](https://github.com/sharkdp/bat)

