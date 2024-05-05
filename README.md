# SQL-Parser

## Introduction

This document describes the code for an SQL parser written in Go. The parser takes an SQL statement as input and performs the following tasks:

- **Tokenization**: Breaks down the SQL string into individual tokens (keywords, identifiers, operators, constants).
- **Syntax Checking**: Verifies the basic structure of the statement for errors like unmatched parentheses or missing semicolons.
- **Parsing**: Validates the token order and structure according to the SQL grammar.
- **Semantic Analysis**: Interprets the parsed structure and assigns meaning based on the provided schema information.

## Dependencies

The code relies on the following dependencies:

- `errors` package for error handling
- `fmt` package for formatting (used for debugging purposes)

## Code Structure

The code is organized into the following parts:

- `parser` package: Contains the core parser logic.
    - `Token` struct: Represents a parsed token with its type and value.
    - `ParsedStmt` interface: Defines methods for accessing information from the parsed statement (query type, tables, columns, conditions).
    - `baseOperation` interface: Base interface for parsed statements to share common methods.
    - `SQLParser` struct: Manages the parsing process and holds the schema information.
    - Various functions for:
        - Parsing SQL statements (`ParseSQL`)
        - Tokenizing the input string (`tokenize`)
        - Performing syntax checks (`syntaxCheck`)
        - Parsing token structure (`parse`)
        - Performing semantic analysis (`semanticAnalysis`)
- `schema.go` file (optional): Contains the `Schema` struct and related functions for loading schema information from a metadata file.

## Key Functions

- `NewSQLParser(schema Schema) *SQLParser`: Creates a new SQL parser instance with the provided schema. Option 1: Schema is loaded here.
- `func (parser *SQLParser) ParseSQL(sql string) (parsedStmt ParsedStmt, warnings []string, err error)`: Parses the given SQL statement and returns the parsed representation, along with any warnings or errors encountered. Optionally accepts a schema name for multi-schema support.
- `func (schema *Schema) Load(filename string) error`: Loads schema information from a metadata file. (Defined in `schema.go`)

## Usage
``` // Option 1 (schema loaded in constructor)
    parser := NewSQLParser(schema) // Assuming schema is already defined

    parsedStmt, warnings, err := parser.ParseSQL(sql)
    if err != nil {
    // Handle error
    }

    // Option 2 (schema loaded explicitly)
    parser := NewSQLParser()
    err := parser.LoadSchema("metadata.json")
    if err != nil {
    // Handle error
    }

    parsedStmt, warnings, err := parser.ParseSQL(sql)
    if err != nil {
    // Handle error
    }
```
-----------------

## To Do (short term urgent)
### semantic analysis
1. **support `*` in select statement.**
   -  **approach:** decide wether to be handled in the parser or in the semantic analysis or in the tokenization.
2. **support mulitple tables in the from clause.**
   -  **approach:** ask the user to decide this column belongs to which table of them (first step to implmenting the joins).
3. **support any conditions not just column with value.**
   -  **approach:** make the condition more generic to support any condition.
4- **support if exist in the drop table statement.**
   -  **approach:** add the support for the if exist in the drop table statement. 




## To Do (long term)
1. Multi-schema support to be implemented by modifying the parser to manage and utilize multiple schema instances.
    - **appraoch:** just parse `use schema` statement and load the schema meta accordingly.
2. Support for more SQL statements and clauses to be added.
    - **appraoch:** add more methods to the `SQLParser` struct for parsing different types of statements.
3. support for more SQL operations to be added like join, group by, order by etc.

