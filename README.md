
## juusprime : Prime Tuplet generation written in golang ##

February 2021:

This is the package code for juusprime which is an engine to generate
Prime Tuplets (Sextuplets and/or Quintuplets and/or Quadruplets). 

juusprime is free software, licensed under the GNU GPL, version 3. It
is written in pure Go, no other dependencies.

Recommended is to use the associated go application I wrote. It is an
interactive linux terminal application which uses this package to
generate tuplets, provide information, and help with the process of
generation and organization. The go code is found at:

https://github.com/Juuliuus/juusprime_app

If you do not use golang, then you can also download a
compiled, ready to run, executable from my website (you may need to
set it as executable before it will run):
Download page (gpg sig): https://www.timepirate.org/downloads.html
App Direct Download: https://www.timepirate.org/downloads/juusprime_app

This project started its life as a proof of concept for prime
structures that I first worked out on paper. I needed to implement the
paper scribbles to see that it actually worked. It worked well enough
that I thought I would share it for those interested in this kind of
thing.

It uses a gestalt-y block approach to tuplet generation, not brute
force. The concepts, as I first outlined with pencil and paper,
indicate that it will be 100% accurate and find 100% of the wanted
tuplets in the region you are searching. I will be working on a pdf
that outlines in detail how the structures are built and used. For now
I include below a "short" summary of the thoughts and theory behind
it.

It uses go's big Int (and where needed, big Float), which means it can
work with any integer that you care to type in. But, as usual for primes,
the bigger the numbers the longer its going to take.

Each run will output a rawdata file, a prettydata file, and an info
file. These are all text files.

-------------------------------------------------------------------------------

### Install ###

Install in the usual Go manner with: 

```
go get github.com/Juuliuus/juusprime
```

### Basic usage ###

For those that prefer to write their own code for this package:

Be sure you are happy with config data file paths by calling
Configure().

Before you can generate Tuplets you need what I refer to as a
"29basis" file. Call GenerateBasis() function to write this file. This
will take a minute or two, the file is about 190Mb.

You can now generate tuplets for any range you desire at any time, call
the GeneratePrimeTupletsInteractive() routine.

-------------------------------------------------------------------------------


## Nutshell overview of the structure and logic ##

With these structures the number line is broken up into 2 types of
blocks of numbers: TNumbers (Template Numbers) and Basis Numbers. It
is recommended to do basis number generation because it is then
assured that no numbers on the number line are left out by accident.

This will be a an overview, there are many details not covered here.

To understand these blocks the Primes are first split into 3 groups, each handled
separately: primes 2-3-5, primes 7-29, and the "primes" 31-59.

### Primes 2, 3, & 5 ###

These form the blocks of the TNumbers (Template Numbers) because they
form sextuplet templates.

It is very easy to see for yourself, mark a list numbered 1 to 120 and
cross out every 2nd, 3rd, and fifth multiple; i.e., a sieve using only
2,3,5.

You will then clearly see at position 25 that a repeating pattern 30
numbers long that shows the sextuplet sitting untouched on the right
of the tempalte and a twin prime on the left. In this project I focus
only on the sextuplet structure, the "rogue" twin on the left is
ignored. 

If you drew the numbers out to 120 you will see the pattern 25-54
(Tnumber 1), 55-84 (TNum 2), and 85-114 (TNum 3). 

It is in TNum 3 that we find the first sextuplet resides: 
97, 101, 103, 107, 109, 113.

NB: The above is the first sexuplet only if you don't want to count the
7,11,13,17,19,23 one.

So actually we need to do nothing with 2,3,5; just accept that every
TNumber always starts life out as a pristine sextuplet waiting to be
realized (or destroyed).

TNumber math is simple:

TNum = (Int + 5) div 30  (div result, not quotient)

Int = (TNum * 30) - 5

The Int= equation above gives the starting integer of that TNumber
range. If you want the ending integer simply add 29.


### Primes 7-29 ###

These form the block I call the "29basis". What's nice about the basis
is that it shows the only places Tuplets =can= occur and it is
re-usable since it recurs after a reasonable length.

Many details here but in summary: It is possible to create a sequence
for each of these primes such that one can calculate what "effect" it
will have at any TNumber by its crossing number at that TNumber. 

An effect is: does it leave the sextuplet above it alone? Or does it
convert the sextuplet to either a left sided quintuplet, or a right
sided quintuplet, or does it destroy it by taking out either one of
the twin primes at the sextuplets center?

prime# 7 for example: Its starting TNumber is TNum 1, because its square
lives in that TNumber. Computing the progressions as it moves
through the seven Templates between 1 and 7 (inclusive) results in  its
natural progression (mod based) of: 3, 1, 6, 4, 2, 0, 5

