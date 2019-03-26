# Predicate

[![Build Status][1]][2] [![Coverage Status][3]][4] [![Go Report Card][5]][6]

Predicate is a simple library for parsing, evaluating and textually representing predicates (boolean functions).

## Install

```sh
git clone git@github.com:lamg/predicate.git
cd predicate/cmd/reduce && go install
```

## Example

Execute the command `reduce` with a predicate as argument, to get a reduced expression. Below appear several examples of the command's execution and its correspondent output shown after `→`.

```
reduce true → true
reduce ¬false → true
reduce ¬true → false
reduce 'true ∧ false' → false
reduce 'false ∧ false' → false
reduce 'false ∨ false' → false
reduce 'false ∨ true' → true
reduce '¬(true ∧ true)' → false
reduce '¬(true ∧ ¬A)' → ¬(¬A)
reduce 'A ∧ A' → A
reduce 'true ⇒ false' → false
reduce 'A ≡ true' → A
reduce 'A ≡ false' → ¬A
reduce 'A ≡ A' → true
reduce 'A ≡ ¬A' → false
reduce 'A ≢ A' → false
reduce 'A ⇐ true' → A
reduce 'A ≢ false' → true
```

## Syntax

The syntax is based on [EWD1300][0] which I have formalized in the following grammar:

```ebnf
predicate = term ('≡'|'≢') term {('≡'|'≢') term}| term.
term = implication | consequence | junction.
implication = junction '⇒' junction {'⇒' junction}.
consequence = junction '⇐' junction {'⇐' junction}.
junction = disjunction | conjunction | factor.
disjunction = factor '∨' factor {'∨' factor}.
conjunction = factor '∧' factor {'∧' factor}.
factor =	[unaryOp] (identifier | '(' predicate ')').
unaryOp = '¬'.
```

## Reduction rules

The procedure `Reduce` applies the following rules while reducing the predicate, i.e. if there's a constant (true, false) in the predicate this rules are applied.

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
