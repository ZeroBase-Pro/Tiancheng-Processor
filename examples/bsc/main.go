package main

import (
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"

	"github.com/PolyhedraZK/ExpanderCompilerCollection"
	"github.com/PolyhedraZK/ExpanderCompilerCollection/test"
)

type Circuit struct {
	// timestamp var
	Pre_timestamp     frontend.Variable
	Current_timestamp frontend.Variable
	Now_time          frontend.Variable

	// block height var
	Pre_height     frontend.Variable
	Current_height frontend.Variable

	// block hash field var
	Parent_hash frontend.Variable `gnark:",public"`
	Prev_hash   frontend.Variable `gnark:",public"`

	// mix hash var
	Mix_hash frontend.Variable

	// uncle hash var
	Uncle_hash frontend.Variable

	// gas var
	Gas_limit frontend.Variable
	Gas_used  frontend.Variable

	// nonce var
	Nonce frontend.Variable

	// diffculty var
	Difficulty frontend.Variable
}

func height_check(api frontend.API, pre_height, current_height frontend.Variable) {
	expect_height := api.Add(pre_height, 1)
	api.AssertIsEqual(current_height, expect_height)
}

func timestamp_check(api frontend.API, pre_timestamp, current_timestamp, now_time frontend.Variable) {
	api.AssertIsLessOrEqual(pre_timestamp, current_timestamp)
	api.AssertIsDifferent(pre_timestamp, current_timestamp)
	api.AssertIsLessOrEqual(current_timestamp, now_time)
	api.AssertIsDifferent(current_timestamp, now_time)
}

func prevhash_check(api frontend.API, parent_hash, prev_hash frontend.Variable) {
	api.AssertIsEqual(parent_hash, prev_hash)
}

func unclehash_check(api frontend.API, uncle_hash frontend.Variable) {
	api.AssertIsEqual(uncle_hash, "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347")
}

func mixhash_check(api frontend.API, mix_hash frontend.Variable) {
	api.AssertIsEqual(mix_hash, "0x0")
}

func gaslimit_check(api frontend.API, gas_limit frontend.Variable, gas_used frontend.Variable) {
	api.AssertIsLessOrEqual(gas_limit, frontend.Variable((1<<63)-1))
	api.AssertIsLessOrEqual(gas_used, gas_limit)
}

func nonce_check(api frontend.API, nonce frontend.Variable) {
	api.AssertIsEqual(nonce, "0x0")
}

func difficulty_check(api frontend.API, difficulty frontend.Variable) {
	diffInTurn := frontend.Variable(2) // Block difficulty for in-turn signatures
	diffNoTurn := frontend.Variable(1) // Block difficulty for out-of-turn signatures
	api.AssertIsEqual(api.Mul(api.Sub(difficulty, diffInTurn), api.Sub(difficulty, diffNoTurn)), 0)
}

// Define declares the circuit's constraints
func (circuit *Circuit) Define(api frontend.API) error {
	height_check(api, circuit.Pre_height, circuit.Current_height)
	timestamp_check(api, circuit.Pre_timestamp, circuit.Current_timestamp, circuit.Now_time)
	prevhash_check(api, circuit.Parent_hash, circuit.Prev_hash)
	unclehash_check(api, circuit.Uncle_hash)
	mixhash_check(api, circuit.Mix_hash)
	gaslimit_check(api, circuit.Gas_limit, circuit.Gas_used)
	nonce_check(api, circuit.Nonce)
	difficulty_check(api, circuit.Difficulty)
	return nil
}

func main() {
	circuit, err := ExpanderCompilerCollection.Compile(ecc.BN254.ScalarField(), &Circuit{})
	if err != nil {
		panic(err)
	}

	c := circuit.GetLayeredCircuit()
	os.WriteFile("circuit.txt", c.Serialize(), 0o644)

	assignment := &Circuit{
		Pre_height:        101,
		Current_height:    102,
		Pre_timestamp:     20240825,
		Current_timestamp: 20240826,
		Now_time:          20240827,
		Parent_hash:       "0xa89c9362b8200c39ec24fadc310391e752f476fa91c407f6d3f5ff2e8bced15c",
		Prev_hash:         "0xa89c9362b8200c39ec24fadc310391e752f476fa91c407f6d3f5ff2e8bced15c",
		Uncle_hash:        "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		Mix_hash:          "0x0",
		Gas_limit:         140000000,
		Gas_used:          8486644,
		Nonce:             "0x0",
		Difficulty:        2,
	}

	inputSolver := circuit.GetInputSolver()
	witness, err := inputSolver.SolveInput(assignment, 1)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("witness.txt", witness.Serialize(), 0o644)
	if err != nil {
		panic(err)
	}

	if !test.CheckCircuit(c, witness) {
		panic("witness error")
	}
}
