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

var (
	primeGTE31AllowedVals = []int64{31, 37, 41, 43, 47, 49, 53, 59}
	isAutomated           = false
)

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

//PrimeHelper struct : holds variables that don't belong properly to a
//prime's properties, but that can be useful in calculations, looping, etc.
type PrimeHelper struct {
	//can use this to control "n" looping
	MaxN *big.Int
}

func getPrimeHelper() *PrimeHelper {
	return &PrimeHelper{
		MaxN: big.NewInt(0),
	}
}

//Stringer for lookups
func (lu *primeGT30Lookup) String() string {
	result := ""
	for i := 0; i < lookUpSize; i++ {
		result = result + fmt.Sprintf("CrNum: %v; mult: %v; effect: %v\n", lu.C[i], lu.Q[i], GetSymbolString(lu.Effect[i], false))
	}
	return result
}

//PrimeGTE31InflationModel : This models a natural progression such that
//one can project forward to any n any position in the natural Progression
//and/or reverse engineer, or de-inflate, an inflated potPrime, used for
//analysis and testing, it also holds effect information; Wait is used as a helper
//var when re-constructing the q var, which is the inflation factor which is then
//based on n level
type PrimeGTE31InflationModel struct {
	Q30     int
	CEffect int
	//rather than have to look ahead to the next index position, I store
	//this value to help the routine and keep it simple
	Wait bool
}

//getPrimeGTE31InflationModel : returns a pointer to an initialized
//PrimeGTE31InflationModel struct
func getPrimeGTE31InflationModel() *PrimeGTE31InflationModel {
	r := &PrimeGTE31InflationModel{
		Q30:     0,
		CEffect: 0,
		Wait:    false,
	}
	return r
}

//PrimeGTE31 : Structure to use for primes greater than or equal to 31;
//the sextuplet program only uses this for primes 31, 37, 41, 43, 47, "49",
//53, and 59; there are no need for others since these can do the checking
//for sextuplets via lookups all the way out to infinity + 1
type PrimeGTE31 struct {
	Helper  *PrimeHelper
	Prime   *primeBase
	LookUp  *primeGT30Lookup
	CQModel []*PrimeGTE31InflationModel

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
		Helper:  getPrimeHelper(),
		Prime:   getPrimeBase(prime),
		LookUp:  getPrimeGT30Lookup(prime),
		CQModel: make([]*PrimeGTE31InflationModel, prime.Int64()),
	}
	for i := 0; i < len(r.CQModel); i++ {
		r.CQModel[i] = getPrimeGTE31InflationModel()
	}

	InitGTE31(r)
	return r
}

//Stringer for PrimeGTE31
func (prime *PrimeGTE31) String() string {
	CQ := ""
	for i := 0; i < len(prime.CQModel); i++ {
		CQ = CQ + fmt.Sprintf("%v:%v:%v:%v ", i, prime.CQModel[i].CEffect, prime.CQModel[i].Q30, prime.CQModel[i].Wait)
	}
	return fmt.Sprintf("%v", prime.Prime) +
		fmt.Sprintf("Has insert before 0: %v\n", prime.hasInsertBefore0) +
		fmt.Sprintf("value squared ends in 1: %v\n", prime.valueSquaredEndsIn1) +
		fmt.Sprintf("CQ lookup slice:\nindex:effect:q-value:wait-value\n%s\n", CQ) +
		fmt.Sprint("Lookup Table:\n") +
		fmt.Sprint(prime.LookUp)
}

