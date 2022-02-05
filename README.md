# Wordle Solver

This small package provides a solver for 'wordle' style puzzles.   It takes an arbitrary list of words and provides the optimal strategy to solve them.

The solution is inspired by Don Knuth's Mastermind algorithm that uses a minmax approach.  For those interested, the link to the original paper is [here](http://www.cs.uni.edu/~wallingf/teaching/cs3530/resources/knuth-mastermind.pdf)

The package has a few components.  In the `words` directory you'll find some word lists.  In particular, the Colline Scrabble Dictionary list is there, as is the subset of that dictionary that is 5 letter words, used for the original Wordle game.

More interestingly, in the `strategies` directory you will find the output of the solver against the ~13000 five letter word list.   The format is a bit wonky but not that hard to understand once you get the hang of it.   Your first guess should be 'SERAI'.  Then you look at the response.   Gray (not present), Yellow (present but in wrong location) and Green ('present and in correct location') are represented by 0, 1, and 2, respectively.   So after your first guess, you find the entry matching 'd1' (depth 1) with the response.   That will give you your next guess.  Once you get that, you find the corresponding 'd2' entry, which will indicate your next guess, and so on.   This strategy, on average solves wordle in 4.3 guesses, across the entire dictionary.  A JSON version of the strategy is also in the same directory.

If you are interested in the tool itself, build it and run

`% ./wordle --usage`

which should give you the various options.  The output goes to stdout.

Enjoy!