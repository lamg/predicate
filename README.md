# Predicate

[![Build Status][1]][2] [![Coverage Status][3]][4] [![Go Report Card][5]][6]

Predicate is a simple library for parsing, evaluating and textually representing predicates (boolean functions).

## Install

```sh
git clone git@github.com:lamg/predicate.git
cd predicate/cmd/reduce && go install
```

## Example

The following table shows some examples of how `reduce` works. Using standard input and output makes easier typing the boolean operators, since you can use Vim's multibyte input method (ex: C-k OR writes ∨), and then pipe the selected text using the visual mode to the `reduce` command, or just store the predicates in a file and then use it as standard input to `reduce` (`reduce < file_with_predicates`).

| Standard input  | Standard output |
|-----------------|-----------------|
| true            | true            |
| ¬false          | true            |
| ¬true           | false           |
| true ∧ false    | false           |
| false ∧ false   | false           |
| false ∨ false   | false           |
| false ∨ true    | true            |
| ¬(true ∧ true)  | false           |
| ¬(true ∧ ¬A)    | ¬(¬A)           |
| A ∧ A           | A               |
| true ⇒ false    | false           |
| A ≡ true        | A               |
| A ≡ false       | ¬A              |
| A ≡ A           | true            |
| A ≡ ¬A          | false           |
| A ≢ A           | false           |
| A ⇐ true        | A               |
| A ≢ false       | A               |

## Syntax

The syntax is based on [EWD1300][0] which I have formalized in the following grammar:

```ebnf
predicate = term {('≡'|'≢') term}.
term = junction ({'⇒' junction} | {'⇐' junction}).
junction = factor ({'∨' factor} | {'∧' factor}).
factor =	[unaryOp] (identifier | '(' predicate ')').
unaryOp = '¬'.
```

## Reduction rules

The procedure `Reduce` applies the following rules while reducing the predicate.

```
¬true ≡ false
¬false ≡ true
A ∨ false ≡ A
A ∧ true ≡ A
A ∨ true ≡ true
A ∧ false ≡ false
A ∨ B ≡ B ∨ A
A ∧ B ≡ B ∧ A
A ≡ true ≡ A
A ≡ false ≡ ¬A
true ⇒ A ≡ A
false ⇒ A ≡ true
A ⇒ true ≡ true
A ⇒ false ≡ ¬A
A ⇐ B ≡ B ⇒ A
A ≢ B ≡ A ≡ ¬B
```

[0]: https://www.cs.utexas.edu/users/EWD/transcriptions/EWD13xx/EWD1300.html
[1]: https://travis-ci.com/lamg/predicate.svg?branch=master
[2]: https://travis-ci.com/lamg/predicate
[3]: https://coveralls.io/repos/github/lamg/predicate/badge.svg?branch=master
[4]: https://coveralls.io/github/lamg/predicate?branch=master
[5]: https://goreportcard.com/badge/github.com/lamg/predicate
[6]: https://goreportcard.com/report/github.com/lamg/predicate
