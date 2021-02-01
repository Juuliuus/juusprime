package juusprime

//juusprime is an engine to find prime Sextuplets and/or Quintuplets and/or Quadruplets
//GNU GPL3
//author: Julius Schön / R. Spicer, 2021

import (
	"fmt"
	"math/big"
	"os"
)

func init() {
	SymbolCount = make(map[int]int)
	SymbolCount[CSextuplet] = -1
	SymbolCount[CLQuint29] = -1
	SymbolCount[CRQuint13] = -1
	SymbolCount[CQuad] = -1
}

//symbols used to represent Tuplets
const (
	cSymSextFull  = "┣━┫"
	cSymSextBrief = " ● "
	cSymLQuint29  = "┣━ "
	cSymRQuint13  = " ━┫"
	cSymQuad      = " ━ "
	cSymDestroyed = " X "
	cInsertSymbol = "o"
	//'possibly for dot n's?? ◎ ■
	//cRawDataFormat = "%v:%v:%v"
	//basisNotification = "-> basis, %s:%v\n"
)

//the integer representation of the Tuplet symbols
const (
	CSextuplet = iota
	CLQuint29
	CRQuint13
	CQuad
	CXNoTrack
	CX17
	CX19
	CX23
	CX25
)

//filter types (ft) for organizing filter choices
const (
	ftAll = iota
	ftSextuplet
	ftLQuint
	ftRQuint
	ftQuints
	ftQuads
	ftCount //used to limit loops
)

//for pre and post filters
const (
	filterSextuplet = 1 << iota
	filterLQuint29
	filterRQuint13
	filterQuad
	filterXNoTrack
)

var (
	SymbolCount map[int]int
	FilterMap   = map[int]uint{
		CSextuplet: filterSextuplet,
		CLQuint29:  filterLQuint29,
		CRQuint13:  filterRQuint13,
		CQuad:      filterQuad,
		CXNoTrack:  filterXNoTrack,
	}
)

//GetFilter29 : Returns the disallow & finalPass filters for
//routines working with LTE 29 Primes
func GetFilter29(filterType int) (disallow, finalPass uint) {
	switch filterType {
	case ftAll: //everything
		//reject these in adder loop
		disallow = filterXNoTrack
		//finalPass then specifies what is allowed to pass through
		finalPass = filterSextuplet | filterLQuint29 | filterRQuint13 | filterQuad
		return
	case ftSextuplet: //sextuplets only
		disallow = filterLQuint29 | filterRQuint13 | filterQuad | filterXNoTrack
		finalPass = filterSextuplet
		return
	case ftLQuint: //L quint 29's only
		disallow = filterRQuint13 | filterQuad | filterXNoTrack
		finalPass = filterLQuint29
		return
	case ftRQuint: //R quint 13's only
		disallow = filterLQuint29 | filterQuad | filterXNoTrack
		finalPass = filterRQuint13
		return
	case ftQuints: //R & L quints only
		disallow = filterQuad | filterXNoTrack
		finalPass = filterRQuint13 | filterLQuint29
		return
	case ftQuads: //Quads only
		disallow = filterXNoTrack
		finalPass = filterQuad
		return
	default:
		fmt.Println("Problem in GetFilter, invalid filterType")
		disallow = 0
		finalPass = 0
		return
	}
	return 0, 0
}

//GetFilter : Returns the pre & post filters for
//routines working with GTE 31 Primes
func GetFilter(filterType int) (pre, post uint) {
	switch filterType {
	case ftAll: //everything
		pre = filterSextuplet | filterLQuint29 | filterRQuint13 | filterQuad
		post = filterXNoTrack
		return
	case ftSextuplet: //sextuplets only
		pre = filterSextuplet
		post = filterLQuint29 | filterRQuint13 | filterQuad | filterXNoTrack
		return
	case ftLQuint: //L quint 29's only
		pre = filterSextuplet | filterLQuint29
		post = filterSextuplet | filterRQuint13 | filterQuad | filterXNoTrack
		return
	case ftRQuint: //R quint 13's only
		pre = filterSextuplet | filterRQuint13
		post = filterSextuplet | filterLQuint29 | filterQuad | filterXNoTrack
		return
	case ftQuints: //R & L quints only
		pre = filterSextuplet | filterRQuint13 | filterLQuint29
		post = filterSextuplet | filterQuad | filterXNoTrack
		return
	case ftQuads: //Quads only
		pre = filterSextuplet | filterRQuint13 | filterLQuint29 | filterQuad
		post = filterSextuplet | filterRQuint13 | filterLQuint29 | filterXNoTrack
		return
	default:
		fmt.Println("Problem in GetFilter, invalid filterType")
		pre = 0
		post = 0
		return
	}
	return 0, 0
}

