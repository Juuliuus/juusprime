package juusprime

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

//juusprime is an engine to find prime Sextuplets and/or Quintuplets and/or Quadruplets
//GNU GPL3
//author: Julius Schön / R. Spicer, 2021
//this is the analysis unit, whimsical funcs() that allow one to crawl among the potPrimes

var (
	p31, p37, p41, p43, p47, p49, p53, p59 *PrimeGTE31
	primes                                 []*PrimeGTE31
)

//GetCritLength : deprecated, moved to PrimeGTE31 method, given fixedN (a chosen, fixed n level),
//and n the n-level you want to compare to, calculate the number of
//Templates between them, prime is a *PrimeGTE31, and
//abs flag is whether to return the absolute value, result is returned in last parameter
func GetCritLength(abs bool, p *PrimeGTE31, fixedN, n, returnHereLen *big.Int) error {
	//TODO remove on major version change
	// d( 2p + cd + 2cN ) d, diff; p, prime value; c=30; N a fixed chosen "n"
	returnHereLen.SetInt64(0)
	if n.Cmp(big0) < 0 {
		return fmt.Errorf("GetCritLength: target n (%v) must be GTE 0", n)
	}

	//d
	iCalcD.Sub(n, fixedN)

	//2p
	iCalcA.Mul(p.Prime.value, big2)

	//cd
	iCalcB.Mul(iCalcD, TemplateLength)

	//2cN
	iCalcC.Mul(fixedN, TemplateLength)
	iCalcC.Mul(iCalcC, big2)

	//sum
	iCalcA.Add(iCalcA, iCalcB)
	iCalcA.Add(iCalcA, iCalcC)

	//d(sum)
	returnHereLen.Mul(iCalcD, iCalcA)
	if abs {
		returnHereLen.Abs(returnHereLen)
	}
	return nil
}

//GetCritLengthByDiff : deprecated, moved to PrimeGTE31 method, return the total number of Templates for
//potPrime p between fromN and toN, toN can be less than fromN, if
//abs is true return the length's absolute value, result returned in last param,
func GetCritLengthByDiff(abs bool, p *PrimeGTE31, N, diff, returnHereLen *big.Int) error {
	//TODO remove on major version change
	returnHereLen.SetInt64(-1)
	if N.Cmp(big0) == -1 {
		return fmt.Errorf("GetCritLengthByDiffWF: desired N (%v) must be 0 or greater", N)
	}

	iCalcA.Add(N, diff)
	if iCalcA.Cmp(big0) == -1 {
		return fmt.Errorf("GetCritLengthByDiffWF: your diff (%v) combined with N (%v) will go below 0", diff, N)
	}
	return GetCritLength(abs, p, N, iCalcA, returnHereLen)

}

//CritLen : get the true critical length (between families), length is returned
//in the return param,  also see CritLenForceGetN()
func CritLen(csid *CritSectID, returnHereLen *big.Int) {
	//TODO: look at this, can it be adjusted for like the within Family:  d( 2p + cd + 2cN )???
	GetLocalPrimes()
	primes[csid.SubN].getFamilyFactoredCritLength(csid.N, returnHereLen)
}

//CritLenForceGetN : complement to CritLen(), given a length (in TNumbers) returns
//in the return param the n needed to get as close as possible to the provided
//given length, e.g., given the length of the 29Baseis (215656441) returns the n
//where the prime's critical length begins exceeding a complete Basis
func CritLenForceGetN(p *PrimeGTE31, len, returnHereN *big.Int) {
	GetLocalPrimes()
	//n for this function is unimportant, only subN, this converts int type for us
	ID := NewCritSectID(big0, p)
	primes[ID.SubN].getFamilyFactoredN(len, returnHereN)
}

//sextPositions : used in the testIsSextuplet() func
var sextPositions = []*big.Int{
	big.NewInt(12),
	big.NewInt(4),
	big.NewInt(2),
	big.NewInt(4),
	big.NewInt(2),
	big.NewInt(4)}

//testIsSextuplet : tests clear TNumbers to see if it actually
//is a sextuplet (ie aligns with a 29basis position) or not
func testIsSextuplet(tNumber *big.Int) bool {
	iCalcA = TNumToInt(tNumber)
	for i := 0; i < len(sextPositions); i++ {
		iCalcA.Add(iCalcA, sextPositions[i])
		if !iCalcA.ProbablyPrime(20) {
			return false
		}
	}
	return true
}

//constant used in PrimeGTE31's to align them
//with a CritSectID's SubN property, 0-7 then
//represents 31-59
const (
	fam31SubN = iota
	fam37SubN
	fam41SubN
	fam43SubN
	fam47SubN
	fam49SubN
	fam53SubN
	fam59SubN
)

