
## juusprime : Prime Tuplet generation written in golang ##

March 2021:

This is the package code for juusprime which is an engine to generate
Prime Tuplets (Sextuplets and/or Quintuplets and/or Quadruplets). 

A pdf file is available in the file list that details the underlying
theory and structures used for the algorithms, only algebra is
needed. A nutshell summary is given below.

New as of 17 March 2021 (v1.1.0):
- Implement new getCrossNumModDirect for 30% speed increase
- Add analysis routines.
- automation with GeneratePrimeTupletsAutomated()

See at bottom for other recent additions.

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
tuplets in the region you are searching. 

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

Support for Automation for producing Tuplets (using basis-#) has been
added so that processing can be done with scripts:

GeneratePrimeTupletsAutomated(auto *AutomationStruct)

Automation has been added to the associated pre-compiled executable (link above)


-------------------------------------------------------------------------------






## Nutshell overview ##

With these structures the number line is broken up into 2 types of
blocks of numbers: TNumbers (Template Numbers) and Basis Numbers. It
is recommended to do basis number generation because it is then
assured that no numbers on the number line are left out by accident.

This is a very short overview. The details are available in the pdf
file which will be updated as new analysis routines are written..

These regions arise out of splitting the prime universe into 3
separate regions and they arise naturally, which helps to turn the
chaos of the primes into groups that useful and easy(ier) to work
with: primes 2-3-5, primes 7-29, and the primes 31-59.

The routines do not factor out or search for primes at all. The blocks
and structures contain patterns that can be combined to "grok" the
entire block.

## Performance ##

For the generation of basis-0 (see pdf) tuplets it takes about 20 minutes to
generate all possible tuplets (Sext's, Quints, and quads) but only
about 40 seconds if you filter for sextuplets only. 

For basis-1000 (the thousand and first basis), filtered by Sextuplets
only, it took a bit under 2 minutes:

Final counts (from TNumber 215656441028 to 215872097468)<br>
(Natural numbers from 6469693230835 to 6476162924064)<br>
(filtered by: Sextuplets only)<br>
177 Sextuplets<br>
0 LQuints<br>
0 RQuints<br>
0 Quads<br>

When generating the one time 29basis file, individual processor cores go to 100%,
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

It found 7 Sextuplets in about 5.5 hours. Here is the prettydata
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


v1.1.0 March 17 2021:
- Implement new getCrossNumModDirect for 30% speed increase
- Add analysis routines.
- automation with GeneratePrimeTupletsAutomated()

v1.0.1 March 5, 2021
- add twin sextuplet check during Tuplet generation
- add automation routines for use through shell scripts, etc.
- Add initial Helper structure to hold useful runtime vars.

v1.0.0 February 2021
