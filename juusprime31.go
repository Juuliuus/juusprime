package juusprime

//juusprime is an engine to find prime Sextuplets and/or Quintuplets and/or Quadruplets
//GNU GPL3
//author: Julius Schön / R. Spicer, 2021

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"text/tabwriter"
)

const lookUpSize = 6 //6 locations where Tuplets can be changed or destroyed

var primeGTE31AllowedVals = []int64{31, 37, 41, 43, 47, 49, 53, 59}

//primeGT30Lookup : Used to keep track of the lookup data for the
//primes greater than or equal to 31
type primeGT30Lookup struct {
	//holds the index into the natural progression which then gives the crossing number
	C []*big.Int
	//the factor that will be used as the potential prime inflates
	Q []*big.Int
	//the effect at the crossing number
	Effect []int
}

//Stringer for lookups
func (lu *primeGT30Lookup) String() string {
	result := ""
	for i := 0; i < lookUpSize; i++ {
		result = result + fmt.Sprintf("CrNum: %v; mult: %v; effect: %v\n", lu.C[i], lu.Q[i], GetSymbolString(lu.Effect[i], false))
	}
	return result
}

//PrimeGTE31 : Structure to use for primes greater than or equal to 31;
//the sextuplet program only uses this for primes 31, 37, 41, 43, 47, "49",
//53, and 59; there are no need for others since these can do the checking
//for sextuplets via lookups all the way out to infinity + 1
type PrimeGTE31 struct {
	Prime  *primeBase
	LookUp *primeGT30Lookup
	//only used in displaying, for humans, the details of the GTE 31 prime
	hasInsertBefore0 bool
	//Knowing this beforehand greatly simplifies the getting N equation.
	valueSquaredEndsIn1 bool
}

//NewPrimeGTE31 : Return a pointer to an initialized PrimeGTE31
//prime is the prime value the variable will have.
func NewPrimeGTE31(prime *big.Int) *PrimeGTE31 {

	if !checkPrimeValues(prime, false) {
		fmt.Println(fmt.Sprintf("Value %v is not a legal prime GTE 31, setting it to '31'", prime))
		prime.SetInt64(31)
	}

	r := &PrimeGTE31{
		Prime:  getPrimeBase(prime),
		LookUp: getPrimeGT30Lookup(prime),
	}
	InitGTE31(r)
	return r
}

//Stringer for PrimeGTE31
func (prime *PrimeGTE31) String() string {
	return fmt.Sprintf("%v", prime.Prime) +
		fmt.Sprintf("Has insert before 0: %v\n", prime.hasInsertBefore0) +
		fmt.Sprintf("value squared ends in 1: %v\n", prime.valueSquaredEndsIn1) +
		fmt.Sprint(prime.LookUp)
}

//GetResultAtCrossNum : tests the GTE 31 primes at the given offset (crossing number) for the
//applicable effect; addResult is changed and will be accumulated in the calling
//function, n is the current n-level (0 based) one is testing.
func (prime *PrimeGTE31) GetResultAtCrossNum(addResult *int, offset, n *big.Int) bool {
	*addResult = CSextuplet

	for i := 0; i < lookUpSize; i++ {
		iCalcA.Mul(n, prime.LookUp.Q[i])
		iCalcA.Add(iCalcA, prime.LookUp.C[i])
		switch iCalcA.Cmp(offset) {
		case -1:
			continue
		case 1:
			return false
		}
		*addResult = prime.LookUp.Effect[i]
		return true
	}
	return false
}

//HumanReadable : Output to f a neatly packaged human readable format
//of the found Tuplet's details; Using this one can translate a .rawdata
//file by sending in the Tnumber and the effect integer, tNum29 is the
//corresponding TNumber from the 29Basis file and must be calculated; and notify
//is any message that needs to be communicate; generally used to inform in the pretty
//file that a basis number change occurred; currently it does check using ProbablyPrime
//if the found primes are probably prime; pretty much useless since they must be! But
//at higher numbers it is worthwhile because it would show in the pretty data if an
//erroneous result (false negative) has ever been seen.
func HumanReadable(tNum *big.Int, symbol *int, tNum29, notify *string, f *os.File) {
	tInt := TNumToInt(tNum)

	checkProbPrime := func(b *big.Int) {
		fmt.Fprint(f, b)
		//according to go guys: 100% accurate at ints less than 2^64
		//and after that should never give false negatives but can have false positives
		if !iCalcA.ProbablyPrime(20) {
			fmt.Fprintln(f, "  <== ProbablyPrime reports false!")
			return
		}
		fmt.Fprintln(f, "")
	}

	if *notify != "" {
		fmt.Fprintln(f, *notify+"\n")
		*notify = ""
	}

	fmt.Fprintf(f, fmt.Sprintf("TNum = %v\n", tNum))
	fmt.Fprintf(f, fmt.Sprintf("BeginsAt : %v\n", tInt))
	fmt.Fprintf(f, fmt.Sprintf("EndsAt : %v\n", TNumLastNatNum(tNum)))
	fmt.Fprintf(f, fmt.Sprintf("[Basis-0-TNum : %s]\n", *tNum29))
	fmt.Fprintf(f, fmt.Sprintf("---primes---   %s (%v)\n", GetSymbolString(*symbol, true), *symbol))
	switch *symbol {
	case CSextuplet, CLQuint29:
		checkProbPrime(iCalcA.Add(tInt, big12))
		checkProbPrime(iCalcA.Add(tInt, big16))
		checkProbPrime(iCalcA.Add(tInt, big18))
		checkProbPrime(iCalcA.Add(tInt, big22))
		checkProbPrime(iCalcA.Add(tInt, big24))
		if *symbol == CSextuplet {
			checkProbPrime(iCalcA.Add(tInt, big28))
			break
		}
		fmt.Fprintln(f, "x")
	case CRQuint13, CQuad:
		fmt.Fprintln(f, "x")
		checkProbPrime(iCalcA.Add(tInt, big16))
		checkProbPrime(iCalcA.Add(tInt, big18))
		checkProbPrime(iCalcA.Add(tInt, big22))
		checkProbPrime(iCalcA.Add(tInt, big24))
		if *symbol == CRQuint13 {
			checkProbPrime(iCalcA.Add(tInt, big28))
			break
		}
		fmt.Fprintln(f, "x")
	default:
		fmt.Fprintf(f, "Symbol '%v' is not valid and/or supported symbol/effect at this time.\n", *symbol)
	}
	fmt.Fprintln(f, "")
}

