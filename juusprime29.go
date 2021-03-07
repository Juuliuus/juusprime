package juusprime

//juusprime is an engine to find prime Sextuplets and/or Quintuplets and/or Quadruplets
//GNU GPL3
//author: Julius Schön / R. Spicer, 2021

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"text/tabwriter"
)

//Outside of experimentation there is no need for primes other than these in this unit
var primeLTE29AllowedVals = []int64{7, 11, 13, 17, 19, 23, 29}

//PrimeIF : interface for the prime struct's and useful methods;
//experimental, needs to be developed
type PrimeIF interface {
	//GetLookupEffect()
	fillNaturalProgression()
}

//primeBase : struct containing the common elements of all juusprime primes
//for the prime sextuplet project
type primeBase struct {
	value              *big.Int
	valueSquared       *big.Int
	startTemplateNum   *big.Int
	naturalProgression []*big.Int
}

//Stringer for primeBase
func (p *primeBase) String() string {
	return fmt.Sprintf("P=%v\nP^2=%v\nT#=%v\nT expanded=%v\n%v\n", p.value, p.valueSquared, p.startTemplateNum, TNumToInt(p.startTemplateNum), p.naturalProgression)
}

//Value : Getter for a prime's value
func (p *primeBase) Value() *big.Int {
	return p.value
}

//ValueSquared : Getter for a prime's valueSquared
func (p *primeBase) ValueSquared() *big.Int {
	return p.valueSquared
}

//StartTemplateNum : Getter for a prime's startTemplateNum
func (p *primeBase) StartTemplateNum() *big.Int {
	return p.startTemplateNum
}

//NaturalProgression : Getter for a prime's naturalProgression
func (p *primeBase) NaturalProgression() []*big.Int {
	return p.naturalProgression
}

//NaturalProgressionAtIndex : Getter for a an individual value of a
//prime's naturalProgression; caller is responsible to keep idx's in range
func (p *primeBase) NaturalProgressionAtIndex(idx int, returnMe *big.Int) error {
	if idx < 0 || idx >= len(p.naturalProgression) {
		returnMe.Set(big0)
		return nil
	}
	returnMe.Set(p.naturalProgression[idx])
	return nil
}

//NaturalProgressionAtIdx : Deprecated, use NaturalProgressionAtIndex; Getter for a an individual value of a
//prime's naturalProgression; caller is responsible to keep idx's in range
func (p *primeBase) NaturalProgressionAtIdx(idx int) (*big.Int, error) {
	if idx < 0 || idx >= len(p.naturalProgression) {
		return getBigInt(big0), fmt.Errorf("idx %v out of bounds", idx)
	}
	return p.naturalProgression[idx], nil
}

//fillNaturalProgression : prime struct method. The natural progression of
//the prime's crossings as it moves through the template numbers.
func (p *primeBase) fillNaturalProgression() {
	if p == nil {
		fmt.Println("TODO: determine if methods implementing interfaces should always have nil check.")
		panic("panicing, don't know what to do about this yet. Would be uncommon.")
	}
	/*
		The point of the primes < 30 is to generate the 29Basis pattern that is repeatedly followed
		out to infinity. Rather than calculate each crossing individually, this below
		allows one to calculate the crossings at any Template# and then "roll" out to any
		other Template#. This also allows looking at any region you like...
		Natural Progression is, in mod, the crossings from 0 to Prime-1 allowing one to
		look up the current effect at the current idx (0 to Prime-1).

		This func is used heavily in the Primes <= 29. It is available in Primes >29 but really only useful
		there for informational purposes. Position finding is done with other methods in primes 31 and greater.

		This starts at the Primes startTNum and increments over all values of the
		prime (0 to prime.value-1) to calculate the order in which the prime crosses TNum boundaries
	*/
	tNum := big.NewInt(0).Set(p.startTemplateNum)
	for i := 0; i < len(p.naturalProgression); i++ {
		p.naturalProgression[i] = CrossingAtTNumMod(p.value, tNum)
		tNum.Add(tNum, big1) //move to next TNumber
	}
}

