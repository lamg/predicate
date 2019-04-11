
import Data.Char

data Predicate = Neg Predicate | 
 Binary String Predicate Predicate | Value String
 deriving Show

predicate :: Parser Predicate
predicate = equivalence +++ term

term :: Parser Predicate
term = implication +++ consequence +++ junction +++ factor

equivalence :: Parser Predicate
equivalence = binaryOp (symbol eq +++ symbol neq) term

implication :: Parser Predicate
implication = binaryOp (symbol impl) junction

consequence :: Parser Predicate
consequence = binaryOp (symbol fll) junction

junction :: Parser Predicate
junction = disjunction +++ conjunction

disjunction :: Parser Predicate
disjunction = binaryOp (symbol orOp) factor

conjunction :: Parser Predicate
conjunction = binaryOp (symbol andOp) factor

factor :: Parser Predicate
factor = negation +++ identifier +++ parenPred

-- the previous region is as clean as the grammar

negation :: Parser Predicate
negation = seqp (symbol neg)
 (\i -> seqp (identifier +++ parenPred)
 (\j -> mreturn (Neg j)))

identifier :: Parser Predicate
identifier = seqp (token ident) (mreturn . Value)

parenPred :: Parser Predicate
parenPred = seqp (symbol "(")
 (\i -> seqp predicate
 (\j -> seqp (symbol ")")
 (\k -> mreturn j)))

binaryOp :: Parser String -> Parser Predicate -> Parser Predicate
binaryOp o p = seqp p
 (\i -> seqp o
 (\j -> seqp p
 (\k -> seqp (assocTail o p)
 (\l -> mreturn $ binaryTail (Binary j i k) l))))

binaryTail :: Predicate -> [(String,Predicate)] -> Predicate
binaryTail (Binary x y z) [(o,p)] = Binary x y (Binary o z p)
binaryTail (Binary x y z) ((o,p):xs) = Binary x y (binaryTail (Binary o z p) xs)

assocTail :: Parser String -> Parser Predicate -> 
 Parser [(String,Predicate)]
assocTail o p = many (seqp o (\i -> seqp p (\j -> mreturn (i,j))))
type Parser a = String -> [(a,String)]

-- the previous region require primitive parsers

mreturn :: a -> Parser a
mreturn v = \inp -> [(v, inp)]

failure :: Parser a
failure = \inp -> []

item :: Parser Char
item = \inp -> case inp of 
  [] -> []
  (x:xs) -> [(x,xs)]
  
parse :: Parser a -> String -> [(a, String)]
parse p inp = p inp

seqp :: Parser a -> (a -> Parser b) -> Parser b 
seqp p f  = \inp -> case parse p inp of
  [] -> []
  [(v,out)] -> parse (f v) out

(+++) :: Parser a -> Parser a -> Parser a
p +++ q = \inp -> case parse p inp of
  [] -> parse q inp
  [(v,out)] -> [(v,out)]
  
sat :: (Char -> Bool) -> Parser Char
sat p = seqp item (\x -> if p x then mreturn x else failure)

char :: Char -> Parser Char
char x = sat (== x)

digit :: Parser Char
digit = sat isDigit

letter :: Parser Char
letter = sat isAlpha

ident :: Parser String
ident = seqp letter 
 (\i -> seqp (many (letter +++ digit))
 (\j -> mreturn (i:j)))

string :: String -> Parser String
string [] = mreturn []
string (x:xs) = seqp (char x) 
 (\i -> seqp (string xs) (\j -> mreturn (x:xs)))

many :: Parser a -> Parser [a]
many p = many1 p +++ mreturn []

many1 :: Parser a -> Parser [a]
many1 p = seqp p (\i -> seqp (many p) (\j -> mreturn (i:j)))

space :: Parser ()
space = seqp (many (sat isSpace)) (\i -> mreturn ())

token :: Parser a -> Parser a
token p = seqp space 
 (\i -> seqp p (\j -> seqp space (\k -> mreturn j)))

symbol :: String -> Parser String
symbol xs = token (string xs)

cond :: Bool -> a -> a -> a
cond ok a b = if ok then a else b

eq :: String
eq = "≡"

neq :: String
neq = "≢"

impl :: String
impl = "⇒"

fll :: String
fll = "⇐"

orOp :: String
orOp = "∨"

andOp :: String
andOp = "∧"

neg :: String
neg = "¬"

-- a ≡ b ≡ c
