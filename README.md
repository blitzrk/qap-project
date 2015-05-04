# Testing

Can be run with `go run main.go matmath.go` from the root dir.

## To Do

Letter transition pairs have two big, generally non-language specific features:

1. Two classes of letters: consonants and vowels
1. The pair transition frequency is related to the product of each letter's frequency

The exceptions to these rules are language-specific, so we ignore them in our attempt to create a random dataset for something language-like.

However, we still must make some assumptions to fill in the gaps for our data generation algorithm. Vowels usually have a higher mean frequency of appearance than 
consonants. 