//getNandSubN : analysis use, returns in params the true
//crit sect identifiers N and subN belonging to the range
//tTarget lives in
func getNandSubN(tTarget, n *big.Int, subN *int64) {
	GetLocalPrimes()

	var highestPrime int64
	highestPrime = -1
	maxN := big.NewInt(-1)
	GetNfromTNum(tTarget, p31, maxN)

	//find which prime is the last to have 31's N, ie. which subN is it?
	for i := range primes {
		GetNfromTNum(tTarget, primes[i], iCalcFuncResult)
		if primes[i].Helper.MaxN.Cmp(maxN) == 0 {
			highestPrime = int64(i)
			//fmt.Println("highestPrime", highestPrime)
		}
	}
	*subN = highestPrime
	n.Set(primes[highestPrime].Helper.MaxN)
}

//AnalyzeCritSectByTNumber : Given a TNumber will analyze its parent CritSectID for
//clear channels and/or sextuplets, details shows a lot of information, the calculated
//crit length will be inclusive, caller be aware: the "to" var here is processed as LTE(<=)
func AnalyzeCritSectByTNumber(tTarget *big.Int, details bool) {

	GetLocalPrimes()
	from, to := big.NewInt(0), big.NewInt(0)
	maxN := big.NewInt(-1)

	var highestPrime int64
	highestPrime = -1

	getNandSubN(tTarget, maxN, &highestPrime)

	GetEffectiveTNum(maxN, primes[highestPrime], from)

	//example, maxN is at pP(49):
	//31: n=3
	//37: n=3
	//41: n=3
	//43: n=3
	//47: n=3
	//49: n=3
	//53: n=2
	//59: n=2
	//highestPrime from above will be 49, we want the effective start of this into "from"
	//want to find the start of the next pP: 53 to go into "to"
	//we need to increase 53's N so that it gives back the next effective start after 49 at n=3
	//also properly handles the movement from a n=3 59 to a n=4 31

	highestPrime = (highestPrime + 1) % 8 //keep it in range of 0-7
	primes[highestPrime].Helper.MaxN.Add(primes[highestPrime].Helper.MaxN, big1)
	//fmt.Println("highestPrime2=", highestPrime)

	GetEffectiveTNum(primes[highestPrime].Helper.MaxN, primes[highestPrime], to)

	//the crit length from "from" to "to" is inclusive
	to.Sub(to, big1)

	//CountClearTNumbers(from, to, <selectedN's>) //interesting! use lower levels for analysis maybe
	CountClearTNumbers(from, to, maxN, details)
}