//GeneratePrimeTupletsInteractive : The user interactive func to begin
//the generation of juusprimes
func GeneratePrimeTupletsInteractive() {

	ctrl := NewGenPrimesStruct()

	rFile, chosen := Choose29BasisFile()
	if !chosen {
		return
	}
	ctrl.FullPathto29RawFile = rFile

	var (
		input       = ""
		wasCanceled = false
	)

	//get operation mode
	fmt.Println("Enter the operation mode you want:")
	fmt.Println("1 - by Basis  (recommended)")
	fmt.Println("2 - by Template Number (TNum)")
	fmt.Println("3 - by Natural Numbers")

	if input, wasCanceled = GetUserInputInteger("Enter Mode:", "1", "x"); wasCanceled {
		return
	}

	switch input {
	case "1":
		ctrl.OpMode = omBasis
		if input, wasCanceled = GetUserInputInteger("Enter Basis Number (0 based):", "0", "x"); wasCanceled {
			return
		}
		fmt.Sscan(input, ctrl.BasisNum)
	case "2":
		ctrl.OpMode = omTNum
		/*
		   //from testing, wanted to see that bigInt's are working, this gave 7 Sextuplets!! Took 7 hours
		   if input, wasCanceled = GetUserInputInteger("Enter From TNumber:", "18546453926011000028", "x"); wasCanceled {
		     return
		   }
		   fmt.Sscan(input, ctrl.From)
		   of input, wasCanceled = GetUserInputInteger("Enter To TNumber:", "18546453926100000027", "x"); wasCanceled {
		     return
		   }
		*/
		if input, wasCanceled = GetUserInputInteger("Enter From TNumber:", "28", "x"); wasCanceled {
			return
		}
		fmt.Sscan(input, ctrl.From)
		if input, wasCanceled = GetUserInputInteger("Enter To TNumber:", "215656468", "x"); wasCanceled {
			return
		}
		fmt.Sscan(input, ctrl.To)
	case "3":
		ctrl.OpMode = omNatNum
		if input, wasCanceled = GetUserInputInteger("Enter From Natural Number (must be GTE 835):", "835", "x"); wasCanceled {
			return
		}
		fmt.Sscan(input, ctrl.From)
		if input, wasCanceled = GetUserInputInteger("Enter To Natural Number:", "1000000", "x"); wasCanceled {
			return
		}
		fmt.Sscan(input, ctrl.To)
	default:
		fmt.Println("Invalid Operation mode")
		return
	}

	//get filter
	fmt.Println("")
	fmt.Println("<hint: No filter long time, Sextuplets less time (eg. for basis-0: 32 mins vs. 1 min.)>")
	fmt.Println("Enter the filter type you want:")
	fmt.Println("1 - No filter")
	fmt.Println("2 - Sextuplets only")
	fmt.Println("3 - Left Quints only")
	fmt.Println("4 - Right Quints only")
	fmt.Println("5 - Left & Right Quints only")
	fmt.Println("6 - Quadruplets only")
	fmt.Println("")

	if input, wasCanceled = GetUserInputInteger("Enter desired filter:", "1", "x"); wasCanceled {
		return
	}
	switch input {
	//case "1": ctrl.FilterType = ftAll
	case "2":
		ctrl.FilterType = ftSextuplet
	case "3":
		ctrl.FilterType = ftLQuint
	case "4":
		ctrl.FilterType = ftRQuint
	case "5":
		ctrl.FilterType = ftQuints
	case "6":
		ctrl.FilterType = ftQuads
	default:
		ctrl.FilterType = ftAll
	}

	ctrl.DefaultPath = DataPath

	GeneratePrimeTuplets(ctrl)
}

