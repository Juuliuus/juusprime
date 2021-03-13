package juusprime

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

//juusprime is an engine to find prime Sextuplets and/or Quintuplets and/or Quadruplets
//GNU GPL3
//author: Julius Schön / R. Spicer, 2021

//this file is in development, not much here right now

//GetCritLengthPositiveWF : WF=witthin family, given N (a chosen, fixed n level),
//and diff (the difference between N and your desired n), calculate the number of
//Templates between them in positive direction, ie. toward infinity; prime is the
//potPrime base value (31, 37,...), result is returned in last parameter
func GetCritLengthPositiveWF(p *PrimeGTE31, N, diff, returnHereLen *big.Int) error {
	// d( 2p + dc + 2cN ) d, diff; p, prime value; c=30; N a fixed chosen "n"
	returnHereLen.SetInt64(0)
	if diff.Cmp(big1) < 0 {
		return errors.New("GetCritLengthPositive: difference (diff) parameter must be greater than or equal to 1")
	}
	iCalcA.Mul(p.Prime.value, big2)

	iCalcB.Mul(diff, TemplateLength)

	iCalcC.Mul(N, TemplateLength)
	iCalcC.Mul(iCalcC, big2)

	iCalcA.Add(iCalcA, iCalcB)
	iCalcA.Add(iCalcA, iCalcC)

	returnHereLen.Mul(diff, iCalcA)
	return nil
}

//GetCritLengthNegativeWF : WF=witthin family, given N (a chosen, fixed n level),
//and diff (the difference between N and your desired n), calculate the number of
//Templates between them in negative direction, ie. toward zero; prime is the
//potPrime base value (31, 37,...), result is returned in last parameter
func GetCritLengthNegativeWF(p *PrimeGTE31, N, diff, returnHereLen *big.Int) error {
	// d( 2p - dc + 2cN ) d, diff; p, prime value; c= 30; N a fixed chosen "n"
	returnHereLen.SetInt64(0)
	if diff.Cmp(big1) < 0 {
		return errors.New("GetCritLengthNegative: difference parameter (diff) must be greater than or equal to 1")
	}
	iCalcD.Sub(N, diff)
	if iCalcD.Cmp(big0) == -1 {
		return fmt.Errorf("GetCritLengthNegative: N (%v) minus the diff (%v) goes below zero", N, diff)
	}
	iCalcA.Mul(p.Prime.value, big2)

	iCalcB.Mul(diff, TemplateLength)

	iCalcC.Mul(N, TemplateLength)
	iCalcC.Mul(iCalcC, big2)

	iCalcA.Sub(iCalcA, iCalcB)
	iCalcA.Add(iCalcA, iCalcC)

	returnHereLen.Mul(diff, iCalcA)
	return nil
}

//GetCritLengthWF : within family; return the total number of Templates for
//potPrime p between fromN and toN, toN can be less than fromN.
func GetCritLengthWF(p *PrimeGTE31, fromN, toN, returnHereLen *big.Int) error {
	returnHereLen.SetInt64(0)
	if fromN.Cmp(toN) == 0 {
		return fmt.Errorf("GetCritLengthWF: from n (%v) and to n (%v) can not be equal", fromN, toN)
	}
	if fromN.Cmp(big0) == -1 || toN.Cmp(big0) == -1 {
		return fmt.Errorf("GetCritLengthWF: Given n values '%v' and '%v', n's must be 0 or greater", fromN, toN)
	}
	diff := big.NewInt(0).Sub(fromN, toN)
	diff.Abs(diff)
	if fromN.Cmp(toN) == -1 {
		return GetCritLengthPositiveWF(p, fromN, diff, returnHereLen)
	}
	return GetCritLengthNegativeWF(p, fromN, diff, returnHereLen)
}

//GetCritLengthByDiffWF : within family; return the total number of Templates for
//potPrime p between fromN and toN, toN can be less than fromN.
func GetCritLengthByDiffWF(p *PrimeGTE31, N, diff, returnHereLen *big.Int) error {
	returnHereLen.SetInt64(0)
	if N.Cmp(big0) == -1 {
		return fmt.Errorf("GetCritLengthByDiffWF: desired N (%v) must be 0 or greater", N)
	}
	if diff.Cmp(big0) == 0 {
		return fmt.Errorf("GetCritLengthByDiffWF: diff (%v) can be positive or negative, but not 0", diff)
	}
	iCalcA.Add(N, diff)
	if iCalcA.Cmp(big0) == -1 {
		return fmt.Errorf("GetCritLengthByDiffWF: your diff (%v) combined with N (%v) will go below 0", diff, N)
	}
	localDiff := big.NewInt(0).Set(diff)
	localDiff.Abs(localDiff)
	if diff.Cmp(big0) == 1 {
		return GetCritLengthPositiveWF(p, N, localDiff, returnHereLen)
	}
	return GetCritLengthNegativeWF(p, N, localDiff, returnHereLen)
}