//AnalyzeCritSectByCSIDRange : Given a from and to CritSectID will analyze in detail all
//TNumbers in that (inclusive) range, be aware this can be a LOT of information, details shows
//very detailed information, primesOnly will reject non-Prime pP's and lump together until the
//next Prime pP, doAdd assumes the "toID" is to be added to the "fromID" instead of being
//a static toID, this always prints a summary and will print a ".B." represenation where "."
//is a critical section where there are clear channels, and "B" is a critical section that is "black",
//i.e., no clear channels
func AnalyzeCritSectByCSIDRange(fromID, toID *CritSectID, primesOnly, doAdd, details bool) {

	GetLocalPrimes()
	black, tot, comp, sext := 0, 0, 0, 0
	var sl []int

	from, to := big.NewInt(0), big.NewInt(0)
	maxN := big.NewInt(-1)
	skipList := ""
	//will contain a string representation of the critical sections
	//"." means there was an open channel, even if it was not a sextuplet
	//"B" means the crit sect was "black", no open channels
	var sb strings.Builder

	//finalID is break check to exit infinite loop
	//when doAdd is true the toID should be added to the fromID
	//otherwise it is a static csID that is the end of the range
	finalID := NewCritSectID(big0, p31)
	switch doAdd {
	case true:
		finalID.SetFromID(fromID)
		finalID.Add(toID)
	default:
		if fromID.Compare(toID) == 1 {
			fmt.Println(fmt.Sprintf("GetTrueCritSectionFromCritSectIDs: toID %v is less than fromID %v", toID, fromID))
			return
		}
		finalID.SetFromID(toID)
	}

	batch := fmt.Sprintf("from ID %v to ID %v", fromID, finalID)
	if primesOnly {
		batch = batch + "(+)"
	}

	//always inc finalID to set it to ID =after= the last requested ID
	//leave this operation BELOW the batch setting above or the numbers
	//will be misleading
	finalID.Increment(1)

	//locID is the "current" csID, it walks the values
	locID := NewCritSectID(big0, p31)
	locID.SetFromID(fromID)

	//calcID is always looking forward to the next csID(s)
	calcID := NewCritSectID(big0, p31)
	calcID.SetFromID(locID)

	fmt.Println("********************************************************************")
	fmt.Println(fmt.Sprintf("******* CritSect by ID: %s", batch))
	fmt.Println("")

	iter := -1
	testPrime := big.NewInt(0)

	GetEffectiveTNumSimple(locID.N, primes[locID.SubN], from)
	for {
		iter = 0
		GetNfromTNum(from, primes[locID.SubN], maxN)
		calcID.Increment(1) //look forward to next boundary

		if primesOnly {
			//throw out non-Prime pP's
			for {
				primes[calcID.SubN].MemberAtN(calcID.N, testPrime)
				if testPrime.ProbablyPrime(5) {
					break
				}
				skipList = skipList + " " + fmt.Sprintf("%v", testPrime)
				calcID.Increment(1)
				iter++
			}
		}
		if skipList != "" {
			fmt.Println("--Skipped non-Primes:", skipList)
		}
		fmt.Println(fmt.Sprintf("%v to %v", locID, calcID))

		GetEffectiveTNumSimple(calcID.N, primes[calcID.SubN], to)
		to.Sub(to, big1) //make it inclusive

		//meat and potatoes call
		//sl will be []int with length 3: counts of sextuplets, clears, total
		sl = CountClearTNumbers(from, to, maxN, details)

		sext = sext + sl[0]
		switch sl[0] {
		case 0:
		default:
			sb.WriteString(strconv.Itoa(sl[0]))
		}
		comp = comp + sl[1]

		switch sl[2] {
		case 0:
			black++
			sb.WriteString("B")
		default:
			sb.WriteString(".")
		}
		tot++

		skipList = ""
		from.Add(to, big1)

		locID.Increment(1 + iter)

		if locID.N.Cmp(finalID.N) > -1 && locID.SubN >= finalID.SubN {
			break
		}
	}
	fmt.Println("dot representation:")
	fmt.Println(sb.String())
	fmt.Println("")
	fmt.Println("===Summary:")
	fmt.Println(fmt.Sprintf("%s", batch))
	fmt.Println(fmt.Sprintf("primesOnly=%v, doAdd: %v", primesOnly, doAdd))
	fmt.Println(fmt.Sprintf("#CS: %v   black=%v  sext=%v, comp=%v  combinedS&C=%v", tot, black, sext, comp, comp+sext))

	return
}

//CritSectID : abbrev csID, Unique identifier for each Critical Section in format "n:subN",
//N is the usual n, subN is an index from 0 to 7 corresponding with the 8 potPrimes
//31 to 59, while there are associated "math" routines the CS ID's have no 0 or
//negative value, ie., 0:0 is equivalent to 1 and stands for n=0, pP=31
type CritSectID struct {
	N *big.Int
	//actually relates to the idx, 0-7, of the 8 pP familes
	//but for traversing true crit sects a "sub"N more sense
	SubN int64
}

//NewCritSectID : returns a pointer to a CritSectID struct, n >= 0,
//p is any PrimeGTE31, each of which have a constant "subN"
func NewCritSectID(n *big.Int, p *PrimeGTE31) *CritSectID {
	if n.Cmp(big0) == -1 {
		n.Set(big0)
	}
	return &CritSectID{
		N:    big.NewInt(0).Set(n),
		SubN: p.subN,
	}
}

//Stringer for CritSectID
func (id *CritSectID) String() string {
	return id.AsString()
}

//AsString : central func to return a string
//representation of the CritSectID
func (id *CritSectID) AsString() string {
	return fmt.Sprintf("%v:%v", id.N, id.SubN)
}

//Increment : increments the csID by value, value must >= 0
func (id *CritSectID) Increment(value int) error {
	if value < 0 {
		return fmt.Errorf("Increment: value %v is less than 0, decrement not allowed, convert to integers first then subtract", value)
	}
	for i := 1; i <= value; i++ {
		if id.SubN == 7 {
			id.N.Add(id.N, big1)
		}
		id.SubN = (id.SubN + 1) % 8
	}
	return nil
}

//ID2Int : This gives back the integer representation of the csID allowing a unique
//integer for each csID, also math can be used and a csID can then be set from
//an the resulting integer, see SetFromInt
func (id *CritSectID) ID2Int() *big.Int {
	result := big.NewInt(0).Set(id.N)
	subn := big.NewInt(0).SetInt64(id.SubN)
	result.Mul(result, big8)
	result.Add(result, subn)
	return result.Add(result, big1) //switch from 0 base to 1 base
}