//GetNaturalProgressionIdx : internal use only; used only in specific
//situations for LTE 29's
func (p *primeBase) GetNaturalProgressionIdx(tNumber *big.Int) int {
	cNum := CrossingAtTNumMod(p.value, tNumber)
	for i := 0; i < int(p.value.Int64()); i++ {
		if cNum.Cmp(p.naturalProgression[i]) == 0 {
			return i
		}
	}
	return -1
}

//getCrossingNumber : internal use only and of limited usefulness;
//used only in showing "raw" details for GTE 31s; from older pascal code
//could be refactored?
func (p *primeBase) getCrossingNumber(naturalProgressionIdx int) *big.Int {
	return p.naturalProgression[naturalProgressionIdx]
}

//PrimeLTE29 : Structure to use for primes 7 to less than or equal to 29.
//Primes LTE 29 are useful only in generating the sequence showing all
//possible prime sextuplets (the 29Basis)
type PrimeLTE29 struct {
	Prime          *primeBase
	CrossNumEffect map[*big.Int]int
}

//Stringer for PrimeLTE29
func (prime *PrimeLTE29) String() string {
	return fmt.Sprintf("%v", prime.Prime) + "crossNumEffect:\n" + fmt.Sprint(prime.CrossNumEffect) + "\n"
}

//checkPrimeValues : internal function to control allowed values
func checkPrimeValues(prime *big.Int, fromLTE29 bool) bool {
	isLegal := false

	switch fromLTE29 {
	case true:
		for i := 0; i < len(primeLTE29AllowedVals); i++ {
			if prime.Int64() == primeLTE29AllowedVals[i] {
				isLegal = true
				break
			}
		}
	case false:
		for i := 0; i < len(primeGTE31AllowedVals); i++ {
			if prime.Int64() == primeGTE31AllowedVals[i] {
				isLegal = true
				break
			}
		}
	}
	return isLegal
}

//getPrimeBase : common function returning a pointer to an
//intialized primeBase struct. prime is the value the Prime struct will use;
//used internally when creating the various juusprime Primes
func getPrimeBase(prime *big.Int) *primeBase {
	r := &primeBase{
		value:              getBigInt(prime),
		valueSquared:       big.NewInt(0).Mul(prime, prime),
		naturalProgression: make([]*big.Int, prime.Int64()),
	}
	r.startTemplateNum = IntToTNum(getBigInt(r.valueSquared))
	//r.fillNaturalProgression()
	fillNatProg(r)
	return r
}

//NewPrimeLTE29 : Return a pointer to an initialized PrimeLTE29
//prime is the prime value the variable will have.
func NewPrimeLTE29(prime *big.Int) *PrimeLTE29 {

	if !checkPrimeValues(prime, true) {
		fmt.Println(fmt.Sprintf("Value %v is not a legal prime LTE 29, setting it to '7'", prime))
		prime.SetInt64(7)
	}

	r := &PrimeLTE29{
		Prime:          getPrimeBase(prime),
		CrossNumEffect: make(map[*big.Int]int),
	}
	r.fillEffectsMap()
	return r
}

//fillEffectsMap : internal function used to make lookup of an effect
//at a crossing number more intuitive
func (prime *PrimeLTE29) fillEffectsMap() {
	for i := 0; i < len(prime.Prime.naturalProgression); i++ {
		prime.CrossNumEffect[prime.Prime.naturalProgression[i]] = prime.GetEffect(prime.Prime.naturalProgression[i])
	}
}

/*
//Repurpose : prime is the prime value for the Prime struct. This
//re-purposes an existing one. Experimental, not currently used.
func (p *primeBase) Repurpose (prime *big.Int) {
	p.value = prime
	p.valueSquared = big.NewInt(0).Mul(prime, prime)
	p.naturalProgression = make([]*big.Int, prime.Int64())
	p.startTemplateNum = IntToTNum(getBigInt(p.valueSquared))
	p.fillNaturalProgression()
}
*/