In this case we see that its third value  will cross TNum 3. The
effect prime 7 has when it crosses with 6 is none, ie., it leaves the
sextuplet untouched, and that is why the sextuplet in TNum 3 exists.

It turns out any crossing of prime 7, other than mod 6, will either destroy
the sextuplet or change it to a Quintuplet (crossing 0 leaves a left
sided quintuplet, crossing 5 leaves a right sided quintuplet, and
crossings 1,2,3,4 destroy the sextuplet). Factoid: From this you can
infer that any difference between TNumbers that contain sextuplets, or
members of differing sextuplets will always be evenly divisible by 7,
which is the case.

Similar structures are built for each of the other primes 11-29. Then
starting at TNumber 28 (which is where prime 29 begins having an
effect on the following integers) one can use these cycles to simply
walk through all possible combinations of the 7-29 cycles until the
overall cycle starts repeating: the 29Basis.

This is not as bad as it sounds! The size will be 7-29 primorial
(ie. 7 * 11 * 13 * 17 * 19 * 23 * 29). This results in a block,
starting at TNum 28 that stretches out to TNum 215656468 (inclusive)
(the length is 215656441 TNumbers). This encompasses a range of
integers from 835 to 6469694064. This is basis-0.

Basis-1 covers TNums from 215656469 to 431312909, and so on.

When this 29Basis file is written out we now have a map to the only
possible locations of the Tuplets, and since it can be re-used when it
cycles we can look at any basis number we want. The engine also allows
generating from specified TNumbers or integers, but basis generation
is recommended because it is assured that all possible Tuplets are
tested.


### "Primes" 31-59 ###

With the generation of the 29basis file a lot of the work of
checking Tuplets is already done, and we know exactly where to look.

Now we can do the final check to see if the possible Sextuplet (or
quint or quad) in the 29basis file becomes a real Sextuplet or is
altered or destroyed.

To do this we only need 8 "primes". Primes is quoted because they are
formed a bit differently than what you expect. We add 1, 7, 11, 13,
17, 19, 23 and 29 to 30. This results in the "primes":

31, 37, 41, 43, 47, 49, 53, 59

You may be surprised that 49 is a "prime", but for this structure it
actually is! Perhaps here is a good time to call these what they
actually are: potential primes (potPrimes). Here's why.

We now construct similar cycles for these primes just as we did for
the 7-29 primes. It turns out to be a much simpler process and they
expect regular behaviors, patterns.

So it also turns out that the potPrimes are very amenable to equations
that can be used to "look up" their effects at any TNumber. They are
regularly irregular or irregularly regular, can't decide which is
better.

The process we need to do in generating Tuplets is simply look up a
possible location in the 29Basis file, which will be by TNumber, and
then check to see if any potPrimes have an effect at that TNumber, or
allow the Tuplet to pass on untouched, and this done by manipulating
their look up tables.

How many potPrimes need to be checked at any given TNumber is
calculable (I call it n, it starts at 0). n0 covers all 8 of the
potPrimes above. But, it could be that n is 1 and suddenly I've run
out of potPrimes...

This is where the potPrimes are regularly irregular. It turns out that
61, the 9th pot prime (n=1), has exactly the same structure as 31 with
the exception that its length is 30 greater. 91, n-2, has length 60
greater...

To handle that all one has to do is to "expand" the look up tables
(like the Universe expands). The look up tables, btw, are constants
and must be figured out with pen and paper.

And so the proper definition of a PotPrime, P, is: 
P = Pvalue + 30n

For n's = 0,1,2 and potPrimes 31, 49

31: 31, 61, 91

49: 49, 79, 109

Since we're checking by blocks, and using simple fast calculations to
do look ups and using "Tuplet math", the entire thing is fast and
relatively simple.

We are checking non-prime potPrimes, yes, but that means we leave
nothing out by accident. It's the best scenario I have right now,
though I continue to think about it. I have tried a couple of schemes
to programatically weed out non-prime potPrimes and all it does is
significantly slow the engine down.

As I said, there are undiscussed details to all this, but in a
nutshell that how this engine works.


## Performance ##