//SetFromInt : Set the csID by an integer value, value is a corresponding integer
//of any csID, obtain it with ID2Int()
func (id *CritSectID) SetFromInt(value *big.Int) error {
	if value.Cmp(big1) == -1 {
		id.N.Set(big0)
		id.SubN = 0
		return fmt.Errorf("SetFromInt: value %v is less than 1", value)
	}
	val := big.NewInt(0).Set(value)

	//reset to 0 based
	val.Sub(val, big1)

	id.N.Div(val, big8)
	id.SubN = val.Mod(val, big8).Int64()
	return nil
}

//SetFromID : Sets, copies, the id's values to those of csID
func (id *CritSectID) SetFromID(csID *CritSectID) error {
	if csID == nil {
		return fmt.Errorf("SetFromID: csID is nil")
	}
	id.N.Set(csID.N)
	id.SubN = csID.SubN
	return nil
}

//SetFromTNumber : Sets id's vars to the appropriate values based on tNumber
func (id *CritSectID) SetFromTNumber(tNumber *big.Int) {
	getNandSubN(tNumber, id.N, &id.SubN)
}

//SetFromString : sets id's vars from a string representation of format "n:subN", eg.  0:0, 123:7,
//n >=0, subN is 0-7, only returns error if an invalid format otherwise adjusts to valid values
func (id *CritSectID) SetFromString(value string) error {
	sl := strings.Split(value, ":")
	if len(sl) != 2 {
		id.N.Set(big0)
		id.SubN = 0
		return fmt.Errorf("SetFromString: value %v is invalid, expecting #:#", value)
	}

	_, err := fmt.Sscan(sl[0], id.N)
	if err != nil {
		return err
	}
	if id.N.Cmp(big0) == -1 {
		id.N.Set(big0)
		fmt.Println("SetFromString: invalid N, setting to 0")
	}

	res, err := strconv.Atoi(sl[1])
	if err != nil {
		return err
	}
	id.SubN = int64(res)

	if id.SubN < 0 || id.SubN > 7 {
		id.SubN = 0
		fmt.Println("SetFromString: invalid SubN, setting to 0")
	}
	return nil
}

//GetTNumberRange : generally analysis use only, returns in params the corresponding
//from/to TNumber and its (inclusive) length (L)
func (id *CritSectID) GetTNumberRange(from, to, L *big.Int) {
	GetLocalPrimes()

	next := NewCritSectID(big0, p31)

	CritLen(id, L)
	GetEffectiveTNumSimple(id.N, primes[id.SubN], from)
	next.SetFromID(id)
	next.Increment(1)
	GetEffectiveTNumSimple(next.N, primes[next.SubN], to)
	to.Sub(to, big1) //adjust to be inclusive
}

//Compare : compares id with csID, returns -1 if <, 0 if =, +1 if greater
func (id *CritSectID) Compare(csID *CritSectID) int {
	return id.ID2Int().Cmp(csID.ID2Int())
}

//Add : adds csID to id
func (id *CritSectID) Add(csID *CritSectID) {
	idInt := id.ID2Int()
	idInt.Add(idInt, csID.ID2Int())
	id.SetFromInt(idInt)
}

//Subtract : subtracts csID from id, csID's have no 0, if the subtraction results
//in 0 or less id is set to 0:0 (equiv to 1), caller is reponsible to compare
//and only subtract if csID is < id
func (id *CritSectID) Subtract(csID *CritSectID) error {
	idInt := id.ID2Int()
	csidInt := csID.ID2Int()

	idInt.Sub(idInt, csidInt)
	if idInt.Cmp(big1) == -1 {
		idInt.Set(big1)
	}
	id.SetFromInt(idInt)
	if csidInt.Cmp(idInt) > -1 {
		return fmt.Errorf("ID %v is GTE this ID %v (setting to 0:0 equivalent to '1')", csID, id)
	}
	return nil
}