//fillNatProg : Interface version, currently just for testing and development;
//It currently has no valid use case but I am dipping a toe in the interface
//waters..
func fillNatProg(p PrimeIF) {
	p.fillNaturalProgression()
}

//GetEffect : Given the prime's crossing number return the effect the
//crossing has for that Template Number. These are "constants" in the
//sense that for each prime the patterns repeat forever; These numbers came
//from on paper analysis; We assumes it is being compared to the 2,3,5 template
//and, so, against a Sextuplet
func (prime *PrimeLTE29) GetEffect(crossingNum *big.Int) int {
	switch prime.Prime.value.Int64() {
	case 7:
		//  3, 1, 6, 4, 2, 0, 5, 7's natural progression
		switch crossingNum.Int64() {
		case 0:
			return CLQuint29
		case 1, 2, 3, 4:
			return CXNoTrack
			//these constants were made in the case that I wanted to track these
			//Tuplet killers individually, but have so far found no need. However,
			//if needed, it is important they have the cases as shown.
			//1 : result := cX_23;
			//2 : result := cX_17;
			//3 : result := cX_25;
			//4 : result := cX_19;
		case 5:
			return CRQuint13
		}
	case 11:
		switch crossingNum.Int64() {
		case 0, 2, 5, 7:
			return CXNoTrack
			//0 : result := cX_23;
			//2 : result := cX_25;
			//5 : result := cX_17;
			//7 : result := cX_19;
		case 1:
			return CRQuint13
		case 6:
			return CLQuint29
		}
	case 13:
		switch crossingNum.Int64() {
		case 3, 5, 9, 11:
			return CXNoTrack
			//9 : result := cX_23;
			//11 : result := cX_25;
			//3 : result := cX_17;
			//5 : result := cX_19;
		case 12:
			return CRQuint13
		case 2:
			return CLQuint29
		}
	case 17:
		switch crossingNum.Int64() {
		case 1, 5, 7, 16:
			return CXNoTrack
			//5 : result := cX_23;
			//7 : result := cX_25;
			//16 : result := cX_17;
			//1 : result := cX_19;
		case 12:
			return CRQuint13
		case 11:
			return CLQuint29
		}
	case 19:
		switch crossingNum.Int64() {
		case 3, 5, 16, 18:
			return CXNoTrack
			//3 : result := cX_23;
			//5 : result := cX_25;
			//16 : result := cX_17;
			//18 : result := cX_19;
		case 12:
			return CRQuint13
		case 9:
			return CLQuint29
		}
	case 23:
		switch crossingNum.Int64() {
		case 1, 16, 18, 22:
			return CXNoTrack
			//22 : result := cX_23;
			//1 : result := cX_25;
			//16 : result := cX_17;
			//18 : result := cX_19;
		case 12:
			return CRQuint13
		case 5:
			return CLQuint29
		}
	case 29:
		switch crossingNum.Int64() {
		case 16, 18, 22, 24:
			return CXNoTrack
			//22 : result := cX_23;
			//24 : result := cX_25;
			//16 : result := cX_17;
			//18 : result := cX_19;
		case 12:
			return CRQuint13
		case 28:
			return CLQuint29
		}
	}
	//all other indexes allow Sextuplets to exist
	return CSextuplet
}

//ShowDetails : produce the details in mod & counting numbers
//of how the prime will effect the Prime Templates
func (prime *PrimeLTE29) ShowDetails(withPausing bool) {
	fmt.Println(prime)
	if withPausing {
		waitForInput()
	}
	fmt.Println("Natural Progression details:\n=====")
	fmt.Println("Prime:", prime.Prime.value)
	fmt.Println("TNum = # (expanded); Crossing: MOD crossing (COUNTING crossing)")
	fmt.Println("And finally the effect at that crossing (blank means no effect)\n\n--------")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	tNum := big.NewInt(0)

	for i := prime.Prime.startTemplateNum.Int64(); i < prime.Prime.startTemplateNum.Int64()+prime.Prime.value.Int64(); i++ {
		tNum.SetInt64(i)
		crossingNum := CrossingAtTNumMod(prime.Prime.value, tNum)

		fmt.Fprintf(w, "TNum = %v (%v)\tCrossing: %v\t(%v)\t= %v\n",
			tNum,
			TNumToInt(tNum),
			crossingNum,                          //mod
			big.NewInt(0).Add(crossingNum, big1), //counting
			GetEffectDisplay(prime.GetEffect(crossingNum)))
	}
	w.Flush()
	fmt.Println("========= END ", prime.Prime.value)
	fmt.Println("")
	if withPausing {
		waitForInput()
	}
	prime.ShowRawDetails()
}

