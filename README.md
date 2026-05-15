# Token Debt Analyzer (MVC-Scanner)

Token Debt Analyzer is a command-line tool designed to calculate the "Minimum Viable Context" (MVC) of source code files in your project. It helps identify files that carry a massive "token debt"—meaning they have too many dependencies or are too large themselves to be easily processed by Large Language Models (LLMs) with limited context windows.

## Concept

The tool works by analyzing the Abstract Syntax Tree (AST) of your code to discover its true dependencies (imports, requires, etc.). 
1. **AST Parsing:** It uses `go-tree-sitter` to parse files and precisely extract dependencies. Supported languages are Java, JavaScript, TypeScript, and Python.
2. **Token Counting:** It uses `tiktoken` (specifically the `cl100k_base` model used by OpenAI) to calculate the token count of each file.
3. **MVC Calculation:** The "Total Context" (MVC) of a file is the sum of its own tokens plus the tokens of all its direct dependencies.

> **Note on Depth Limit:** The dependency analysis is currently limited to a depth of **1**. This intentional constraint keeps the tool usable and highly discriminating, making it much easier to pinpoint architectural hotspots and immediate technical debt without exploding the context graph.

## Supported Languages

As of today, the analyzer fully supports extracting dependencies via AST parsing for the following languages:
- **Java** (`.java`)
- **JavaScript** (`.js`)
- **TypeScript** (`.ts`)
- **Python** (`.py`)

## Why Go?

This tool was built in Go for several key reasons:
- **Performance & Speed:** Go compiles to native machine code, making file traversal and text processing incredibly fast.
- **AST Parsing:** The `go-tree-sitter` bindings are robust and fast, allowing for efficient parsing of multiple languages.
- **Concurrency Potential:** Go's goroutines make it trivial to introduce multi-threading for scanning massive codebases in the future.
- **Portability:** Go allows compiling the tool into a single, static executable binary that can run anywhere without needing a runtime (like Node.js or Python) installed.

## How to Compile

Make sure you have Go installed (version 1.17+).

1. Clone the repository and navigate to the project root.
2. Download dependencies:
   ```bash
   go mod download
   ```
3. Build the executable:
   ```bash
   go build -o token-debt-analyzer main.go
   ```

## Running Tests

To run the unit tests for the analyzer package, simply execute:
```bash
go test ./...
```

## How to Use

Run the compiled executable and provide the path to the directory you want to scan:

```bash
./token-debt-analyzer scan /path/to/your/code
```

### Options

- `--output`: Specifies the output format. Available options: `text` (default), `json`, `csv`.
  ```bash
  ./token-debt-analyzer scan --output json .
  ```
- `--fail-on`: Fails the CI pipeline (exits with code 1) if at least one file reaches the specified grade or worse (e.g., `D` or `F`).
  ```bash
  ./token-debt-analyzer scan --fail-on D .
  ```

## Example Output

By default, the tool outputs a neat table grouped by grade in your terminal:

```text
========================================================================
TOKEN DEBT ANALYZER - SCAN RESULTS
Model: cl100k_base | Depth: 1 | Files scanned: 12
========================================================================

GRADE   FILE                                  SELF TOKENS   TOTAL TOKENS   DEPS   TOP CONTRIBUTOR
A       src/utils/helpers.ts                  120           120            0      -
B       src/components/Button.tsx             450           1250           2      src/theme/colors.ts (800)
C       src/services/UserService.ts           1200          25000          5      src/models/User.ts (15000)
F       src/controllers/MainController.ts     5000          120000         15     src/services/GodObject.ts (80000)

========================================================================
```

### Grades Explanation
- **A (Optimal):** < 2,000 tokens
- **B (Standard):** < 10,000 tokens
- **C (Heavy):** < 30,000 tokens
- **D (Warning):** < 100,000 tokens
- **F (Critical):** ≥ 100,000 tokens