//GeneratePrimeTuplets : The engine func that generates juusprime Tuplets; it can
//be called with a GenPrimesStruct already filled in if you are comfortable
//with that; otherwise recommended to use GeneratePrimeTupletsInteractive
func GeneratePrimeTuplets(ctrl *GenPrimesStruct) {
	var (
		inResult, addResult int
	)

	if !FileExists(ctrl.FullPathto29RawFile) {
		fmt.Println(fmt.Sprintf("Path '%s' to 29bais rawdata is invalid, quitting.", ctrl.FullPathto29RawFile))
		return
	}

	outputFileName := ctrl.Prepare()

	//must call Prepare above before accessing ctrl's fields
	if ctrl.BasisNum.Cmp(big0) < 0 {
		fmt.Println("Basis can not be less than 0, ie., 0 for the first basis, 1 for the 2nd and so on.")
		return
	}
	if FileExists(outputFileName) {
		if approved := GetUserConfirmation(fmt.Sprintf("ATTENTION, File: \n%v\nalready exists. Do you want to overwrite it?", outputFileName), "y"); !approved {
			return
		}
	}

	tTarget := big.NewInt(0)
	nControl := big.NewInt(0)
	toTest := big.NewInt(0)
	n31 := big.NewInt(0)
	n37 := big.NewInt(0)
	n41 := big.NewInt(0)
	n43 := big.NewInt(0)
	n47 := big.NewInt(0)
	n49 := big.NewInt(0)
	n53 := big.NewInt(0)
	n59 := big.NewInt(0)

	p31 := NewPrimeGTE31(big.NewInt(31))
	p37 := NewPrimeGTE31(big.NewInt(37))
	p41 := NewPrimeGTE31(big.NewInt(41))
	p43 := NewPrimeGTE31(big.NewInt(43))
	p47 := NewPrimeGTE31(big.NewInt(47))
	p49 := NewPrimeGTE31(big.NewInt(49))
	p53 := NewPrimeGTE31(big.NewInt(53))
	p59 := NewPrimeGTE31(big.NewInt(59))

	preFilter, postFilter := GetFilter(ctrl.FilterType)

	rawF, err := FileOpen(outputFileName, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(rawF)

	prettyF, err := FileOpen(ChangeFileExt(outputFileName, fileExtPretty), false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(prettyF)

	rawData, err := FileOpen(ctrl.FullPathto29RawFile, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(rawData)

	rawScan := bufio.NewScanner(rawData)

	destroys := func(prime *PrimeGTE31, curN, finalN *big.Int) bool {
		//if one wanted to intercept the potential prime to reject it if it
		//is not prime, one would do that here. I've tried various schemes
		//already and all they do is slow down the process
		if tTarget.Cmp(prime.Prime.startTemplateNum) > -1 && curN.Cmp(finalN) < 1 {
			GetCrossNumMod(tTarget, curN, prime, toTest)
			//this test uses the primes lookups, calcs there will be at least 1, at most 6, so its quick(ish)
			if prime.GetResultAtCrossNum(&addResult, toTest, curN) {
				inResult = AddSymbols(&inResult, &addResult)
				//filtering: this post filter needs to be done twice. Here to stop unnecessary processing,
				//but also after the proc. loop ends, because some incoming tNumbers will never go through
				//this conditional if nothing effects it and the original passes.
				//An example: We want LQuints only, incoming has to allow LQuints and
				//Sextuplets, because the sextuplet may be whittled down to an LQuint. but if the sextuplet
				//passes through without anything effecting it, it will show up in the results. So that
				//final filter look has to take place for those cases.
				return FilterMap[inResult]&postFilter != 0
			}
		}
		return false
	}

	basisTrack := big.NewInt(0).Mul(basisLen, ctrl.BasisNum)
	basis29TNumberStr := "" //for the pretty print

	doScan := func(t *big.Int, e *int) bool {
		// TODO: need to also validate the file for format, this could be a real pill
		//if they throw in some random text file!
		//rawdata is in pairs on separate lines, hence 2 Scans
		r := rawScan.Scan()
		fmt.Sscan(rawScan.Text(), t)
		basis29TNumberStr = fmt.Sprintf("%v", t)
		t.Add(t, basisTrack)
		r = rawScan.Scan()
		fmt.Sscan(rawScan.Text(), e)
		return r
	}

	ClearSymbolCounts()

	//20 breaks up progress into approx. 5% chunks, 100 would be 1% chunks
	displayProgress := DisplayProgressBig(ctrl.From, ctrl.To, 100)
	//necessary because the range of TNums does not match range of possibilities in 29Basis file
	displayProgressLastPos := big.NewInt(0).Set(ctrl.From)
	startTime := DisplayProgressBookend(fmt.Sprintf("Generating prime tuplets between TNums %v & %v (filtered by: %s)",
		ctrl.From,
		ctrl.To,
		GetFilterDesc(ctrl.FilterType)), true)
	ignoredCounter := 0 //let user know it is scanning for the first possibility

	//if the file is empty and no tTarget is scanned this will break the loop immediately.
	tTarget.Add(ctrl.To, big1)
	basisNotification := fmt.Sprintf("BASIS:%v", ctrl.BasisNum)

	fmt.Fprintln(prettyF, fmt.Sprintf("TNumbers from %v to %v", ctrl.From, ctrl.To))
	fmt.Fprintln(prettyF, fmt.Sprintf("(Natural #'s from %v to %v)", TNumToInt(ctrl.From), TNumLastNatNum(ctrl.To)))
	fmt.Fprintln(prettyF, fmt.Sprintf("29Basis file used: %s", rawData.Name()))
	fmt.Fprintln(prettyF, fmt.Sprintf("filtered by: %v", GetFilterDesc(ctrl.FilterType)))

	//main loop
	for {
		if !doScan(tTarget, &inResult) {
			//crossing basis boundary, reset for the next basis
			rawData.Seek(0, 0)
			rawScan = bufio.NewScanner(rawData)
			ctrl.BasisNum.Add(ctrl.BasisNum, big1)
			basisTrack.Mul(basisLen, ctrl.BasisNum)
			basisNotification = fmt.Sprintf("BASIS WRAPPED:%v", ctrl.BasisNum)
			doScan(tTarget, &inResult)
		}

		if tTarget.Cmp(ctrl.From) == -1 {
			ignoredCounter++
			//show that something is happening if the file is being scanned to find the from TNum
			if ignoredCounter%100000 == 0 {
				fmt.Print(".")
			}
			continue
		}

		if tTarget.Cmp(ctrl.To) == 1 {
			break //only location that breaks the loop
		}
		displayProgress(tTarget, displayProgressLastPos)

		//Filter incoming symbols. After processing these, there are two other sites
		//where a post filter is needed
		if FilterMap[inResult]&preFilter == 0 {
			continue
		}

		GetNfromTNum(tTarget, p31, n31) //n31 should be the maximum of n
		GetNfromTNum(tTarget, p37, n37)
		GetNfromTNum(tTarget, p41, n41)
		GetNfromTNum(tTarget, p43, n43)
		GetNfromTNum(tTarget, p47, n47)
		GetNfromTNum(tTarget, p49, n49)
		GetNfromTNum(tTarget, p53, n53)
		GetNfromTNum(tTarget, p59, n59)

		nControl.SetInt64(0)

		//internal loop
		for nControl.Cmp(n31) < 1 {

			if destroys(p31, nControl, n31) {
				break
			}
			if destroys(p37, nControl, n37) {
				break
			}
			if destroys(p41, nControl, n41) {
				break
			}
			if destroys(p43, nControl, n43) {
				break
			}
			if destroys(p47, nControl, n47) {
				break
			}
			if destroys(p49, nControl, n49) {
				break
			}
			if destroys(p53, nControl, n53) {
				break
			}
			if destroys(p59, nControl, n59) {
				break
			}

			nControl.Add(nControl, big1)
		} //nControl.Cmp(n31) < 1

		if FilterMap[inResult]&postFilter != 0 {
			continue
		}

		fmt.Fprintln(rawF, tTarget)
		fmt.Fprintln(rawF, inResult)
		HumanReadable(tTarget, &inResult, &basis29TNumberStr, &basisNotification, prettyF)

		//accumulate statistics of what we've found
		SymbolCount[inResult]++

	} //for{(ever)}

	fmt.Println("Duration:", DisplayProgressBookend("Done", false).Sub(startTime))
	fmt.Println("")

	infoF, err := FileOpen(ChangeFileExt(outputFileName, fileExtInfo), false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(infoF)

	ShowSymbolCounts(ctrl.From, ctrl.To, ctrl.FilterType, infoF)
	fmt.Fprintln(infoF, fmt.Sprintf("29Basis file used: %s", rawData.Name()))
	fmt.Fprintln(infoF, cInfoDataFormat)
	fmt.Fprintln(infoF, cInfoSymbols)
	ShowSymbolFileDesignations(infoF)

	ShowSymbolCounts(ctrl.From, ctrl.To, ctrl.FilterType, os.Stdout)
	fmt.Println("")
	fmt.Println(fmt.Sprintf("The files:\n%s\n%s\n%s\nhave been generated.", rawF.Name(), prettyF.Name(), infoF.Name()))
}

//ShowDetails : print to screen all the GTE 31 details
func (prime *PrimeGTE31) ShowDetails(withPausing bool) {
	fmt.Println(fmt.Sprintf("P = %v", prime.Prime.value))
	fmt.Println(fmt.Sprintf("P^2 = %v", big.NewInt(0).Mul(prime.Prime.value, prime.Prime.value)))
	fmt.Println(fmt.Sprintf("T# = %v", prime.Prime.startTemplateNum))
	fmt.Println(fmt.Sprintf("T expanded = %v", TNumToInt(prime.Prime.startTemplateNum)))
	fmt.Println(prime.Prime.naturalProgression)
	fmt.Println("\nLook Up matrix:")
	fmt.Println(prime.LookUp)
	if withPausing {
		waitForInput()
	}
	prime.showNaturalProgression()
	if withPausing {
		waitForInput()
	}
	prime.showNaturalProgressionSimple()
	if withPausing {
		waitForInput()
	}
	prime.showRawDetails()
}

//showRawDetails : display the basic data of the prime structure
//and the "raw" crossing effects. Natural progressions are derived from
//this data; the func is complicated, much easier to do on paper! But it
//is a good check that the printed data matches the paper data
func (prime *PrimeGTE31) showRawDetails() {

	fmt.Println(fmt.Sprintf("P %v Raw Data ('%s' symbol means jumping over a Template)", prime.Prime.value, cInsertSymbol))
	fmt.Println("---------------------- mod based")

	inserts := big.NewInt(0).Mod(prime.Prime.value, TemplateLength)
	inserts.Sub(inserts, big1)
	insertChk := int(inserts.Int64())

	primeInt := int(prime.Prime.value.Int64())

	symbolStr := ""
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\n", "CrossNum", "effect")
	spacer := 0

	for i := 0; i < int(TemplateLength.Int64()); i++ {
		if spacer%5 == 0 {
			fmt.Fprintf(w, "%v\t%v\n", "-----", "------")
		}
		spacer++
		cNum := i % primeInt
		cResult := prime.getEffect(cNum)
		symbolStr = GetSymbolString(cResult, false)
		if i <= insertChk {
			symbolStr = cInsertSymbol + " " + symbolStr
		}
		fmt.Fprintf(w, "%v :\t%s\n", cNum, symbolStr)
	}
	w.Flush()
	fmt.Println("----------------------")
	fmt.Println("\n--- end of: P", prime.Prime.value, "Raw Data")
	fmt.Println("")
}

//showNaturalProgressionSimple : Show the natural progression for the GTE 31;
//This is the simple(r) format as opposed to the detailed one.
func (prime *PrimeGTE31) showNaturalProgressionSimple() {
	fmt.Println("P", prime.Prime.value, "Natural Progression, simple")
	fmt.Println("---------------------- Mod based")
	cNum := big.NewInt(0)
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\t%v\n", "NatProgIdx", "Cr.Num", "effect")
	spacer := 0
	for i := 0; i < int(prime.Prime.value.Int64()); i++ {
		if spacer%5 == 0 {
			fmt.Fprintf(w, "%v\t%v\t%v\n", "----------", "------", "------")
		}
		spacer++
		cNum = prime.Prime.getCrossingNumber(i)
		if cNum.Cmp(TemplateLength) > -1 {
			fmt.Fprintf(w, "Idx %v\t%v\tSKIP %s\n", i, cNum, cInsertSymbol)
		} else {
			fmt.Fprintf(w, "Idx %v\t%v\t%v\n", i, cNum, GetSymbolString(prime.getEffect(int(cNum.Int64())), false))
		}
	}
	w.Flush()
	fmt.Println("=================")
	fmt.Println("\n--- end of: P", prime.Prime.value, "Natural Progression, simple")
	fmt.Println("")
}

//ShowNaturalProgression : shows the natural progression of GTE 31's;
//This is the detailed format as opposed to the simpler one;
func (prime *PrimeGTE31) showNaturalProgression() {
	//Those with an insert at 0/1 (natural progression-wise): the insert needs to be placed at the
	//end of the sequence, this happens naturally for all the calcs, and for most of the detail displays, but
	//not this NatProg display, so need to check if we are at the "last" TNumber but still
	//have one actual count left. This requirement is for displaying the NatProg structure only,
	//the calculations have this included.

	var (
		//the various pieces are individually made and then assembled later
		insertAppend, r30, tN, output string
	)

	crossingNumMod := big.NewInt(0)
	tExp := big.NewInt(0)
	tNum := big.NewInt(0)

	i := big.NewInt(prime.Prime.startTemplateNum.Int64() - 1)
	loopControl := big.NewInt(prime.Prime.startTemplateNum.Int64() + prime.Prime.value.Int64() - 1)
	loopControlMinus1 := big.NewInt(0).Sub(loopControl, big1)

	//countActual keeps track of where we are when skips are inserted
	countActual := -1
	//count30 keeps track of the movement through the Template length (length 30)
	count30 := -1

	fmt.Println("P", prime.Prime.value, "Natural Progression, detailed")
	fmt.Println("---------------------- mod based")

	incrementCounters := func() {
		countActual++
		i.Add(i, big1)
		tNum.Set(i)
		tExp = TNumToInt(tNum)
		crossingNumMod = CrossingAtTNumExpandedMod(prime.Prime.value, tExp)
	}

	getNonSkip := func() {
		tN = fmt.Sprintf("%v (%v)   ", tNum, tExp)
		output = fmt.Sprintf("%v:%v", crossingNumMod, GetSymbolString(prime.getEffect(int(crossingNumMod.Int64())), false))
		insertAppend = fmt.Sprintf(" (%v)", countActual)
	}

	getSkip := func() {
		tempTN := fmt.Sprintf("%v (%v)   ", tNum, tExp)
		tempOutput := fmt.Sprintf("%v SKIP: %v", crossingNumMod, cInsertSymbol)
		incrementCounters()
		tN = fmt.Sprintf("%v%v\n\t+ %v (%v)   ", tempTN, tempOutput, tNum, tExp)
		output = fmt.Sprintf("%v:%v", crossingNumMod, GetSymbolString(prime.getEffect(int(crossingNumMod.Int64())), false))
		insertAppend = fmt.Sprintf(" (%v)", countActual)
	}

	appendInsert := func() {
		countActual++
		i.Add(i, big1)
		tExp = TNumToInt(i)
		insertAppend = fmt.Sprintf(" + %v (%v) SKIP: %v (%v)", i, tExp, cInsertSymbol, countActual)
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", " ", "TNum", "cr#", "cum. count")

	for i.Cmp(loopControl) == -1 {

		//keep track of the length of the GTE Prime (from 31 to 59)
		incrementCounters()
		count30++
		r30 = fmt.Sprintf("%v  ", count30)

		if crossingNumMod.Cmp(TemplateLength) > -1 {
			getSkip()
		} else {
			getNonSkip()
		}

		if prime.hasInsertBefore0 && i.Cmp(loopControlMinus1) == 0 {
			appendInsert()
		}

		fmt.Fprintf(w, "%s\tTNum %s\t%s\t%s\n", r30, tN, output, insertAppend)
		//fmt.Println(r30 + tN + output + insertAppend)
	}
	w.Flush()
	fmt.Println("\n--- end of: P", prime.Prime.value, "Natural Progression, detailed")
	fmt.Println("")
}

//getPrimeGT30Lookup : return a pointer to an
//intialized primeGT30Lookup struct
func getPrimeGT30Lookup(prime *big.Int) *primeGT30Lookup {
	r := &primeGT30Lookup{
		C:      make([]*big.Int, lookUpSize),
		Q:      make([]*big.Int, lookUpSize),
		Effect: make([]int, lookUpSize),
	}
	return r
}

//getEffect : Only use for this func is for
//details printing of the basic (non-inflated) GTE 31 prime.
func (prime *PrimeGTE31) getEffect(crossingNum int) int {
	//The mod based effect is the same for all primes > 30
	//this proc was used heavily in the primes < 30 but here it is not.
	//this is left as a check and for maybe future use, but the workhorse of
	//the primes > 30 are the lookup tables which take into account the inflation
	//of the natural progression as one moves to successive p + 30n iterations of
	//the particular prime root of the object ( 31 (61, 91...), 37 (67, 97...), etc.
	switch crossingNum {
	case 12:
		return CRQuint13
	case 16, 18, 22, 24:
		return CXNoTrack
		//16 : result := cX_17;
		//18 : result := cX_19;
		//22 : result := cX_23;
		//24 : result := cX_25;
	case 28:
		return CLQuint29
	}
	return CSextuplet
}

//InitGTE31 : fill in the appropriate data for the particular GTE 31 prime; in essence
//these are "constants" associated with that GTE 31 prime, they are calculated with
//pen and paper
func InitGTE31(prime *PrimeGTE31) {
	switch prime.Prime.value.Int64() {
	case 31:
		prime.hasInsertBefore0 = false
		prime.valueSquaredEndsIn1 = true
		for i := 0; i < len(prime.LookUp.C); i++ {
			switch i {
			case 0:
				prime.LookUp.C[i] = big.NewInt(6)
				prime.LookUp.Q[i] = big.NewInt(6)
				prime.LookUp.Effect[i] = CRQuint13
			case 1:
				prime.LookUp.C[i] = big.NewInt(10)
				prime.LookUp.Q[i] = big.NewInt(10)
				prime.LookUp.Effect[i] = CXNoTrack //cX_17;
			case 2:
				prime.LookUp.C[i] = big.NewInt(12)
				prime.LookUp.Q[i] = big.NewInt(12)
				prime.LookUp.Effect[i] = CXNoTrack //cX_19;
			case 3:
				prime.LookUp.C[i] = big.NewInt(16)
				prime.LookUp.Q[i] = big.NewInt(16)
				prime.LookUp.Effect[i] = CXNoTrack //cX_23;
			case 4:
				prime.LookUp.C[i] = big.NewInt(18)
				prime.LookUp.Q[i] = big.NewInt(18)
				prime.LookUp.Effect[i] = CXNoTrack //cX_25;
			case 5:
				prime.LookUp.C[i] = big.NewInt(22)
				prime.LookUp.Q[i] = big.NewInt(22)
				prime.LookUp.Effect[i] = CLQuint29
			}
		}
	case 37:
		prime.hasInsertBefore0 = false
		prime.valueSquaredEndsIn1 = false
		for i := 0; i < len(prime.LookUp.C); i++ {
			switch i {
			case 0:
				prime.LookUp.C[i] = big.NewInt(0)
				prime.LookUp.Q[i] = big.NewInt(0)
				prime.LookUp.Effect[i] = CXNoTrack //cX_25;
			case 1:
				prime.LookUp.C[i] = big.NewInt(5)
				prime.LookUp.Q[i] = big.NewInt(4)
				prime.LookUp.Effect[i] = CXNoTrack //cX_23;
			case 2:
				prime.LookUp.C[i] = big.NewInt(15)
				prime.LookUp.Q[i] = big.NewInt(12)
				prime.LookUp.Effect[i] = CXNoTrack //cX_19;
			case 3:
				prime.LookUp.C[i] = big.NewInt(20)
				prime.LookUp.Q[i] = big.NewInt(16)
				prime.LookUp.Effect[i] = CXNoTrack //cX_17;
			case 4:
				prime.LookUp.C[i] = big.NewInt(27)
				prime.LookUp.Q[i] = big.NewInt(22)
				prime.LookUp.Effect[i] = CLQuint29
			case 5:
				prime.LookUp.C[i] = big.NewInt(30)
				prime.LookUp.Q[i] = big.NewInt(24)
				prime.LookUp.Effect[i] = CRQuint13
			}
		}
	case 41:
		prime.hasInsertBefore0 = true
		prime.valueSquaredEndsIn1 = true
		for i := 0; i < len(prime.LookUp.C); i++ {
			switch i {
			case 0:
				prime.LookUp.C[i] = big.NewInt(2)
				prime.LookUp.Q[i] = big.NewInt(2)
				prime.LookUp.Effect[i] = CLQuint29
			case 1:
				prime.LookUp.C[i] = big.NewInt(8)
				prime.LookUp.Q[i] = big.NewInt(6)
				prime.LookUp.Effect[i] = CRQuint13
			case 2:
				prime.LookUp.C[i] = big.NewInt(16)
				prime.LookUp.Q[i] = big.NewInt(12)
				prime.LookUp.Effect[i] = CXNoTrack //cX_19;
			case 3:
				prime.LookUp.C[i] = big.NewInt(24)
				prime.LookUp.Q[i] = big.NewInt(18)
				prime.LookUp.Effect[i] = CXNoTrack //cX_25;
			case 4:
				prime.LookUp.C[i] = big.NewInt(27)
				prime.LookUp.Q[i] = big.NewInt(20)
				prime.LookUp.Effect[i] = CXNoTrack //cX_17;
			case 5:
				prime.LookUp.C[i] = big.NewInt(35)
				prime.LookUp.Q[i] = big.NewInt(26)
				prime.LookUp.Effect[i] = CXNoTrack //cX_23;
			}
		}
	case 43:
		prime.hasInsertBefore0 = false
		prime.valueSquaredEndsIn1 = false
		for i := 0; i < len(prime.LookUp.C); i++ {
			switch i {
			case 0:
				prime.LookUp.C[i] = big.NewInt(0)
				prime.LookUp.Q[i] = big.NewInt(0)
				prime.LookUp.Effect[i] = CXNoTrack //cX_25;
			case 1:
				prime.LookUp.C[i] = big.NewInt(6)
				prime.LookUp.Q[i] = big.NewInt(4)
				prime.LookUp.Effect[i] = CXNoTrack //cX_17;
			case 2:
				prime.LookUp.C[i] = big.NewInt(9)
				prime.LookUp.Q[i] = big.NewInt(6)
				prime.LookUp.Effect[i] = CRQuint13
			case 3:
				prime.LookUp.C[i] = big.NewInt(23)
				prime.LookUp.Q[i] = big.NewInt(16)
				prime.LookUp.Effect[i] = CXNoTrack //cX_23;
			case 4:
				prime.LookUp.C[i] = big.NewInt(26)
				prime.LookUp.Q[i] = big.NewInt(18)
				prime.LookUp.Effect[i] = CXNoTrack //cX_19;
			case 5:
				prime.LookUp.C[i] = big.NewInt(40)
				prime.LookUp.Q[i] = big.NewInt(28)
				prime.LookUp.Effect[i] = CLQuint29
			}
		}
	case 47:
		prime.hasInsertBefore0 = false
		prime.valueSquaredEndsIn1 = false
		for i := 0; i < len(prime.LookUp.C); i++ {
			switch i {
			case 0:
				prime.LookUp.C[i] = big.NewInt(0)
				prime.LookUp.Q[i] = big.NewInt(0)
				prime.LookUp.Effect[i] = CXNoTrack //cX_25;
			case 1:
				prime.LookUp.C[i] = big.NewInt(3)
				prime.LookUp.Q[i] = big.NewInt(2)
				prime.LookUp.Effect[i] = CLQuint29
			case 2:
				prime.LookUp.C[i] = big.NewInt(19)
				prime.LookUp.Q[i] = big.NewInt(12)
				prime.LookUp.Effect[i] = CXNoTrack //cX_19;
			case 3:
				prime.LookUp.C[i] = big.NewInt(22)
				prime.LookUp.Q[i] = big.NewInt(14)
				prime.LookUp.Effect[i] = CXNoTrack //cX_23;
			case 4:
				prime.LookUp.C[i] = big.NewInt(38)
				prime.LookUp.Q[i] = big.NewInt(24)
				prime.LookUp.Effect[i] = CRQuint13
			case 5:
				prime.LookUp.C[i] = big.NewInt(41)
				prime.LookUp.Q[i] = big.NewInt(26)
				prime.LookUp.Effect[i] = CXNoTrack //cX_17;
			}
		}
	case 49:
		prime.hasInsertBefore0 = true
		prime.valueSquaredEndsIn1 = true
		for i := 0; i < len(prime.LookUp.C); i++ {
			switch i {
			case 0:
				prime.LookUp.C[i] = big.NewInt(6)
				prime.LookUp.Q[i] = big.NewInt(4)
				prime.LookUp.Effect[i] = CXNoTrack //cX_23;
			case 1:
				prime.LookUp.C[i] = big.NewInt(16)
				prime.LookUp.Q[i] = big.NewInt(10)
				prime.LookUp.Effect[i] = CXNoTrack //cX_17;
			case 2:
				prime.LookUp.C[i] = big.NewInt(19)
				prime.LookUp.Q[i] = big.NewInt(12)
				prime.LookUp.Effect[i] = CXNoTrack //cX_25;
			case 3:
				prime.LookUp.C[i] = big.NewInt(29)
				prime.LookUp.Q[i] = big.NewInt(18)
				prime.LookUp.Effect[i] = CXNoTrack //cX_19;
			case 4:
				prime.LookUp.C[i] = big.NewInt(39)
				prime.LookUp.Q[i] = big.NewInt(24)
				prime.LookUp.Effect[i] = CRQuint13
			case 5:
				prime.LookUp.C[i] = big.NewInt(45)
				prime.LookUp.Q[i] = big.NewInt(28)
				prime.LookUp.Effect[i] = CLQuint29
			}
		}
	case 53:
		prime.hasInsertBefore0 = false
		prime.valueSquaredEndsIn1 = false
		for i := 0; i < len(prime.LookUp.C); i++ {
			switch i {
			case 0:
				prime.LookUp.C[i] = big.NewInt(0)
				prime.LookUp.Q[i] = big.NewInt(0)
				prime.LookUp.Effect[i] = CXNoTrack //cX_25;
			case 1:
				prime.LookUp.C[i] = big.NewInt(11)
				prime.LookUp.Q[i] = big.NewInt(6)
				prime.LookUp.Effect[i] = CRQuint13
			case 2:
				prime.LookUp.C[i] = big.NewInt(14)
				prime.LookUp.Q[i] = big.NewInt(8)
				prime.LookUp.Effect[i] = CLQuint29
			case 3:
				prime.LookUp.C[i] = big.NewInt(25)
				prime.LookUp.Q[i] = big.NewInt(14)
				prime.LookUp.Effect[i] = CXNoTrack //cX_17;
			case 4:
				prime.LookUp.C[i] = big.NewInt(32)
				prime.LookUp.Q[i] = big.NewInt(18)
				prime.LookUp.Effect[i] = CXNoTrack //cX_19;
			case 5:
				prime.LookUp.C[i] = big.NewInt(46)
				prime.LookUp.Q[i] = big.NewInt(26)
				prime.LookUp.Effect[i] = CXNoTrack //cX_23;
			}
		}
	case 59:
		prime.hasInsertBefore0 = true
		prime.valueSquaredEndsIn1 = true
		for i := 0; i < len(prime.LookUp.C); i++ {
			switch i {
			case 0:
				prime.LookUp.C[i] = big.NewInt(15)
				prime.LookUp.Q[i] = big.NewInt(8)
				prime.LookUp.Effect[i] = CLQuint29
			case 1:
				prime.LookUp.C[i] = big.NewInt(23)
				prime.LookUp.Q[i] = big.NewInt(12)
				prime.LookUp.Effect[i] = CXNoTrack //cX_25;
			case 2:
				prime.LookUp.C[i] = big.NewInt(27)
				prime.LookUp.Q[i] = big.NewInt(14)
				prime.LookUp.Effect[i] = CXNoTrack //cX_23;
			case 3:
				prime.LookUp.C[i] = big.NewInt(35)
				prime.LookUp.Q[i] = big.NewInt(18)
				prime.LookUp.Effect[i] = CXNoTrack //cX_19;
			case 4:
				prime.LookUp.C[i] = big.NewInt(39)
				prime.LookUp.Q[i] = big.NewInt(20)
				prime.LookUp.Effect[i] = CXNoTrack //cX_17;
			case 5:
				prime.LookUp.C[i] = big.NewInt(47)
				prime.LookUp.Q[i] = big.NewInt(24)
				prime.LookUp.Effect[i] = CRQuint13
			}
		}
	}
}