//GeneratePrimes7to23 : Trivial output, but important as sanity check
//that the routines give back the correct effects; The primes 7 to 23
//cannot be generated using the general methods I've written, this func
//prints them out for completeness, but they are not for use in any further
//routines or analysis
func GeneratePrimes7to23() {

	var effectPtr int
	addResult := -1
	type controlStruct struct {
		idx  int
		size int
	}

	rootpath := filepath.Join(DataPath, "juusprimes_1_27%s")
	rawName := fmt.Sprintf(rootpath, fileExtRaw23)
	prettyName := fmt.Sprintf(rootpath, fileExtPretty)
	infoName := fmt.Sprintf(rootpath, fileExtInfo)

	//	rootpath := filepath.Join(rootpath, "primes_1_27."+fileExtRaw)
	rawF, err := FileOpen(rawName, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(rawF)

	prettyF, err := FileOpen(prettyName, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(prettyF)

	infoF, err := FileOpen(infoName, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(infoF)

	//closure; the starting position of a primes effects are staggered
	//so this func takes that into account.
	adder := func(prime *PrimeLTE29) func(*big.Int) {
		cs := controlStruct{idx: 0, size: int(prime.Prime.value.Int64())}
		return func(comp *big.Int) {
			if comp.Cmp(prime.Prime.startTemplateNum) >= 0 {
				//get prime's crossing number at this TNumber and add its effect'
				effectPtr = prime.CrossNumEffect[prime.Prime.naturalProgression[cs.idx]]
				cs.idx = (cs.idx + 1) % cs.size //increment but keep it in mod range of 0..6 for 7, 0..10 for 11, etc.
				addResult = AddSymbols(&addResult, &effectPtr)
			}
		}
	}

	i7 := adder(NewPrimeLTE29(big.NewInt(7)))
	i11 := adder(NewPrimeLTE29(big.NewInt(11)))
	i13 := adder(NewPrimeLTE29(big.NewInt(13)))
	i17 := adder(NewPrimeLTE29(big.NewInt(17)))
	i19 := adder(NewPrimeLTE29(big.NewInt(19)))
	i23 := adder(NewPrimeLTE29(big.NewInt(23)))
	tNum29Str := "n/a"
	basisNumStr := "n/a"

	//get the pattern up to, but not including, prime 29.
	//prime 29 starts at TNumber 28...so iterate TNums 1 to 27
	for i := 1; i < 28; i++ {
		comp := big.NewInt(int64(i))
		//for each TNum start with Sextuplet (235 temlate) and let prime's TNum effect try to destroy it
		addResult = CSextuplet
		i7(comp)
		i11(comp)
		i13(comp)
		i17(comp)
		i19(comp)
		i23(comp)
		if !SymbolIsOfInterest(addResult) {
			continue
		}

		fmt.Fprintln(rawF, i)
		fmt.Fprintln(rawF, addResult)
		HumanReadable(big.NewInt(int64(i)), &addResult, &tNum29Str, &basisNumStr, prettyF)
	}
	fmt.Fprintln(infoF, cInfo23Text)
	fmt.Fprintln(infoF, cInfoSymbols)
	fmt.Println(fmt.Sprintf("The files:\n%s\n%s\n%s\nhave been generated.", rawName, prettyName, infoName))
}

//GenerateBasis : the default and recommended 29 basis file generation
//routine for new/inexperienced users
func GenerateBasis() {

	ctrl := NewGenPrimesStruct()
	ctrl.FilterType = ftAll

	ctrl.From.Set(basisBegin)
	ctrl.To.Set(basisEnd)

	ctrl.fileNameBase = filepath.Join(Basis29Path, fmt.Sprintf("%s_%v_%v_%s%s", filePrefix29, ctrl.From, ctrl.To, GetFilterAbbrev(ctrl.FilterType), fileExtRaw29))

	if FileExists(ctrl.fileNameBase) {
		if approved := GetUserConfirmation(fmt.Sprintf("ATTENTION, File: \n%v\nalready exists. Do you want to overwrite it?", ctrl.fileNameBase), "y"); !approved {
			return
		}
	}

	f, err := FileOpen(ctrl.fileNameBase, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(f)

	GenerateBasis29(ctrl, f)

}

//GenerateBasisInteractive : Interactive func to take user input and generate 29Basis
//files; generally one only needs the "no filter" version, but one is allowed
//to choose other varieties; recommended to always use default ranges unless
//you are completely comfortable with juusprimes and filtering
func GenerateBasisInteractive() {

	fmt.Println(Basis29Msg + "\n")
	if approved := GetUserConfirmation("Do you want to continue?", "y"); !approved {
		return
	}

	var (
		input       = ""
		wasCanceled = false
	)

	ctrl := NewGenPrimesStruct()

	fmt.Println("Enter the filter type you want:")
	fmt.Println("1 - No filter (with default TNum range, same as 'basis' in Generation menu)")
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

	if input, wasCanceled = GetUserInputInteger("Enter from TNum (recommend default!):", "28", "x"); wasCanceled {
		return
	}
	fmt.Sscan(input, ctrl.From)

	if input, wasCanceled = GetUserInputInteger("Enter to TNum (recommend default!):", "215656468", "x"); wasCanceled {
		return
	}
	fmt.Sscan(input, ctrl.To)

	ctrl.fileNameBase = filepath.Join(Basis29Path, fmt.Sprintf("%s_%v_%v_%s%s", filePrefix29, ctrl.From, ctrl.To, GetFilterAbbrev(ctrl.FilterType), fileExtRaw29))

	if FileExists(ctrl.fileNameBase) {
		if approved := GetUserConfirmation(fmt.Sprintf("ATTENTION, File: \n%v\nalready exists. Do you want to overwrite it?", ctrl.fileNameBase), "y"); !approved {
			return
		}
	}

	fmt.Printf("\n\n-----\nWill generate Basis 29 file starting from TNumber %v and ending at TNumber %v, filtered by: %s\n", ctrl.From, ctrl.To, GetFilterDesc(ctrl.FilterType))
	if approved := GetUserConfirmation("file will be written to: \n"+
		ctrl.fileNameBase+"\n\nis this what you want?",
		"y"); !approved {
		return
	}

	f, err := FileOpen(ctrl.fileNameBase, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(f)

	GenerateBasis29(ctrl, f)

}

//GenerateBasis29 : The engine func that generates 29Basis files; it can
//be called with a GenPrimesStruct already filled in if you are comfortable
//with that; otherwise recommended to use GenerateBasisInteractive()
func GenerateBasis29(ctrl *GenPrimesStruct, f *os.File) {

	if ctrl.From == nil || ctrl.To == nil {
		fmt.Println("Basis29: invalid GenPrimeStruct, from and/or to TNum not assigned")
		return
	}

	if ctrl.From.Cmp(basisBegin) < 0 {
		fmt.Println(fmt.Sprintf("From value %v is less than 28. It must be 28 or greater.", ctrl.From))
		return
	}

	if ctrl.To.Cmp(ctrl.From) < 1 {
		fmt.Println(fmt.Sprintf("To value %v is <= from value %v", ctrl.To, ctrl.From))
		return
	}

	adderFilter, finalFilter := GetFilter29(ctrl.FilterType)

	type controlStruct struct {
		idx  int
		size int
		p    *PrimeLTE29
	}

	var (
		controlStructs       []*controlStruct
		addResult, effectPtr int
	)
	addResult = -1

	//a function to check sanity: indexes, in a complete 29basis, should return to initial values minus one
	ListIndexes := func(isAtStart bool) {
		adjust := 0
		switch isAtStart {
		case true:
			fmt.Println("Starting:")
			adjust = 1
		case false:
			fmt.Println("Ending:")
		}
		for _, cs := range controlStructs {
			fmt.Println(fmt.Sprintf("%v: %v", cs.size, cs.idx+adjust))
		}
		fmt.Println("")
	}

	getControlStruct := func(prime *PrimeLTE29) *controlStruct {
		result := &controlStruct{
			//in order to bail out (continue) if effect is "X" (ie, possibility killed/destroyed) then
			//Need to increment at TOP of the loop, hence the " - 1";
			idx:  prime.Prime.GetNaturalProgressionIdx(ctrl.From) - 1,
			size: int(prime.Prime.value.Int64()),
			p:    prime,
		}
		controlStructs = append(controlStructs, result)
		return result
	}

	increment := func() {
		for _, cs := range controlStructs {
			cs.idx = (cs.idx + 1) % cs.size //increment but keep it in mod range of 0..6 for 7, 0..10 for 11, etc.
		}
	}

	adder := func(cs *controlStruct) func() bool {
		return func() bool {
			effectPtr = cs.p.CrossNumEffect[cs.p.Prime.naturalProgression[cs.idx]]
			addResult = AddSymbols(&addResult, &effectPtr)
			return FilterMap[addResult]&adderFilter == 0
		}
	}

	i7 := adder(getControlStruct(NewPrimeLTE29(big.NewInt(7))))
	i11 := adder(getControlStruct(NewPrimeLTE29(big.NewInt(11))))
	i13 := adder(getControlStruct(NewPrimeLTE29(big.NewInt(13))))
	i17 := adder(getControlStruct(NewPrimeLTE29(big.NewInt(17))))
	i19 := adder(getControlStruct(NewPrimeLTE29(big.NewInt(19))))
	i23 := adder(getControlStruct(NewPrimeLTE29(big.NewInt(23))))
	i29 := adder(getControlStruct(NewPrimeLTE29(big.NewInt(29))))

	ClearSymbolCounts()
	//when passing through the entire 29 basis all the counters
	//should come around to the same position minus 1, listIndexes allows us to
	//check this
	ListIndexes(true)

	displayProgress := DisplayProgress(ctrl.From.Int64(), ctrl.To.Int64(), 40)
	startTime := DisplayProgressBookend(fmt.Sprintf("Generating 29 basis from TNumber %v to %v, filtered by %s", ctrl.From, ctrl.To, GetFilterDesc(ctrl.FilterType)), true)

	for i := ctrl.From.Int64(); i <= ctrl.To.Int64(); i++ {

		displayProgress()
		increment()

		//for the generation of first 29 primes we always are comparing to a sextuplet,
		//the 2,3,5 pattern, no prefilter needed
		addResult = CSextuplet

		if !i7() {
			continue
		}
		if !i11() {
			continue
		}
		if !i13() {
			continue
		}
		if !i17() {
			continue
		}
		if !i19() {
			continue
		}
		if !i23() {
			continue
		}
		if !i29() {
			continue
		}

		if FilterMap[addResult]&finalFilter == 0 {
			continue
		}

		//accumulate statistics of what we've found
		SymbolCount[addResult]++

		fmt.Fprintln(f, i)
		fmt.Fprintln(f, addResult)
		//29Basis files have no pretty print, no real need and kinda useless
	}

	fmt.Println("Duration:", DisplayProgressBookend("Done", false).Sub(startTime))
	fmt.Println("")
	ListIndexes(false)

	infoF, err := FileOpen(ChangeFileExt(ctrl.fileNameBase, fileExtInfo), false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(infoF)

	ShowSymbolCounts(ctrl.From, ctrl.To, ctrl.FilterType, os.Stdout)

	ShowSymbolCounts(ctrl.From, ctrl.To, ctrl.FilterType, infoF)
	fmt.Fprintln(infoF, cInfoDataFormat)
	fmt.Fprintln(infoF, cInfoSymbols)
	ShowSymbolFileDesignations(infoF)
	fmt.Fprintln(infoF, cInfoBasis)
	fmt.Fprintln(infoF, cInfoBasisExtra)
	fmt.Println("")
	fmt.Println(fmt.Sprintf("The files:\n%s\n%s\nhave been generated.", f.Name(), infoF.Name()))

}

//Basis29RangeTNum : func that returns the start/end values, in TNumbers,
//of a 29Basis range; not used, primarily here so that an interested person
//can see how they were determined; The pattern then begins again and
//repeats to infinity
func Basis29RangeTNum() (startTNum, endTNum int) {
	startTNum = 28 //constant, 28 is the starttemplate of Prime 29
	endTNum = (7 * 11 * 13 * 17 * 19 * 23 * 29) + startTNum - 1
	return
}

//ShowRawDetails : display the basic data of the prime structure
//and the "raw" crossing effects. Natural progressions are derived from
//this data.
func (prime *PrimeLTE29) ShowRawDetails() {
	fmt.Println("raw crossing details:\n=====")
	fmt.Println("P  =", prime.Prime.value)
	fmt.Println("T# =", prime.Prime.startTemplateNum)
	fmt.Println("T expanded =", TNumToInt(prime.Prime.startTemplateNum))
	fmt.Println("These are raw results with crossing #'s in mod form.")
	fmt.Println("The natural progression can be derived from this data.")
	fmt.Println("")
	cNum := big.NewInt(0)
	for i := 0; i < int(prime.Prime.value.Int64()); i++ {
		cNum.Mod(big.NewInt(int64(i)), prime.Prime.value)
		fmt.Println(fmt.Sprintf("Cross#: %v   Symbol: %s", cNum, GetSymbolString(prime.GetEffect(cNum), false)))
	}
	fmt.Println("---------------------- END")
	fmt.Println("")
}

//CrossingAtTNumCountingNumber : Given the Prime and TemplateNum, return
//the crossing number in Counting number format (starting at 1); only for
//humans preferring counting numbers rather than mod results, do not use
//in calculations!
func CrossingAtTNumCountingNumber(prime, tNum *big.Int) *big.Int {
	iCalcFuncResult = CrossingAtTNumMod(prime, tNum)
	return big.NewInt(0).Add(iCalcFuncResult, big1)
}

//CrossingAtTNumExpandedCountingNumber : Given the Prime and TemplateNum expanded to int, return
//the crossing number in Counting number format (starting at 1); only for
//humans preferring counting numbers rather than mod results, do not use
//in calculations!
func CrossingAtTNumExpandedCountingNumber(prime, tNumExpanded *big.Int) *big.Int {
	iCalcFuncResult = CrossingAtTNumExpandedMod(prime, tNumExpanded)
	return big.NewInt(0).Add(iCalcFuncResult, big1)
}

//CrossingAtTNumExpandedMod : Given the Prime and a TemplateNum expanded to number line integer,
//return the crossing number in mod format (starting at 0); Key function for this package
func CrossingAtTNumExpandedMod(prime, tNumExpanded *big.Int) *big.Int {
	//	return (prime - (tNumExpanded % prime)) % prime
	iCalcA.Mod(tNumExpanded, prime)
	iCalcA.Sub(prime, iCalcA)
	return big.NewInt(0).Mod(iCalcA, prime)
}

//CrossingAtTNumMod : Given the Prime and TemplateNum return
//the crossing number; Key function for this package; "prime" is
//not necessarily a prime #, but a potentially prime #
func CrossingAtTNumMod(prime, tNum *big.Int) *big.Int {
	//	return (prime - (TNumToInt(tNum) % prime)) % prime
	iCalcFuncResult = TNumToInt(tNum)
	iCalcFuncResult.Mod(iCalcFuncResult, prime)
	iCalcFuncResult.Sub(prime, iCalcFuncResult)
	return big.NewInt(0).Mod(iCalcFuncResult, prime)
}

/*
alternate input from user
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")

	name, _ := reader.ReadString('\n')
	fmt.Printf("Hello %s\n", name)
*/