For the generation of basis-0 tuplets it takes about half an hour to
generate all possible tuplets (Sext's, Quints, and quads) but only
about 1 minute if you filter for sextuplets only. 

For basis-1000 (the thousand and first basis), filtered by Sextuplets
only, it took a bit under 3 minutes:

Final counts (from TNumber 215656441028 to 215872097468)<br>
(Natural numbers from 6469693230835 to 6476162924064)<br>
(filtered by: Sextuplets only)<br>
177 Sextuplets<br>
0 LQuints<br>
0 RQuints<br>
0 Quads<br>

When generating a basis file, individual processor cores go to 100%,
but this basis generation process only takes about a minute.

However, at least on my machine, the processor usage does not peg out
when generating tuplets. And I'm not sure why. Perhaps the go language
ability to compile to spread work across processors has kicked in?

In any case, on my dual-core quad processor, usage remains at about
the 20% level, leaving the computer very responsive to other programs
while generating tuplets.

I did do, for testing purposes, a rather large range of TNumbers (it
is not an entire basis, just a portion) in basis-86000000000

This puts the sextuplets in the 10^20 range. 

It found 7 Sextuplets in about 7.5 hours. Here is the prettydata
output:

TNumbers from 18546453926011000028 to 18546453926100000027<br>
(Natural #'s from 556393617780330000835 to 556393617783000000834)<br>
filtered by: Sextuplets only<br>
BASIS:86000000000<br>

TNum = 18546453926016026965<br>
BeginsAt : 556393617780480808945<br>
EndsAt : 556393617780480808974<br>
[Basis-0-TNum : 16026965]<br>
---primes---   ┣━┫ (0)<br>
556393617780480808957<br>
556393617780480808961<br>
556393617780480808963<br>
556393617780480808967<br>
556393617780480808969<br>
556393617780480808973<br>

TNum = 18546453926049112178<br>
BeginsAt : 556393617781473365335<br>
EndsAt : 556393617781473365364<br>
[Basis-0-TNum : 49112178]<br>
---primes---   ┣━┫ (0)<br>
556393617781473365347<br>
556393617781473365351<br>
556393617781473365353<br>
556393617781473365357<br>
556393617781473365359<br>
556393617781473365363<br>

TNum = 18546453926052983318<br>
BeginsAt : 556393617781589499535<br>
EndsAt : 556393617781589499564<br>
[Basis-0-TNum : 52983318]<br>
---primes---   ┣━┫ (0)<br>
556393617781589499547<br>
556393617781589499551<br>
556393617781589499553<br>
556393617781589499557<br>
556393617781589499559<br>
556393617781589499563<br>

TNum = 18546453926062497662<br>
BeginsAt : 556393617781874929855<br>
EndsAt : 556393617781874929884<br>
[Basis-0-TNum : 62497662]<br>
---primes---   ┣━┫ (0)<br>
556393617781874929867<br>
556393617781874929871<br>
556393617781874929873<br>
556393617781874929877<br>
556393617781874929879<br>
556393617781874929883<br>

TNum = 18546453926079720686<br>
BeginsAt : 556393617782391620575<br>
EndsAt : 556393617782391620604<br>
[Basis-0-TNum : 79720686]<br>
---primes---   ┣━┫ (0)<br>
556393617782391620587<br>
556393617782391620591<br>
556393617782391620593<br>
556393617782391620597<br>
556393617782391620599<br>
556393617782391620603<br>

TNum = 18546453926082821469<br>
BeginsAt : 556393617782484644065<br>
EndsAt : 556393617782484644094<br>
[Basis-0-TNum : 82821469]<br>
---primes---   ┣━┫ (0)<br>
556393617782484644077<br>
556393617782484644081<br>
556393617782484644083<br>
556393617782484644087<br>
556393617782484644089<br>
556393617782484644093<br>

TNum = 18546453926092051529<br>
BeginsAt : 556393617782761545865<br>
EndsAt : 556393617782761545894<br>
[Basis-0-TNum : 92051529]<br>
---primes---   ┣━┫ (0)<br>
556393617782761545877<br>
556393617782761545881<br>
556393617782761545883<br>
556393617782761545887<br>
556393617782761545889<br>
556393617782761545893<br>

FYI, The 29Basis file contains only possible locations, but these are
also the only locations we need to search. The statistics for the
possibilities in the 29basis file are:

1956955 Sextuplets<br>
5010341 LQuints<br>
5010341 RQuints<br>
5528488 Quads<br>

From this one can see that filtering for what you are interested in is
a good idea if you want to save time. A Sextuplets only search only
has to check 2 million TNumbers, and the Sextuplets are destroyed
faster by the potPrime look up tables.

-------------------------------------------------------------------------------

I'm a biologist and computer programmer, not a mathematician. But I
have made every effort, and done many tests, to insure the structures
I build, and the results, are correct. Please let me know if you find
any problems.

-------------------------------------------------------------------------------
### History ###

v1.0.0 February 2021
