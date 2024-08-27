# Tiga Processor

<img src="https://github.com/ZeroBase-Pro/Tiga-Processor/blob/main/tiga.png" width="300" >


The Tiga processor is a BNB Chain blockchain hash header verification tool written using Gnark. This circuit is mainly used in scenarios where ZK cross-chain bridges need to verify the original information of blocks. Diga references PolyhedraZK's ExpanderCompilerCollection for circuit compilation.

## Workflow

Our realization of gnark-based BSC light client(or zklightclient for BSC chain) is in ``./examples/BSC/main.go``


### Complie and Output the Circuit 

To compile the circuit, and ouput the optimized layered circuit, as well as the witness, just run:

```
go run ./examples/BSC/main.go
```

You can find the ``circuit.txt`` and ``witness.txt`` in the ``./examples/BSC/`` directory.


### Prove the Circuit

To prove the circuit,it's recommanded that you put the ``circuit.txt`` and ``witness.txt`` in the ``./data`` directory. Then run:

```
RUSTFLAGS="-C target-cpu=native" cargo run --bin expander-exec --release -- prove ./data/circuit.txt ./data/witness.txt ./data/out.bin
```

It will produce the a binary called ``out.bin`` for verification.

### Verify the Circuit

To verify the circuit, in the same directory above, run:

```
RUSTFLAGS="-C target-cpu=native" cargo run --bin expander-exec --release -- verify ./data/circuit.txt ./data/witness.txt ./data/out.bin
```

## Logical Introduction(Expander Version)

os: Used for file operations.

github.com/consensys/gnark-crypto/ecc and github.com/consensys/gnark/frontend: Used for building circuits and zero-knowledge proofs.

github.com/PolyhedraZK/ExpanderCompilerCollection and github.com/PolyhedraZK/ExpanderCompilerCollection/test: Used for compiling and testing circuits.

### Define Circuit Structure

The Circuit struct contains various fields that represent information within a blockchain, such as timestamp, block height, block hash, mix hash, uncle hash, gas limit and usage, nonce, and difficulty.
Each field is of type frontend.Variable, which is used to represent variables within the circuit.

### Check Functions

height_check: Verifies that the current block height is equal to the previous block height plus one.

timestamp_check: Ensures that timestamps are in ascending order.

prevhash_check: Verifies that the parent hash of the current block matches the hash of the previous block.

unclehash_check: Checks if the uncle hash is equal to a specific value.

mixhash_check: Verifies if the mix hash is zero.

gaslimit_check: Ensures that the gas limit and usage are within a reasonable range.

nonce_check: Checks if the nonce is zero.

difficulty_check: Ensures that the block difficulty meets the rules.

### Define Circuit Constraints

Uses frontend.API to perform various constraint checks. By calling the previously defined check functions, the circuit is ensured to meet all constraint conditions.

### Main Function

Compile the Circuit: Uses ExpanderCompilerCollection.Compile to compile the Circuit circuit and saves the result to a circuit.txt file.

Set Inputs: Creates an instance of Circuit and assigns values to each field, which serves as input to the solver.

Generate Witness: Uses the solver to generate the circuit's witness and saves it to a witness.txt file.

Verify the Circuit: Uses the test function test.CheckCircuit to verify if the generated witness is correct.

## Gnark Only Version

Due to the fact that the Expander version is verified off chain, we offer only the GnRH version deployed on the BNB test network.

Contract address : 

PK&VK: /gnarkversion/key.zip


