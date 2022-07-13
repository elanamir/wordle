# WordleBot - Wordle style puzzle solver

This small package provides a solver for 'wordle' style puzzles.   It takes an arbitrary list of words and provides the optimal strategy to solve them.   There are configuration options to load dictionaries with different word lengths to generate strategies, as well as way to use the strategy in an interactive bot.

The solution was initially inspired by Don Knuth's Mastermind algorithm that uses a minmax approach.  For those interested, the link to the original paper is [here](http://www.cs.uni.edu/~wallingf/teaching/cs3530/resources/knuth-mastermind.pdf).    Instead of the minmax, it uses the entropy of the remaining dictionary distribution to develop a strategy that is as efficient as possible.   


The package has a few components.  In the `words` directory you'll find some word lists.  In particular, the Collins Scrabble Dictionary list is there, as is the subset of that dictionary that is 5 letter words.  I'm not sure, but I think this was used for the original Wordle game.   Using this dictionary, the solver generates a strategy that gets to a solution in an average of 4.07 guesses.

More interestingly, in the `strategies` directory you will find the output of the solver against the ~13000 five letter word list.   The .txt format can be followed by a human - it's a bit wonky but not that hard to understand once you get the hang of it.   Your first guess is at the top ('TARES' in the CollinsDictionary strategy).  Then you look at the response.   Gray (not present), Yellow (present but in wrong location) and Green (present and in correct location) are represented by 0, 1, and 2, respectively.   So after your first guess, you find the entry matching 'd1' (depth 1) with the response.   That will give you your next guess.  Once you get that, you find the corresponding 'd2' entry, which will indicate your next guess, and so on.   

Alternatively, a JSON version of the strategy is also in the same directory, which you can ingest to create an interactive bot.   The command to run is:

`% ./wordlebot -cmdline -strategyfile=strategies/CollinsStrategy.json`

The bot will walk you through the guesses.  You enter the responses using the code above.  For example, `gray gray green yellow green` would be entered `00212` as a response.

Aside from that: 

`% ./wordlebot -usage`

which should give you the various options.  The output goes to stdout.

Enjoy!