//GetFilterAbbrev : abbreviation for filter types, usually used
//in file names
func GetFilterAbbrev(filterType int) string {
	switch filterType {
	case ftAll: //everything
		return "6L5R5Q"
	case ftSextuplet: //sextuplets only
		return "6"
	case ftLQuint: //L quint 29's only
		return "L5"
	case ftRQuint: //R quint 13's only
		return "R5"
	case ftQuints: //R & L quints only
		return "L5R5"
	case ftQuads: //Quads only
		return "Q"
	default:
		fmt.Println("Problem in GetFilterAbbrev, invalid filterType")
		return "xx"
	}
	return "xx"
}

//GetFilterDesc : Human readable description for filter types, usually used
//in information and messages to the user
func GetFilterDesc(filterType int) string {
	switch filterType {
	case ftAll:
		return "no filter"
	case ftSextuplet:
		return "Sextuplets only"
	case ftLQuint:
		return "Left Quints only"
	case ftRQuint:
		return "Right Quints only"
	case ftQuints:
		return "Left and Right Quints only"
	case ftQuads:
		return "Quads only"
	default:
		fmt.Println("Problem in GetFilterDesc, invalid filterType")
		return "Unknown"
	}
	return "Unknown"
}

func HelpOutputFiles() {
	fmt.Println("\njuusprime file information:")
	fmt.Println("\nfile names:")
	fmt.Println("\nname_basis-#_tFrom_tTo_filter.ext")
	fmt.Println("\nname is a descriptor for the various types of files.")
	fmt.Println("basis-#: # is replaced by the basis number the file starts out in.")
	fmt.Println("tFrom is the from TNumber")
	fmt.Println("tTo is the to TNumber")
	fmt.Println("filter is a short description showing which tuplets the file contains.")
	fmt.Println("ext is the extension, either: .rawdata, .prettydata, or .info")
	fmt.Println("")
	waitForInput()

	fmt.Println("=====")
	fmt.Println("for example: juusprimes_basis-0_28_215656468_6L5R5Q.rawdata")
	fmt.Println("=====")

	fmt.Println("\nIt is a tuplet outpfile (juusprimes)")
	fmt.Println("the basis it starts in is 0")
	fmt.Println("it starts at TNumber 28")
	fmt.Println("it ends at TNumber 215656468")
	fmt.Println("it contains sextuplets, right & left quintuplets, and quadruplets")
	fmt.Println("it is a rawdata file.")
	fmt.Println("")
	ShowSymbolFileDesignations(os.Stdout)

	//fmt.Println(cInfoDataFormat)
}

//HelpSymbolsMath : Display how symbol "math" is done
func HelpSymbolsMath() {
	InterimResult := -1

	doDisplay := func(first, second int) {
		InterimResult = AddSymbols(&first, &second)
		fmt.Println(fmt.Sprintf("%s + %s = %s", GetSymbolString(first, true), GetSymbolString(second, true), GetSymbolString(InterimResult, true)))
	}

	doDisplay(CSextuplet, CSextuplet)
	doDisplay(CSextuplet, CLQuint29)
	doDisplay(CSextuplet, CRQuint13)
	doDisplay(CSextuplet, CQuad)
	doDisplay(CSextuplet, CXNoTrack)
	fmt.Println("------")
	doDisplay(CLQuint29, CSextuplet)
	doDisplay(CLQuint29, CLQuint29)
	doDisplay(CLQuint29, CRQuint13)
	doDisplay(CLQuint29, CQuad)
	doDisplay(CLQuint29, CXNoTrack)
	fmt.Println("------")
	doDisplay(CRQuint13, CSextuplet)
	doDisplay(CRQuint13, CLQuint29)
	doDisplay(CRQuint13, CRQuint13)
	doDisplay(CRQuint13, CQuad)
	doDisplay(CRQuint13, CXNoTrack)
	fmt.Println("------")
	doDisplay(CQuad, CSextuplet)
	doDisplay(CQuad, CLQuint29)
	doDisplay(CQuad, CRQuint13)
	doDisplay(CQuad, CQuad)
	doDisplay(CQuad, CXNoTrack)
	fmt.Println("------")
	doDisplay(CXNoTrack, CSextuplet)
	doDisplay(CXNoTrack, CLQuint29)
	doDisplay(CXNoTrack, CRQuint13)
	doDisplay(CXNoTrack, CQuad)
	doDisplay(CXNoTrack, CXNoTrack)
	fmt.Println("------")
}