//displayFullCritLengths : internal use, testing, checking
//Will show first from 0 to fixedN-1, then from fixedN+1 to 0
//the number of Templates between each diff of fixedN +/- n
//and will also display the appropriate offsets of each n level to the fixedN
//border; careful, large fixedN means lots and lots of output
func displayFullCritLengths(pP *PrimeGTE31, fixedN *big.Int) {
	iter := big.NewInt(0).Set(fixedN)

	N := big.NewInt(0)
	nctrl := big.NewInt(0)
	diff := big.NewInt(1)
	res := big.NewInt(0)

	mod := big.NewInt(0)
	pp := big.NewInt(0)

	iCalcA.Add(N, iter)
	iCalcA.Add(iCalcA, big1)

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 1)

	pP.MemberAtN(iter, pp)
	fmt.Println(fmt.Sprintf("distances in TNumbers between n's for potPrime %v at N=%v (%v)",
		pP.Prime.value, iter, pp))

	pP.MemberAtN(iCalcA, pp)
	fmt.Println(fmt.Sprintf("from N=%v (%v) to N+%v (%v) [inclusive]", N, pP.Prime.value, iCalcA, pp))

	fmt.Println(fmt.Sprintf("Offset (o) is for potprime value %v", pP.Prime.value))

	for nctrl.Cmp(iter) < 1 {
		nctrl.Add(nctrl, big1)
		diff.Sub(N, nctrl)
		diff.Abs(diff)
		GetCritLengthPositiveWF(pP, N, diff, res)

		pP.MemberAtN(nctrl, pp)
		mod.Mod(res, pP.Prime.Value())
		mod.Sub(mod, big1)

		fmt.Fprintf(w, "N=%v (%v) to n=%v (%v) length =\t %v\t   o=%v\n", N, pP.Prime.Value(), nctrl, pp, res, mod)
	}
	w.Flush()

	fmt.Println("-------------")
	pP.MemberAtN(iter, pp)
	fmt.Println(fmt.Sprintf("from N=%v (%v) to N-%v (%v)", iter, pp, iter, pP.Prime.Value()))

	fmt.Println("Offset (o) is for current potprime at n at the border of N")

	N.Set(iter)
	nctrl.Set(N)

	pP.MemberAtN(N, pp)
	frompp := big.NewInt(0).Set(pp)

	cmp := big.NewInt(0).Sub(N, iter)
	for nctrl.Cmp(cmp) > 0 {
		nctrl.Sub(nctrl, big1)
		diff.Sub(N, nctrl)
		diff.Abs(diff)
		err := GetCritLengthNegativeWF(pP, N, diff, res)
		if err != nil {
			fmt.Println(err)
			break
		}
		pP.MemberAtN(nctrl, pp)
		mod.Mod(res, pp)
		mod.Sub(mod, big1)

		fmt.Fprintf(w, "n=%v (%v) to n=%v (%v) length =\t %v\t   o (%v)=%v\n",
			N, frompp, nctrl, pp, res, pp, mod)
	}
	w.Flush()
}