//CountClearTNumbers : fromTNum/toTNum are inclusive, detailed gives a lot
//of output, returned []int has 3 rows: Sextuplet, Clear, and total counts, this
//func tests each TNumber to see if it has a clear channel, if so adds to sextuplet total
//if a true sextuplet, otherwise goes into "unused" clear channels, caller be aware -
//- toTNum is processed as LTE (<=) since it is inclusive
func CountClearTNumbers(fromTNum, toTNum, maxN *big.Int, detailed bool) []int {
	GetLocalPrimes()
	cnt := 0

	locN := big.NewInt(0)
	offset := big.NewInt(0)
	tTarget := big.NewInt(0).Set(fromTNum)
	var (
		comp []string
		sext []string
	)

	addResult := 0
	failure := false

	var (
		progPrecision int64
		startTime     time.Time
	)
	tnRange := big.NewInt(0).Sub(toTNum, fromTNum)
	tnRange.Add(tnRange, big1)
	progPrecision = 5
	if tnRange.Cmp(big.NewInt(10000)) > -1 {
		progPrecision = 100
	}
	displayProgress := DisplayProgressBig(fromTNum, toTNum, progPrecision)
	displayProgressLastPos := big.NewInt(0).Set(fromTNum)

	if detailed {
		fmt.Println("====  CountClearTNumbers  ====")
		startTime = DisplayProgressBookend("", true)
	}

	failed := func(p *PrimeGTE31) bool {
		if locN.Cmp(p.Helper.MaxN) == 1 {
			//fmt.Println("skipped n", locN, "p", p.Prime.value)
			return false
		}
		//older alternative -> GetCrossNumModDirect(tTarget, locN, p, offset)
		GetCrossNumModSimple(tTarget, locN, p, offset)
		if p.GetResultAtCrossNum(&addResult, offset, locN) {
			failure = true
			return true
		}
		return false
	}

	for tTarget.Cmp(toTNum) < 1 {
		if detailed {
			displayProgress(tTarget, displayProgressLastPos)
		}
		failure = false
		locN.Set(big0)

		//GetNfromTNum modified to set prime.Helper.MaxN directly
		//so here we can toss the result param
		GetNfromTNum(tTarget, p31, iCalcFuncResult)
		GetNfromTNum(tTarget, p37, iCalcFuncResult)
		GetNfromTNum(tTarget, p41, iCalcFuncResult)
		GetNfromTNum(tTarget, p43, iCalcFuncResult)
		GetNfromTNum(tTarget, p47, iCalcFuncResult)
		GetNfromTNum(tTarget, p49, iCalcFuncResult)
		GetNfromTNum(tTarget, p53, iCalcFuncResult)
		GetNfromTNum(tTarget, p59, iCalcFuncResult)

		for locN.Cmp(maxN) < 1 {
			if failed(p31) {
				break
			}
			if failed(p37) {
				break
			}
			if failed(p41) {
				break
			}
			if failed(p43) {
				break
			}
			if failed(p47) {
				break
			}
			if failed(p49) {
				break
			}
			if failed(p53) {
				break
			}
			if failed(p59) {
				break
			}
			locN.Add(locN, big1)
		}

		if !failure {
			cnt++
			switch testIsSextuplet(tTarget) {
			case true:
				sext = append(sext, fmt.Sprintf("%v", tTarget))
			default:
				comp = append(comp, fmt.Sprintf("%v", tTarget))
			}
		}

		tTarget.Add(tTarget, big1)
	}
	sl := []int{len(sext), len(comp), cnt}
	if detailed {
		fmt.Println("Duration:", DisplayProgressBookend("Done", false).Sub(startTime))
		fmt.Println("---------------------")
		fmt.Println("N Depths:")
		for i := range primes {
			fmt.Println(fmt.Sprintf("%v: n=%v", primes[i].Prime.value, primes[i].Helper.MaxN))
		}

		fmt.Println("TNumber Clear (false Sext) ", sl[1], ":\n", comp)
		fmt.Println("TNumber: Sextuplet(s) ", sl[0], ":\n", sext)
		fmt.Println(fmt.Sprintf(":: count=%v (%v to %v, %v TNumbers searched)", sl[2], fromTNum, toTNum, tnRange))

		fmt.Println("")
	}
	return sl
}

//GetPrimeGTE31Slice : helper for analysis routines, gives
//back a slice of all Primes GTE 31, len = 8, so [0]-[7] is 31-59
func GetPrimeGTE31Slice() []*PrimeGTE31 {
	return []*PrimeGTE31{
		P31(),
		P37(),
		P41(),
		P43(),
		P47(),
		P49(),
		P53(),
		P59(),
	}
}

//GetLocalPrimes : helper for analysis routines, initializes the
//p## *PrimeGTE31 variables for direct use, can be called any number of times
func GetLocalPrimes() {
	if p31 == nil {
		p31 = P31()
		p37 = P37()
		p41 = P41()
		p43 = P43()
		p47 = P47()
		p49 = P49()
		p53 = P53()
		p59 = P59()
	}
	if primes == nil {
		primes = []*PrimeGTE31{p31, p37, p41, p43, p47, p49, p53, p59}
	}
}

//P31 : Call to get a *PrimeGTE31 for pP(31) for whatever purpose you like
func P31() *PrimeGTE31 {
	return NewPrimeGTE31(big.NewInt(31))
}