//HelpSymbols : Display the explanation and symbols used for Tuplets
func HelpSymbols() {
	fmt.Println("Symbols and info")
	fmt.Println("----------------")
	fmt.Println("(integer representation is what you see in raw files)")
	fmt.Println("")
	fmt.Println(fmt.Sprintf("%v or %v : both equal Sextuplet, i.e. its a 'pass'", cSymSextFull, cSymSextBrief))
	fmt.Println(fmt.Sprintf("Integer representation: %v", CSextuplet))
	fmt.Println("")
	fmt.Println(fmt.Sprintf("%v : Left Handed Quintuplet", cSymLQuint29))
	fmt.Println(fmt.Sprintf("Integer representation: %v", CLQuint29))
	fmt.Println("")
	fmt.Println(fmt.Sprintf("%v : Right Handed Quintuplet", cSymRQuint13))
	fmt.Println(fmt.Sprintf("Integer representation: %v", CRQuint13))
	fmt.Println("")
	fmt.Println(fmt.Sprintf("%v : Quadruplet", cSymQuad))
	fmt.Println(fmt.Sprintf("Integer representation: %v", CQuad))
	fmt.Println("")
	fmt.Println(fmt.Sprintf("%v : Structure Destroyed", cSymDestroyed))
	fmt.Println("(ie., crosses 17,19,23,or 25 (mod 16,18,22,or 24))")
	fmt.Println(fmt.Sprintf("Integer representation: %v", CXNoTrack))
	fmt.Println("")
	fmt.Println(fmt.Sprintf("%v : inserted pass, prime has jumped a complete 30 template", cInsertSymbol))
	fmt.Println("")
}

//GetEffectDisplay : Return a user friendly string for symbol.
//used where only the effect is shown.
func GetEffectDisplay(symbol int) string {
	switch symbol {
	case CSextuplet, CQuad:
		return "" //"┣━┫", "";
	case CLQuint29:
		return cSymLQuint29 // "┣━ "
	case CRQuint13:
		return cSymRQuint13 // " ━┫"
	}
	return cSymDestroyed // "X"
}

//AddSymbols : Return the "addition" of effects that determine whether
//the two symbols add to a significant result
func AddSymbols(symbol1, symbol2 *int) int {

	switch *symbol1 {
	case CSextuplet:
		switch *symbol2 {
		case CSextuplet:
			return CSextuplet
		case CLQuint29:
			return CLQuint29
		case CRQuint13:
			return CRQuint13
		case CQuad:
			return CQuad
			//          cX17, cX19, cX23, cX25,
			//          cXNoTrack : result := ' X ';
		}
	case CLQuint29:
		switch *symbol2 {
		case CSextuplet:
			return CLQuint29
		case CLQuint29:
			return CLQuint29
		case CRQuint13:
			return CQuad
		case CQuad:
			return CQuad
			//          cX17, cX19, cX23, cX25,
			//          cXNoTrack : result := ' X ';
		}
	case CRQuint13:
		switch *symbol2 {
		case CSextuplet:
			return CRQuint13
		case CLQuint29:
			return CQuad
		case CRQuint13:
			return CRQuint13
		case CQuad:
			return CQuad
			//          cX17, cX19, cX23, cX25,
			//          cXNoTrack : result := ' X ';
		}
	case CQuad:
		switch *symbol2 {
		case CSextuplet, CLQuint29, CRQuint13, CQuad:
			return CQuad
			//          cX17, cX19, cX23, cX25,
			//          cXNoTrack : result := ' X ';
		}
	}
	return CXNoTrack
}

//SymbolIsOfInterest : basic filtering func; rejects any
//thing that is not a Tuple; only used in GeneratePrimes7to23()
func SymbolIsOfInterest(symbol int) bool {
	if symbol <= CQuad {
		return true
	}
	return false
}

