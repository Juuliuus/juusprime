package juusprime

//juusprime is an engine to find prime Sextuplets and/or Quintuplets and/or Quadruplets
//GNU GPL3
//author: Julius Schön / R. Spicer, 2021

import (
	"crypto/md5"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

func init() {
	//all of these are numbers that are used over and over again, my thoughts are that
	//it is better to re-use rather than a storm of big.NewInt calls.
	TemplateLength = big.NewInt(30)
	big0 = big.NewInt(0)
	big1 = big.NewInt(1)
	big2 = big.NewInt(2)
	big5 = big.NewInt(5)
	big8 = big.NewInt(8)
	big19 = big.NewInt(19)

	big12 = big.NewInt(12)
	big16 = big.NewInt(16)
	big18 = big.NewInt(18)
	big22 = big.NewInt(22)
	big24 = big.NewInt(24)
	big28 = big.NewInt(28)
	big29 = big.NewInt(29)

	basisBegin = big.NewInt(28)
	basisLen = big.NewInt(215656441)
	basisEnd = big.NewInt(215656468)

	iCalcFuncResult, iCalcA, iCalcB, iCalcC, iCalcD = big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	//all big.NewFloats default to precision 53, precision is adjusted on the fly where it is needed
	fCalcFuncResult, fCalcA, fCalcB, fCalcC = big.NewFloat(0), big.NewFloat(0), big.NewFloat(0), big.NewFloat(0)
	bigFTemplateLength = big.NewFloat(0).SetInt(TemplateLength)

	bigF1 = big.NewFloat(0).SetInt64(1)
	bigF19 = big.NewFloat(0).SetInt64(19)
	currPrecision = 53
}

const (
	filePrefix29  = "29basis"
	fileExtRaw23  = ".rawdata23"
	fileExtRaw29  = ".rawdata29"
	fileExtRaw    = ".rawdata"
	fileExtPretty = ".prettydata"
	fileExtInfo   = ".info"
)

var (
	iCalcFuncResult, iCalcA, iCalcB, iCalcC, iCalcD *big.Int
	fCalcFuncResult, fCalcA, fCalcB, fCalcC         *big.Float
	bigFTemplateLength, bigF1, bigF19               *big.Float
	currPrecision                                   uint
	big0, big1, big2, big5, big8, big19             *big.Int
	big12, big16, big18, big22, big24, big28, big29 *big.Int
	//TemplateLength : This a constant used a lot in big.Int calcs, it is the length of a Template (30)
	TemplateLength                 *big.Int
	basisBegin, basisLen, basisEnd *big.Int
)

const (
	omBasis = iota
	omTNum
	omNatNum
)

//GenPrimesStruct : A structure that can be filled in to control
//prime generation (both LTE 29's and GTE 31's, though the usage is
//slightly different), for GTE 31's it will call its Prepare() method
//which adjusts parameters based on OpMode and returns a proper filename
type GenPrimesStruct struct {
	BasisNum            *big.Int
	DefaultPath         string
	FilterType          int
	fileNameBase        string
	From                *big.Int
	FullPathto29RawFile string //only used for GTE 31s
	To                  *big.Int
	OpMode              int
	//omBasis = iota
	//omTNum
	//omNatNum
}

//Prepare : Method for GenPrimesStruct which takes the OpMode
//and prepares proper settings, it returns a proper filename and fullpath,
//Always used for GTE 31's, LTE 29's are special and handled slightly differently,
//Generally no need to call this, it is done automatically when generating Tuplets
func (ctrl *GenPrimesStruct) Prepare() (fileName string) {
	switch ctrl.OpMode {
	case omBasis:
		BasisToTNumRange(ctrl.BasisNum, ctrl.From, ctrl.To)
	case omTNum:
		TNumToBasisNum(ctrl.From, ctrl.BasisNum)
	case omNatNum:
		ctrl.From.Set(IntToTNum(ctrl.From))
		ctrl.To.Set(IntToTNum(ctrl.To))
		TNumToBasisNum(ctrl.From, ctrl.BasisNum)
	}

	//one can manipulate settings at this point for testing

	//section below tests moving between basis #s
	//ctrl.BasisNum = 0
	//ctrl.OpMode = omTNum
	//ctrl.To.SetInt64(215648043)
	//ctrl.To.SetInt64(215671939)

	//ctrl.To.SetInt64(28) //original testing
	//ctrl.To.SetInt64(65145) //original testing which gives a good number of Sext's
	//ctrl.To.SetInt64(215656468)  //end of basis 0

	fileName = filepath.Join(ctrl.DefaultPath,
		fmt.Sprintf(
			ctrl.fileNameBase,
			AdjustTNumsForFilename(ctrl.BasisNum),
			AdjustTNumsForFilename(ctrl.From),
			AdjustTNumsForFilename(ctrl.To),
			GetFilterAbbrev(ctrl.FilterType),
			fileExtRaw))
	return
}

//NewGenPrimesStruct : returns a pointer to GenPrimesStruct which is
//used to control prime GTE 31's generation
func NewGenPrimesStruct() *GenPrimesStruct {
	return &GenPrimesStruct{
		BasisNum:            big.NewInt(0),
		DefaultPath:         os.Getenv("HOME"),
		FilterType:          ftAll,
		fileNameBase:        "juusprimes_basis-%v_%v_%v_%v%s",
		FullPathto29RawFile: "",
		From:                big.NewInt(0),
		To:                  big.NewInt(0),
		OpMode:              omBasis,
	}
}

//AdjustTNumsForFilename : used to keep filename lengths under control,
//if the incoming TNumber or basis (n) is bigger than 9999999999999999999999999999999
//it returns an md5 hash, otherwise returns the string representation of the
//incoming number
func AdjustTNumsForFilename(n *big.Int) string {
	//same(ish) length as an md5 sum, so limits to a max(ish) filename length
	//         27ccb0eea8a706c4c34a16891f84e7b
	maxstr := "9999999999999999999999999999999"
	max := big.NewInt(0)
	fmt.Sscan(maxstr, max)
	returnStr := fmt.Sprintf("%v", n)
	if n.Cmp(max) > 0 {
		return fmt.Sprintf("%x", md5.Sum([]byte(returnStr)))
	}
	return returnStr
}

//GetEffectiveTNum : "effective" TNumber is the TNumber where the square of
//the potential prime (p) resides, ie. where it begins its effect on the following
//TNumbers. It is the Start Template Number for the calculated potential GTE 31 primes
func GetEffectiveTNum(n *big.Int, p *PrimeGTE31, returnHereTNum *big.Int) {
	//  result := StartTNum + ( 2 * Prime.Value * n ) + ( cTemplateLength * sqr( n ) )
	//note: this is seemingly more complicated than it needs to be because
	//the effective tnumber is nothing more that the potential prime's value squared
	//converted to a TNumber. I use this instead because it uses the generic "n" based
	//math equations which tie the system together, the other would hide this information.
	iCalcB.Mul(big2, p.Prime.value)
	iCalcB.Mul(iCalcB, n)

	iCalcC.Exp(n, big2, nil)
	iCalcC.Mul(TemplateLength, iCalcC)

	iCalcB.Add(p.Prime.startTemplateNum, iCalcB)
	returnHereTNum.Add(iCalcB, iCalcC)
}

//GetEffectiveTNumSimple : lighter alternative to GetEffectiveTNum:
// see note in GetEffectiveTNum
func GetEffectiveTNumSimple(n *big.Int, p *PrimeGTE31, returnHereTNum *big.Int) {
	//  result := IntToTNum( (p.value + 30n)^2 )
	p.MemberAtN(n, iCalcFuncResult)
	returnHereTNum.Set(IntToTNum(iCalcFuncResult.Mul(iCalcFuncResult, iCalcFuncResult)))
}

//GetNfromTNumComplicated : Given a TNum & PrimeGTE31 return the the n value,
//ie. how many potential primes must be tested, Validation that givenTNum is
//equal to or greater than p's effective start TNumber is the responsibility of the
//calling func(), to complete the calculation. This comes from the first
//derivation of the formula and is complicated in many ways (lots of roots),
//a simpler form was found, see GetNfromTNum, result is returned in param and in p.Helper.MaxN
func GetNfromTNumComplicated(givenTNum *big.Int, p *PrimeGTE31, returnedHereN *big.Int) {
	//superceded by GetNfromTNum(), see notes below

	//example: floor( (sqrt( 194091003877655 - 93 + ( 53^2 <=2809> / 30 ) ) - ( 53 / sqrt( 30 ) ))  / sqrt( 30 ))
	//givenTNum: 194091003877655, Prime.startTNum: 93, Prime.value: 53, Prime.valuesquared: 53^2
	//answer should be: 2543558
	//this calculations require big.Float's, attention to precision!

	currPrecision = uint(givenTNum.BitLen()) + 33          //I want at least 10 extra decimal digits
	fCalcA.SetPrec(currPrecision).Sqrt(bigFTemplateLength) //sqrt(30) RootTLength : Root of Template Length

	//SQRTvalue := ( sqrt( givenTNum - startTNum + PvalueSquared ) );
	fCalcB.SetPrec(currPrecision).Sub(
		big.NewFloat(0).SetPrec(currPrecision).SetInt(givenTNum),
		big.NewFloat(0).SetPrec(currPrecision).SetInt(p.Prime.startTemplateNum))

	//PvalueSquared := ( valueSquared / cTemplateLength );
	fCalcC.SetPrec(currPrecision).Quo(
		big.NewFloat(0).SetPrec(currPrecision).SetInt(p.Prime.valueSquared),
		bigFTemplateLength)

	fCalcB.Add(fCalcB, fCalcC)
	fCalcB.Sqrt(fCalcB)

	//Pvalue := ( value / RootTLength );
	fCalcC.Quo(big.NewFloat(0).SetInt(p.Prime.value), fCalcA)
	//( SQRTvalue - Pvalue )
	fCalcB.Sub(fCalcB, fCalcC)

	//result := floor( result / RootTLength );
	fCalcB.Quo(fCalcB, fCalcA)

	r, _ := fCalcB.Int(nil)
	//fmt.Println(fCalcB, r)
	returnedHereN.Set(r)
	//return r
	p.Helper.MaxN.Set(r)

	//This has been superseded by a simpler algorithm after I re-did the algebra.
	//This is the first equation I had and is algebraically exact. It is very much
	//more complicated in that it has lots of roots and is very finicky, one has
	//to be very careful passing the correct precision.
	//I've done many tests (also in test unit) with SpeedCrunch calc and both of These
	//functions, they all agree to very high accuracy...as long as one pays attention to
	//correct precision.
}

//GetNfromIntComplicated : the compliment to GetNfromTNumComplicated where an Int
//(regular number line value) is used, rather than an already known Template Number.
//rNum is a "real" number, ie. a number line integer, not a Template Number. p is
//a PrimeGTE31
func GetNfromIntComplicated(rNum *big.Int, p *PrimeGTE31, returnedHereN *big.Int) {
	GetNfromTNumComplicated(IntToTNum(rNum), p, returnedHereN)
}

//GetNfromTNum : Given a TNum and a PrimeGTE31 the n value, ie how many potential
//primes must be tested to complete the calculation, Validation that givenTNum is
//equal to or greater than p's effective start TNumber is the responsibility of the
//calling func(), This is later derived algebraic equation
//that turns out to be far simpler and just as accurate, The first equation was quite
//complicated, see GetNfromTNumComplicated, result is returned in param and in p.Helper.MaxN
func GetNfromTNum(givenTNum *big.Int, p *PrimeGTE31, returnedHereN *big.Int) {
	//supercedes, because its simpler/faster, GetNfromTNumComplicated
	//floor( (sqrt( ( givenTNum * cTemplateLength ) + (1 or 19) ) - value) / 30 )
	//example:( sqrt( ( 194091003877655 * 30 )  + 19 ) - 53 ) / 30
	//givenTNum: 194091003877655, Prime.value: 53, Prime.valueSquaredEndsIn1 is false while squared value ends in 9
	//answer should be: 2543558
	//this calculations require big.Float's, attention to precision!

	currPrecision = uint(givenTNum.BitLen()) + 33 //I want at least 10 exact decimal digits

	//interim-result: ( givenTNum * cTemplateLength )
	fCalcC.SetPrec(currPrecision).SetInt(givenTNum)
	fCalcA.SetPrec(currPrecision).Mul(fCalcC, bigFTemplateLength)

	if p.valueSquaredEndsIn1 {
		//interim-result + 1
		fCalcA.Add(fCalcA, bigF1)
	} else {
		//interim-result + 19
		fCalcA.Add(fCalcA, bigF19)
	}
	//sqrt( interim-result )
	fCalcA.Sqrt(fCalcA)

	//floor( ( interim-result - value ) / cTemplateLength )
	fCalcC.SetPrec(currPrecision).SetInt(p.Prime.value)
	fCalcA.Sub(fCalcA, fCalcC)

	fCalcA.Quo(fCalcA, bigFTemplateLength)
	r, _ := fCalcA.Int(nil)
	returnedHereN.Set(r)
	p.Helper.MaxN.Set(r)
}

//GetNfromInt : the compliment to GetNfromTNum where an Int
//(regular number line value) is used, rather than an already known Template Number.
//rNum is a "real" number, ie. a number line integer, not a Template Number, p is
//is a potential prime of type PrimeGTE31
func GetNfromInt(rNum *big.Int, p *PrimeGTE31, returnedHereN *big.Int) {
	GetNfromTNum(IntToTNum(rNum), p, returnedHereN)
}

//GetCrossNumModDirect : An alternative to GetCrossNumMod original,
//it uses a direct "mod" operation which then adjusts back to TNumbers crossing,
//basically "unwrapping" the expansions, result is
//returned in parameter, A bit less calculation and requires no floor function
//since effectiveTNum is not required, also see GetCrossNumModSimple,
//This DOES NOT return errors if you send in too small a TNumber because this routine
//is used very often in loops, caller needs to validate the TNumbers submitted are >= the p's startTemplateNumber
func GetCrossNumModDirect(givenTNum, n *big.Int, p *PrimeGTE31, returnHereCrossNumMod *big.Int) {
	//R is p.value mod 30, ie. 1, 7, 11, 13, 17, 19, 23, or 29
	//givenTNum mod (p+30n) - ((n+1)R) - K
	//
	//K is constant for each potPrime = 0, 1, 4, 5, 9, 12, 17, 28
	//K is calculated from p's offset into its starting TNumber + C - p.value

	// get p+30n into iCalcC
	p.MemberAtN(n, iCalcC)
	//fmt.Println(iCalcC, iCalcFuncResult)

	// givenTNum mod (p+30n) into iCalcA
	iCalcA.Mod(givenTNum, iCalcC)

	// n+1 into iCalcB
	iCalcB.Add(n, big1)
	// (n+1)R
	iCalcB.Mul(iCalcB, p.Prime.mod30)

	//givenTNum mod (p+30n) - ((n+1)R)
	iCalcA.Sub(iCalcA, iCalcB)

	//subtract K
	returnHereCrossNumMod.Sub(iCalcA, p.Prime.modConst)
	//returnHereCrossNumMod.Sub(iCalcA, iCalcB)
	//if answer is negative then it is the complementary portion of pP
	if returnHereCrossNumMod.Cmp(big0) == -1 {
		returnHereCrossNumMod.Add(returnHereCrossNumMod, iCalcC)
	}

}

//GetCrossNumModSimple : Another alternative to GetCrossNumMod original, but even
//slightly faster, found by inspection rather than derivation,
//also see GetCrossNumModDirect,  This DOES NOT return errors if you send
//in too small a TNumber because this routine is used very often in loops, caller needs
//to validate the TNumbers submitted are >= the p's startTemplateNumber
func GetCrossNumModSimple(givenTNum, n *big.Int, p *PrimeGTE31, returnHereCrossNumMod *big.Int) {
	//F is p.value mod 30, ie. 1, 7, 11, 13, 17, 19, 23, or 29
	//p.sMinusp is another constant = p.startTemplateNumber - p.value
	//
	//(giveTNum - p.sMinusp - (F*n)) mod pP(n)

	// get pP(n) {p+30n} into iCalcC
	p.MemberAtN(n, iCalcC)

	//(F*n)
	iCalcB.Mul(p.Prime.mod30, n)

	iCalcA.Sub(givenTNum, iCalcB)
	iCalcA.Sub(iCalcA, p.Prime.sMinusp)

	returnHereCrossNumMod.Mod(iCalcA, iCalcC)
}

//GetCrossNumMod : returns in pointer param the offset (crossing number) for the GTE 31 Prime
//at level n at the specified given target TNumber, NOTE - After the release I realized that this
//function is very poorly named and misleading, It does NOT return a crossing number, it returns
//the relative offset of givenTNum into the prime's Natural Progression. Renaming the function
//would require a major version number change because it would break existing code, this is the
//original CrossNumMod function built using TNumber mod algorithms,
//also see the newer GetCrossNumModDirect or GetCrossNumModSimple funcs which are simpler and faster,
//This DOES NOT return errors if you send in too small a TNumber because this routine
//is used very often in loops, caller needs to validate the TNumbers submitted are >= the p's startTemplateNumber
func GetCrossNumMod(givenTNum, n *big.Int, p *PrimeGTE31, returnHereCrossNumMod *big.Int) {
	//givenTNum - (floor(givenTNum-effectiveTNum / p+30n) (p+30n) ) + effectiveTNum
	GetEffectiveTNum(n, p, iCalcFuncResult)

	//PotentialPrime := p + ( cTemplateLength * n ) [i.e. p+30n]
	iCalcC.Mul(TemplateLength, n)
	iCalcC.Add(p.Prime.value, iCalcC)

	//DivResult := ( givenTNum - effectiveTNum ) div PotentialPrime;
	iCalcA.Sub(givenTNum, iCalcFuncResult)
	iCalcA.Div(iCalcA, iCalcC)
	//ToSubtract := ( DivResult * PotentialPrime ) + effectiveTNum;
	iCalcA.Mul(iCalcA, iCalcC)
	iCalcA.Add(iCalcA, iCalcFuncResult)
	//result := givenTNum - ToSubtract;
	returnHereCrossNumMod.Sub(givenTNum, iCalcA)
}

//TNumToInt : given any Template Number return the number
//line integer where the template begins. Key function for this package.
func TNumToInt(tNum *big.Int) *big.Int {
	//	return (tNum * 30) - 5
	iCalcA.Mul(tNum, TemplateLength)
	return big.NewInt(0).Sub(iCalcA, big5)
}

//IntToTNum : given any integer return the TemplateNum where
//it will be found. Key function for this package.
func IntToTNum(rNum *big.Int) *big.Int {
	//	return (rNum + 5) / 30
	iCalcA.Add(rNum, big5)
	return big.NewInt(0).Div(iCalcA, TemplateLength)
}

//BasisToTNumRange : basis is 0 based, results are set directly to passed
//in param pointers, results are the beginning and ending TNumbers for the
//specified basis
func BasisToTNumRange(basis, returnHereBegin, returnHereEnd *big.Int) {
	//tNumBeg = basisBeg + (basisLen * basisNum)
	//tNumEnd = basisEnd + (basisLen * basisNum)
	iCalcA.Mul(basis, basisLen)
	returnHereBegin.Add(basisBegin, iCalcA)
	returnHereEnd.Add(basisEnd, iCalcA)
}

//TNumToBasisNum : Given a TNumber, return the basis Number
//it will be found in to the passed pointer param
func TNumToBasisNum(tNum, returnHereBasisNum *big.Int) error {
	//basisNum = (tNum - basisBeg) div basisLen
	iCalcA.Sub(tNum, basisBegin)
	returnHereBasisNum.Set(iCalcA.Div(iCalcA, basisLen))
	if tNum.Cmp(basisBegin) < 0 {
		return fmt.Errorf("TNumToBasisNum: TNumber '%v' is less than beg. basis # '%v', meaningless result.", tNum, basisBegin)
	}
	return nil
}

//IntToBasisNum : Given an integer, return the basis Number
//it will be found in to the passed in pointer param
func IntToBasisNum(anInt, returnHereBasisNum *big.Int) error {
	//basisNum = (tNum - basisBegin) div basisLen
	TNumToBasisNum(IntToTNum(anInt), returnHereBasisNum)
	if anInt.Cmp(big.NewInt(835)) < 0 {
		return fmt.Errorf("IntToBasisNum: int '%v' is less than beg. basis int '%d', meaningless result.", anInt, 835)
	}
	return nil
}

//getBigInt : helper func to get fresh initialized *big.Int,
//particularly useful when passing big.Int's as params & it
//is =NOT= limited to int64.
func getBigInt(initTo *big.Int) *big.Int {
	return big.NewInt(0).Add(big0, initTo)
}

//ShowMe : testing and internal use, good to use
//when you want to see unexported variable results without a lot
//of bother
func ShowMe(choice int, p *PrimeGTE31, b *big.Int) error {

	switch choice {
	case 1:
		if p == nil {
			return fmt.Errorf("ShowMe: Bad prime sent in")
		}
		if b.Cmp(big1) == -1 {
			return fmt.Errorf("fixed N (%v) must GTE 1", b)
		}
		p.DisplayFullCritLengths(b)
	default:
		fmt.Println("ShowMe: bad choice sent in.")

	}
	return nil
}

/*
original-ish notes from pascal version re: N and very large TNumber


SpeedCrunch testing of very large TNumber and resulting N....

Research into N
===============

A large TNumber and the checked sextuplet it belongs to:

TNumber: 194091003877655  ┣━┫
---------
5822730116329657
5822730116329661
5822730116329663
5822730116329667
5822730116329669
5822730116329673

using potential primes of type 53 (a 9'er, ie, squares end in 9) is the check into
N calculations (also checked with potential primes of 31 (a 1'er). they were also exact):

easy N method
( sqrt( ( 194091003877655 * 30 )  + 19 ) - 53 ) / 30
2543558.75536837777049820408

This N generates:
Ttar = ( 30 * 2543558^2 ) + ( 62 * 2543558 ) + 93
194090776681609

P = 53 + ( 30 * 2543558 ) = 76306793

76306793^2 = 5822726657944849
T = ( 5822726657944849 + 5 ) / 30 = 194090888598161.8

so this is the square T:
194090888598161
easy N
( sqrt( ( 194090888598161 * 30 )  + 19 ) - 53 ) / 30
2543558


Hard N:

Using T from above:
( sqrt( 194091003877655 - 93 + ( 53^2 / 30 ) ) - ( 53 / sqrt( 30 ) ) ) / sqrt( 30 )
Hard N : 2543558.7553683777 7049820408
golangH: 2543558.7553683777 7049820407848511219158913194961228194112612748908185548954047
Easy N : 2543558.7553683777 7049820408
golangE: 2543558.7553683777 7049820407848511219158913194961228194112612748908185548954047

Using T from The squared P:
( sqrt(  194090888598161 - 93 + ( 53^2 / 30 ) ) - ( 53 / sqrt( 30 ) ) ) / sqrt( 30 )
Hard N : 2543558
Easy N : 2543558
So both are correct and exact. Below easy and hard are tested with the
TNumber 1 over and 1 under. Again exactly matching. The two equations are equivalent
and produce exactly the same numbers...given the precision of the calculator is ok.

However, since it could be that very far out there is some rounding error in
decimal precision that "could" cause it to round down 1 too many (all these functions
are "floor"ed), then to be certain that N is high enough I add 1 to N and then check
that the TNumber being tested is in range (ie., less than the T(effective) of that
potential prime. Yes, belt AND suspenders...just in case. But with precise enough decimal,
the adding of 1 to N is not needed. <I believe we now have precise enough precision Dec. 2020>)

Here are "easy" N's 1 T over and 1 T under square P TNumber...
1 over:
( sqrt( ( 194090888598162 * 30 )  + 19 ) - 53 ) / 30
2543558.00000000655249657786

1 under:
( sqrt( ( 194090888598160 * 30 )  + 19 ) - 53 ) / 30
2543557.99999999344750342214


And here are the Hard N's , 1 over and under:
1 over:
( sqrt( 194090888598162 - 93 + ( 53^2 / 30 ) ) - ( 53 / sqrt( 30 ) ) ) / sqrt( 30 )
2543558.00000000655249657786 hard
2543558.00000000655249657786 easy

1 under:
( sqrt(  194090888598160 - 93 + ( 53^2 / 30 ) ) - ( 53 / sqrt( 30 ) ) ) / sqrt( 30 )
2543557.99999999344750342214 hard
2543557.99999999344750342214 easy

These calculations were done in SpeedCrunch on Debian 9, so the equations
above go in that calculator and are not pascal expressions.

*/