//GetPrimeGTE31Slice : helper for analysis routines, gives
//back a slice of all Primes GTE 31
func GetPrimeGTE31Slice() []*PrimeGTE31 {
	return []*PrimeGTE31{
		NewPrimeGTE31(big.NewInt(31)),
		NewPrimeGTE31(big.NewInt(37)),
		NewPrimeGTE31(big.NewInt(41)),
		NewPrimeGTE31(big.NewInt(43)),
		NewPrimeGTE31(big.NewInt(47)),
		NewPrimeGTE31(big.NewInt(49)),
		NewPrimeGTE31(big.NewInt(53)),
		NewPrimeGTE31(big.NewInt(59))}
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

	if input, wasCanceled = GetUserInput("Show broken/effected only?", "false", "x"); wasCanceled {
		return
	}
	switch strings.ToUpper(input) {
	case "T", "TRUE", "Y", "YES":
		broken = true
	default:
		broken = false
	}

	if input, wasCanceled = GetUserInput("Show prime pP's only?", "true", "x"); wasCanceled {
		return
	}
	switch strings.ToUpper(input) {
	case "T", "TRUE", "Y", "YES":
		primesOnly = true
	default:
		primesOnly = false
	}

	if input, wasCanceled = GetUserInput("Pause for review  after each potPrime?", "true", "x"); wasCanceled {
		return
	}
	switch strings.ToUpper(input) {
	case "T", "TRUE", "Y", "YES":
		doPause = true
	default:
		doPause = false
	}

	if input, wasCanceled = GetUserInput("Enter TNumbers comma separated:", "535, 647", "x"); wasCanceled {
		return
	}

	sl := strings.Split(input, ",")
	fmt.Println("analysing ", sl)

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 1)
	wFmt := "%s\t  %s\t %v \t %s\n"

	effectStr := ""
	wStr := ""
	primeStr := ""
	printIt := false
	wasBroken := true
	big32 := big.NewInt(32)

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
				//The two funcs below are equivalent, I'm prefering Direct because it
				//doesn't need an effectiveTNumber to work and saves a few processor cycles.
				//GetCrossNumMod(giveTNum, n, primes[i], offset)
				GetCrossNumModDirect(giveTNum, n, primes[i], offset)

				pass = "ok"
				printIt = !broken
				effectStr = "-"
				if primes[i].GetResultAtCrossNum(&addResult, offset, n) {
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
					fmt.Fprintf(w, fmt.Sprintf(wFmt, wStr, effectStr, offset, pass))
				}
				n.Add(n, big1)
			}

			if broken {
				if !wasBroken {
					fmt.Fprintf(w, fmt.Sprintf(wFmt, " ok", "", "", ""))
				} else {
					//fmt.Fprintf(w, fmt.Sprintf(wFmt, " ==>  EFFECTED  <==", "", "", ""))
					fmt.Fprintf(w, fmt.Sprintf(wFmt, "", "", "", "                           ==>  EFFECTED  <=="))
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

//CheckTwinSextuplet : testing, check rawdata files for twin
//twin Sextuplets.
func CheckTwinSextuplet(filename string) {

	var (
		inResult int
	)

	if !FileExists(filename) {
		fmt.Println(fmt.Sprintf("Path '%s' to 29bais rawdata is invalid, quitting.", filename))
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
	// TODO: fix this to check over bases!!
	from, to := big.NewInt(28), big.NewInt(215656468)
	//from, to := big.NewInt(215656469), big.NewInt(431312909)
	//from, to := big.NewInt(28), big.NewInt(7731)
	big7 := big.NewInt(7)

	displayProgress := DisplayProgressBig(from, to, 100)
	displayProgressLastPos := big.NewInt(0).Set(from)

	//if the file is empty and no tTarget is scanned this will break the loop immediately.
	tTarget := big.NewInt(0).Add(to, big1)

	tCompare := big.NewInt(0)
	count := 0

	//main loop
	for {
		if !doScan(tTarget, &inResult) {
			fmt.Println("doscan broke it.")
			fmt.Println(fmt.Sprintf("%v %v", tTarget, tCompare))
			break
		}

		if tTarget.Cmp(from) == -1 {
			continue
		}

		if tTarget.Cmp(to) == 1 {
			fmt.Println("tTarget.Cmp  broke it.")
			fmt.Println(fmt.Sprintf("%v %v", tTarget, tCompare))
			break
		}

		if tCompare.Cmp(big0) == 0 {
			//need to have gotten the second number
			tCompare.Set(tTarget)
			continue
		}

		iCalcA.Sub(tTarget, tCompare)
		iCalcA.Div(iCalcA, big7)

		//		fmt.Println(fmt.Sprintf("===> Compare gives: %v", iCalcA.Cmp(big5)))

		//if its <= 5 print em.
		//if iCalcA.Cmp(big5) == 0 || iCalcA.Cmp(big5) == -1
		if iCalcA.Cmp(big1) == 0 {
			count++
			fmt.Println(tTarget)
			fmt.Println(tCompare)
			fmt.Println("")
			/*
				iCalcB.Sub(tTarget, tCompare)
				iCalcB.Div(iCalcB, big7)
				fmt.Println(fmt.Sprintf("divided gives: %v", iCalcB))
				fmt.Println(fmt.Sprintf("Found %v: %v %v", count, tTarget, tCompare))
			*/
		}
		tCompare.Set(tTarget)

		displayProgress(tTarget, displayProgressLastPos)

	} //for{(ever)}

	fmt.Println("finished:", count)
}