//GetSymbolDisplay : Return a user friendly string for symbol.
//used where Symbol is displayed either screen or file.
func GetSymbolDisplay(symbol int) string {
	switch symbol {
	case CSextuplet:
		return cSymSextFull // "┣━┫"
	case CLQuint29:
		return cSymLQuint29 // "┣━ "
	case CRQuint13:
		return cSymRQuint13 // " ━┫"
	case CQuad:
		return cSymQuad // " ━ "
	}
	return ""
}

//GetSymbolString : Similar to GetSymbolDisplay. Return a user friendly string
//for symbol, with the addition of a less busy string for Sextuplets which
//is used if fullSymbol is false
func GetSymbolString(symbol int, fullSymbol bool) string {
	switch symbol {
	case CSextuplet:
		if fullSymbol {
			return cSymSextFull // "┣━┫"
		}
		return cSymSextBrief // " ● "
	case CLQuint29:
		return cSymLQuint29 // "┣━ "
	case CRQuint13:
		return cSymRQuint13 // " ━┫"
	case CQuad:
		return cSymQuad // " ━ "
	case CX17, CX19, CX23, CX25, CXNoTrack:
		return cSymDestroyed // " X "
	}
	return "UNK"
}

//ClearSymbolCounts : Clears the SymbolCount map which
//is used to accumulate results
func ClearSymbolCounts() {
	SymbolCount[CSextuplet] = 0
	SymbolCount[CLQuint29] = 0
	SymbolCount[CRQuint13] = 0
	SymbolCount[CQuad] = 0
}

//TNumLastNatNum : helper func that calculates the last natural number in
//a TNumber range
func TNumLastNatNum(tNum *big.Int) *big.Int {
	result := big.NewInt(0).Set(TNumToInt(tNum))
	return result.Add(result, big29)
}

//ShowSymbolCounts: Print out the SymbolCount map accumulated results
func ShowSymbolCounts(from, to *big.Int, filter int, f *os.File) {
	fmt.Fprintln(f, fmt.Sprintf("\nFinal counts (from TNumber %v to %v)", from, to))
	fmt.Fprintln(f, fmt.Sprintf("(Natural numbers from %v to %v)", TNumToInt(from), TNumLastNatNum(to)))
	fmt.Fprintln(f, fmt.Sprintf("(filtered by: %s)", GetFilterDesc(filter)))

	fmt.Fprintln(f, fmt.Sprintf("%v %s", SymbolCount[CSextuplet], "Sextuplets"))
	fmt.Fprintln(f, fmt.Sprintf("%v %s", SymbolCount[CLQuint29], "LQuints"))
	fmt.Fprintln(f, fmt.Sprintf("%v %s", SymbolCount[CRQuint13], "RQuints"))
	fmt.Fprintln(f, fmt.Sprintf("%v %s\n", SymbolCount[CQuad], "Quads"))

	sumInteresting := 0
	for cnt := range SymbolCount {
		sumInteresting = sumInteresting + SymbolCount[cnt]
	}
	fmt.Fprintln(f, fmt.Sprintf("Sum of found Symbols: %v\n", sumInteresting))

	//fmt.Println("\nSanity Check of SymbolCounts::")
	//fmt.Println(SymbolCount)
}

//ShowSymbolFileDesignations : a routine that assembles the pieces that are
//printed to the info file so people can understand the symbol strings in file names
func ShowSymbolFileDesignations(f *os.File) {
	fmt.Fprintln(f, cInfoFilter)
	for i := 0; i < ftCount; i++ {
		fmt.Fprintln(f, fmt.Sprintf(cInfoFilterBase, GetFilterAbbrev(i), GetFilterDesc(i)))
	}
	fmt.Fprintln(f, "")
}