//P37 : Call to get a *PrimeGTE31 for pP(37) for whatever purpose you like
func P37() *PrimeGTE31 {
	return NewPrimeGTE31(big.NewInt(37))
}

//P41 : Call to get a *PrimeGTE31 for pP(41) for whatever purpose you like
func P41() *PrimeGTE31 {
	return NewPrimeGTE31(big.NewInt(41))
}

//P43 : Call to get a *PrimeGTE31 for pP(43) for whatever purpose you like
func P43() *PrimeGTE31 {
	return NewPrimeGTE31(big.NewInt(43))
}

//P47 : Call to get a *PrimeGTE31 for pP(47) for whatever purpose you like
func P47() *PrimeGTE31 {
	return NewPrimeGTE31(big.NewInt(47))
}

//P49 : Call to get a *PrimeGTE31 for pP(49) for whatever purpose you like
func P49() *PrimeGTE31 {
	return NewPrimeGTE31(big.NewInt(49))
}

//P53 : Call to get a *PrimeGTE31 for pP(53) for whatever purpose you like
func P53() *PrimeGTE31 {
	return NewPrimeGTE31(big.NewInt(53))
}

//P59 : Call to get a *PrimeGTE31 for pP(59) for whatever purpose you like
func P59() *PrimeGTE31 {
	return NewPrimeGTE31(big.NewInt(59))
}

//getResultAtCrossNumANALYSIS : for analysis only, an alternative routine that
//provides results in an analysis format,
//tests the GTE 31 primes at the given offset (crossing number) for the
//applicable effect, addResult is changed and will be accumulated in the calling
//function, n is the current level (0 based) one is testing, int result is used as a flag
func (prime *PrimeGTE31) getResultAtCrossNumANALYSIS(addResult *int, offset, n *big.Int) (bool, int) {
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
			return false, (i * -1) - 1
		}
		*addResult = prime.LookUp.Effect[i]
		return true, i
	}
	return false, -99
}

