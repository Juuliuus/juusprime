package juusprime

//juusprime is an engine to find prime Sextuplets and/or Quintuplets and/or Quadruplets
//GNU GPL3
//author: Julius Schön / R. Spicer, 2021

import (
	"fmt"
	"math/big"
	"testing"
)

//TestNaturalProgression : perhaps trivial but if something got messed up in
//generating the natural progression it indicates even wider problems and
//certain calculation errors.
func TestNaturalProgression(t *testing.T) {
	//Natural progression is "easily" calculated by hand, and so they could be
	//set by literally assigning the paper calculations. However, having the
	//primes generate it themselves provides a proof that the setup and support
	//funcs have had no inadvertent bug invasions.
	fmt.Println("testing Natural Progression...")

	testNatProg := func(natProg []int64, p *PrimeLTE29) {
		if len(natProg) != len(p.Prime.naturalProgression) {
			t.Errorf("P%v natProg length (%d) not equal to constant length (%d).\n%v\n%v : Constant", p.Prime.value, len(p.Prime.naturalProgression), len(natProg), p.Prime.naturalProgression, natProg)
		} else {
			for i := range p.Prime.naturalProgression {
				if p.Prime.naturalProgression[i].Cmp(big.NewInt(natProg[i])) != 0 {
					t.Errorf("P%v natProg does not match at index %d (counting#: %d).\n%v\n%v : Constant", p.Prime.value, i, i+1, p.Prime.naturalProgression, natProg)
					break
				}
			}
		}
	}

	var prime *PrimeLTE29

	prime = NewPrimeLTE29(big.NewInt(7))
	p7NatProg := []int64{3, 1, 6, 4, 2, 0, 5}
	//p7NatProg := []int64{3, 1, 6, 4, 2, 0, 5, 77} //too long
	testNatProg(p7NatProg, prime)

	prime = NewPrimeLTE29(big.NewInt(11))
	p11NatProg := []int64{6, 9, 1, 4, 7, 10, 2, 5, 8, 0, 3}
	//p11NatProg := []int64{6, 9, 3, 4, 7, 10, 2, 5, 8, 0} //too short
	testNatProg(p11NatProg, prime)

	prime = NewPrimeLTE29(big.NewInt(13))
	p13NatProg := []int64{11, 7, 3, 12, 8, 4, 0, 9, 5, 1, 10, 6, 2}
	//p13NatProg := []int64{11, 7, 3, 12, 8, 4, 0, 9, 5, 1, 10, 6, 999} //wrong sequence #
	testNatProg(p13NatProg, prime)

	prime = NewPrimeLTE29(big.NewInt(17))
	//prime = NewPrime29(big.NewInt(15)) //wrong prime
	p17NatProg := []int64{7, 11, 15, 2, 6, 10, 14, 1, 5, 9, 13, 0, 4, 8, 12, 16, 3}
	testNatProg(p17NatProg, prime)

	prime = NewPrimeLTE29(big.NewInt(19))
	p19NatProg := []int64{6, 14, 3, 11, 0, 8, 16, 5, 13, 2, 10, 18, 7, 15, 4, 12, 1, 9, 17}
	testNatProg(p19NatProg, prime)

	prime = NewPrimeLTE29(big.NewInt(23))
	p23NatProg := []int64{1, 17, 10, 3, 19, 12, 5, 21, 14, 7, 0, 16, 9, 2, 18, 11, 4, 20, 13, 6, 22, 15, 8}
	testNatProg(p23NatProg, prime)

	prime = NewPrimeLTE29(big.NewInt(29))
	p29NatProg := []int64{6, 5, 4, 3, 2, 1, 0, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7}
	testNatProg(p29NatProg, prime)

}

type mystruct struct {
	myint int
}

func (m *mystruct) fillNaturalProgression() {
	fmt.Println("TestInterface: I set it to 42")
	m.myint = 42
}

//TestInterface : a first attempt to see if interfaces can
//be useful in this package
func TestInterface(t *testing.T) {

	my := &mystruct{
		myint: 2,
	}

	fillNatProg(my)
	if my.myint != 42 {
		t.Errorf("my = %v", my)
	}
}