//GetResultAtCrossNum : tests the GTE 31 primes at the given offset (crossing number) for the
//applicable effect; addResult is changed and will be accumulated in the calling
//function, n is the current n-level (0 based) one is testing, offset is calculated by
//GetCrossNumModDirect() or GetCrossNumMod()
func (prime *PrimeGTE31) GetResultAtCrossNum(addResult *int, offset, n *big.Int) bool {
	//if the function does not find  that there is an effect at offset & n then it is a
	//pass and CSextuplet needs to be returned
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

//GetQbyReverseInflation : takes n and offset from other routines and returns,
//in the parameter, a q appropriate for that n; the q is used to "deconstruct" an
//inflated potPrime so that the offset can be re-traced to its actual natural
//progression position so that its effect can be retrieved; this is a VERY INEFFICIENT
//func, and is meant for testing and analyis only; it has been used to search for sextuplets
//and returns correct results;
//Lookup tables are the efficient way to see whether an offset is has an effect
//of interest, if offset is in an inflated region then returns -1 in param
func (prime *PrimeGTE31) GetQbyReverseInflation(n, offset, returnHereQ *big.Int) error {
	//see notes at bottom of func

	prime.MemberAtN(n, iCalcD)
	if n.Cmp(big0) == -1 || offset.Cmp(iCalcD) != -1 || offset.Cmp(big0) == -1 {
		return fmt.Errorf("invalid params, n (%v) must be >= 0 and/or offset (%v) must be >=0 and less than prime value", n, offset)
	}

	returnHereQ.Set(big0)
	cnt := big.NewInt(-1)

	runningQ := big.NewInt(0)
	inflated := n.Cmp(big0) == 1

	for i := 0; i < len(prime.CQModel); i++ {
		//move through arrary, each position counts. zero based.
		cnt.Add(cnt, big1)

		//landed on inflated crossing exactly
		if cnt.Cmp(offset) == 0 {
			returnHereQ.Set(runningQ)
			return nil
		}

		if prime.CQModel[i].Wait {
			//if current iteration members have same Q value do NOT add inflation!
			//this is hard coded along with Q30 pos so we don't need to look ahead here...
			continue
		}

		if inflated {
			//cnt is incremented here too
			cnt.Add(cnt, n)
			//keep track of how many inflation spaces, ie. q for this n
			runningQ.Add(runningQ, n)
			if cnt.Cmp(offset) > -1 {
				//ignore, effect is CSextuplet because its in inflated space
				//this flag must be checked by calling func
				returnHereQ.SetInt64(-1)
				return nil
			}
		}

	}
	return nil

	//this routine was developed to be able to analyze offsets to ANY position into and/or out of
	//an inflated potPrime. Long ago I had thought this would be un-doable,
	//but it turns out it is not. It is just supremely inefficient. The reason it is
	//difficult is that, in an inflated potPrime, one needs 3 variables, c, q and n, to be able to get
	//back to the natural progression crossing, and we only have access to 2 vars: offset and n.

	//Normal testing is done by lookup tables that can project out the position we are
	//interested in to find out if an effect has happened, this direction is "easy"
	//The reverse process of coming from an inflated potPrime back to its origin is "hard"

	//This procedure gives back one of the unknown vars: q (already adjusted for n and so ready to
	//use immediately without multiplication), the calculation must be very precise and must keep track
	//of two running totals and handling "skip" spaces and "inflation" spaces properly.

	//an entire struct was developed to keep track of this and to help this calculation.
	//below is a diagrammatic represention, Numbers in parens, which can be single or double. If double
	//the first number is a "skip" space and there is NO inflation between these pairs of numbers.
	//I stands for the inflation spaces which be from 0 to infinity. I is tied directly to n. It is the
	//sum of these I's that is returned as the q.

	//(0) I (1,2) I (3) I (4,5) I (6,7) I (8)...

	//The diagram above is exactly the natural progression when n, and therefore I, is zero.

}

//GetResultAtCrossNumByReverseInflation : Alternative to GetResultAtCrossNum(); Used for
//testing and analysis because it is very inefficient and slow; It uses GetQbyReverseInflation()
//which re-contructs the inflation at n to get the q needed to take the offset back to its
//natural progression index to get its effect; addResult is changed and will be accumulated in the calling
//function, n is the current n-level (0 based) one is testing, offset is calculated by
//GetCrossNumModDirect() or GetCrossNumMod()
func (prime *PrimeGTE31) GetResultAtCrossNumByReverseInflation(addResult *int, offset, n *big.Int) bool {

	const (
		attn  = "ATTN: Result is meaningless!"
		where = "func GetResultAtCrossNumByReverseInflation()"
	)
	//if the function does not find  that there is an effect at offset & n then it is a
	//pass and CSextuplet needs to be returned
	*addResult = CSextuplet

	//since this func is used for analysis and playing around, quite important to validate params
	if n.Cmp(big0) == -1 {
		fmt.Println("")
		fmt.Println(where)
		fmt.Println(fmt.Sprintf("%s n=%v, n must be GTE 0", attn, n))
		return false
	}

	prime.MemberAtN(n, iCalcD)
	if offset.Cmp(iCalcD) != -1 || offset.Cmp(big0) == -1 {
		fmt.Println("")
		fmt.Println(where)
		fmt.Println(fmt.Sprintf("%s offset %v must be >=0 and less than potPrime %v", attn, offset, iCalcD))
		return false
	}

	prime.GetQbyReverseInflation(n, offset, iCalcA)

	//GetQbyReverseInflation func returns -1 if offset is in inflation space
	if iCalcA.Cmp(big0) == -1 {
		return false
	}

	//The returned q is already adjusted for n, only need to subtract
	iCalcB.Sub(offset, iCalcA)
	*addResult = prime.CQModel[iCalcB.Int64()].CEffect
	//fmt.Println("p=", prime.Prime.value, "n=", n, offset, iCalcA, *addResult)

	return *addResult > CSextuplet
}

//MemberAtN : return the member of the potPrime family at
//n; e.g. family 31, n=0 return 31, n=1 return 61, n=2 return 91, etc.
func (p *PrimeGTE31) MemberAtN(n, returnMember *big.Int) {
	returnMember.Mul(n, TemplateLength)
	returnMember.Add(returnMember, p.Prime.Value())
}

//displayResultsAtCrossNum : Internal use testing/debugging
func (prime *PrimeGTE31) displayResultsAtCrossNum(n *big.Int) string {
	result := ""
	for i := 0; i < lookUpSize; i++ {
		iCalcA.Mul(n, prime.LookUp.Q[i])
		iCalcA.Add(iCalcA, prime.LookUp.C[i])
		result = result + fmt.Sprintf("n:%v  C=%v Q=%v Calc=%v    E=%v\n",
			n, prime.LookUp.C[i], prime.LookUp.Q[i], iCalcA, prime.LookUp.Effect[i])
	}
	return result
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
	if !isAutomated && FileExists(outputFileName) {
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
			//original func, slower GetCrossNumMod(tTarget, curN, prime, toTest)
			GetCrossNumModDirect(tTarget, curN, prime, toTest)
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
			//You can also use GetResultAtCrossNumByReverseInflation() which is great for
			//testing and analysis, but extremely inefficient. Of theoretical importance but
			//not for practical use.
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

	//add vars for testing for twin sextuplets
	twinCheck := big.NewInt(0)
	big7 := big.NewInt(7)
	twinCheckStr := ""

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

		//twin sextuplet check for TNum to be 7 difference
		if inResult == CSextuplet {
			iCalcA.Sub(tTarget, twinCheck)
			/*
				if true {
					_ = big7
					fmt.Println("Twin Sextuplets found, TNumbers:", twinCheck, tTarget, iCalcA)
					twinCheckStr = twinCheckStr + fmt.Sprintf("\n%v %v %v", twinCheck, tTarget, iCalcA)
				}
			*/
			if iCalcA.Cmp(big7) == 0 {
				fmt.Println("Twin Sextuplets found, TNumbers:", twinCheck, tTarget)
				twinCheckStr = twinCheckStr + fmt.Sprintf("\n%v %v", twinCheck, tTarget)
			}
			twinCheck.Set(tTarget)
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
	ShowTwinSextupletResults(&twinCheckStr, infoF)
	fmt.Fprintln(infoF, fmt.Sprintf("29Basis file used: %s", rawData.Name()))
	fmt.Fprintln(infoF, cInfoDataFormat)
	fmt.Fprintln(infoF, cInfoSymbols)
	ShowSymbolFileDesignations(infoF)

	ShowSymbolCounts(ctrl.From, ctrl.To, ctrl.FilterType, os.Stdout)
	ShowTwinSextupletResults(&twinCheckStr, os.Stdout)
	fmt.Println("")
	fmt.Println(fmt.Sprintf("The files:\n%s\n%s\n%s\nhave been generated.", rawF.Name(), prettyF.Name(), infoF.Name()))
}

//AutomationStruct : used to pass automation params into
//automation routines
type AutomationStruct struct {
	BasisFile, OutputPath    string
	FromBasisNum, ToBasisNum string
	Filter                   int
	//Overwrite                bool
}

//GetNewAutomationStructure : return a pointer to an AutomationStruct
func GetNewAutomationStructure() *AutomationStruct {
	return &AutomationStruct{
		BasisFile:    "",
		OutputPath:   "",
		Filter:       1,
		FromBasisNum: "0",
		ToBasisNum:   "0",
		//Overwrite:    false,
	}
}

//GeneratePrimeTupletsAutomated : For use with automation through code or
//shell scripts, etc. Pass in a filled automation structure.
func GeneratePrimeTupletsAutomated(auto *AutomationStruct) int {

	if !FileExists(auto.BasisFile) {
		fmt.Println(fmt.Sprintf("Automation: Path '%s' to 29basis rawdata is invalid.", auto.BasisFile))
		return 1
	}

	if !FileExists(auto.OutputPath) {
		fmt.Println(fmt.Sprintf("Automation: Output Path '%s' is invalid.", auto.OutputPath))
		return 1
	}

	locFrom := big.NewInt(-1)
	locTo := big.NewInt(-2)

	fmt.Sscan(auto.FromBasisNum, locFrom)
	if locFrom.Cmp(big0) == -1 {
		fmt.Println(fmt.Sprintf("Automation: From Basis Num '%v' must >= 0.", locFrom))
		return 1
	}

	fmt.Sscan(auto.ToBasisNum, locTo)
	if locTo.Cmp(locFrom) == -1 {
		fmt.Println(fmt.Sprintf("Automation: To Basis Num '%v' must greater than or eqaul to From Basis Num '%v'.", locTo, locFrom))
		return 1
	}

	if auto.Filter >= ftCount {
		fmt.Println(fmt.Sprintf("Automation: Filter choice '%v' is invalid.", auto.Filter))
		return 1
	}

	isAutomated = true
	// TODO: Overwriting!!! a bit painful, must call prepare repeatedly....
	//	outputFileName := ctrl.Prepare()
	//	if !auto.Overwrite {
	//	}

	ctrl := NewGenPrimesStruct()

	ctrl.FullPathto29RawFile = auto.BasisFile
	ctrl.OpMode = omBasis
	ctrl.FilterType = auto.Filter
	ctrl.DefaultPath = auto.OutputPath

	for locFrom.Cmp(locTo) < 1 {
		ctrl.BasisNum.Set(locFrom)
		fmt.Println("")
		fmt.Println("=============================")
		fmt.Println(fmt.Sprintf("--- Processing basis-%v ---", locFrom))
		fmt.Println("=============================")
		fmt.Println("")
		GeneratePrimeTuplets(ctrl)
		locFrom.Add(locFrom, big1)
	}
	return 0
}

//ShowDetails : print to screen all the GTE 31 details
func (prime *PrimeGTE31) ShowDetails(withPausing bool) {
	fmt.Println(prime)
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
//these are "constants" associated with that GTE 31 prime, they can be calculated with
//pen and paper
func InitGTE31(prime *PrimeGTE31) {
	//constants before the switch are used to unwind effects of
	//TNumber mod arithmetic to normal mod arithmetic

	//store the pP mod 30 value for normal arithmetic mod routines, 1 for 31, 7 for 37, etc.
	prime.Prime.mod30.Mod(prime.Prime.value, TemplateLength)

	//get the primes offset into its own Starting Template
	prime.Prime.modOffset.Mod(prime.Prime.startTemplateNum, prime.Prime.value)

	//constant that reflects the initial Template offsets
	//C + modOffset - prime's value
	//= modOffset-mod30 (same as 30 + offset - (30+mod30))
	//prime.Prime.modConst.Add(TemplateLength, prime.Prime.modOffset)
	//prime.Prime.modConst.Sub(prime.Prime.modConst, prime.Prime.value)
	prime.Prime.modConst.Sub(prime.Prime.modOffset, prime.Prime.mod30)

	fillCQ := func(inflate []int) {
		val := 0
		idx := 0
		for i := 0; i < len(prime.CQModel); i++ {
			prime.CQModel[i].Q30 = val
			/*
				if inflate[idx] == 99 {
					continue
				}
			*/
			if inflate[idx] == i {
				prime.CQModel[i].Wait = true
				idx++
				continue
			}
			val++
		}
	}

	switch prime.Prime.value.Int64() {
	case 31:
		prime.hasInsertBefore0 = false
		prime.valueSquaredEndsIn1 = true

		//See case 37 for explanation of this
		sl := []int{24, 999}
		fillCQ(sl)
		prime.CQModel[6].CEffect = CRQuint13
		prime.CQModel[10].CEffect = CXNoTrack //cX_17;
		prime.CQModel[12].CEffect = CXNoTrack //cX_19;
		prime.CQModel[16].CEffect = CXNoTrack //cX_23;
		prime.CQModel[18].CEffect = CXNoTrack //cX_25;
		prime.CQModel[22].CEffect = CLQuint29

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

		//This structure is used for reverse inflation Analysis
		//rather than type all the structures by hand, one only needs
		//to indicate to fillCQ() which postiions in the natural progression
		//hang together as a group, that is, one is a "skip" space and those
		//pairs (and sometimes triplets) of numbers have the same q value.
		sl := []int{1, 7, 12, 17, 22, 28, 33, 999}
		fillCQ(sl)
		prime.CQModel[0].CEffect = CXNoTrack  //cX_25;
		prime.CQModel[5].CEffect = CXNoTrack  //cX_23;
		prime.CQModel[15].CEffect = CXNoTrack //cX_19;
		prime.CQModel[20].CEffect = CXNoTrack //cX_17;
		prime.CQModel[27].CEffect = CLQuint29
		prime.CQModel[30].CEffect = CRQuint13

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

		sl := []int{3, 6, 10, 14, 18, 21, 25, 29, 32, 36, 39, 999}
		fillCQ(sl)
		prime.CQModel[2].CEffect = CLQuint29
		prime.CQModel[8].CEffect = CRQuint13
		prime.CQModel[16].CEffect = CXNoTrack //cX_19;
		prime.CQModel[24].CEffect = CXNoTrack //cX_25;
		prime.CQModel[27].CEffect = CXNoTrack //cX_17;
		prime.CQModel[35].CEffect = CXNoTrack //cX_23;

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

		sl := []int{1, 4, 8, 11, 14, 17, 21, 24, 27, 31, 34, 37, 41, 999}
		fillCQ(sl)
		prime.CQModel[0].CEffect = CXNoTrack //cX_25;
		prime.CQModel[6].CEffect = CXNoTrack //cX_17;
		prime.CQModel[9].CEffect = CRQuint13
		prime.CQModel[23].CEffect = CXNoTrack //cX_23;
		prime.CQModel[26].CEffect = CXNoTrack //cX_19;
		prime.CQModel[40].CEffect = CLQuint29

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

		sl := []int{1, 4, 6, 9, 12, 15, 17, 20, 23, 26, 28, 31, 34, 37, 40, 42, 45, 999}
		fillCQ(sl)
		prime.CQModel[0].CEffect = CXNoTrack //cX_25;
		prime.CQModel[3].CEffect = CLQuint29
		prime.CQModel[19].CEffect = CXNoTrack //cX_19;
		prime.CQModel[22].CEffect = CXNoTrack //cX_23;
		prime.CQModel[38].CEffect = CRQuint13
		prime.CQModel[41].CEffect = CXNoTrack //cX_17;

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

		sl := []int{2, 4, 7, 9, 12, 15, 17, 20, 22, 25, 28, 30, 33, 35, 38, 40, 43, 46, 47, 999}
		fillCQ(sl)
		prime.CQModel[6].CEffect = CXNoTrack  //cX_23;
		prime.CQModel[16].CEffect = CXNoTrack //cX_17;
		prime.CQModel[19].CEffect = CXNoTrack //cX_25;
		prime.CQModel[29].CEffect = CXNoTrack //cX_19;
		prime.CQModel[39].CEffect = CRQuint13
		prime.CQModel[45].CEffect = CLQuint29

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

		sl := []int{1, 3, 5, 8, 10, 12, 15, 17, 19, 21, 24, 26, 28, 31, 33, 35, 38, 40, 42, 45, 47, 49, 51, 999}
		fillCQ(sl)
		prime.CQModel[0].CEffect = CXNoTrack //cX_25;
		prime.CQModel[11].CEffect = CRQuint13
		prime.CQModel[14].CEffect = CLQuint29
		prime.CQModel[25].CEffect = CXNoTrack //cX_17;
		prime.CQModel[32].CEffect = CXNoTrack //cX_19;
		prime.CQModel[46].CEffect = CXNoTrack //cX_23;

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

		sl := []int{1, 3, 5, 7, 9, 11, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 57, 999}
		fillCQ(sl)
		prime.CQModel[15].CEffect = CLQuint29
		prime.CQModel[23].CEffect = CXNoTrack //cX_25;
		prime.CQModel[27].CEffect = CXNoTrack //cX_23;
		prime.CQModel[35].CEffect = CXNoTrack //cX_19;
		prime.CQModel[39].CEffect = CXNoTrack //cX_17;
		prime.CQModel[47].CEffect = CRQuint13

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