//AnalyzeTNumbersInteractive : interactive full analysis of any
//TNumbers you like >= 32, =careful= large TNumbers could take a while
//and can produce LOTS of output.
func AnalyzeTNumbersInteractive() {
	var (
		input       string
		wasCanceled bool
	)
	offset := big.NewInt(0)
	//offsetSimple := big.NewInt(0)

	primes := GetPrimeGTE31Slice()

	giveTNum := big.NewInt(0)
	effectiveP := big.NewInt(0)
	addResult := CSextuplet
	pass := ""
	n := big.NewInt(0)
	big0 := big.NewInt(0)
	big1 := big.NewInt(1)
	broken := false
	primesOnly := true
	doPause := true

	fmt.Println("")
	fmt.Println("")
	fmt.Println("===============================")
	fmt.Println("==> Be aware: These routines assume that any TNumbers entered have already been")
	fmt.Println("processedby the primes 7-29! If that is not the case you =CAN BE FOOLED= into ")
	fmt.Println("thinking a TNumber you entered is a Tuplet. For this reason, recommended is ")
	fmt.Println("to enter TNumbers from a 29bais catalog.")
	fmt.Println("===============================")
	fmt.Println("")
	fmt.Println("")

	if wasCanceled = GetUserBoolChoice("Show broken/effected only?", "false", &broken); wasCanceled {
		return
	}

	if wasCanceled = GetUserBoolChoice("Show prime pP's only?", "true", &primesOnly); wasCanceled {
		return
	}

	if wasCanceled = GetUserBoolChoice("Pause for review  after each potPrime?", "true", &doPause); wasCanceled {
		return
	}

	if input, wasCanceled = GetUserInput("Enter TNumbers comma separated:", "535, 647", "x"); wasCanceled {
		return
	}

	sl := strings.Split(input, ",")
	fmt.Println("analysing ", sl)

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 1)
	wFmt := "%s\t  %s\t %v \t %s \t  (%s  %s %s)\n"

	effectStr := ""
	wStr := ""
	primeStr := ""
	printIt := false
	wasBroken := true
	big32 := big.NewInt(32)
	idx := -1
	found := false
	C, Q, I := "", "", ""

	for x := range sl {
		fmt.Sscan(sl[x], giveTNum)
		if giveTNum.Cmp(big32) < 0 {
			fmt.Println(fmt.Sprintf("Skipping %v, TNumbers must be >= %v", giveTNum, big32))
			continue
		}

		GetNfromTNum(giveTNum, primes[0], primes[0].Helper.MaxN)
		GetNfromTNum(giveTNum, primes[1], primes[1].Helper.MaxN)
		GetNfromTNum(giveTNum, primes[2], primes[2].Helper.MaxN)
		GetNfromTNum(giveTNum, primes[3], primes[3].Helper.MaxN)
		GetNfromTNum(giveTNum, primes[4], primes[4].Helper.MaxN)
		GetNfromTNum(giveTNum, primes[5], primes[5].Helper.MaxN)
		GetNfromTNum(giveTNum, primes[6], primes[6].Helper.MaxN)
		GetNfromTNum(giveTNum, primes[7], primes[7].Helper.MaxN)
		// Debug:
		//offsetCheck := big.NewInt(0)
		fmt.Println("")
		fmt.Println("==== Target:", giveTNum, "MaxN:", primes[0].Helper.MaxN)
		for i := range primes {
			fmt.Println("==== Begin  -- p", primes[i].Prime.Value(), "  max N=", primes[i].Helper.MaxN, "  target: ", giveTNum)
			n.Set(big0)
			runningCnt := 0
			wasBroken = false
			for n.Cmp(primes[i].Helper.MaxN) < 1 {

				primeStr = ""
				primes[i].MemberAtN(n, effectiveP)
				prime := effectiveP.ProbablyPrime(20)
				switch prime {
				case false:
					primeStr = "not prime"
				default:
					primeStr = ""
				}
				if primesOnly && !prime {
					n.Add(n, big1)
					continue
				}

				runningCnt++
				//The three funcs below are equivalent, I'm prefering Direct or simple because they
				//don't need an effectiveTNumber to work and saves a few processor cycles.
				//GetCrossNumModDirect(giveTNum, n, primes[i], offset)
				//new
				GetCrossNumModSimple(giveTNum, n, primes[i], offset)
				// Debug:
				//GetCrossNumMod(giveTNum, n, primes[i], offsetCheck)

				pass = "ok"
				printIt = !broken
				effectStr = "-"
				if found, idx = primes[i].getResultAtCrossNumANALYSIS(&addResult, offset, n); found {
					pass = "altered"
					if addResult == 4 {
						pass = "broken!"
					}
					wasBroken = true
					printIt = true
					effectStr = strconv.Itoa(addResult)
				}

				if printIt {
					wStr = fmt.Sprintf("%v: n=%v  %v %s", runningCnt, n, effectiveP, primeStr)
					C, Q = "", ""
					switch idx {
					case -99:
						I = "early"
					case -1, -2, -3, -4, -5, -6:
						I = strconv.Itoa(idx)
					default:
						C = fmt.Sprintf("%v", primes[i].LookUp.C[idx])
						Q = fmt.Sprintf("%v", primes[i].LookUp.Q[idx])
						I = strconv.Itoa(idx + 1) //I think of the lookups in terms of 1 based!
					}
					fmt.Fprintf(w, fmt.Sprintf(wFmt, wStr, effectStr, offset, pass, I, C, Q))
					// Debug:
					//fmt.Fprintf(w, fmt.Sprintf(wFmt, wStr, "", offsetCheck, "--="))
				}
				n.Add(n, big1)
			}

			if broken {
				if !wasBroken {
					fmt.Fprintf(w, fmt.Sprintf(wFmt, " ok", "", "", "", "", "", ""))
				} else {
					//fmt.Fprintf(w, fmt.Sprintf(wFmt, " ==>  EFFECTED  <==", "", "", ""))
					fmt.Fprintf(w, fmt.Sprintf(wFmt, "", "", "", "", "", "", "                           ==>  EFFECTED  <=="))
				}
			}
			w.Flush()

			fmt.Println("==== End  -- p", primes[i].Prime.Value(), "  max N=", primes[i].Helper.MaxN, "  target: ", giveTNum)
			fmt.Println("")
			fmt.Println("")
			if !broken && doPause {
				waitForInput()
			}
		}
		if broken && doPause {
			waitForInput()
		}
	}
}

//InflatePrimeGTE31Position : Given any idx <= prime.value-1 use n to calculate
//its expanded position. This is similar to the Lookup Table where c + qn is used,
//but this can be for any position in the prime's cycle using the CQModel slice,
//used for analysis, offset result returned in parameter, effect at that pos in result
func InflatePrimeGTE31Position(prime *PrimeGTE31, idx int, n, returnHereOffset *big.Int) (effect int) {
	if idx >= len(prime.CQModel) {
		fmt.Println(fmt.Sprintf("InflatePrimeGTE31Position: idx %v is out of range 0-%v", idx, len(prime.CQModel)-1))
		effect = -1
		returnHereOffset.SetInt64(-1)
		return
	}
	effect = prime.CQModel[idx].CEffect
	//c + qn, c is 0-based idx into the prime natural progression, q is a constant from nat. progression
	iCalcA.SetInt64(int64(idx))
	iCalcB.SetInt64(int64(prime.CQModel[idx].Q30))
	returnHereOffset.Mul(iCalcB, n)
	returnHereOffset.Add(returnHereOffset, iCalcA)
	return
}