//TestFindingN : The calculation of "n" is extremely important in this
//package, it must be accurate and rounded/floor'd properly. There are
//two methods to do this an original complicated but algebraically accurate
//func and a re-worked function that uses a shortcut. This tests that These
//funcs continue to agree to high accuracy.
func TestFindingN(t *testing.T) {

	//A big TNumber with its sextuplet:
	//TNumber: 194091003877655  ┣━┫
	//---------
	//5822730116329657
	//5822730116329661
	//5822730116329663
	//5822730116329667
	//5822730116329669
	//5822730116329673

	//example: floor( (sqrt( 194091003877655 - 93 + ( 53^2 <=2809> / 30 ) ) - ( 53 / sqrt( 30 ) ))  / sqrt( 30 ))
	//givenTNum: 194091003877655, startTNum: 93, 53, 53^2
	//answer should be: 2543558
	//BTW: these 3 examples were tested, as you see, with the two GetN func's and with SpeedCrunch
	//in all cases the results were exactly the same at very high precision.
	fmt.Println("Test Getting-N using the two Get N func's")

	var check, i, j *big.Int
	var templ int64
	templ = 194091003877655
	check = big.NewInt(0).SetInt64(2543558)
	i, j = big.NewInt(0), big.NewInt(0)
	p53 := NewPrimeGTE31(big.NewInt(53))

	//GetNfromTNumComplicated(big.NewInt(templ), big.NewInt(93), big.NewInt(53), big.NewInt(2809), i)
	////j = GetNfromTNum(big.NewInt(templ), big.NewInt(53), false)
	//GetNfromTNum(big.NewInt(templ), big.NewInt(53), false, j)
	GetNfromTNumComplicated(big.NewInt(templ), p53, i)
	//j = GetNfromTNum(big.NewInt(templ), big.NewInt(53), false)
	GetNfromTNum(big.NewInt(templ), p53, j)

	if i.Cmp(check) != 0 {
		t.Errorf("Compl. GetN, Template #: %v. Wanted %v, got %v", templ, check, i)
	}
	if j.Cmp(check) != 0 {
		t.Errorf("Compl. GetN, Template #: %v. Wanted %v, got %v", templ, check, j)
	}

	check = big.NewInt(0)
	fmt.Sscan("2543560522035041559030", check)
	var bigTempl *big.Int
	bigTempl = big.NewInt(0)
	//todo use bigger primes than 53, try 53 + (100*30)
	fmt.Sscan("194091003877655194091003877655194091003877655", bigTempl)
	GetNfromTNumComplicated(getBigInt(bigTempl), p53, i)
	GetNfromTNum(getBigInt(bigTempl), p53, j)

	if i.Cmp(check) != 0 {
		t.Errorf("Compl. GetN, Template #: %v. Wanted %v, got %v", bigTempl, check, i)
	}
	if j.Cmp(check) != 0 {
		t.Errorf("Compl. GetN, Template #: %v. Wanted %v, got %v", bigTempl, check, j)
	}

	check = big.NewInt(0)
	fmt.Sscan("80434446161176326808970547259116604596006499", check)
	bigTempl = big.NewInt(0)
	fmt.Sscan("194091003877655194091003877655194091003877655194091003877655194091003877655194091003877657", bigTempl)
	GetNfromTNumComplicated(getBigInt(bigTempl), p53, i)
	GetNfromTNum(getBigInt(bigTempl), p53, j)

	if i.Cmp(check) != 0 {
		t.Errorf("Compl. GetN, Template #: %v. Wanted %v, got %v", bigTempl, check, i)
	}
	if j.Cmp(check) != 0 {
		t.Errorf("Compl. GetN, Template #: %v. Wanted %v, got %v", bigTempl, check, j)
	}

}

//TestCritLengthRoutines : tests that the crit length routines return
//the same results that a manual subtraction of effective TNumbers produces,
//for use within a PotPrime family
func TestCritLengthRoutines(t *testing.T) {
	fmt.Println("testing crit length routines...")
	//use the getcritlength functions to calculate all lengths between
	//fixed N and the n's represented by diff. Check the calculation is correct
	//using manual subtaction of corresponding effective TNumbers.
	//iter controls the looping, first we loop from 0 to N-1, then reverse the
	//looping from N-1 to 0

	iter := big.NewInt(25)
	var potprime *PrimeGTE31

	N := big.NewInt(0)
	nctrl := big.NewInt(-1)
	//diff := big.NewInt(1)
	effective := big.NewInt(0)
	test := big.NewInt(0)
	current := big.NewInt(0)
	res := big.NewInt(0)
	cmp := big.NewInt(0)

	potprimes := []int64{31, 37, 41, 43, 47, 49, 53, 59}

	for p := range potprimes {
		potprime = NewPrimeGTE31(big.NewInt(potprimes[p]))

		N.SetInt64(0)
		GetEffectiveTNum(N, potprime, current)
		nctrl.SetInt64(0)

		for nctrl.Cmp(iter) < 1 {
			nctrl.Add(nctrl, big1)
			//diff.Sub(N, nctrl)
			//diff.Abs(diff)
			//GetCritLengthPositiveWF(potprime.Prime.Value(), N, diff, res)
			GetCritLength(true, potprime, N, nctrl, res)

			GetEffectiveTNum(nctrl, potprime, effective)
			test.Sub(effective, current)
			//fmt.Println(res, test, effective)
			if test.Cmp(res) != 0 {
				t.Errorf("CritLengthPos error, Wanted %v, got %v", res, test)
			}

		}

		//fmt.Println("-------------")

		//We need to go the next N so that it backtraces
		//perfectly the ascending series we just did.
		N.Add(iter, big1)

		nctrl.Set(N)
		GetEffectiveTNum(N, potprime, current)

		cmp.Sub(N, iter)
		for nctrl.Cmp(cmp) > -1 {
			nctrl.Sub(nctrl, big1)
			if nctrl.Cmp(big0) == -1 {
				//this will never happen in this function because it
				//is controlled. This here if someone copies this code for
				//general use, this break would be important.
				break
			}
			//diff.Sub(N, nctrl)
			//diff.Abs(diff)
			//GetCritLengthNegativeWF(potprime.Prime.Value(), N, diff, res)
			GetCritLength(true, potprime, N, nctrl, res)

			GetEffectiveTNum(nctrl, potprime, effective)
			test.Sub(current, effective)
			//fmt.Println(res, test, current)
			if test.Cmp(res) != 0 {
				t.Errorf("CritLengthNeg error, Wanted %v, got %v", res, test)
			}

		}
	}
}
