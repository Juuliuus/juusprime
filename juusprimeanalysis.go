package juusprime

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
)

//juusprime is an engine to find prime Sextuplets and/or Quintuplets and/or Quadruplets
//GNU GPL3
//author: Julius Schön / R. Spicer, 2021

//this file is in development, not much here right now

//GetCritLengthPositive : todo add to calcs, check this! if this only compares 61 to 31 then its
//kinda useless, need to check the other primes rather, particularly of
//interest if the next PotPrime is NOT prime, then the crit section defaults
//to the next PotPrime that is prime.
func GetCritLengthPositive(prime, n, diff, returnHereLen *big.Int) error {
	// d( 2p + dc + 2cn ) d, diff; p, prime value; c, 30; n, N for a template number??
	if diff.Cmp(big1) < 0 {
		return errors.New("GetCritLengthPositive: difference (diff) parameter must be at greater than or equal to 1")
	}
	iCalcA.Mul(prime, big2)
	iCalcB.Mul(diff, TemplateLength)
	iCalcC.Mul(n, TemplateLength)
	iCalcC.Mul(iCalcC, big2)
	iCalcD.Add(iCalcA, iCalcB)
	iCalcD.Add(iCalcA, iCalcC)
	returnHereLen.Mul(diff, iCalcD)
	//return big.NewInt(0).Mul(diff, iCalcD)
	return nil
}

//GetCritLengthNegative : untested, under development
func GetCritLengthNegative(prime, n, diff, returnHereLen *big.Int) error {
	// d( 2p - dc + 2cn ) d, diff; p, prime value; c, 30; n, N for a template number??
	if diff.Cmp(big1) < 0 {
		return errors.New("GetCritLengthNegative: difference parameter (diff) must be at greater than or equal to 1")
	}
	iCalcA.Mul(prime, big2)
	iCalcB.Mul(diff, TemplateLength)
	iCalcC.Mul(n, TemplateLength)
	iCalcC.Mul(iCalcC, big2)
	iCalcD.Sub(iCalcA, iCalcB)
	iCalcD.Add(iCalcA, iCalcC)
	returnHereLen.Mul(diff, iCalcD)
	//return big.NewInt(0).Mul(diff, iCalcD)
	return nil
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
