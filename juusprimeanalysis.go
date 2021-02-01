package juusprime

import (
	"errors"
	"math/big"
)

//juusprime is an engine to find prime Sextuplets and/or Quintuplets and/or Quadruplets
//GNU GPL3
//author: Julius Schön / R. Spicer, 2021

//this file is in development, not much here right now

//todo add to calcs, check this! if this only compares 61 to 31 then its
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