//BuildInflationMap : Builds, from first principles, the inflated structure for p at
//n by repeadedly adding the values, essentially equivalent to doing the inflation
//by hand, used primarily as a check to other inflation routines
func BuildInflationMap(p *PrimeGTE31, n *big.Int) {
	realC := big.NewInt(-1)
	tNum := big.NewInt(0)
	lastTNum := big.NewInt(0)
	tDiff := big.NewInt(0)
	nPlus1 := big.NewInt(0).Set(n)
	nPlus1.Add(nPlus1, big1)

	//start gets our pP at n (31,61,91,...etc for pp(31) and so on)
	//starts initial value becomes the constant to be added
	start := big.NewInt(0)
	p.MemberAtN(n, start)
	addByHand := big.NewInt(0).Set(start)
	//then we square start so that we will begin at that pP's
	//effective starting TNumber
	start.Mul(start, start)

	lastTNum = IntToTNum(start)

	nCntrl := int(n.Int64())
	msg := fmt.Sprintf("Inflation Map for prime %v: n=%v -->  pP=%v", p.Prime.value, n, addByHand)

	const (
		inflate = "TNumbers in dashes (eg. -123-) are spaces added by inflation"
		skip    = "TNumbers in with > (eg. > 123) are the natural progression skip spaces"
	)

	fmt.Println("")
	fmt.Println(inflate)
	fmt.Println(skip)
	fmt.Println("=================")
	fmt.Println(msg)
	fmt.Println("=================")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 1)
	wFmt := "%s\t  %v\t %v \n"

	//need to print out past the end of the section to force the
	//last skip space results to show, so we iterate to 31 instead of 30
	//so there will be a final row where q again is equal to 0
	for i := 0; i < 31; i++ {
		realC.Add(realC, big1)
		tNum = IntToTNum(start)
		tDiff.Sub(tNum, lastTNum)

		//determines if a natural skip occurred
		if tDiff.Cmp(nPlus1) == 1 {
			fmt.Fprintf(w, fmt.Sprintf(wFmt, "", fmt.Sprintf("> %v", tDiff.Sub(tNum, big1)), fmt.Sprintf("c=%v", realC)))
			realC.Add(realC, big1)
		}

		if i > 29 {
			w.Flush()
			fmt.Println("--cycle complete--")
			break
		}

		//print the q line
		fmt.Fprintf(w, fmt.Sprintf(wFmt, fmt.Sprintf("q = %v", i), tNum, fmt.Sprintf("c=%v", realC)))
		start.Add(start, addByHand)
		lastTNum.Set(tNum)

		//add inflation spaces
		for j := 1; j <= nCntrl; j++ {
			realC.Add(realC, big1)
			fmt.Fprintf(w, fmt.Sprintf(wFmt, "", fmt.Sprintf("-%v-", tNum.Add(tNum, big1)), fmt.Sprintf("c=%v", realC)))
		}

	}
	fmt.Println("=================")
	fmt.Println(msg)
	fmt.Println("=================")
	fmt.Println(inflate)
	fmt.Println(skip)
	fmt.Println("")

}

//CheckTwinSextuplet : testing, check rawdata files for twin
//Sextuplets, can also serve as template for other quick tests one wants.
func CheckTwinSextuplet(filename string, out *os.File) {

	var inResult int

	if !FileExists(filename) {
		fmt.Println(fmt.Sprintf("Path '%s' to rawdata is invalid, quitting.", filename))
		return
	}

	rawData, err := FileOpen(filename, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(rawData)

	rawScan := bufio.NewScanner(rawData)

	doScan := func(t *big.Int, e *int) bool {
		r := rawScan.Scan()
		fmt.Sscan(rawScan.Text(), t)
		r = rawScan.Scan()
		fmt.Sscan(rawScan.Text(), e)
		return r
	}

	big7 := big.NewInt(7)
	tCompare := big.NewInt(0)
	count := 0
	tTarget := big.NewInt(0)

	for {
		if !doScan(tTarget, &inResult) {
			fmt.Println("End of file")
			break
		}

		if inResult != CSextuplet {
			continue
		}

		iCalcA.Sub(tTarget, tCompare)
		if iCalcA.Cmp(big7) == 0 {
			count++
			fmt.Fprintln(out, fmt.Sprintf("%v: %v %v", count, tCompare, tTarget))
		}
		tCompare.Set(tTarget)

	} //for{(ever)}

	fmt.Println("finished:", count)
}