//some large string constants for printing into various information files
const (
	cInfoFilter     = "Filter designations in file names:\n"
	cInfoFilterBase = "%s = %s"
	Basis29Msg      = `
You've chosen to print out a custom 29 basis file. Perhaps
you are testing things, want to filter it for some reason, or 
are choosing odd ranges for some reason.

In the case you are not completely familiar with the purpose of basis files
be aware that routines for juusprime generation expect a normal default file.

If you want special basis files for some reason you need to be completely familiar
with filtering and usage. If you are not you can expect odd or incomplete results
if used or filtered improperly.`
	cInfo23Text = `The associated file is the outcome of the surviving sextuplets and quads in the range of template numbers 1 to 27.

These are separate because they needed to be handled separately from the basis 29 calculations, as well as the generation of prime sextuplets, quintuplets, and quadruplets. 

Technically they "belong" at the start of the 29 rawdata, but only technically. These are exceptions to the generality of the prime structures and are not meant to be used in any way in the building of the structures. They are for completeness and simply show the results up to and including TNumber 27.

If you feel the urge to append these to the beginning of the 29Basis rawdata, =don't= do it!

 `
	cInfoSymbols = `The symbols for the sextuplets and so on are of two types, a plain integer for computing purposes (in the raw files) and a pretty symbol (in the pretty files) to make them easy to spot in a list meant for humans to read.

0 = Sextuplet  = ┣━┫

1 = LQuint_29  = ┣━

2 = RQuint_13  =  ━┫

3 = Quadruplet =  ━	


As to file names if the From TNumber or the To TNumber is greater than 9999999999999999999999999999999 they will be converted to md5 sums so that file names remain readable, under the filename length limit, and to prevent file name clashes.

 `
	cInfoDataFormat = `The rawdata file is linear, in TNumber order, and uses a pair of lines to describe the location and the structure at that location. The first of the pair (odd numbered lines) is the Template Number, the second of the pair (even numbered lines) is the internal symbol for the tuplet structure. Look in the prettydata file for human readable info.

The pretty data files will show TNumbers, expanded TNumbers, corresponding 29Basis TNumber, and the primes found.

info files contain pertinent information regarding the associated rawdata file.

TNumbers (template numbers) are an abstraction of the number line into chunks of 30 numbers. The first Template (TNum 1) starts at number line position 25, the next TNumber (2) is 30 away, ie. at 55, and so on. To convert a TNumber to the number line position which starts the Template use: number = (TNumber * 30) - 5. To get the TNumber for any natural number use TNumber = (number + 5) div 30

`
	cInfoBasis = `The default 29basis file (no filters, and from TNumber 28 to TNumber 215656468) contains all (meaning Sextuplets, Left and Right quintuplets, and quadruplets) the "potential" locations of these prime structures. These are the only locations that are possible throughout the entire number line. If you have generated other ranges for the 29basis, be aware that only the default settings can be meaningfully used, in general, when generating the juusprimes. But it may be useful, when you desire, to generate a basis file containing only possible sextuplets by filtering by sextuplets, that file can also be used, but be aware that it will not contain "natural" quints or quads. Filtering can be used in a similar fashion for all types. Be aware that if you used filtered 29 basis files, they must be in the default range (TNumber 28 to TNumber 215656468) unless you also restrict your juusprimes generation to the appropriate range of your custom 29basis file.

The default basis file starts from template number 28 and goes out for a "basis length" of 215656441 to template number 215656468. The "basis length" of 215656441 comes from [7-29]Pactorial (ie., 7*11*13*17*19*23*29), which gives all the possible combinations of natural progression crossing numbers and their effects on a pure Sextuplet.

The pattern then repeats from template number 215656469 for a length of 215656441, and so on, out to infinity. One can then use this raw data against primes 31 and greater to find out where the sextuplets, quints, and quads survive, at any position in the number line.

if you generate a filtered 29basis file (for example sextuplets only) then be aware that that file is good for ONLY sextuplets, =AND= that if you filter the processing file by anything other than no filter or sextuplet filter you won't get any results.

`
	cInfoBasisExtra = `Following are the starting and ending offsets of the primes less than 29. Since one must go through =ALL= combinations, one comes back to the exact same offsets, minus 1 each, proving that we are at the end of a repeating cycle.

Starting (Tnumber 28):
7: 6
11: 2
13: 10
17: 2
19: 16
23: 11
29: 0

Ending (Tnumber 215656468):
7: 5
11: 1
13: 9
17: 1
19: 15
23: 10
29: 28

One can also see at the next TNumber that the "combination" returns to the exact same as at TNumber 28. The cycle starts again.

Starting (Tnumber 215656469):
7: 6
11: 2
13: 10
17: 2
19: 16
23: 11
29: 0

BTW, the numbers above are indexes into the respective LTE 29 prime's Natural Progression array, which then can be looked up to find the crossing number, which then can be used to get the "effect" of the prime at that index. An "effect" is whether it allows a sext/quint/quad to exist, or does it modify or destroy it.

`
)
